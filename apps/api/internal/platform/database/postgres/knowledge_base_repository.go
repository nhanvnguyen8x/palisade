package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nhanvnguyen8x/palisade/internal/features/knowledgebase/domain"
)

type KnowledgeBaseRepository struct {
	db *pgxpool.Pool
}

func NewKnowledgeBaseRepository(db *pgxpool.Pool) *KnowledgeBaseRepository {
	return &KnowledgeBaseRepository{db: db}
}

func (r *KnowledgeBaseRepository) Create(ctx context.Context, kb *domain.KnowledgeBase) error {
	const query = `
INSERT INTO knowledge_bases (id, workspace_id, name, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
`

	_, err := r.db.Exec(
		ctx,
		query,
		kb.ID,
		kb.WorkspaceID,
		kb.Name,
		kb.Description,
		kb.CreatedAt,
		kb.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert knowledge base: %w", err)
	}

	return nil
}

func (r *KnowledgeBaseRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.KnowledgeBase, error) {
	const query = `
SELECT id, workspace_id, name, description, created_at, updated_at
FROM knowledge_bases
WHERE id = $1
`

	var kb domain.KnowledgeBase

	err := r.db.QueryRow(ctx, query, id).Scan(
		&kb.ID,
		&kb.WorkspaceID,
		&kb.Name,
		&kb.Description,
		&kb.CreatedAt,
		&kb.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrKnowledgeBaseNotFound
		}

		return nil, fmt.Errorf("find knowledge base: %w", err)
	}

	return &kb, nil
}

func (r *KnowledgeBaseRepository) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*domain.KnowledgeBase, error) {
	const query = `
SELECT id, workspace_id, name, description, created_at, updated_at
FROM knowledge_bases
WHERE workspace_id = $1
ORDER BY created_at DESC
`

	rows, err := r.db.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("list knowledge bases: %w", err)
	}
	defer rows.Close()

	var bases []*domain.KnowledgeBase

	for rows.Next() {
		var kb domain.KnowledgeBase

		if err := rows.Scan(
			&kb.ID,
			&kb.WorkspaceID,
			&kb.Name,
			&kb.Description,
			&kb.CreatedAt,
			&kb.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan knowledge base: %w", err)
		}

		bases = append(bases, &kb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate knowledge bases: %w", err)
	}

	return bases, nil
}

func (r *KnowledgeBaseRepository) WorkspaceExists(ctx context.Context, id uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM workspaces WHERE id = $1)`

	var exists bool

	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check workspace: %w", err)
	}

	return exists, nil
}
