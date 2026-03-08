from fastapi import FastAPI

from app.api.routes import router as classify_router
from app.api.health import router as health_router

app = FastAPI(
    title="Mushroom ML Service",
    version="1.0.0"
)

app.include_router(classify_router)
app.include_router(health_router)