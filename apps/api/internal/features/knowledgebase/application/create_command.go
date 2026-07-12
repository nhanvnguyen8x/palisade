package application

import "github.com/google/uuid"

type CreateCommand struct {
	WorkspaceID uuid.UUID
	Name        string
	Description string
}
