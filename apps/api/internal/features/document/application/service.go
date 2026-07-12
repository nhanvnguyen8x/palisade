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
	"github.com/nhanvnguyen8x/palisade/internal/features/document/domain"
	"github.com/nhanvnguyen8x/palisade/internal/features/document/domain/events"
	"github.com/nhanvnguyen8x/palisade/internal/platform/event"
	"github.com/nhanvnguyen8x/palisade/internal/platform/storage/object"
)

type Service struct {
	repository domain.Repository
	storage    object.Storage
	publisher  event.Publisher
}

func NewService(
	repository domain.Repository,
	storage object.Storage,
	publisher event.Publisher,
) *Service {
	return &Service{
		repository: repository,
		storage:    storage,
		publisher:  publisher,
	}
}

func (s *Service) Upload(
	ctx context.Context,
	cmd UploadCommand,
) (*UploadResult, error) {

	checksum, reader, err := calculateChecksum(cmd.Reader)
	if err != nil {
		return nil, err
	}

	document := domain.NewDocument(
		cmd.KnowledgeSourceID,
		cmd.WorkspaceID,
		cmd.FileName,
		cmd.ContentType,
		cmd.Size,
		buildStorageKey(cmd.WorkspaceID),
		checksum,
	)

	_, err = s.storage.Upload(ctx, object.UploadRequest{
		Key:         document.StorageKey,
		Reader:      reader,
		Size:        cmd.Size,
		ContentType: cmd.ContentType,
	})

	if err != nil {
		return nil, fmt.Errorf("upload object: %w", err)
	}

	if err := s.repository.Create(ctx, document); err != nil {
		return nil, fmt.Errorf("create document: %w", err)
	}

	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, events.DocumentUploaded{
			DocumentID:  document.ID,
			WorkspaceID: document.WorkspaceID,
			Occurred:    time.Now().UTC(),
		})
	}

	return &UploadResult{
		ID:        document.ID,
		Status:    string(document.Status),
		FileName:  document.FileName,
		CreatedAt: document.CreatedAt,
	}, nil
}

func buildStorageKey(workspaceID uuid.UUID) string {
	return fmt.Sprintf(
		"workspaces/%s/%s/original",
		workspaceID,
		uuid.New(),
	)
}

func calculateChecksum(
	reader io.Reader,
) (string, io.Reader, error) {

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", nil, err
	}

	sum := sha256.Sum256(data)

	return hex.EncodeToString(sum[:]),
		bytes.NewReader(data),
		nil
}
