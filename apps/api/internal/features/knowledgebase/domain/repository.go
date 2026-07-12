package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, kb *KnowledgeBase) error
	FindByID(ctx context.Context, id uuid.UUID) (*KnowledgeBase, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*KnowledgeBase, error)
}
