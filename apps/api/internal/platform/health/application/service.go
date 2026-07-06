package application

import (
	"context"

	"github.com/nhanvnguyen8x/palisade/internal/platform/health/domain"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetHealth(ctx context.Context) (*domain.HealthReport, error) {
	return &domain.HealthReport{}, nil
}
