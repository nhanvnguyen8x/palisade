package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, source *KnowledgeSource) error
	FindByID(ctx context.Context, id uuid.UUID) (*KnowledgeSource, error)
	Update(ctx context.Context, source *KnowledgeSource) error
}
