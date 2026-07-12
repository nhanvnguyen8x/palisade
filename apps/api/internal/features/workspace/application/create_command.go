package application

import "github.com/google/uuid"

type CreateCommand struct {
	OrganizationID uuid.UUID
	Name           string
}
