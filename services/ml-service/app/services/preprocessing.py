import io
from PIL import Image
import torch
from torchvision import transforms


transform = transforms.Compose(
    [
        transforms.Resize((224, 224)),
        transforms.ToTensor(),
    ]
)


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
    tensor = transform(image)
    tensor = tensor.unsqueeze(0)
    return tensor