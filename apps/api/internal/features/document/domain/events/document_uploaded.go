package events

import (
	"time"

	"github.com/google/uuid"
)

type DocumentUploaded struct {
	DocumentID uuid.UUID

	WorkspaceID uuid.UUID

	Occurred time.Time
}

func (e DocumentUploaded) Name() string {
	return "document.uploaded"
}

func (e DocumentUploaded) OccurredAt() time.Time {
	return e.Occurred
}
