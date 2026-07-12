package bootstrap

import (
	"context"
	"fmt"

	"github.com/nhanvnguyen8x/palisade/internal/config"
	"github.com/nhanvnguyen8x/palisade/internal/platform/storage/object"
	"github.com/nhanvnguyen8x/palisade/internal/platform/storage/object/minio"
)

func NewObjectStorage(ctx context.Context, cfg config.Storage) (object.Storage, error) {
	client, err := minio.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}

	storage := minio.NewStorage(client, cfg.Bucket)
	if err := storage.Initialize(ctx); err != nil {
		return nil, fmt.Errorf("minio init: %w", err)
	}

	return storage, nil
}
