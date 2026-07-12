package minio

import (
	"context"
	"fmt"

	minioSDK "github.com/minio/minio-go/v7"
)

func (s *Storage) Initialize(ctx context.Context) error {

	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("check bucket existence: %w", err)
	}

	if exists {
		return nil
	}

	err = s.client.MakeBucket(ctx, s.bucket, minioSDK.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("create bucket %q: %w", s.bucket, err)
	}

	return nil
}
