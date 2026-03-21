import io
from PIL import Image
import torch
from torchvision import transforms
from app.core.config import settings


# ImageNet normalization for EfficientNet
IMAGENET_MEAN = [0.485, 0.456, 0.406]
IMAGENET_STD = [0.229, 0.224, 0.225]


def get_transform():
    """Get preprocessing transform with correct image size."""
    return transforms.Compose([
        transforms.Resize((settings.IMAGE_SIZE, settings.IMAGE_SIZE)),
        transforms.ToTensor(),
        transforms.Normalize(mean=IMAGENET_MEAN, std=IMAGENET_STD),
    ])


def load_image_from_bytes(data: bytes) -> Image.Image:
    try:
        image = Image.open(io.BytesIO(data))
        image.verify()
        image = Image.open(io.BytesIO(data))
        image = image.convert("RGB")
        return image
    except Exception:
        raise ValueError("Invalid image file")


def preprocess_image(image: Image.Image) -> torch.Tensor:
    transform = get_transform()
    tensor = transform(image)
    tensor = tensor.unsqueeze(0)
    return tensor
