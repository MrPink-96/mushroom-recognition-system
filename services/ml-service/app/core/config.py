from pydantic import BaseSettings


class Settings(BaseSettings):

    MODEL_PATH: str = "models/mushroom_classifier.pt"

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