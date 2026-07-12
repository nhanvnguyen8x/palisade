from __future__ import annotations

import uuid
from typing import Iterable


class MilvusStore:
    def __init__(
        self,
        host: str,
        port: int,
        collection_name: str,
        dim: int,
    ) -> None:
        self.host = host
        self.port = port
        self.collection_name = collection_name
        self.dim = dim

        try:
            from pymilvus import (
                Collection,
                CollectionSchema,
                DataType,
                FieldSchema,
                utility,
                connections,
            )

            self._pymilvus = {
                "Collection": Collection,
                "CollectionSchema": CollectionSchema,
                "DataType": DataType,
                "FieldSchema": FieldSchema,
                "utility": utility,
                "connections": connections,
            }
        except Exception as e:  # pragma: no cover
            raise RuntimeError(
                "pymilvus is required for MILVUS_ENABLED=true. Install requirements.txt"
            ) from e

        self._connected = False
        self._collection = None

    def _connect(self) -> None:
        if self._connected:
            return
        self._pymilvus["connections"].connect("default", host=self.host, port=self.port)
        self._connected = True

    def ensure_collection(self) -> None:
        self._connect()
        utility = self._pymilvus["utility"]
        Collection = self._pymilvus["Collection"]
        CollectionSchema = self._pymilvus["CollectionSchema"]
        DataType = self._pymilvus["DataType"]
        FieldSchema = self._pymilvus["FieldSchema"]

        if utility.has_collection(self.collection_name):
            self._collection = Collection(self.collection_name)
            return

        id_field = FieldSchema(
            name="id",
            dtype=DataType.INT64,
            is_primary=True,
            auto_id=False,
        )
        vector_field = FieldSchema(
            name="embedding",
            dtype=DataType.FLOAT_VECTOR,
            dim=self.dim,
        )
        schema = CollectionSchema(fields=[id_field, vector_field], description="knowledge embeddings")

        self._collection = Collection(
            name=self.collection_name,
            schema=schema,
        )

    def _chunk_id_to_int64(self, chunk_id: uuid.UUID) -> int:
        # Make it deterministic and fit into INT64 range.
        return int.from_bytes(chunk_id.bytes[:8], "little", signed=False) % (2**63 - 1)

    def insert_vectors(
        self,
        chunk_ids: Iterable[uuid.UUID],
        vectors: list[list[float]],
    ) -> None:
        if not self._collection:
            self.ensure_collection()

        ids = [self._chunk_id_to_int64(cid) for cid in chunk_ids]
        self._collection.insert([ids, vectors])
        self._collection.flush()

