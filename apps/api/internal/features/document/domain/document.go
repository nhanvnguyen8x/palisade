package domain

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID                uuid.UUID
	KnowledgeSourceID uuid.UUID
	WorkspaceID       uuid.UUID
	FileName          string
	ContentType       string
	Size              int64
	StorageKey        string
	Checksum          string
	Status            Status
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewDocument(
	knowledgeSourceID uuid.UUID,
	workspaceID uuid.UUID,
	filename string,
	contentType string,
	size int64,
	storageKey string,
	checksum string,
) *Document {
	now := time.Now().UTC()

	return &Document{
		ID:                uuid.New(),
		KnowledgeSourceID: knowledgeSourceID,
		WorkspaceID:       workspaceID,
		FileName:          filename,
		ContentType:       contentType,
		Size:              size,
		StorageKey:        storageKey,
		Checksum:          checksum,
		Status:            StatusUploaded,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}
