package transport

import "github.com/nhanvnguyen8x/palisade/internal/features/ingestionjob/domain"

func ToIngestionJobResponse(job *domain.IngestionJob) IngestionJobResponse {
	return IngestionJobResponse{
		ID:                job.ID.String(),
		KnowledgeSourceID: job.KnowledgeSourceID.String(),
		Status:            string(job.Status),
		ErrorMessage:      job.ErrorMessage,
		StartedAt:         job.StartedAt,
		CompletedAt:       job.CompletedAt,
		CreatedAt:         job.CreatedAt,
		UpdatedAt:         job.UpdatedAt,
	}
}

func ToIngestionJobListResponse(jobs []*domain.IngestionJob) []IngestionJobResponse {
	result := make([]IngestionJobResponse, 0, len(jobs))
	for _, job := range jobs {
		result = append(result, ToIngestionJobResponse(job))
	}

	return result
}
