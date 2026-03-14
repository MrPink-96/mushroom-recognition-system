import os
import shutil
import random
import gc
import numpy as np
import pandas as pd
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import Dataset, DataLoader
from torchvision import transforms
import cv2
import timm
from tqdm import tqdm
from sklearn.preprocessing import LabelEncoder
import matplotlib.pyplot as plt
from pytorch_grad_cam import GradCAM
from pytorch_grad_cam.utils.model_targets import ClassifierOutputTarget
from pytorch_grad_cam.utils.image import show_cam_on_image
import torch.nn.functional as F
import torch.multiprocessing as mp
mp.set_start_method("spawn", force=True)


# -------------------- Конфигурация --------------------
DEVICE = torch.device("cuda" if torch.cuda.is_available() else "cpu")
IMAGE_SIZE = 320
BATCH_SIZE = 10
ACCUMULATION_STEPS = 2
NUM_WORKERS = 4
PIN_MEMORY = torch.cuda.is_available()
PREFETCH_FACTOR = 3 # 4 if NUM_WORKERS > 0 else None

EPOCHS_STAGE1 = 10      # было 5
EPOCHS_STAGE2 = 30      # было 15
EPOCHS_STAGE3 = 30      # было 15
PATIENCE = 10           # было 5

LR_STAGE1 = 1e-3
LR_STAGE2 = 1e-4
LR_STAGE3 = 5e-6  # было 1e-5

# Пункт 6: MixUp/CutMix 70/30 (без None)
MIXUP_CUTMIX_PROB = 0.3  # 30% CutMix, 70% MixUp
MIXUP_ALPHA = 0.4
CUTMIX_ALPHA = 1.0

EMA_DECAY = 0.9999

TRAIN_CSV = "dataset/train.csv"
VAL_CSV = "dataset/val.csv"
UNFREEZE_BLOCK_INDICES = [4, 5, 6]

torch.set_float32_matmul_precision("high")
torch.backends.cudnn.benchmark = True


def set_seed(seed=42):
    random.seed(seed)
    np.random.seed(seed)
    torch.manual_seed(seed)
    torch.cuda.manual_seed_all(seed)


set_seed(42)

def print_device_info():
    print("=" * 60)
    print("🔧 DEVICE INFO")
    print("=" * 60)

    print(f"torch.cuda.is_available(): {torch.cuda.is_available()}")
    print(f"Selected DEVICE: {DEVICE}")

    if torch.cuda.is_available():
        print(f"GPU: {torch.cuda.get_device_name(0)}")
        print(f"GPU memory: {torch.cuda.get_device_properties(0).total_memory / 1e9:.2f} GB")
    else:
        print("⚠️ CUDA GPU not available — training will run on CPU")

    print("=" * 60)



# -------------------- EMA --------------------
class ModelEMA:
    def __init__(self, model, decay=0.9999):
        self.model = model
        self.decay = decay
        self.shadow = {k: v.clone().detach() for k, v in model.state_dict().items()}
        self.backup = {}

    def update(self):
        for name, param in self.model.state_dict().items():
            self.shadow[name] = (1 - self.decay) * param.detach() + self.decay * self.shadow[name]

    def apply_shadow(self):
        self.backup = {k: v.clone() for k, v in self.model.state_dict().items()}
        self.model.load_state_dict(self.shadow)

    def restore(self):
        self.model.load_state_dict(self.backup)
        self.backup = {}


# -------------------- Dataset --------------------
class MushroomDataset(Dataset):
    def __init__(self, df, transform=None):
        self.df = df
        self.transform = transform

    def __len__(self):
        return len(self.df)

    def __getitem__(self, idx):
        img_path = self.df.iloc[idx]["image_path"]
        label = self.df.iloc[idx]["label"]
        if not os.path.exists(img_path):
            raise FileNotFoundError(f"Image not found: {img_path}")
        image = cv2.imread(img_path)
        if image is None:
            raise RuntimeError(f"Failed to load: {img_path}")
        image = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
        if self.transform:
            image = self.transform(image)
        return image, label


# -------------------- Progressive resizing --------------------
def build_train_transform(image_size):
    return transforms.Compose([
        transforms.ToPILImage(),
        transforms.RandomResizedCrop(image_size, scale=(0.6, 1.0)),  # было (0.8, 1.0)
        transforms.RandomHorizontalFlip(p=0.5),
        transforms.RandomRotation(15),
        transforms.ColorJitter(brightness=0.4, contrast=0.4, saturation=0.4), # усильте
        transforms.RandomApply([transforms.GaussianBlur(kernel_size=3)], p=0.3), # (brightness=0.2, contrast=0.2, saturation=0.2)
        transforms.ToTensor(),
        transforms.Normalize(mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225]),
        transforms.RandomErasing(p=0.5, scale=(0.02, 0.3), ratio=(0.3, 3.3))  # p=0.5 вместо 0.25
    ])


