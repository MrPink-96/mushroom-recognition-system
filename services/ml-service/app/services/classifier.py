import time
from PIL import Image
from app.core.model_loader import ModelLoader
from app.services.preprocessing import preprocess_image
from app.services.detection import detect_mushroom
from app.ml.inference import run_inference
from app.ml.postprocessing import get_top_predictions


def classify_image(image: Image.Image):

    start = time.time()

    detected_image = detect_mushroom(image)

    tensor = preprocess_image(detected_image)

    model = ModelLoader.get_model()

    logits = run_inference(model, tensor)

    top_predictions = get_top_predictions(logits)

    prediction = top_predictions[0]

    inference_time = (time.time() - start) * 1000

    return {
        "prediction": prediction,
        "top_predictions": top_predictions,
        "inference_time_ms": inference_time,
    }