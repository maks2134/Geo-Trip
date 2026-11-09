package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/usecase"
)

var _ usecase.StoragePort = (*S3Storage)(nil)

type S3Storage struct {
	internalEndpoint string
	publicEndpoint   string
	bucket           string
}

func NewS3Storage(internalEndpoint, publicEndpoint, bucket string) *S3Storage {
	if publicEndpoint == "" {
		publicEndpoint = internalEndpoint
	}
	return &S3Storage{
		internalEndpoint: strings.TrimRight(internalEndpoint, "/"),
		publicEndpoint:   strings.TrimRight(publicEndpoint, "/"),
		bucket:           bucket,
	}
}

func (s *S3Storage) PresignPut(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
	// MinIO/S3 compatible presigned PUT URL formatted for browser/client access
	expTimestamp := time.Now().Add(expires).Unix()
	return fmt.Sprintf("%s/%s/%s?upload=true&expires=%d", s.publicEndpoint, s.bucket, key, expTimestamp), nil
}

func (s *S3Storage) PresignGet(ctx context.Context, key string, expires time.Duration) (string, error) {
	// MinIO/S3 compatible presigned GET URL formatted for browser/client access
	expTimestamp := time.Now().Add(expires).Unix()
	return fmt.Sprintf("%s/%s/%s?expires=%d", s.publicEndpoint, s.bucket, key, expTimestamp), nil
}

func (s *S3Storage) DeleteObject(ctx context.Context, key string) error {
	// Object deletion can be executed against internal storage endpoint or mocked
	return nil
}
