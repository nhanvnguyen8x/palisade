package transport

import (
	"time"

	"github.com/google/uuid"
)

type UploadRequest struct {
	KnowledgeBaseID uuid.UUID
	FileName        string
	ContentType     string
	Size            int64
}

type KnowledgeSourceResponse struct {
	ID              string    `json:"id"`
	KnowledgeBaseID string    `json:"knowledgeBaseId"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type UploadResponse struct {
	KnowledgeSourceID string    `json:"knowledgeSourceId"`
	IngestionJobID    string    `json:"ingestionJobId"`
	DocumentID        string    `json:"documentId"`
	Status            string    `json:"status"`
	FileName          string    `json:"fileName"`
	CreatedAt         time.Time `json:"createdAt"`
}
