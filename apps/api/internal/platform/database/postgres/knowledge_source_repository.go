package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nhanvnguyen8x/palisade/internal/features/knowledgesource/domain"
)

type KnowledgeSourceRepository struct {
	db *pgxpool.Pool
}

func NewKnowledgeSourceRepository(db *pgxpool.Pool) *KnowledgeSourceRepository {
	return &KnowledgeSourceRepository{db: db}
}

func (r *KnowledgeSourceRepository) Create(ctx context.Context, source *domain.KnowledgeSource) error {
	const query = `
INSERT INTO knowledge_sources (id, knowledge_base_id, type, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
`

	_, err := r.db.Exec(
		ctx,
		query,
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

	return nil
}

func (r *KnowledgeSourceRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.KnowledgeSource, error) {
	const query = `
SELECT id, knowledge_base_id, type, status, created_at, updated_at
FROM knowledge_sources
WHERE id = $1
`

	var source domain.KnowledgeSource

	err := r.db.QueryRow(ctx, query, id).Scan(
		&source.ID,
		&source.KnowledgeBaseID,
		&source.Type,
		&source.Status,
		&source.CreatedAt,
		&source.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrKnowledgeSourceNotFound
		}

		return nil, fmt.Errorf("find knowledge source: %w", err)
	}

	return &source, nil
}

func (r *KnowledgeSourceRepository) Update(ctx context.Context, source *domain.KnowledgeSource) error {
	const query = `
UPDATE knowledge_sources
SET status = $2, updated_at = $3
WHERE id = $1
`

	_, err := r.db.Exec(ctx, query, source.ID, source.Status, source.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update knowledge source: %w", err)
	}

	return nil
}
