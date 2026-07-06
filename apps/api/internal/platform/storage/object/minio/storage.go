package minio

import "github.com/minio/minio-go/v7"

type Storage struct {
	client *minio.Client
	bucket string
}

func NewStorage(client *minio.Client, bucket string,
) *Storage {
	return &Storage{}
}
