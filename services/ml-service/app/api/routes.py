from fastapi import APIRouter, UploadFile, File, HTTPException
import logging

from app.core.config import settings
from app.services.preprocessing import load_image_from_bytes
from app.services.classifier import predict_image

logger = logging.getLogger(__name__)
router = APIRouter()


async def read_upload(file: UploadFile) -> bytes:
    """Validate and read uploaded file."""
    if file.content_type not in settings.ALLOWED_IMAGE_TYPES:
        raise HTTPException(415, "Unsupported media type")

    data = await file.read()

    if len(data) > settings.MAX_FILE_SIZE:
        raise HTTPException(413, "File too large")

    return data


@router.post("/predict")
async def predict(file: UploadFile = File(...)):
    """
    Predict mushroom species from uploaded image.
    Returns top-k predictions with probabilities.
    """
    data = await read_upload(file)

    try:
        image = load_image_from_bytes(data)
    except ValueError as e:
        raise HTTPException(400, str(e))

    if image.width > settings.MAX_IMAGE_WIDTH:
        raise HTTPException(400, "Image width too large")
    if image.height > settings.MAX_IMAGE_HEIGHT:
        raise HTTPException(400, "Image height too large")

    try:
        result = predict_image(image)
    except RuntimeError as e:
        logger.error(f"Prediction failed: {e}")
        raise HTTPException(503, "Model not available")
    except Exception as e:
        logger.error(f"Unexpected error during prediction: {e}")
        raise HTTPException(500, "Prediction failed")

    return result
