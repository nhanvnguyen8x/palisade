package domain

import (
	"time"

	"github.com/google/uuid"
)

type IngestionJob struct {
	ID                uuid.UUID
	KnowledgeSourceID uuid.UUID
	Status            JobStatus
	ErrorMessage      *string
	StartedAt         *time.Time
	CompletedAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewIngestionJob(knowledgeSourceID uuid.UUID) *IngestionJob {
	now := time.Now().UTC()

	return &IngestionJob{
		ID:                uuid.New(),
		KnowledgeSourceID: knowledgeSourceID,
		Status:            JobStatusPending,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func (j *IngestionJob) MarkProcessing() {
	now := time.Now().UTC()
	j.Status = JobStatusProcessing
	j.StartedAt = &now
	j.UpdatedAt = now
}

func (j *IngestionJob) MarkCompleted() {
	now := time.Now().UTC()
	j.Status = JobStatusCompleted
	j.CompletedAt = &now
	j.UpdatedAt = now
}

func (j *IngestionJob) MarkFailed(message string) {
	now := time.Now().UTC()
	j.Status = JobStatusFailed
	j.ErrorMessage = &message
	j.CompletedAt = &now
	j.UpdatedAt = now
}
