import torch
from tqdm import tqdm


def train_epoch(model, loader, optimizer, loss_fn, device):

    model.train()

    total_loss = 0

    for images, labels in tqdm(loader):

        images = images.to(device)
        labels = labels.to(device)

        optimizer.zero_grad()

        outputs = model(images)

        loss = loss_fn(outputs, labels)

        loss.backward()

        optimizer.step()

        total_loss += loss.item()

    return total_loss / len(loader