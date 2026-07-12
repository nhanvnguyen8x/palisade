package application

import (
	"time"

	"github.com/google/uuid"
)

type UploadResult struct {
	ID uuid.UUID

	Status string

	FileName string

	CreatedAt time.Time
}
