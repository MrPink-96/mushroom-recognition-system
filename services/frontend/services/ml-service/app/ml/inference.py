import torch

def run_inference(model, tensor: torch.Tensor):
    with torch.no_grad():
        logits = model(tensor)
    return logits