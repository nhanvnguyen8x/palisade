package domain

import (
	"time"

	"github.com/google/uuid"
)

type Workspace struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewWorkspace(organizationID uuid.UUID, name string) *Workspace {
	now := time.Now().UTC()

	return &Workspace{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		Name:           name,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
