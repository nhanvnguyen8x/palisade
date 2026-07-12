package application

import (
	"context"
	"strings"

	"github.com/google/uuid"
	orgdomain "github.com/nhanvnguyen8x/palisade/internal/features/organization/domain"
	"github.com/nhanvnguyen8x/palisade/internal/features/workspace/domain"
)

type OrganizationChecker interface {
	FindByID(ctx context.Context, id uuid.UUID) (*orgdomain.Organization, error)
}

type Service struct {
	repository          domain.Repository
	organizationChecker OrganizationChecker
}

func NewService(repository domain.Repository, organizationChecker OrganizationChecker) *Service {
	return &Service{
		repository:          repository,
		organizationChecker: organizationChecker,
	}
}

func (s *Service) Create(ctx context.Context, cmd CreateCommand) (*domain.Workspace, error) {
	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		return nil, domain.ErrInvalidName
	}

	if _, err := s.organizationChecker.FindByID(ctx, cmd.OrganizationID); err != nil {
		return nil, domain.ErrOrganizationNotFound
	}

	workspace := domain.NewWorkspace(cmd.OrganizationID, name)

	if err := s.repository.Create(ctx, workspace); err != nil {
		return nil, err
	}

	return workspace, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Workspace, error) {
	return s.repository.FindByID(ctx, id)
}
