import torch
from typing import List, Dict, Any
from app.core.config import settings
from app.core.model_loader import ModelLoader


def get_top_predictions(logits: torch.Tensor) -> List[Dict[str, Any]]:
    """
    Convert model logits to top-k predictions with class labels.
    """
    probs = torch.softmax(logits, dim=1)
    values, indices = torch.topk(probs, settings.TOP_K)

    results = []
    for idx, prob in zip(indices[0], values[0]):
        class_idx = int(idx)
        label = ModelLoader.get_class_name(class_idx)
        results.append({
            "class_id": class_idx,
            "label": label,
            "confidence": round(float(prob), 6),
        })

    return results
