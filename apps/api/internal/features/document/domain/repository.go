package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, document *Document) error
	Update(ctx context.Context, document *Document) error
	FindByID(ctx context.Context, id uuid.UUID) (*Document, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
