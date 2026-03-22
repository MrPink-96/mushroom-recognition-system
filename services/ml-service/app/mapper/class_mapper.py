import json
import os

class ClassMapper:
    def __init__(self, path: str):
        self._map = self._load(path)

    def _load(self, path: str) -> dict[int, int]:
        with open(path, "r") as f:
            raw = json.load(f)

        return {int(k): int(v) for k, v in raw.items()}

    def to_species_id(self, class_id: int) -> int | None:
        return self._map.get(class_id)