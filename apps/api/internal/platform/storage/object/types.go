package object

import "io"

type UploadRequest struct {
	Key         string
	Reader      io.Reader
	Size        int64
	ContentType string
	Metadata    map[string]string
}

type UploadResult struct {
	Key       string
	ETag      string
	VersionID string
}
