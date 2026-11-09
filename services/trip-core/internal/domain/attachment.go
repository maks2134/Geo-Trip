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
	OwnerType   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

const (
	MaxPhotoSizeBytes = 10 * 1024 * 1024 // 10 MB for photos (JPEG, PNG, WEBP)
	MaxDocSizeBytes   = 25 * 1024 * 1024 // 25 MB for documents (PDF)
)

var allowedPhotoTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

func ValidateAttachmentMeta(contentType string, size int64) error {
	ct := strings.ToLower(strings.TrimSpace(contentType))
	if ct == "application/pdf" {
		if size > MaxDocSizeBytes {
			return ErrFileTooLarge
		}
		return nil
	}
	if allowedPhotoTypes[ct] {
		if size > MaxPhotoSizeBytes {
			return ErrFileTooLarge
		}
		return nil
	}
	return ErrInvalidMimeType
}
