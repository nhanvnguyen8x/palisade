from __future__ import annotations

import uuid
from dataclasses import dataclass
from typing import Any

import psycopg

from .config import Config


@dataclass(frozen=True)
class ClaimedJob:
    job_id: uuid.UUID
    knowledge_source_id: uuid.UUID


class Postgres:
    def __init__(self, cfg: Config) -> None:
        self.cfg = cfg

    def connect(self) -> psycopg.Connection[Any]:
        # autocommit=False so we can claim/transition statuses safely.
        return psycopg.connect(self.cfg.database_url, autocommit=False)

    def claim_pending_jobs(
        self, conn: psycopg.Connection[Any], limit: int
    ) -> list[ClaimedJob]:
        sql = """
        WITH claimed AS (
            UPDATE ingestion_jobs
            SET
                status = 'PROCESSING',
                started_at = NOW(),
                updated_at = NOW()
            WHERE id IN (
                SELECT id
                FROM ingestion_jobs
                WHERE status = 'PENDING'
                ORDER BY created_at ASC
                LIMIT %s
            )
            RETURNING id, knowledge_source_id
        ),
        updated_ks AS (
            UPDATE knowledge_sources ks
            SET
                status = 'PROCESSING',
                updated_at = NOW()
            FROM claimed
            WHERE ks.id = claimed.knowledge_source_id
            RETURNING ks.id
        )
        SELECT claimed.id, claimed.knowledge_source_id
        FROM claimed
        ;
        """

        rows = conn.execute(sql, (limit,)).fetchall()
        return [
            ClaimedJob(job_id=row[0], knowledge_source_id=row[1])  # type: ignore[misc]
            for row in rows
        ]

    def fetch_document_for_source(
        self, conn: psycopg.Connection[Any], knowledge_source_id: uuid.UUID
    ) -> dict[str, Any]:
        sql = """
        SELECT
            d.workspace_id,
            d.filename,
            d.content_type,
            d.storage_key
        FROM knowledge_sources ks
        JOIN documents d ON d.knowledge_source_id = ks.id
        WHERE ks.id = %s
        LIMIT 1
        ;
        """
        row = conn.execute(sql, (knowledge_source_id,)).fetchone()
        if not row:
            raise RuntimeError(
                f"document not found for knowledge_source_id={knowledge_source_id}"
            )

        return {
            "workspace_id": row[0],
            "filename": row[1],
            "content_type": row[2],
            "storage_key": row[3],
        }

    def insert_chunks_and_embeddings(
        self,
        conn: psycopg.Connection[Any],
        knowledge_source_id: uuid.UUID,
        chunks: list[dict[str, Any]],
        vectors: list[list[float]],
        model: str,
    ) -> None:
        # Expect 1:1 mapping: chunks[i] -> vectors[i]
        if len(chunks) != len(vectors):
            raise ValueError("chunks and vectors length mismatch")

        # Dev safety:
        # - If embedding algorithm changes during development, we must avoid
        #   mixing old/new vectors for the same source.
        # - Also prevents UNIQUE conflicts on chunks(knowledge_source_id, chunk_index).
        conn.execute(
            """
            DELETE FROM embeddings
            USING chunks
            WHERE embeddings.chunk_id = chunks.id
              AND chunks.knowledge_source_id = %s
            """,
            (knowledge_source_id,),
        )

        conn.execute(
            """
            DELETE FROM chunks
            WHERE knowledge_source_id = %s
            """,
            (knowledge_source_id,),
        )

        # Insert all chunks, then embeddings. Keep it simple for now (dev).
        for i, chunk in enumerate(chunks):
            chunk_id = chunk["chunk_id"]

            conn.execute(
                """
                INSERT INTO chunks (
                    id,
                    knowledge_source_id,
                    chunk_index,
                    page_number,
                    text,
                    created_at,
                    updated_at
                )
                VALUES ($1,$2,$3,$4,$5,NOW(),NOW())
                """,
                (
                    chunk_id,
                    knowledge_source_id,
                    chunk["chunk_index"],
                    chunk.get("page_number"),
                    chunk["text"],
                ),
            )

            conn.execute(
                """
                INSERT INTO embeddings (
                    id,
                    chunk_id,
                    vector,
                    model,
                    created_at
                )
                VALUES ($1,$2,$3,$4,NOW())
                """,
                (
                    chunk["embedding_id"],
                    chunk_id,
                    vectors[i],
                    model,
                ),
            )

    def mark_completed(
        self,
        conn: psycopg.Connection[Any],
        job_id: uuid.UUID,
        knowledge_source_id: uuid.UUID,
    ) -> None:
        conn.execute(
            """
            UPDATE ingestion_jobs
            SET
                status = 'COMPLETED',
                completed_at = NOW(),
                updated_at = NOW(),
                error_message = NULL
            WHERE id = %s
            """,
            (job_id,),
        )

        conn.execute(
            """
            UPDATE knowledge_sources
            SET
                status = 'READY',
                updated_at = NOW()
            WHERE id = %s
            """,
            (knowledge_source_id,),
        )

    def mark_failed(
        self,
        conn: psycopg.Connection[Any],
        job_id: uuid.UUID,
        knowledge_source_id: uuid.UUID,
        error_message: str,
    ) -> None:
        conn.execute(
            """
            UPDATE ingestion_jobs
            SET
                status = 'FAILED',
                completed_at = NOW(),
                updated_at = NOW(),
                error_message = %s
            WHERE id = %s
            """,
            (error_message, job_id),
        )

        conn.execute(
            """
            UPDATE knowledge_sources
            SET
                status = 'FAILED',
                updated_at = NOW()
            WHERE id = %s
            """,
            (knowledge_source_id,),
        )

