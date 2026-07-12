package transport

import "github.com/google/uuid"

type ChatRequest struct {
	KnowledgeBaseID uuid.UUID `json:"knowledgeBaseId" binding:"required"`
	Question        string    `json:"question" binding:"required"`
	TopK            int       `json:"topK"`
}

