import logging
from fastapi import FastAPI
from app.api.routes import router as classify_router
from app.api.health import router as health_router
from app.core.model_loader import ModelLoader

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="Mushroom ML Service",
    version="1.0.0"
)


@app.on_event("startup")
def startup():
    logger.info("Starting ML service")
    ModelLoader.load_model()

@app.on_event("shutdown")
def shutdown():
    logger.info("Shutting down ML service")
    model = ModelLoader._model
    if model is not None:
        del model
        ModelLoader._model = None

app.include_router(classify_router)
app.include_router(health_router)