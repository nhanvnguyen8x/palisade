package domain

import "errors"

var (
	ErrWorkspaceNotFound      = errors.New("workspace not found")
	ErrInvalidName            = errors.New("invalid workspace name")
	ErrOrganizationNotFound   = errors.New("organization not found")
)
