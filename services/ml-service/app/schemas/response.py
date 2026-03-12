from pydantic import BaseModel


class Prediction(BaseModel):
    species_id: int
    probability: float


class ClassificationResponse(BaseModel):
    prediction: Prediction
    top_predictions: list[Prediction]
    inference_time_ms: float