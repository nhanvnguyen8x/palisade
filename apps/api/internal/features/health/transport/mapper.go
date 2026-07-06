package transport

import (
	"github.com/nhanvnguyen8x/palisade/internal/features/health/domain"
	"github.com/nhanvnguyen8x/palisade/internal/features/health/dto"
)

func ToHealthResponse(report *domain.HealthReport) dto.HealthResponse {
	return dto.HealthResponse{
		Status: string(report.Status),
	}
}
