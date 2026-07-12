package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nhanvnguyen8x/palisade/internal/features/ingestionjob/domain"
)

type IngestionJobRepository struct {
	db *pgxpool.Pool
}

func NewIngestionJobRepository(db *pgxpool.Pool) *IngestionJobRepository {
	return &IngestionJobRepository{db: db}
}

func (r *IngestionJobRepository) Create(ctx context.Context, job *domain.IngestionJob) error {
	const query = `
INSERT INTO ingestion_jobs (
    id, knowledge_source_id, status, error_message,
    started_at, completed_at, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

	_, err := r.db.Exec(
		ctx,
		query,
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

	return nil
}

func (r *IngestionJobRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.IngestionJob, error) {
	const query = `
SELECT id, knowledge_source_id, status, error_message, started_at, completed_at, created_at, updated_at
FROM ingestion_jobs
WHERE id = $1
`

	var job domain.IngestionJob

	err := r.db.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.KnowledgeSourceID,
		&job.Status,
		&job.ErrorMessage,
		&job.StartedAt,
		&job.CompletedAt,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrIngestionJobNotFound
		}

		return nil, fmt.Errorf("find ingestion job: %w", err)
	}

	return &job, nil
}

func (r *IngestionJobRepository) FindPending(ctx context.Context, limit int) ([]*domain.IngestionJob, error) {
	if limit <= 0 {
		limit = 10
	}

	const query = `
SELECT id, knowledge_source_id, status, error_message, started_at, completed_at, created_at, updated_at
FROM ingestion_jobs
WHERE status = $1
ORDER BY created_at ASC
LIMIT $2
`

	rows, err := r.db.Query(ctx, query, domain.JobStatusPending, limit)
	if err != nil {
		return nil, fmt.Errorf("find pending jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*domain.IngestionJob

	for rows.Next() {
		var job domain.IngestionJob

		if err := rows.Scan(
			&job.ID,
			&job.KnowledgeSourceID,
			&job.Status,
			&job.ErrorMessage,
			&job.StartedAt,
			&job.CompletedAt,
			&job.CreatedAt,
			&job.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan ingestion job: %w", err)
		}

		jobs = append(jobs, &job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate ingestion jobs: %w", err)
	}

	return jobs, nil
}

func (r *IngestionJobRepository) Update(ctx context.Context, job *domain.IngestionJob) error {
	const query = `
UPDATE ingestion_jobs
SET status = $2, error_message = $3, started_at = $4, completed_at = $5, updated_at = $6
WHERE id = $1
`

	_, err := r.db.Exec(
		ctx,
		query,
		job.ID,
		job.Status,
		job.ErrorMessage,
		job.StartedAt,
		job.CompletedAt,
		job.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update ingestion job: %w", err)
	}

	return nil
}
