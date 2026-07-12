package domain

import "errors"

var (
	ErrIngestionJobNotFound = errors.New("ingestion job not found")
	ErrInvalidStatus        = errors.New("invalid job status")
)
