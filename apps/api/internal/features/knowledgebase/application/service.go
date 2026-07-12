package application

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/nhanvnguyen8x/palisade/internal/features/knowledgebase/domain"
	wsdomain "github.com/nhanvnguyen8x/palisade/internal/features/workspace/domain"
)

type WorkspaceChecker interface {
	FindByID(ctx context.Context, id uuid.UUID) (*wsdomain.Workspace, error)
}

type Service struct {
	repository       domain.Repository
	workspaceChecker WorkspaceChecker
}

func NewService(repository domain.Repository, workspaceChecker WorkspaceChecker) *Service {
	return &Service{
		repository:       repository,
		workspaceChecker: workspaceChecker,
	}
}

func (s *Service) Create(ctx context.Context, cmd CreateCommand) (*domain.KnowledgeBase, error) {
	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		return nil, domain.ErrInvalidName
	}

	if _, err := s.workspaceChecker.FindByID(ctx, cmd.WorkspaceID); err != nil {
		return nil, domain.ErrWorkspaceNotFound
	}

	kb := domain.NewKnowledgeBase(cmd.WorkspaceID, name, strings.TrimSpace(cmd.Description))

	if err := s.repository.Create(ctx, kb); err != nil {
		return nil, err
	}

	return kb, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.KnowledgeBase, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*domain.KnowledgeBase, error) {
	if _, err := s.workspaceChecker.FindByID(ctx, workspaceID); err != nil {
		return nil, domain.ErrWorkspaceNotFound
	}

	return s.repository.ListByWorkspace(ctx, workspaceID)
}
