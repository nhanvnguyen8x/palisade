package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, job *IngestionJob) error
	FindByID(ctx context.Context, id uuid.UUID) (*IngestionJob, error)
	FindPending(ctx context.Context, limit int) ([]*IngestionJob, error)
	Update(ctx context.Context, job *IngestionJob) error
}
