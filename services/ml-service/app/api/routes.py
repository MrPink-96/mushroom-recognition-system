from fastapi import APIRouter, UploadFile, File

router = APIRouter()

@router.post("/classify")
async def classify(file: UploadFile = File(...)):
    return {
        "prediction": {
            "species_id": 1,
            "probability": 0.95
        },
        "top_predictions": [
            {"species_id": 1, "probability": 0.95},
            {"species_id": 2, "probability": 0.03},
            {"species_id": 3, "probability": 0.01},
            {"species_id": 4, "probability": 0.008},
            {"species_id": 5, "probability": 0.002}
        ]
    }