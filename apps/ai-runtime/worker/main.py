from __future__ import annotations

import time
import uuid

from .chunking import chunk_pages
from .config import Config
from .db import Postgres
from .embedding import EmbeddingProvider
from .milvus_store import MilvusStore
from .minio_client import MinioStorage
from .pdf_parser import extract_text_by_page


def _safe_truncate(s: str, max_len: int = 2000) -> str:
    if len(s) <= max_len:
        return s
    return s[:max_len]


def process_one_job(
    cfg: Config,
    pg: Postgres,
    conn,
    storage: MinioStorage,
    embedder: EmbeddingProvider,
    milvus: MilvusStore | None,
    job_id: uuid.UUID,
    knowledge_source_id: uuid.UUID,
) -> None:
    # One job per DB transaction for data consistency.
    doc = pg.fetch_document_for_source(conn, knowledge_source_id)
    pdf_bytes = storage.download_bytes(cfg.minio_bucket, doc["storage_key"])

    pages = extract_text_by_page(pdf_bytes)
    chunks = chunk_pages(pages, cfg.chunk_max_chars)

    if chunks:
        vectors = embedder.embed_texts([c["text"] for c in chunks])
    else:
        vectors = []

    model_name = "hash-embedding"
    pg.insert_chunks_and_embeddings(
        conn=conn,
        knowledge_source_id=knowledge_source_id,
        chunks=chunks,
        vectors=vectors,
        model=model_name,
    )

    if cfg.milvus_enabled and chunks and milvus is not None:
        milvus.insert_vectors(
            chunk_ids=[c["chunk_id"] for c in chunks],
            vectors=vectors,
        )

    pg.mark_completed(conn, job_id=job_id, knowledge_source_id=knowledge_source_id)


def main() -> None:
    cfg = Config()
    pg = Postgres(cfg)
    storage = MinioStorage(cfg)
    embedder = EmbeddingProvider(dim=cfg.embedding_dim)

    milvus = None
    if cfg.milvus_enabled:
        milvus = MilvusStore(
            host=cfg.milvus_host,
            port=cfg.milvus_port,
            collection_name=cfg.milvus_collection,
            dim=cfg.embedding_dim,
        )

    # Keep one connection; if it fails, reconnect.
    conn = pg.connect()

    while True:
        try:
            claimed_jobs: list[tuple[uuid.UUID, uuid.UUID]] = []
            # Claim jobs in a single transaction (atomic).
            claimed = pg.claim_pending_jobs(conn, cfg.batch_size)
            conn.commit()

            for j in claimed:
                claimed_jobs.append((j.job_id, j.knowledge_source_id))

            if not claimed_jobs:
                time.sleep(cfg.poll_interval_seconds)
                continue

            for job_id, knowledge_source_id in claimed_jobs:
                try:
                    # Start a new transaction for inserts + mark completed.
                    process_conn = conn
                    process_one_job(
                        cfg=cfg,
                        pg=pg,
                        conn=process_conn,
                        storage=storage,
                        embedder=embedder,
                        milvus=milvus,
                        job_id=job_id,
                        knowledge_source_id=knowledge_source_id,
                    )
                    conn.commit()
                except Exception as e:
                    conn.rollback()
                    err = _safe_truncate(str(e))
                    pg.mark_failed(
                        conn=conn,
                        job_id=job_id,
                        knowledge_source_id=knowledge_source_id,
                        error_message=err,
                    )
                    conn.commit()
        except Exception:
            try:
                conn.rollback()
            except Exception:
                pass
            # Reconnect after unexpected connection issues.
            try:
                conn.close()
            except Exception:
                pass
            conn = pg.connect()


if __name__ == "__main__":
    main()

