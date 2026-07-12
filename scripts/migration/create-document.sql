CREATE TABLE documents (

                           id UUID PRIMARY KEY,

                           workspace_id UUID NOT NULL,

                           filename VARCHAR(255) NOT NULL,

                           content_type VARCHAR(128) NOT NULL,

                           size BIGINT NOT NULL,

                           checksum VARCHAR(64) NOT NULL,

                           storage_key TEXT NOT NULL,

                           status VARCHAR(32) NOT NULL,

                           created_at TIMESTAMP NOT NULL,

                           updated_at TIMESTAMP NOT NULL
);