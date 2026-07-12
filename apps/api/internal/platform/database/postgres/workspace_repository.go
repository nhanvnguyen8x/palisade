package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	orgdomain "github.com/nhanvnguyen8x/palisade/internal/features/organization/domain"
	"github.com/nhanvnguyen8x/palisade/internal/features/workspace/domain"
)

type WorkspaceRepository struct {
	db *pgxpool.Pool
}

func NewWorkspaceRepository(db *pgxpool.Pool) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (r *WorkspaceRepository) Create(ctx context.Context, workspace *domain.Workspace) error {
	const query = `
INSERT INTO workspaces (id, organization_id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
`

	_, err := r.db.Exec(
		ctx,
		query,
		workspace.ID,
		workspace.OrganizationID,
		workspace.Name,
		workspace.CreatedAt,
		workspace.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert workspace: %w", err)
	}

	return nil
}

func (r *WorkspaceRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Workspace, error) {
	const query = `
SELECT id, organization_id, name, created_at, updated_at
FROM workspaces
WHERE id = $1
`

	var ws domain.Workspace

	err := r.db.QueryRow(ctx, query, id).Scan(
		&ws.ID,
		&ws.OrganizationID,
		&ws.Name,
		&ws.CreatedAt,
		&ws.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrWorkspaceNotFound
		}

		return nil, fmt.Errorf("find workspace: %w", err)
	}

	return &ws, nil
}

func (r *WorkspaceRepository) OrganizationExists(ctx context.Context, id uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM organizations WHERE id = $1)`

	var exists bool

	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check organization: %w", err)
	}

	return exists, nil
}

func (r *WorkspaceRepository) FindOrganization(ctx context.Context, id uuid.UUID) (*orgdomain.Organization, error) {
	const query = `
SELECT id, name, created_at, updated_at
FROM organizations
WHERE id = $1
`

	var org orgdomain.Organization

	err := r.db.QueryRow(ctx, query, id).Scan(
		&org.ID,
		&org.Name,
		&org.CreatedAt,
		&org.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, orgdomain.ErrOrganizationNotFound
		}

		return nil, fmt.Errorf("find organization: %w", err)
	}

	return &org, nil
}
