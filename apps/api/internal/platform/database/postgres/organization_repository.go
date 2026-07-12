package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nhanvnguyen8x/palisade/internal/features/organization/domain"
)

type OrganizationRepository struct {
	db *pgxpool.Pool
}

func NewOrganizationRepository(db *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	const query = `
INSERT INTO organizations (id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4)
`

	_, err := r.db.Exec(ctx, query, org.ID, org.Name, org.CreatedAt, org.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert organization: %w", err)
	}

	return nil
}

func (r *OrganizationRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	const query = `
SELECT id, name, created_at, updated_at
FROM organizations
WHERE id = $1
`

	var org domain.Organization

	err := r.db.QueryRow(ctx, query, id).Scan(
		&org.ID,
		&org.Name,
		&org.CreatedAt,
		&org.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrganizationNotFound
		}

		return nil, fmt.Errorf("find organization: %w", err)
	}

	return &org, nil
}
