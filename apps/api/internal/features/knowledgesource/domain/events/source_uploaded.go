package events

import (
	"time"

	"github.com/google/uuid"
)

type SourceUploaded struct {
	KnowledgeSourceID uuid.UUID
	KnowledgeBaseID   uuid.UUID
	WorkspaceID       uuid.UUID
	Occurred          time.Time
}

func (e SourceUploaded) Name() string {
	return "knowledge_source.uploaded"
}

func (e SourceUploaded) OccurredAt() time.Time {
	return e.Occurred
}
