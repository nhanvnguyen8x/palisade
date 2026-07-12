package application

import (
	"time"

	"github.com/google/uuid"
)

type UploadResult struct {
	KnowledgeSourceID uuid.UUID
	IngestionJobID    uuid.UUID
	DocumentID        uuid.UUID
	Status            string
	FileName          string
	CreatedAt         time.Time
}
