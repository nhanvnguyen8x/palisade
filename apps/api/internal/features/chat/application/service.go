package application

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/nhanvnguyen8x/palisade/internal/features/chat/domain"
	"github.com/nhanvnguyen8x/palisade/internal/platform/ai/embedding"
)

type Service struct {
	repository domain.Repository
	embeddingDim int
	modelName    string
}

func NewService(repository domain.Repository, embeddingDim int, modelName string) *Service {
	return &Service{
		repository:    repository,
		embeddingDim: embeddingDim,
		modelName:    modelName,
	}
}

type ChatRequest struct {
	KnowledgeBaseID uuid.UUID
	Question        string
	TopK            int
}

type ChatAnswer struct {
	Answer   string `json:"answer"`
	Contexts []domain.RetrievedChunk
}

func (s *Service) Chat(ctx context.Context, req ChatRequest) (*ChatAnswer, error) {
	question := strings.TrimSpace(req.Question)
	if question == "" {
		return nil, fmt.Errorf("question is required")
	}

	topK := req.TopK
	if topK <= 0 {
		topK = 5
	}

	queryVec := embedding.HashEmbedding(question, s.embeddingDim)

	chunks, err := s.repository.ListChunkEmbeddings(ctx, req.KnowledgeBaseID, s.modelName)
	if err != nil {
		return nil, err
	}

	type scored struct {
		chunkID string
		text    string
		score   float64
	}
	scoredChunks := make([]scored, 0, len(chunks))

	for _, c := range chunks {
		// Assume both vectors are normalized -> dot product = cosine similarity.
		if len(c.Vector) == 0 {
			continue
		}
		score := embedding.CosineSimilarity(queryVec, c.Vector)
		scoredChunks = append(scoredChunks, scored{
			chunkID: c.ChunkID.String(),
			text:    c.Text,
			score:   score,
		})
	}

	if len(scoredChunks) == 0 {
		return nil, domain.ErrNoRelevantChunks
	}

	sort.Slice(scoredChunks, func(i, j int) bool {
		return scoredChunks[i].score > scoredChunks[j].score
	})

	if topK > len(scoredChunks) {
		topK = len(scoredChunks)
	}
	scoredChunks = scoredChunks[:topK]

	contextParts := make([]string, 0, len(scoredChunks))
	contexts := make([]domain.RetrievedChunk, 0, len(scoredChunks))
	totalChars := 0
	maxContextChars := 4000

	for _, c := range scoredChunks {
		snippet := c.text
		if len(snippet) > 1200 {
			snippet = snippet[:1200] + "..."
		}

		contextParts = append(contextParts, snippet)
		contexts = append(contexts, domain.RetrievedChunk{
			ChunkID: c.chunkID,
			Text:    snippet,
			Score:   c.score,
		})

		totalChars += len(snippet)
		if totalChars >= maxContextChars {
			break
		}
	}

	contextText := strings.Join(contextParts, "\n\n---\n\n")

	// Dev-only "LLM":
	// Return top contexts as an answer skeleton. Later we can plug a real LLM.
	answer := "Dựa trên tài liệu bạn cung cấp, mình thấy nội dung liên quan như sau:\n\n" + contextText

	return &ChatAnswer{
		Answer:   answer,
		Contexts: contexts,
	}, nil
}

