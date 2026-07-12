package transport

import "time"

type IngestionJobResponse struct {
	ID                string     `json:"id"`
	KnowledgeSourceID string     `json:"knowledgeSourceId"`
	Status            string     `json:"status"`
	ErrorMessage      *string    `json:"errorMessage,omitempty"`
	StartedAt         *time.Time `json:"startedAt,omitempty"`
	CompletedAt       *time.Time `json:"completedAt,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}
