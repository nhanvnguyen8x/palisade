package application

import (
	"fmt"

	"github.com/google/uuid"
)

func parseUUID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid uuid: %w", err)
	}

	return id, nil
}
