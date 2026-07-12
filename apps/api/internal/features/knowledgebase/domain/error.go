package domain

import "errors"

var (
	ErrKnowledgeBaseNotFound = errors.New("knowledge base not found")
	ErrInvalidName           = errors.New("invalid knowledge base name")
	ErrWorkspaceNotFound     = errors.New("workspace not found")
)
