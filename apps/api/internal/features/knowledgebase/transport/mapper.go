package transport

import "github.com/nhanvnguyen8x/palisade/internal/features/knowledgebase/domain"

func ToKnowledgeBaseResponse(kb *domain.KnowledgeBase) KnowledgeBaseResponse {
	return KnowledgeBaseResponse{
		ID:          kb.ID.String(),
		WorkspaceID: kb.WorkspaceID.String(),
		Name:        kb.Name,
		Description: kb.Description,
		CreatedAt:   kb.CreatedAt,
		UpdatedAt:   kb.UpdatedAt,
	}
}

func ToKnowledgeBaseListResponse(bases []*domain.KnowledgeBase) []KnowledgeBaseResponse {
	result := make([]KnowledgeBaseResponse, 0, len(bases))
	for _, kb := range bases {
		result = append(result, ToKnowledgeBaseResponse(kb))
	}

	return result
}
