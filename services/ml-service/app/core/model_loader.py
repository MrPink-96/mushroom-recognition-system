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
    _model_version = None
    
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
        sample_keys = list(checkpoint.keys())[:5]
        return any(
            "weight" in key or "bias" in key or "running_mean" in key
            for key in sample_keys
        )
    
    @classmethod
    def _detect_efficientnet_variant(cls, state_dict) -> str:
        """Detect EfficientNet variant from conv_stem.weight shape."""
        # EfficientNet conv_stem output channels by variant:
        # B0: 32, B1: 32, B2: 32, B3: 40, B4: 48, B5: 48
        variant_map = {
            32: "efficientnet_b0",  # Could be B0/B1/B2, default to B0
            40: "efficientnet_b3",
            48: "efficientnet_b4",  # Could be B4/B5, default to B4
        }
        
        conv_stem_key = "conv_stem.weight"
        if conv_stem_key in state_dict:
            out_channels = state_dict[conv_stem_key].shape[0]
            variant = variant_map.get(out_channels, "efficientnet_b0")
            logger.info(f"Detected {variant} from conv_stem shape: {state_dict[conv_stem_key].shape}")
            return variant
        
        logger.warning("Could not detect variant, defaulting to efficientnet_b0")
        return "efficientnet_b0"
    
    @classmethod
    def _build_model_from_state_dict(cls, state_dict):
        """Build EfficientNet model matching the state_dict architecture."""
        import timm
        import torch.nn as nn
        
        num_classes = settings.NUM_CLASSES
        
        # Auto-detect variant from weights
        variant = cls._detect_efficientnet_variant(state_dict)
        
        # Create base model
        model = timm.create_model(
            variant,
            pretrained=False,
            num_classes=num_classes
        )
        
        # Check if classifier has custom structure (Sequential with Dropout + Linear)
        # This is indicated by keys like "classifier.1.weight" instead of "classifier.weight"
        has_sequential_classifier = "classifier.1.weight" in state_dict
        
        if has_sequential_classifier:
            logger.info("Detected Sequential classifier (Dropout + Linear)")
            in_features = model.classifier.in_features
            model.classifier = nn.Sequential(
                nn.Dropout(0.5),
                nn.Linear(in_features, num_classes)
            )
        
        # Load state dict
        model.load_state_dict(state_dict)
        
        # Save version for later
        cls._model_version = variant
        logger.info(f"Model built: {variant}, classes: {num_classes}")
        
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
    def get_model_version(cls) -> str:
        return cls._model_version or "unknown"
