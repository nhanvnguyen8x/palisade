package application

import (
	"context"

	"github.com/google/uuid"
	jobdomain "github.com/nhanvnguyen8x/palisade/internal/features/ingestionjob/domain"
	ksdomain "github.com/nhanvnguyen8x/palisade/internal/features/knowledgesource/domain"
)

type Service struct {
	jobRepo    jobdomain.Repository
	sourceRepo ksdomain.Repository
}

func NewService(jobRepo jobdomain.Repository, sourceRepo ksdomain.Repository) *Service {
	return &Service{
		jobRepo:    jobRepo,
		sourceRepo: sourceRepo,
	}
}

func (s *Service) ListPending(ctx context.Context, limit int) ([]*jobdomain.IngestionJob, error) {
	return s.jobRepo.FindPending(ctx, limit)
}

func (s *Service) UpdateStatus(ctx context.Context, cmd UpdateStatusCommand) (*jobdomain.IngestionJob, error) {
	job, err := s.jobRepo.FindByID(ctx, cmd.JobID)
	if err != nil {
		return nil, err
	}

	source, err := s.sourceRepo.FindByID(ctx, job.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}

	switch cmd.Status {
	case jobdomain.JobStatusProcessing:
		job.MarkProcessing()
		source.MarkProcessing()
	case jobdomain.JobStatusCompleted:
		job.MarkCompleted()
		source.MarkReady()
	case jobdomain.JobStatusFailed:
		job.MarkFailed(cmd.ErrorMessage)
		source.MarkFailed()
	default:
		return nil, jobdomain.ErrInvalidStatus
	}

	if err := s.jobRepo.Update(ctx, job); err != nil {
		return nil, err
	}

	if err := s.sourceRepo.Update(ctx, source); err != nil {
		return nil, err
	}

	return job, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*jobdomain.IngestionJob, error) {
	return s.jobRepo.FindByID(ctx, id)
}
