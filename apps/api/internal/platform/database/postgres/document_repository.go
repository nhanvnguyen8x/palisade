package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	docdomain "github.com/nhanvnguyen8x/palisade/internal/features/document/domain"
	jobdomain "github.com/nhanvnguyen8x/palisade/internal/features/ingestionjob/domain"
	ksdomain "github.com/nhanvnguyen8x/palisade/internal/features/knowledgesource/domain"
)

type DocumentRepository struct {
	db *pgxpool.Pool
}

func NewDocumentRepository(db *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(ctx context.Context, document *docdomain.Document) error {
	const query = `
INSERT INTO documents (
    id,
    knowledge_source_id,
    workspace_id,
    filename,
    content_type,
    size,
    checksum,
    storage_key,
    status,
    created_at,
    updated_at
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
`

	_, err := r.db.Exec(
		ctx,
		query,
		document.ID,
		document.KnowledgeSourceID,
		document.WorkspaceID,
		document.FileName,
		document.ContentType,
		document.Size,
		document.Checksum,
		document.StorageKey,
		document.Status,
		document.CreatedAt,
		document.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert document: %w", err)
	}

	return nil
}

func (r *DocumentRepository) FindByID(ctx context.Context, id uuid.UUID) (*docdomain.Document, error) {
	const query = `
SELECT
    id,
    knowledge_source_id,
    workspace_id,
    filename,
    content_type,
    size,
    checksum,
    storage_key,
    status,
    created_at,
    updated_at
FROM documents
WHERE id = $1
`

	var d docdomain.Document

	err := r.db.QueryRow(ctx, query, id).Scan(
		&d.ID,
		&d.KnowledgeSourceID,
		&d.WorkspaceID,
		&d.FileName,
		&d.ContentType,
		&d.Size,
		&d.Checksum,
		&d.StorageKey,
		&d.Status,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, docdomain.ErrDocumentNotFound
		}

		return nil, fmt.Errorf("find document: %w", err)
	}

	return &d, nil
}

func (r *DocumentRepository) Update(ctx context.Context, document *docdomain.Document) error {
	const query = `
UPDATE documents SET
    filename = $2,
    content_type = $3,
    size = $4,
    checksum = $5,
    storage_key = $6,
    status = $7,
    updated_at = $8
WHERE id = $1
`

	_, err := r.db.Exec(
		ctx,
		query,
		document.ID,
		document.FileName,
		document.ContentType,
		document.Size,
		document.Checksum,
		document.StorageKey,
		document.Status,
		document.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update document: %w", err)
	}

	return nil
}

func (r *DocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM documents WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}

	return nil
}

func SaveKnowledgeSourceUpload(
	ctx context.Context,
	pool *pgxpool.Pool,
	source *ksdomain.KnowledgeSource,
	document *docdomain.Document,
	job *jobdomain.IngestionJob,
) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const sourceQuery = `
INSERT INTO knowledge_sources (id, knowledge_base_id, type, status, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6)
`
	_, err = tx.Exec(
		ctx,
		sourceQuery,
		source.ID,
		source.KnowledgeBaseID,
		source.Type,
		source.Status,
		source.CreatedAt,
		source.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert knowledge source: %w", err)
	}

	const docQuery = `
INSERT INTO documents (
    id, knowledge_source_id, workspace_id, filename, content_type,
    size, checksum, storage_key, status, created_at, updated_at
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
`
	_, err = tx.Exec(
		ctx,
		docQuery,
		document.ID,
		document.KnowledgeSourceID,
		document.WorkspaceID,
		document.FileName,
		document.ContentType,
		document.Size,
		document.Checksum,
		document.StorageKey,
		document.Status,
		document.CreatedAt,
		document.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert document: %w", err)
	}

	const jobQuery = `
INSERT INTO ingestion_jobs (
    id, knowledge_source_id, status, error_message,
    started_at, completed_at, created_at, updated_at
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
`
	_, err = tx.Exec(
		ctx,
		jobQuery,
		job.ID,
		job.KnowledgeSourceID,
		job.Status,
		job.ErrorMessage,
		job.StartedAt,
		job.CompletedAt,
		job.CreatedAt,
		job.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert ingestion job: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