def build_val_transform(image_size):
    return transforms.Compose([
        transforms.ToPILImage(),
        transforms.Resize(int(image_size * 1.15)),
        transforms.CenterCrop(image_size),
        transforms.ToTensor(),
        transforms.Normalize(mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225])
    ])


# -------------------- MixUp + CutMix --------------------
def mixup_data(x, y, alpha=MIXUP_ALPHA):
    if alpha > 0:
        lam = np.random.beta(alpha, alpha)
    else:
        lam = 1
    batch_size = x.size(0)
    index = torch.randperm(batch_size).to(x.device)
    mixed_x = lam * x + (1 - lam) * x[index]
    return mixed_x, y, y[index], lam


def cutmix_data(x, y, alpha=CUTMIX_ALPHA):
    if alpha > 0:
        lam = np.random.beta(alpha, alpha)
    else:
        lam = 1
    batch_size = x.size(0)
    index = torch.randperm(batch_size).to(x.device)
    bbx1, bby1, bbx2, bby2 = rand_bbox(x.size(), lam)
    x = x.clone()
    x[:, :, bby1:bby2, bbx1:bbx2] = x[index, :, bby1:bby2, bbx1:bbx2]
    lam = 1 - ((bbx2 - bbx1) * (bby2 - bby1) / (x.size(2) * x.size(3)))
    return x, y, y[index], lam


