package application

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	docdomain "github.com/nhanvnguyen8x/palisade/internal/features/document/domain"
	jobdomain "github.com/nhanvnguyen8x/palisade/internal/features/ingestionjob/domain"
	kbdomain "github.com/nhanvnguyen8x/palisade/internal/features/knowledgebase/domain"
	"github.com/nhanvnguyen8x/palisade/internal/features/knowledgesource/domain"
	"github.com/nhanvnguyen8x/palisade/internal/features/knowledgesource/domain/events"
	"github.com/nhanvnguyen8x/palisade/internal/platform/database/postgres"
	"github.com/nhanvnguyen8x/palisade/internal/platform/event"
	"github.com/nhanvnguyen8x/palisade/internal/platform/storage/object"
)

type KnowledgeBaseChecker interface {
	FindByID(ctx context.Context, id uuid.UUID) (*kbdomain.KnowledgeBase, error)
}

type Service struct {
	pool        *pgxpool.Pool
	kbChecker   KnowledgeBaseChecker
	sourceRepo  domain.Repository
	storage     object.Storage
	publisher   event.Publisher
}

func NewService(
	pool *pgxpool.Pool,
	kbChecker KnowledgeBaseChecker,
	sourceRepo domain.Repository,
	storage object.Storage,
	publisher event.Publisher,
) *Service {
	return &Service{
		pool:       pool,
		kbChecker:  kbChecker,
		sourceRepo: sourceRepo,
		storage:    storage,
		publisher:  publisher,
	}
}

func (s *Service) Upload(ctx context.Context, cmd UploadCommand) (*UploadResult, error) {
	kb, err := s.kbChecker.FindByID(ctx, cmd.KnowledgeBaseID)
	if err != nil {
		return nil, domain.ErrKnowledgeBaseNotFound
	}

	checksum, reader, err := calculateChecksum(cmd.Reader)
	if err != nil {
		return nil, err
	}

	source := domain.NewKnowledgeSource(cmd.KnowledgeBaseID, domain.SourceTypeDocument)
	storageKey := buildStorageKey(kb.WorkspaceID, source.ID)

	_, err = s.storage.Upload(ctx, object.UploadRequest{
		Key:         storageKey,
		Reader:      reader,
		Size:        cmd.Size,
		ContentType: cmd.ContentType,
	})
	if err != nil {
		return nil, fmt.Errorf("upload object: %w", err)
	}

	source.MarkUploaded()

	document := docdomain.NewDocument(
		source.ID,
		kb.WorkspaceID,
		cmd.FileName,
		cmd.ContentType,
		cmd.Size,
		storageKey,
		checksum,
	)

	job := jobdomain.NewIngestionJob(source.ID)

	if err := postgres.SaveKnowledgeSourceUpload(ctx, s.pool, source, document, job); err != nil {
		return nil, fmt.Errorf("save upload: %w", err)
	}

	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, events.SourceUploaded{
			KnowledgeSourceID: source.ID,
			KnowledgeBaseID:   kb.ID,
			WorkspaceID:       kb.WorkspaceID,
			Occurred:          time.Now().UTC(),
		})
	}

	return &UploadResult{
		KnowledgeSourceID: source.ID,
		IngestionJobID:    job.ID,
		DocumentID:        document.ID,
		Status:            string(source.Status),
		FileName:          cmd.FileName,
		CreatedAt:         source.CreatedAt,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.KnowledgeSource, error) {
	return s.sourceRepo.FindByID(ctx, id)
}

func buildStorageKey(workspaceID, sourceID uuid.UUID) string {
	return fmt.Sprintf(
		"workspaces/%s/sources/%s/original",
		workspaceID,
		sourceID,
	)
}

func calculateChecksum(reader io.Reader) (string, io.Reader, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", nil, err
	}

	sum := sha256.Sum256(data)

	return hex.EncodeToString(sum[:]), bytes.NewReader(data), nil
}
