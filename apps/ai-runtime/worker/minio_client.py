from __future__ import annotations

import io
from dataclasses import dataclass

from minio import Minio

from .config import Config


@dataclass(frozen=True)
class ObjectRef:
    key: str


class MinioStorage:
    def __init__(self, cfg: Config) -> None:
        self.cfg = cfg
        self.client = Minio(
            cfg.minio_endpoint,
            access_key=cfg.minio_access_key,
            secret_key=cfg.minio_secret_key,
            secure=cfg.minio_use_ssl,
        )

    def download_bytes(self, bucket: str, key: str) -> bytes:
        # MinIO get_object returns an object-like stream.
        response = self.client.get_object(bucket, key)
        try:
            data = response.read()
            return data
        finally:
            response.close()
            response.release_conn()

