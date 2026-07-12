package application

import (
	"io"

	"github.com/google/uuid"
)

type UploadCommand struct {
	KnowledgeSourceID uuid.UUID
	WorkspaceID       uuid.UUID

	FileName string

	ContentType string

	Size int64

	Reader io.Reader
}
