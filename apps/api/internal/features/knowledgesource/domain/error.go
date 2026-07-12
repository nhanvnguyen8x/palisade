package domain

import "errors"

var (
	ErrKnowledgeSourceNotFound = errors.New("knowledge source not found")
	ErrKnowledgeBaseNotFound = errors.New("knowledge base not found")
	ErrInvalidFile           = errors.New("invalid file")
)
