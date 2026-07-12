package domain

import "errors"

var (
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrInvalidName          = errors.New("invalid organization name")
)
