from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    # Model paths
    MODEL_PATH: str = "models/best_model_ema_final.pth"
    CLASSES_PATH: str = "models/classes.npy"

    # Model config
    MODEL_VERSION: str = "efficientnet-b3-v1"
    IMAGE_SIZE: int = 320
    NUM_CLASSES: int = 0  # Будет задан из classes.npy

    # File limits
    MAX_FILE_SIZE: int = 100 * 1024 * 1024  # 100 MB
    MAX_IMAGE_WIDTH: int = 4096
    MAX_IMAGE_HEIGHT: int = 4096

    ALLOWED_IMAGE_TYPES: set[str] = {
        "image/jpeg",
        "image/png",
        "image/jpg",
    }

    TOP_K: int = 5


settings = Settings()
