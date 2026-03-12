from fastapi import APIRouter, UploadFile, File, HTTPException
from app.core.config import settings
from app.services.preprocessing import load_image_from_bytes
from app.services.classifier import classify_image


router = APIRouter()


async def read_upload(file: UploadFile) -> bytes:
    if file.content_type not in settings.ALLOWED_IMAGE_TYPES:
        raise HTTPException(415, "Unsupported media type")
    data = await file.read()
    if len(data) > settings.MAX_FILE_SIZE:
        raise HTTPException(413, "File too large")

    return data


@router.post("/classify")
async def classify(file: UploadFile = File(...)):
    data = await read_upload(file)
    image = load_image_from_bytes(data)

    if image.width > settings.MAX_IMAGE_WIDTH:
        raise HTTPException(400, "Image width too large")
    if image.height > settings.MAX_IMAGE_HEIGHT:
        raise HTTPException(400, "Image height too large")
    result = classify_image(image)

    return result