package transport

import "github.com/nhanvnguyen8x/palisade/internal/features/chat/domain"

type ChatResponse struct {
	Answer   string              `json:"answer"`
	Contexts []domain.RetrievedChunk `json:"contexts"`
}

