package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, workspace *Workspace) error
	FindByID(ctx context.Context, id uuid.UUID) (*Workspace, error)
}
