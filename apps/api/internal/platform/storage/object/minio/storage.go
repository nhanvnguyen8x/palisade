package minio

import (
	"context"
	"fmt"
	"io"

	objectstorage "github.com/nhanvnguyen8x/palisade/internal/platform/storage/object"

	minioSDK "github.com/minio/minio-go/v7"
)

type Storage struct {
	client *minioSDK.Client
	bucket string
}

func NewStorage(client *minioSDK.Client, bucket string) *Storage {
	return &Storage{
		client: client,
		bucket: bucket,
	}
}

func (s *Storage) Upload(ctx context.Context, req objectstorage.UploadRequest) (*objectstorage.UploadResult, error) {
	info, err := s.client.PutObject(
		ctx,
		s.bucket,
		req.Key,
		req.Reader,
		req.Size,
		minioSDK.PutObjectOptions{
			ContentType:  req.ContentType,
			UserMetadata: req.Metadata,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("upload object: %w", err)
	}

	return &objectstorage.UploadResult{
		Key:       info.Key,
		ETag:      info.ETag,
		VersionID: info.VersionID,
	}, nil
}

func (s *Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	reader, err := s.client.GetObject(ctx, s.bucket, key, minioSDK.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("download object: %w", err)
	}

	return reader, nil
}
func (s *Storage) Delete(ctx context.Context, key string) error {
	err := s.client.RemoveObject(ctx, s.bucket, key, minioSDK.RemoveObjectOptions{})

	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}

	return nil
}

func (s *Storage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucket, key, minioSDK.StatObjectOptions{})
	if err == nil {
		return true, nil
	}

	resp := minioSDK.ToErrorResponse(err)
	if resp.Code == "NoSuchKey" {
		return false, nil
	}

	return false, fmt.Errorf("stat object: %w", err)
}

func (s *Storage) Health(ctx context.Context) error {
	_, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("check bucket: %w", err)
	}

	return nil
}
