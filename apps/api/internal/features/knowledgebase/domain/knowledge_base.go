package domain

import (
	"time"

	"github.com/google/uuid"
)

type KnowledgeBase struct {
	ID          uuid.UUID
	WorkspaceID uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewKnowledgeBase(workspaceID uuid.UUID, name, description string) *KnowledgeBase {
	now := time.Now().UTC()

	return &KnowledgeBase{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
