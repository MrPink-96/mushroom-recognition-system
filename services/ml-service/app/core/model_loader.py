import torch
import logging

from app.core.config import settings

logger = logging.getLogger(__name__)


class ModelLoader:
    _model = None
    @classmethod
    def load_model(cls):
        if cls._model is not None:
            return cls._model
        logger.info("Loading ML model")

        try:
            model = torch.load(settings.MODEL_PATH, map_location="cpu")
            model.eval()
            cls._model = model
            logger.info("Model loaded successfully")
        except Exception as e:
            logger.error("Failed to load model")
            raise RuntimeError("Model loading failed") from e

        return cls._model

    @classmethod
    def get_model(cls):
        if cls._model is None:
            raise RuntimeError("Model is not loaded")
        return cls._model