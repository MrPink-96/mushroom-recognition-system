from pydantic import BaseModel
from typing import List


class Prediction(BaseModel):
    species_id: int
    label: str
    probability: float


class ClassificationResponse(BaseModel):
    model_version: str
    predictions: List[Prediction]
    inference_time_ms: float
