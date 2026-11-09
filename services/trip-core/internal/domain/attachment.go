package domain

import (
	"strings"
	"time"
)

type AttachmentVisibility int32

const (
	AttachmentVisibilityUnspecified AttachmentVisibility = 0
	AttachmentVisibilityPrivate     AttachmentVisibility = 1
	AttachmentVisibilityPublic      AttachmentVisibility = 2
)

type Attachment struct {
	ID          string
	TripID      string
	PlaceID     *string
	DayID       *string
	Filename    string
	ContentType string
	FileSize    int64
	StorageKey  string
	DownloadURL string
	Visibility  AttachmentVisibility
	IsReady     bool
	CreatedAt   time.Time
}

var allowedMimeTypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/png":       true,
	"image/webp":      true,
}

const MaxAttachmentSizeBytes = 25 * 1024 * 1024 // 25 MB

func ValidateAttachmentMeta(contentType string, size int64) error {
	ct := strings.ToLower(strings.TrimSpace(contentType))
	if !allowedMimeTypes[ct] {
		return ErrInvalidMimeType
	}
	if size > MaxAttachmentSizeBytes {
		return ErrFileTooLarge
	}
	return nil
}