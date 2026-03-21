import torch
from app.core.config import settings


def get_top_predictions(logits: torch.Tensor):
    probs = torch.softmax(logits, dim=1)
    values, indices = torch.topk(probs, settings.TOP_K)
    results = []
    for idx, prob in zip(indices[0], values[0]):
        results.append(
            {
                "species_id": int(idx),
                "probability": float(prob),
            }
        )
    return results