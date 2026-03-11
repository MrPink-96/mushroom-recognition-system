import timm
import torch.nn as nn


def build_model(num_classes):

    model = timm.create_model(
        "efficientnet_b3",
        pretrained=True
    )

    in_features = model.classifier.in_features

    model.classifier = nn.Sequential(
        nn.Dropout(0.3),
        nn.Linear(in_features, num_classes)
    )

    return model