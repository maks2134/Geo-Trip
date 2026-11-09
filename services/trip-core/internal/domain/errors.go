package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrForbidden       = errors.New("forbidden")
	ErrValidation      = errors.New("validation failed")
	ErrConflict        = errors.New("conflict")
	ErrInvalidMimeType = errors.New("unsupported MIME type: allowed formats are PDF, JPEG, PNG, WEBP")
	ErrFileTooLarge    = errors.New("file exceeds maximum allowed size")
	ErrShareRevoked    = errors.New("public share link has been revoked")
)