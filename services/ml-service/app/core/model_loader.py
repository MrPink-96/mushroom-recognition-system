import torch
import numpy as np
import logging
from pathlib import Path

from app.core.config import settings

logger = logging.getLogger(__name__)


class ModelLoader:
    _model = None
    _classes = None
    _device = None
    
    @classmethod
    def load_model(cls):
        if cls._model is not None:
            return cls._model
        
        logger.info("Loading ML model")
        
        # Determine device
        cls._device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        logger.info(f"Using device: {cls._device}")
        
        try:
            # Load classes first
            cls._load_classes()
            
            # Load model checkpoint
            checkpoint = torch.load(
                settings.MODEL_PATH, 
                map_location=cls._device,
                weights_only=False
            )
            
            # Handle different checkpoint formats
            if isinstance(checkpoint, dict):
                if "model" in checkpoint:
                    # Full model saved
                    model = checkpoint["model"]
                elif "state_dict" in checkpoint:
                    # State dict wrapped in checkpoint
                    model = cls._build_model_from_state_dict(checkpoint["state_dict"])
                elif cls._is_state_dict(checkpoint):
                    # Raw state dict (OrderedDict with layer weights)
                    logger.info("Detected raw state_dict format")
                    model = cls._build_model_from_state_dict(checkpoint)
                else:
                    # Log keys for debugging
                    logger.error(f"Unknown checkpoint keys: {list(checkpoint.keys())[:10]}")
                    raise ValueError("Unknown checkpoint format")
            else:
                model = checkpoint
            
            model = model.to(cls._device)
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
        # State dict keys typically look like: "conv_stem.weight", "bn1.weight", "classifier.weight"
        sample_keys = list(checkpoint.keys())[:5]
        return any(
            "weight" in key or "bias" in key or "running_mean" in key
            for key in sample_keys
        )
    
    @classmethod
    def _build_model_from_state_dict(cls, state_dict):
        """Build EfficientNet-B4 model and load state dict."""
        import timm
        
        num_classes = settings.NUM_CLASSES
        model = timm.create_model(
            "efficientnet_b4",
            pretrained=False,
            num_classes=num_classes
        )
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
