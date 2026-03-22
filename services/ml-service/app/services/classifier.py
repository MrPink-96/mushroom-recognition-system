import time
from PIL import Image
from typing import Dict, Any

from app.core.config import settings
from app.core.model_loader import ModelLoader
from app.services.preprocessing import preprocess_image
from app.services.detection import detect_mushroom
from app.ml.inference import run_inference
from app.ml.postprocessing import get_top_predictions
from app.mapper.class_mapper import ClassMapper

mapper = ClassMapper("mapping.json")

def predict_image(image: Image.Image) -> Dict[str, Any]:
    """
    Full prediction pipeline:
    1. Detection (stub for now)
    2. Preprocessing
    3. Inference
    4. Postprocessing
    """
    start = time.time()

    # Обнаружение и выделение гриба (пока не реализовано)
    detected_image = detect_mushroom(image)

    # Предобработка: изменение размера, нормализация
    tensor = preprocess_image(detected_image)

    # Перемещение тензора на то же устройство, где находится модель
    device = ModelLoader.get_device()
    tensor = tensor.to(device)

    # Получение модели и запуск инференса
    model = ModelLoader.get_model()
    logits = run_inference(model, tensor)

    # Получение топ-N предсказаний с названиями классов
    model_predictions = get_top_predictions(logits)

    predictions = []

    for p in model_predictions:
        class_id = p["class_id"]
        confidence = p["confidence"]
        label = p["label"]

        if class_id is None or confidence is None:
            continue

        species_id = mapper.to_species_id(class_id)

        if species_id is None:
            continue

        predictions.append({
            "species_id": species_id,
            "confidence": confidence,
            "label": label,
        })

    inference_time = (time.time() - start) * 1000

    return {
        "model_version": settings.MODEL_VERSION,
        "predictions": predictions,
        "inference_time_ms": round(inference_time, 2)
    }
