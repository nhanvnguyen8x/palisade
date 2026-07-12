package object

import "errors"

var (
	ErrObjectNotFound      = errors.New("object not found")
	ErrObjectAlreadyExists = errors.New("object already exists")
	ErrBucketNotFound      = errors.New("bucket not found")
	ErrInvalidObjectKey    = errors.New("invalid object key")
	ErrStorageUnavailable  = errors.New("storage unavailable")
)
