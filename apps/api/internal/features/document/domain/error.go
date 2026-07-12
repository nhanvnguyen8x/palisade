package domain

import "errors"

var (
	ErrDocumentNotFound   = errors.New("document not found")
	ErrInvalidFileName    = errors.New("invalid file name")
	ErrInvalidWorkspace   = errors.New("invalid workspace")
	ErrInvalidContentType = errors.New("invalid content type")
)