def rand_bbox(size, lam):
    W, H = size[3], size[2]
    cut_rat = np.sqrt(1. - lam)
    cut_w = int(W * cut_rat)
    cut_h = int(H * cut_rat)
    cx = np.random.randint(W)
    cy = np.random.randint(H)
    bbx1 = np.clip(cx - cut_w // 2, 0, W)
    bby1 = np.clip(cy - cut_h // 2, 0, H)
    bbx2 = np.clip(cx + cut_w // 2, 0, W)
    bby2 = np.clip(cy + cut_h // 2, 0, H)
    return bbx1, bby1, bbx2, bby2


def mixup_cutmix_loss(loss_fn, pred, y_a, y_b, lam):
    return lam * loss_fn(pred, y_a) + (1 - lam) * loss_fn(pred, y_b)


# -------------------- Метрики --------------------
def accuracy_topk(output, target, topk=None):
    with torch.inference_mode():
        num_classes = output.size(1)
        if topk is None:
            topk = (1, 3) if num_classes < 20 else (1, 5)  # ← Автовыбор
        topk = tuple(min(k, num_classes) for k in topk)
        maxk = max(topk)

        _, pred = output.topk(maxk, 1, True, True)
        pred = pred.t()
        correct = pred.eq(target.view(1, -1).expand_as(pred))

        return [correct[:k].reshape(-1).float().sum(0).item() for k in topk]


# -------------------- TTA Validation (Пункт 5) --------------------
def validate_tta(model, loader, loss_fn, image_size):
    """4-view детерминированный TTA: original + hflip + scale + hflip+scale"""
    model.eval()
    total_loss = 0.0
    top1, top5 = 0, 0
    total_samples = 0

    with torch.inference_mode():
        for images, labels in tqdm(loader, desc="Validation TTA (4 views)"):
            # Пункт 8: Channels Last для батчей
            images = images.to(DEVICE, memory_format=torch.channels_last, non_blocking=True)
            labels = labels.to(DEVICE, non_blocking=True)

            batch_size = images.size(0)
            num_classes = model.classifier[-1].out_features

            preds = torch.zeros((batch_size, num_classes), device=DEVICE)

            # 1. Оригинал
            preds += model(images)
            # 2. Horizontal flip
            preds += model(images.flip(-1))

            # 3. Детерминированный scale (1.05×) + center crop
            resized = F.interpolate(images, scale_factor=1.05, mode='bilinear', align_corners=False)
            h, w = resized.shape[2], resized.shape[3]
            ch = cw = image_size
            start_h = (h - ch) // 2
            start_w = (w - cw) // 2
            scaled = resized[:, :, start_h:start_h + ch, start_w:start_w + cw]
            preds += model(scaled)
            # 4. Horizontal flip на масштабированной версии
            preds += model(scaled.flip(-1))

            preds /= 4.0

            loss = loss_fn(preds, labels)
            total_loss += loss.item()
            a1, a5 = accuracy_topk(preds, labels)
            top1 += a1
            top5 += a5
            total_samples += labels.size(0)

    return total_loss / len(loader), top1 / total_samples * 100, top5 / total_samples * 100


# -------------------- Grad-CAM --------------------
def run_gradcam(model, val_loader, stage_name, num_images=20, image_size=IMAGE_SIZE):
    model.eval()
    target_layers = [model.conv_head]
    cam = GradCAM(model=model, target_layers=target_layers)

    os.makedirs(f"gradcam_{stage_name}", exist_ok=True)
    indices = random.sample(range(len(val_loader.dataset)), min(num_images, len(val_loader.dataset)))

    for i, idx in enumerate(indices):
        try:
            img_tensor, label = val_loader.dataset[idx]
            img_tensor = img_tensor.unsqueeze(0).to(DEVICE)
            img_tensor.requires_grad = True

            targets = [ClassifierOutputTarget(label.item())]

            grayscale_cam = cam(input_tensor=img_tensor, targets=targets)
            grayscale_cam = grayscale_cam[0, :]

            img = img_tensor.squeeze(0).detach().cpu().permute(1, 2, 0).numpy()
            img = (img * np.array([0.229, 0.224, 0.225]) + np.array([0.485, 0.456, 0.406])).clip(0, 1)

            visualization = show_cam_on_image(img, grayscale_cam, use_rgb=True)

            plt.figure(figsize=(6, 6))
            plt.imshow(visualization)
            plt.title(f"Class: {label.item()}")
            plt.axis('off')
            plt.savefig(f"gradcam_{stage_name}/img_{i:03d}_class_{label.item()}.png", bbox_inches='tight', dpi=120)
            plt.close()

        except Exception as e:
            print(f"Grad-CAM error on image {i}: {e}")
            continue


# -------------------- Модель --------------------
def build_model(num_classes):
    model = timm.create_model(
        "efficientnet_b3",
        pretrained=True,
        drop_rate=0.3,        # было 0.2
        drop_path_rate=0.3    # было 0.2
    )
    in_features = model.classifier.in_features
    model.classifier = nn.Sequential(
        nn.Dropout(0.5),      # было 0.3
        nn.Linear(in_features, num_classes)
    )
    return model
"""def build_model(num_classes):
    model = timm.create_model(
        "efficientnet_b3",
        pretrained=True,
        drop_rate=0.2,
        drop_path_rate=0.2
    )
    in_features = model.classifier.in_features
    # Пункт 10: Оставляем dropout в classifier
    model.classifier = nn.Sequential(
        nn.Dropout(0.3),
        nn.Linear(in_features, num_classes)
    )
    return model
"""

# -------------------- Train epoch --------------------
def train_epoch(model, loader, optimizer, loss_fn, scaler, ema):
    model.train()
    total_loss = 0.0
    total_samples = 0
    optimizer.zero_grad(set_to_none=True)

    for i, (images, labels) in enumerate(tqdm(loader, desc="Training")):
        images = images.to(DEVICE, memory_format=torch.channels_last, non_blocking=True)
        labels = labels.to(DEVICE, non_blocking=True)

        # 70% MixUp / 30% CutMix
        if random.random() < 0.3:
            images, y_a, y_b, lam = cutmix_data(images, labels)
        else:
            images, y_a, y_b, lam = mixup_data(images, labels)

        with torch.autocast(device_type="cuda", dtype=torch.float16):
            outputs = model(images)
            loss_raw = mixup_cutmix_loss(loss_fn, outputs, y_a, y_b, lam)
            loss = loss_raw / ACCUMULATION_STEPS

        scaler.scale(loss).backward()

        if (i + 1) % ACCUMULATION_STEPS == 0 or (i + 1) == len(loader):
            scaler.unscale_(optimizer)
            torch.nn.utils.clip_grad_norm_(model.parameters(), 1.0)
            scaler.step(optimizer)
            scaler.update()
            optimizer.zero_grad(set_to_none=True)
            if ema is not None:
                ema.update()

        total_loss += loss_raw.item() * images.size(0)
        total_samples += images.size(0)

    return total_loss / total_samples


# -------------------- Validate (без TTA) --------------------
def validate(model, loader, loss_fn):
    model.eval()
    total_loss = 0.0
    top1, top5 = 0, 0
    total_samples = 0
    with torch.inference_mode():
        for images, labels in tqdm(loader, desc="Validation"):
            # Пункт 8: Channels Last для батчей
            images = images.to(DEVICE, memory_format=torch.channels_last)
            labels = labels.to(DEVICE)
            outputs = model(images)
            loss = loss_fn(outputs, labels)
            total_loss += loss.item()
            a1, a5 = accuracy_topk(outputs, labels)
            top1 += a1
            top5 += a5
            total_samples += labels.size(0)
    return total_loss / len(loader), top1 / total_samples * 100, top5 / total_samples * 100

########################


def test_spatial_bias(model, val_loader, device):
    """Проверка: смотрит ли модель на объект или на координаты"""
    model.eval()

    correct_count = 0
    total_count = 0

    with torch.inference_mode():
        for batch_idx, (img, _) in enumerate(val_loader):
            if batch_idx >= 5:  # Проверяем 5 батчей
                break

            img = img.to(device)
            preds = model(img).softmax(dim=1)
            classes = preds.argmax(dim=1)

            img_flipped = img.flip(-1)
            preds_f = model(img_flipped).softmax(dim=1)
            classes_f = preds_f.argmax(dim=1)

            matches = (classes == classes_f).sum().item()
            correct_count += matches
            total_count += img.size(0)

        accuracy = correct_count / total_count * 100
        print(f"\n🔍 Spatial Bias Test ({total_count} images):")
        print(f"  ✅ Инвариантность к flip: {accuracy:.1f}%")
        if accuracy < 90:
            print("  ⚠️  Модель частично зависит от положения!")
        else:
            print("  ✅ Отлично, модель инвариантна к положению")

# -------------------- АДАПТИВНЫЕ НАСТРОЙКИ ДЛЯ 10-169 КЛАССОВ --------------------
def get_adaptive_config(train_df, num_classes_target=10):
    """Автоматически подбирает настройки под размер датасета"""

    # Фильтрация классов (если нужно меньше 169)
    all_classes = train_df["label"].unique()
    if len(all_classes) > num_classes_target:
        selected = random.sample(list(all_classes), num_classes_target)
        train_df = train_df[train_df["label"].isin(selected)].reset_index(drop=True)
        print(f"🎯 Выбрано {num_classes_target} классов из {len(all_classes)}")
    else:
        selected = list(all_classes)
        print(f"🎯 Все {len(all_classes)} классов будут использованы")

    # Статистика
    class_counts = train_df["label"].value_counts()
    avg_per_class = len(train_df) / len(selected)
    min_per_class = class_counts.min()
    max_per_class = class_counts.max()
    imbalance_ratio = max_per_class / min_per_class

    print(
        f"📊 Статистика: avg={avg_per_class:.0f}, min={min_per_class}, max={max_per_class}, imbalance={imbalance_ratio:.2f}x")

    # 🔥 АДАПТИВНЫЕ ЭПОХИ
    if avg_per_class >= 1500:
        epochs = {"s1": 15, "s2": 40, "s3": 60}
        patience = 20
        lr_stage3 = 2e-6
    elif avg_per_class >= 800:
        epochs = {"s1": 12, "s2": 35, "s3": 45}
        patience = 15
        lr_stage3 = 3e-6
    elif avg_per_class >= 400:
        epochs = {"s1": 10, "s2": 30, "s3": 35}
        patience = 12
        lr_stage3 = 5e-6
    else:  # 250-400
        epochs = {"s1": 10, "s2": 30, "s3": 30}
        patience = 10
        lr_stage3 = 5e-6

    # 🔥 АДАПТИВНЫЙ BATCH_SIZE (проверка GPU)
    if torch.cuda.is_available():
        gpu_mem = torch.cuda.get_device_properties(0).total_memory / 1e9
        batch_size = 32 if gpu_mem > 20 else (16 if gpu_mem > 12 else 8)
    else:
        batch_size = 8

    # 🔥 УСИЛЕНИЕ WEIGHTS при сильном дисбалансе
    if imbalance_ratio > 5:
        weight_boost = 1.5  # Усиливаем веса редких классов
    else:
        weight_boost = 1.0

    return {
        "selected_classes": selected,
        "train_df": train_df,
        "epochs": epochs,
        "patience": patience,
        "lr_stage3": lr_stage3,
        "batch_size": batch_size,
        "weight_boost": weight_boost,
        "avg_per_class": avg_per_class,
        "imbalance_ratio": imbalance_ratio
    }


def create_loaders(train_ds, val_ds, batch_size, num_workers, pin_memory, prefetch_factor):
    return (
        DataLoader(train_ds, batch_size=batch_size, shuffle=True, num_workers=num_workers,
                   pin_memory=pin_memory,  pin_memory_device="cuda" if pin_memory else "", prefetch_factor=prefetch_factor, drop_last=True,
                   persistent_workers=True if num_workers > 0 else False),
        DataLoader(val_ds, batch_size=batch_size, shuffle=False, num_workers=num_workers,
                   pin_memory=pin_memory,  pin_memory_device="cuda" if pin_memory else "", prefetch_factor=prefetch_factor,
                   persistent_workers=True if num_workers > 0 else False,
                   drop_last=False)
    )

########################

# -------------------- Main --------------------
def main():
    print_device_info()
    for f in [
        "best_model_stage1.pth",
        "best_model_stage2.pth",
        "best_model_final.pth",
    ]:
        if os.path.exists(f):
            os.remove(f)

    # Очистка старых Grad-CAM
    if os.path.exists("gradcam_stage1"):
        shutil.rmtree("gradcam_stage1")
    if os.path.exists("gradcam_stage2"):
        shutil.rmtree("gradcam_stage2")
    if os.path.exists("gradcam_stage3"):
        shutil.rmtree("gradcam_stage3")

    global IMAGE_SIZE, BATCH_SIZE, EPOCHS_STAGE1, EPOCHS_STAGE2, EPOCHS_STAGE3, PATIENCE, LR_STAGE3

    # Загрузка данных
    train_df = pd.read_csv(TRAIN_CSV)
    val_df = pd.read_csv(VAL_CSV)

    # 🔍 ПРОВЕРКА ИСХОДНЫХ CSV ФАЙЛОВ
    print("\n" + "=" * 60)
    print("📊 АНАЛИЗ ИСХОДНЫХ CSV ФАЙЛОВ")
    print("=" * 60)
    print(f"✅ Train: {len(train_df)} изображений, {train_df['label'].nunique()} классов")
    print(f"✅ Val: {len(val_df)} изображений, {val_df['label'].nunique()} классов")

    # Проверка перекрытия классов
    train_classes = set(train_df['label'].unique())
    val_classes = set(val_df['label'].unique())
    common = train_classes & val_classes
    print(f"✅ Общие классы: {len(common)} из {len(train_classes | val_classes)}")
    print("=" * 60 + "\n")

    # 🔥 РЕШЕНИЕ: Использовать исходный сплит (ваши CSV уже разделены!)
    USE_ORIGINAL_SPLIT = True  # ← True для production, False только для экспериментов

    # 🔥 ПОЛУЧАЕМ АДАПТИВНЫЕ НАСТРОЙКИ
    # Укажите желаемое количество классов: 10, 50, 100 или 169
    num_classes_target = 10 # ← ИЗМЕНИТЕ ЭТО ЗНАЧЕНИЕ ПРИ НЕОБХОДИМОСТИ
    config = get_adaptive_config(train_df, num_classes_target=num_classes_target)

    # 🔥 Quick test override (если нужно)
    QUICK_TEST = False #True  # ← False для полноценного обучения
    # 🔥 SUPER QUICK TEST: подвыборка для мгновенной проверки
    SUPER_QUICK = False #True  # ← True для теста за 5-10 минут

    if QUICK_TEST:
        config["epochs"] = {"s1": 1, "s2": 2, "s3": 2}
        config["patience"] = 2
        print("⚡ QUICK TEST MODE")

    # Применяем настройки из config
    train_df = config["train_df"]
    val_df = val_df[val_df["label"].isin(config["selected_classes"])].reset_index(drop=True)
    selected_classes = config["selected_classes"]

    # Обновляем глобальные гиперпараметры
    EPOCHS_STAGE1 = config["epochs"]["s1"]
    EPOCHS_STAGE2 = config["epochs"]["s2"]
    EPOCHS_STAGE3 = config["epochs"]["s3"]
    PATIENCE = config["patience"]
    LR_STAGE3 = config["lr_stage3"]
    BATCH_SIZE = config["batch_size"]

    print(f"\n⚙️ Применены адаптивные настройки:")
    print(f"   📊 Классов: {len(selected_classes)}")
    print(f"   📈 Epochs: S1={EPOCHS_STAGE1}, S2={EPOCHS_STAGE2}, S3={EPOCHS_STAGE3}")
    print(f"   ⏱️  PATIENCE={PATIENCE}")
    print(f"   🎯 LR_STAGE3={LR_STAGE3}")
    print(f"   📦 BATCH_SIZE={BATCH_SIZE}")
    print(f"   ⚖️  Weight boost: {config['weight_boost']}x (imbalance={config['imbalance_ratio']:.2f}x)")

    # 🔥 СТРАТИФИЦИРОВАННЫЙ СПЛИТ (ТОЛЬКО ЕСЛИ USE_ORIGINAL_SPLIT = False)
    if not USE_ORIGINAL_SPLIT:
        print("\n🔄 Создаём стратифицированный сплит...")
        from sklearn.model_selection import train_test_split
        val_ratio = 0.15
        train_list, val_list = [], []
        for class_name in selected_classes:
            class_train = train_df[train_df["label"] == class_name]
            class_val = val_df[val_df["label"] == class_name]
            class_all = pd.concat([class_train, class_val], ignore_index=True)
            if len(class_all) >= 10:
                tr, vl = train_test_split(
                    class_all, test_size=val_ratio, random_state=42,
                    stratify=class_all["label"] if len(class_all["label"].unique()) > 1 else None
                )
                train_list.append(tr)
                val_list.append(vl)
            else:
                train_list.append(class_all)
                val_list.append(class_all.sample(min(3, len(class_all)), random_state=42, replace=len(class_all) < 3))
        train_df = pd.concat(train_list, ignore_index=True)
        val_df = pd.concat(val_list, ignore_index=True)
        print(f"✅ Stratified split: Train={len(train_df)}, Val={len(val_df)}")
    else:
        print("\n✅ Используем исходное разделение из CSV файлов")
        print(f"   Train: {len(train_df)}, Val: {len(val_df)}")


    if SUPER_QUICK and QUICK_TEST:
        sample_frac = 0.005
        train_df = train_df.sample(frac=sample_frac, random_state=42).reset_index(drop=True)
        val_df = val_df.sample(frac=sample_frac, random_state=42).reset_index(drop=True)
        print(f"⚡ SUPER QUICK: {sample_frac * 100:.0f}% данных → {len(train_df)} train, {len(val_df)} val")

    # 🔍 Валидация датасета перед обучением
    print("\n" + "=" * 60)
    print("🔍 DATASET VALIDATION")
    print("=" * 60)

    unique_classes = train_df["label"].nunique()
    print(f"✅ Количество классов в train: {unique_classes}")

    class_counts = train_df["label"].value_counts()
    min_count = class_counts.min()
    max_count = class_counts.max()
    ratio = max_count / min_count
    print(f"✅ Мин. изображений/класс: {min_count}")
    print(f"✅ Макс. изображений/класс: {max_count}")
    print(f"✅ Дисбаланс: {ratio:.2f}x {'✅ OK' if ratio < 3 else '⚠️ Высокий'}")

    # Проверка путей к файлам (выборочно)
    missing = 0
    for _, row in train_df.head(20).iterrows():
        if not os.path.exists(row["image_path"]):
            missing += 1
    if missing > 0:
        print(f"⚠️  Не найдено {missing} изображений из первых 20")
    else:
        print("✅ Все пути к изображениям валидны")

    print("=" * 60 + "\n")

    # -------------------- LabelEncoder --------------------
    le = LabelEncoder()
    train_df["label"] = le.fit_transform(train_df["label"])
    val_df["label"] = le.transform(val_df["label"])
    num_classes = len(le.classes_)
    np.save("classes.npy", le.classes_)

    # -------------------- Class weights с усилением при дисбалансе --------------------
    class_counts = train_df["label"].value_counts().sort_index()
    class_weights = 1.0 / class_counts

    # 🔥 УСИЛЕНИЕ ПРИ ДИСБАЛАНСЕ
    if config["imbalance_ratio"] > 3:
        class_weights = class_weights ** config["weight_boost"]
        class_weights = class_weights / class_weights.mean()

    # Создаём тензор весов
    class_weights_tensor = torch.zeros(num_classes, device=DEVICE)
    for idx, w in class_weights.items():
        class_weights_tensor[int(idx)] = w

    print(f"✅ Class weights: min={class_weights_tensor.min():.4f}, max={class_weights_tensor.max():.4f}")
    class_weights_tensor = class_weights_tensor.to(DEVICE)

    # -------------------- Dataset и DataLoader --------------------
    train_dataset = MushroomDataset(train_df, transform=None)
    val_dataset = MushroomDataset(val_df, transform=None)

    # -------------------- Модель и Loss --------------------
    model = build_model(num_classes).to(DEVICE)
    model = model.to(memory_format=torch.channels_last)

    loss_fn = nn.CrossEntropyLoss(
        label_smoothing=0.02,
        weight=class_weights_tensor
    )

    best_val_acc = 0.0
    patience_counter = 0

    print(f"Train samples: {len(train_dataset)}")
    print(f"Val samples: {len(val_dataset)}")
    print(f"Class weights range: [{class_weights_tensor.min():.3f}, {class_weights_tensor.max():.3f}]")

    # ==================== STAGE 1 ====================
    print("\n=== STAGE 1: Classifier only ===")
    IMAGE_SIZE = 224

    train_dataset.transform = build_train_transform(IMAGE_SIZE)
    val_dataset.transform = build_val_transform(IMAGE_SIZE)

    for param in model.parameters():
        param.requires_grad = False
    for param in model.classifier.parameters():
        param.requires_grad = True

    optimizer = optim.AdamW(
        filter(lambda p: p.requires_grad, model.parameters()),
        lr=LR_STAGE1,
        weight_decay=1e-4
    )
    scheduler = optim.lr_scheduler.CosineAnnealingWarmRestarts(
        optimizer,
        T_0=EPOCHS_STAGE1,
        eta_min=1e-6
    )

    scaler = torch.amp.GradScaler("cuda", enabled=(DEVICE.type == "cuda"))
    ema = ModelEMA(model, decay=EMA_DECAY)

    # 🔥 Пересоздаём DataLoader
    train_loader, val_loader = create_loaders(train_dataset, val_dataset, BATCH_SIZE, NUM_WORKERS, PIN_MEMORY,
                                              PREFETCH_FACTOR)

    for epoch in range(EPOCHS_STAGE1):
        train_loss = train_epoch(model, train_loader, optimizer, loss_fn, scaler, ema)
        scheduler.step()

        val_loss, top1, top5 = validate(model, val_loader, loss_fn)
        ema.apply_shadow()
        _, ema_top1, ema_top5 = validate_tta(model, val_loader, loss_fn, IMAGE_SIZE)
        ema.restore()

        print(f"Stage1 {epoch + 1}/{EPOCHS_STAGE1} | Train Loss: {train_loss:.4f} | "
              f"Val Loss: {val_loss:.4f} | "
              f"Top-1 (no TTA): {top1:.2f}%  Top-5 (no TTA): {top5:.2f}% | "
              f"EMA + TTA Top-1: {ema_top1:.2f}%  EMA + TTA Top-5: {ema_top5:.2f}%")

        if ema_top1 > best_val_acc:
            best_val_acc = ema_top1
            ema.apply_shadow()
            torch.save(model.state_dict(), "best_model_stage1.pth")
            ema.restore()
            patience_counter = 0
        else:
            patience_counter += 1
            if patience_counter >= PATIENCE:
                print("Early stopping in Stage 1")
                break

    print("\nGenerating Grad-CAM visualizations for Stage 1...")
    run_gradcam(model, val_loader, "stage1", num_images=20, image_size=IMAGE_SIZE)

    gc.collect()
    torch.cuda.empty_cache()

    # ==================== STAGE 2 ====================

    best_val_acc = 0
    patience_counter = 0
    print("\n=== STAGE 2: Unfreeze top blocks ===")
    IMAGE_SIZE = 288

    train_dataset.transform = build_train_transform(IMAGE_SIZE)
    val_dataset.transform = build_val_transform(IMAGE_SIZE)

    for param in model.parameters():
        param.requires_grad = False
    for param in model.classifier.parameters():
        param.requires_grad = True
    for idx in UNFREEZE_BLOCK_INDICES:
        if idx < len(model.blocks):
            for param in model.blocks[idx].parameters():
                param.requires_grad = True

    if os.path.exists("best_model_stage1.pth"):
        model.load_state_dict(torch.load("best_model_stage1.pth", map_location=DEVICE))
        model = model.to(memory_format=torch.channels_last)

    ema = ModelEMA(model, decay=EMA_DECAY)

    optimizer = optim.AdamW(
        filter(lambda p: p.requires_grad, model.parameters()),
        lr=LR_STAGE2,
        weight_decay=1e-4
    )
    scheduler = optim.lr_scheduler.CosineAnnealingWarmRestarts(
        optimizer,
        T_0=EPOCHS_STAGE2 // 2, # EPOCHS_STAGE2,
        eta_min=1e-6
    )

    scaler = torch.amp.GradScaler("cuda", enabled=(DEVICE.type == "cuda"))


    # 🔥 Пересоздаём DataLoader
    train_loader, val_loader = create_loaders(train_dataset, val_dataset, BATCH_SIZE, NUM_WORKERS, PIN_MEMORY,
                                              PREFETCH_FACTOR)

    for epoch in range(EPOCHS_STAGE2):
        train_loss = train_epoch(model, train_loader, optimizer, loss_fn, scaler, ema)
        scheduler.step()

        val_loss, top1, top5 = validate(model, val_loader, loss_fn)
        ema.apply_shadow()
        _, ema_top1, ema_top5 = validate_tta(model, val_loader, loss_fn, IMAGE_SIZE)
        ema.restore()

        print(f"Stage2 {epoch + 1}/{EPOCHS_STAGE2} | Train Loss: {train_loss:.4f} | "
              f"Val Loss: {val_loss:.4f} | "
              f"Top-1 (no TTA): {top1:.2f}%  Top-5 (no TTA): {top5:.2f}% | "
              f"EMA + TTA Top-1: {ema_top1:.2f}%  EMA + TTA Top-5: {ema_top5:.2f}%")

        if ema_top1 > best_val_acc:
            best_val_acc = ema_top1
            ema.apply_shadow()
            torch.save(model.state_dict(), "best_model_stage2.pth")
            ema.restore()
            patience_counter = 0
        else:
            patience_counter += 1
            if patience_counter >= PATIENCE:
                print("Early stopping in Stage 2")
                break

    print("\nGenerating Grad-CAM visualizations for Stage 2...")
    run_gradcam(model, val_loader, "stage2", num_images=20, image_size=IMAGE_SIZE)

    gc.collect()
    torch.cuda.empty_cache()

    # ==================== STAGE 3 ====================

    best_val_acc = 0
    patience_counter = 0
    print("\n=== STAGE 3: Full fine-tuning ===")
    IMAGE_SIZE = 320

    train_dataset.transform = build_train_transform(IMAGE_SIZE)
    val_dataset.transform = build_val_transform(IMAGE_SIZE)

    for param in model.parameters():
        param.requires_grad = True

    if os.path.exists("best_model_stage2.pth"):
        model.load_state_dict(torch.load("best_model_stage2.pth", map_location=DEVICE))
        model = model.to(memory_format=torch.channels_last)

    ema = ModelEMA(model, decay=EMA_DECAY)

    optimizer = optim.AdamW(
        model.parameters(),
        lr=LR_STAGE3,
        weight_decay=1e-4
    )
    scheduler = optim.lr_scheduler.CosineAnnealingWarmRestarts(
        optimizer,
        #T_0=EPOCHS_STAGE3 // 2, # EPOCHS_STAGE3,
        T_0=10,
        T_mult=2,
        eta_min=1e-6
    )

    scaler = torch.amp.GradScaler("cuda", enabled=(DEVICE.type == "cuda"))


    # 🔥 Пересоздаём DataLoader
    train_loader, val_loader = create_loaders(train_dataset, val_dataset, BATCH_SIZE, NUM_WORKERS, PIN_MEMORY,
                                              PREFETCH_FACTOR)

    for epoch in range(EPOCHS_STAGE3):
        train_loss = train_epoch(model, train_loader, optimizer, loss_fn, scaler, ema)
        scheduler.step()

        val_loss, top1, top5 = validate(model, val_loader, loss_fn)
        ema.apply_shadow()
        _, ema_top1, ema_top5 = validate_tta(model, val_loader, loss_fn, IMAGE_SIZE)
        ema.restore()

        print(f"Stage3 {epoch + 1}/{EPOCHS_STAGE3} | Train Loss: {train_loss:.4f} | "
              f"Val Loss: {val_loss:.4f} | "
              f"Top-1 (no TTA): {top1:.2f}%  Top-5 (no TTA): {top5:.2f}% | "
              f"EMA + TTA Top-1: {ema_top1:.2f}%  EMA + TTA Top-5: {ema_top5:.2f}%")

        if ema_top1 > best_val_acc:
            best_val_acc = ema_top1
            ema.apply_shadow()
            torch.save(model.state_dict(), "best_model_final.pth")
            ema.restore()
            patience_counter = 0
        else:
            patience_counter += 1
            if patience_counter >= PATIENCE:
                print("Early stopping in Stage 3")
                break

    print("\nGenerating Grad-CAM visualizations for Stage 3...")
    run_gradcam(model, val_loader, "stage3", num_images=20, image_size=IMAGE_SIZE)

    # Финальное сохранение
    ema.apply_shadow()
    torch.save(model.state_dict(), "best_model_ema_final.pth")
    torch.save({
        "model_state_dict": model.state_dict(),
        "classes": le.classes_,
        "image_size": IMAGE_SIZE
    }, "model_bundle.pth")
    ema.restore()

    test_spatial_bias(model, val_loader, DEVICE)

    print(f"\n{'=' * 60}")
    print(f"✅ ОБУЧЕНИЕ ЗАВЕРШЕНО!")
    print(f"{'=' * 60}")
    print(f"📊 Лучшая EMA + TTA Top-1: {best_val_acc:.2f}%")
    print(f"\n📁 Файлы сохранены в текущей директории:")
    print(f"   • best_model_stage1.pth")
    print(f"   • best_model_stage2.pth")
    print(f"   • best_model_final.pth")
    print(f"   • best_model_ema_final.pth")
    print(f"   • model_bundle.pth")
    print(f"   • classes.npy")
    print(f"   • gradcam_stage1/, gradcam_stage2/, gradcam_stage3/")
    print(f"{'=' * 60}")

    gc.collect()
    torch.cuda.empty_cache()



if __name__ == "__main__":
    main()


