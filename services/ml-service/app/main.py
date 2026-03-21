import logging
from contextlib import asynccontextmanager
from fastapi import FastAPI

from app.api.routes import router as predict_router
from app.api.health import router as health_router
from app.core.model_loader import ModelLoader

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifespan context manager for startup/shutdown."""
    # Startup
    logger.info("Starting ML service")
    try:
        ModelLoader.load_model()
        logger.info("ML service ready")
    except Exception as e:
        logger.error(f"Failed to initialize ML service: {e}")
        raise
    
    yield
    
    # Shutdown
    logger.info("Shutting down ML service")
    if ModelLoader._model is not None:
        del ModelLoader._model
        ModelLoader._model = None


app = FastAPI(
    title="Mushroom ML Service",
    version="1.0.0",
    lifespan=lifespan
)

app.include_router(predict_router)
app.include_router(health_router)
