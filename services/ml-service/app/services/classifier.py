import time
from PIL import Image
from app.core.model_loader import ModelLoader
from app.services.preprocessing import preprocess_image
from app.services.detection import detect_mushroom
from app.ml.inference import run_inference
from app.ml.postprocessing import get_top_predictions


def predict_image(image: Image.Image):

    start = time.time()

    detected_image = detect_mushroom(image)

    tensor = preprocess_image(detected_image)

    model = ModelLoader.get_model()

    logits = run_inference(model, tensor)

    top_predictions = get_top_predictions(logits)


    inference_time = (time.time() - start) * 1000

    return {

    "model_version": "efficientnet-b4-v1",
        "predictions": [
    { "species_id": 12, "label": "Amanita muscaria", "probability": 0.93 },
    { "species_id": 15, "label": "Amanita muscaria","probability": 0.04 },
    { "species_id": 7,  "label": "Amanita muscaria","probability": 0.02 },
    { "species_id": 9,  "label": "Amanita muscaria","probability": 0.008 },
    { "species_id": 3,  "label": "Amanita muscaria", "probability": 0.002 }
  ],
  "inference_time_ms": 42
}