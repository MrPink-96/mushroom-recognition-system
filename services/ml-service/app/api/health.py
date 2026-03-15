from fastapi import APIRouter
from app.core.model_loader import ModelLoader

router = APIRouter()


@router.get("/health")
def health():
    model_loaded = ModelLoader._model is not None
    return {
        "status": "ok",
        "model_loaded": model_loaded,
    }