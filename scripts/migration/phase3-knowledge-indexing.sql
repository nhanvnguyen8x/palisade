-- Knowledge indexing (Phase 3)
-- Stores parsed chunks + embeddings for each Knowledge Source.

CREATE TABLE IF NOT EXISTS chunks (
    id UUID PRIMARY KEY,
    knowledge_source_id UUID NOT NULL REFERENCES knowledge_sources(id),
    chunk_index INT NOT NULL,
    page_number INT,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE (knowledge_source_id, chunk_index)
);

CREATE TABLE IF NOT EXISTS embeddings (
    id UUID PRIMARY KEY,
    chunk_id UUID NOT NULL REFERENCES chunks(id),
    vector DOUBLE PRECISION[] NOT NULL,
    model VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    UNIQUE (chunk_id)
);

CREATE INDEX IF NOT EXISTS idx_chunks_ks_id ON chunks(knowledge_source_id);
CREATE INDEX IF NOT EXISTS idx_embeddings_chunk_id ON embeddings(chunk_id);

