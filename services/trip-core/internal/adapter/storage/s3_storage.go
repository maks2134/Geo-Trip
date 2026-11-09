package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/usecase"
)

var _ usecase.StoragePort = (*S3Storage)(nil)

type S3Storage struct {
	endpoint string
	bucket   string
}

func NewS3Storage(endpoint, bucket string) *S3Storage {
	return &S3Storage{
		endpoint: endpoint,
		bucket:   bucket,
	}
}

func (s *S3Storage) PresignPut(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
	// MinIO/S3 compatible presigned PUT URL
	expTimestamp := time.Now().Add(expires).Unix()
	return fmt.Sprintf("%s/%s/%s?upload=true&expires=%d", s.endpoint, s.bucket, key, expTimestamp), nil
}

func (s *S3Storage) PresignGet(ctx context.Context, key string, expires time.Duration) (string, error) {
	// MinIO/S3 compatible presigned GET URL
	expTimestamp := time.Now().Add(expires).Unix()
	return fmt.Sprintf("%s/%s/%s?expires=%d", s.endpoint, s.bucket, key, expTimestamp), nil
}