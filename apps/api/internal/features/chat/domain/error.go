package domain

import "errors"

var (
	ErrKnowledgeBaseNotFound = errors.New("knowledge base not found")
	ErrNoRelevantChunks      = errors.New("no relevant chunks")
)

