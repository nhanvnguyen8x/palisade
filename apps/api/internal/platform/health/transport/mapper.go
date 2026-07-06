package transport

import (
	"github.com/nhanvnguyen8x/palisade/internal/platform/health/domain"
	"github.com/nhanvnguyen8x/palisade/internal/platform/health/dto"
)

func ToHealthResponse(report *domain.HealthReport) dto.HealthResponse {
	return dto.HealthResponse{
		Status: string(report.Status),
	}
}
