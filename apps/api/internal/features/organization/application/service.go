package application

import (
	"context"
	"strings"

	"github.com/nhanvnguyen8x/palisade/internal/features/organization/domain"
)

type Service struct {
	repository domain.Repository
}

func NewService(repository domain.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, cmd CreateCommand) (*domain.Organization, error) {
	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		return nil, domain.ErrInvalidName
	}

	org := domain.NewOrganization(name)

	if err := s.repository.Create(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*domain.Organization, error) {
	parsed, err := parseUUID(id)
	if err != nil {
		return nil, domain.ErrOrganizationNotFound
	}

	return s.repository.FindByID(ctx, parsed)
}
