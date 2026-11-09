package repo

import (
	"context"
	"sync"
	"time"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
)

type InMemoryTokenRepository struct {
	mu     sync.RWMutex
	byHash map[string]*domain.RefreshTokenRecord
	byUser map[string][]string // userID -> hashes
}

func NewInMemoryTokenRepository() *InMemoryTokenRepository {
	return &InMemoryTokenRepository{
		byHash: make(map[string]*domain.RefreshTokenRecord),
		byUser: make(map[string][]string),
	}
}

func (r *InMemoryTokenRepository) Save(ctx context.Context, record *domain.RefreshTokenRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copyRec := *record
	r.byHash[record.TokenHash] = &copyRec
	r.byUser[record.UserID] = append(r.byUser[record.UserID], record.TokenHash)
	return nil
}

func (r *InMemoryTokenRepository) GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshTokenRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, ok := r.byHash[tokenHash]
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	copyRec := *record
	return &copyRec, nil
}

func (r *InMemoryTokenRepository) Revoke(ctx context.Context, tokenHash string, revokedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.byHash[tokenHash]
	if !ok {
		return domain.ErrInvalidToken
	}
	record.Revoke(revokedAt)
	return nil
}

func (r *InMemoryTokenRepository) RevokeAllForUser(ctx context.Context, userID string, revokedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	hashes := r.byUser[userID]
	for _, h := range hashes {
		if rec, ok := r.byHash[h]; ok {
			rec.Revoke(revokedAt)
		}
	}
	return nil
}