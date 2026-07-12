from __future__ import annotations

import hashlib
from typing import List

import numpy as np


class EmbeddingProvider:
    def __init__(self, dim: int) -> None:
        self.dim = dim

    def _embed_one(self, text: str) -> list[float]:
        # Hash-based deterministic embedding (no RNG):
        # - Compatible to implement the same logic in Go for retrieval.
        # - Returns a unit-normalized vector.
        vec = np.zeros((self.dim,), dtype=np.float32)

        # Use per-dimension hashes for stability.
        for i in range(self.dim):
            h = hashlib.sha256(f"{text}\0{i}".encode("utf-8")).digest()
            # Map to [0,1)
            u = int.from_bytes(h[:4], "little", signed=False)
            vec[i] = (u % 10_000_000) / 10_000_000.0

        norm = float(np.linalg.norm(vec) + 1e-12)
        vec = vec / norm
        return [float(x) for x in vec]

    def embed_texts(self, texts: list[str]) -> list[list[float]]:
        return [self._embed_one(text) for text in texts]

