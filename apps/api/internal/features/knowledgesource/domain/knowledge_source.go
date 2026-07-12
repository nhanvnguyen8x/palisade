package domain

import (
	"time"

	"github.com/google/uuid"
)

type KnowledgeSource struct {
	ID              uuid.UUID
	KnowledgeBaseID uuid.UUID
	Type            SourceType
	Status          Status
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewKnowledgeSource(knowledgeBaseID uuid.UUID, sourceType SourceType) *KnowledgeSource {
	now := time.Now().UTC()

	return &KnowledgeSource{
		ID:              uuid.New(),
		KnowledgeBaseID: knowledgeBaseID,
		Type:            sourceType,
		Status:          StatusUploading,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (ks *KnowledgeSource) MarkUploaded() {
	ks.Status = StatusUploaded
	ks.UpdatedAt = time.Now().UTC()
}

func (ks *KnowledgeSource) MarkProcessing() {
	ks.Status = StatusProcessing
	ks.UpdatedAt = time.Now().UTC()
}

func (ks *KnowledgeSource) MarkReady() {
	ks.Status = StatusReady
	ks.UpdatedAt = time.Now().UTC()
}

func (ks *KnowledgeSource) MarkFailed() {
	ks.Status = StatusFailed
	ks.UpdatedAt = time.Now().UTC()
}
