from torchvision import transforms
from config import IMAGE_SIZE


def get_train_transforms():

    return transforms.Compose([
        transforms.ToPILImage(),

        transforms.RandomResizedCrop(
            IMAGE_SIZE,
            scale=(0.8, 1.0)
        ),

        transforms.RandomHorizontalFlip(),

        transforms.RandomRotation(15),

        transforms.ColorJitter(
            brightness=0.2,
            contrast=0.2,
            saturation=0.2
        ),

        transforms.GaussianBlur(3),

        transforms.ToTensor(),

        transforms.Normalize(
            mean=[0.485,0.456,0.406],
            std=[0.229,0.224,0.225]
        )
    ])


def get_val_transforms():

    return transforms.Compose([
        transforms.ToPILImage(),

        transforms.Resize((IMAGE_SIZE, IMAGE_SIZE)),

        transforms.ToTensor(),

        transforms.Normalize(
            mean=[0.485,0.456,0.406],
            std=[0.229,0.224,0.225]
        )
    ])