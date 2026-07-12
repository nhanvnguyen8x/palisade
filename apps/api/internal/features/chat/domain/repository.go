package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// ListChunkEmbeddings returns all chunk vectors for a knowledge base.
	// The service will rank by cosine similarity in-memory for this dev build.
	ListChunkEmbeddings(
		ctx context.Context,
		knowledgeBaseID uuid.UUID,
		model string,
	) ([]ChunkEmbedding, error)
}

