package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	chatdomain "github.com/nhanvnguyen8x/palisade/internal/features/chat/domain"
)

type ChatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) ListChunkEmbeddings(
	ctx context.Context,
	knowledgeBaseID uuid.UUID,
	model string,
) ([]chatdomain.ChunkEmbedding, error) {
	// Dev build: fetch all chunk vectors and rank in-memory.
	// Later we can optimize using Milvus / pgvector.
	const query = `
SELECT
    c.id,
    c.text,
    e.vector
FROM knowledge_sources ks
JOIN chunks c ON c.knowledge_source_id = ks.id
JOIN embeddings e ON e.chunk_id = c.id
WHERE ks.knowledge_base_id = $1
  AND e.model = $2
`

	rows, err := r.db.Query(ctx, query, knowledgeBaseID, model)
	if err != nil {
		return nil, fmt.Errorf("list chunk embeddings: %w", err)
	}
	defer rows.Close()

	var results []chatdomain.ChunkEmbedding

	for rows.Next() {
		var chunkID uuid.UUID
		var text string
		var vector []float64

		if err := rows.Scan(&chunkID, &text, &vector); err != nil {
			return nil, fmt.Errorf("scan chunk embedding: %w", err)
		}

		results = append(results, chatdomain.ChunkEmbedding{
			ChunkID: chunkID,
			Text:    text,
			Vector:  vector,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	if len(results) == 0 {
		// We keep NotFound logic in service/transport based on empty results.
		return nil, nil
	}

	return results, nil
}

