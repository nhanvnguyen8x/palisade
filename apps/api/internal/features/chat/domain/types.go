package domain

import "github.com/google/uuid"

type RetrievedChunk struct {
	ChunkID string  `json:"chunkID"`
	Text    string  `json:"text"`
	Score   float64 `json:"score"`
}

type ChunkEmbedding struct {
	ChunkID uuid.UUID
	Text    string
	Vector  []float64
}

