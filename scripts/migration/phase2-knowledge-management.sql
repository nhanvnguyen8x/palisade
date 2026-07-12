CREATE TABLE organizations (
    id         UUID PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE workspaces (
    id              UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    name            VARCHAR(255) NOT NULL,
    created_at      TIMESTAMP NOT NULL,
    updated_at      TIMESTAMP NOT NULL
);

CREATE TABLE knowledge_bases (
    id           UUID PRIMARY KEY,
    workspace_id UUID NOT NULL REFERENCES workspaces(id),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    created_at   TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP NOT NULL
);

CREATE TABLE knowledge_sources (
    id                UUID PRIMARY KEY,
    knowledge_base_id UUID NOT NULL REFERENCES knowledge_bases(id),
    type              VARCHAR(32) NOT NULL,
    status            VARCHAR(32) NOT NULL,
    created_at        TIMESTAMP NOT NULL,
    updated_at        TIMESTAMP NOT NULL
);

ALTER TABLE documents
    ADD COLUMN IF NOT EXISTS knowledge_source_id UUID REFERENCES knowledge_sources(id);

CREATE TABLE ingestion_jobs (
    id                  UUID PRIMARY KEY,
    knowledge_source_id UUID NOT NULL REFERENCES knowledge_sources(id),
    status              VARCHAR(32) NOT NULL,
    error_message       TEXT,
    started_at          TIMESTAMP,
    completed_at        TIMESTAMP,
    created_at          TIMESTAMP NOT NULL,
    updated_at          TIMESTAMP NOT NULL
);

CREATE INDEX idx_workspaces_org_id ON workspaces(organization_id);
CREATE INDEX idx_knowledge_bases_workspace_id ON knowledge_bases(workspace_id);
CREATE INDEX idx_knowledge_sources_kb_id ON knowledge_sources(knowledge_base_id);
CREATE INDEX idx_documents_knowledge_source_id ON documents(knowledge_source_id);
CREATE INDEX idx_ingestion_jobs_status ON ingestion_jobs(status);
CREATE INDEX idx_ingestion_jobs_ks_id ON ingestion_jobs(knowledge_source_id);
