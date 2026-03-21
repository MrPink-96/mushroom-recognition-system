import torch
import numpy as np
import logging
from pathlib import Path
import timm
import torch.nn as nn

from app.core.config import settings

logger = logging.getLogger(__name__)


class ModelLoader:
    _model = None
    _classes = None
    _device = None
    _image_size = 320

    @classmethod
    def load_model(cls):
        if cls._model is not None:
            return cls._model

        logger.info("Loading ML model")

        # Determine device
        cls._device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        logger.info(f"Using device: {cls._device}")

        try:
            cls._load_classes()

            # Загружаем веса модели из файла
            checkpoint = torch.load(
                settings.MODEL_PATH,
                map_location=cls._device,
                weights_only=True
            )

            #  Обрабатываем различные форматы весов
            if isinstance(checkpoint, dict):
                if "model_state_dict" in checkpoint:
                    # Формат model_bundle.pth (веса + метаданные)
                    logger.info("Detected model_bundle format")
                    state_dict = checkpoint["model_state_dict"]
                elif "state_dict" in checkpoint:
                    # Формат с обёрткой {"state_dict": {...}}
                    logger.info("Detected state_dict wrapper format")
                    state_dict = checkpoint["state_dict"]
                elif cls._is_state_dict(checkpoint):
                    # Простой state_dict
                    logger.info("Detected raw state_dict format")
                    state_dict = checkpoint
                else:
                    logger.error(f"Unknown checkpoint keys: {list(checkpoint.keys())[:10]}")
                    raise ValueError("Unknown checkpoint format")
            else:
                state_dict = checkpoint

            model = cls._build_model_from_state_dict(state_dict)

            # Перемещаем модель на устройство (GPU/CPU)
            model = model.to(cls._device)
            # Это ускоряет работу на GPU NVIDIA
            model = model.to(memory_format=torch.channels_last)
            model.eval()

            cls._model = model
            logger.info("Model loaded successfully")

        except Exception as e:
            logger.error(f"Failed to load model: {e}")
            raise RuntimeError("Model loading failed") from e

        return cls._model

    @classmethod
    def _load_classes(cls):
        """Load class names from numpy file."""
        try:
            cls._classes = np.load(settings.CLASSES_PATH, allow_pickle=True)
            settings.NUM_CLASSES = len(cls._classes)
            logger.info(f"Loaded {len(cls._classes)} classes")
        except Exception as e:
            logger.error(f"Failed to load classes: {e}")
            raise RuntimeError("Classes loading failed") from e

    @classmethod
    def _is_state_dict(cls, checkpoint: dict) -> bool:
        """Check if dict is a raw state_dict (contains layer weight keys)."""
        sample_keys = list(checkpoint.keys())[:5]
        return any(
            "weight" in key or "bias" in key or "running_mean" in key
            for key in sample_keys
        )

    @classmethod
    def _build_model_from_state_dict(cls, state_dict):
        """Build EfficientNet-B3 model and load state dict."""
        import timm

        num_classes = settings.NUM_CLASSES

        model = timm.create_model(
            "efficientnet_b3",
            pretrained=False,
            drop_rate=0.3,
            drop_path_rate=0.3
        )

        in_features = model.classifier.in_features
        model.classifier = nn.Sequential(
            nn.Dropout(0.5),
            nn.Linear(in_features, num_classes)
        )

        # Загрузка весов
        model.load_state_dict(state_dict)
        return model

    @classmethod
    def get_model(cls):
        if cls._model is None:
            raise RuntimeError("Model is not loaded")
        return cls._model

    @classmethod
    def get_classes(cls):
        if cls._classes is None:
            raise RuntimeError("Classes not loaded")
        return cls._classes

    @classmethod
    def get_class_name(cls, idx: int) -> str:
        """Get class name by index."""
        if cls._classes is None:
            raise RuntimeError("Classes not loaded")
        if idx < 0 or idx >= len(cls._classes):
            return "Unknown"
        return str(cls._classes[idx])

    @classmethod
    def get_device(cls):
        return cls._device

    @classmethod
    def get_image_size(cls):
        return cls._image_size