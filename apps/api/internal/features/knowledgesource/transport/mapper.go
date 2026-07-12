package transport

import (
	"github.com/nhanvnguyen8x/palisade/internal/features/knowledgesource/application"
	"github.com/nhanvnguyen8x/palisade/internal/features/knowledgesource/domain"
)

func ToUploadResponse(result *application.UploadResult) UploadResponse {
	return UploadResponse{
		KnowledgeSourceID: result.KnowledgeSourceID.String(),
		IngestionJobID:    result.IngestionJobID.String(),
		DocumentID:        result.DocumentID.String(),
		Status:            result.Status,
		FileName:          result.FileName,
		CreatedAt:         result.CreatedAt,
	}
}

func ToKnowledgeSourceResponse(source *domain.KnowledgeSource) KnowledgeSourceResponse {
	return KnowledgeSourceResponse{
		ID:              source.ID.String(),
		KnowledgeBaseID: source.KnowledgeBaseID.String(),
		Type:            string(source.Type),
		Status:          string(source.Status),
		CreatedAt:       source.CreatedAt,
		UpdatedAt:       source.UpdatedAt,
	}
}
