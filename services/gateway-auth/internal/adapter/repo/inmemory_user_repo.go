package repo

import (
	"context"
	"strings"
	"sync"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
)

type InMemoryUserRepository struct {
	mu     sync.RWMutex
	byID   map[string]*domain.User
	byMail map[string]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		byID:   make(map[string]*domain.User),
		byMail: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	copyUser := *user
	return &copyUser, nil
}

func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clean := strings.ToLower(strings.TrimSpace(email))
	user, ok := r.byMail[clean]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	copyUser := *user
	return &copyUser, nil
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	clean := strings.ToLower(strings.TrimSpace(user.Email))
	if _, ok := r.byMail[clean]; ok {
		return domain.ErrEmailAlreadyExists
	}

	copyUser := *user
	r.byID[user.ID] = &copyUser
	r.byMail[clean] = &copyUser
	return nil
}

func (r *InMemoryUserRepository) Update(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[user.ID]; !ok {
		return domain.ErrUserNotFound
	}

	clean := strings.ToLower(strings.TrimSpace(user.Email))
	copyUser := *user
	r.byID[user.ID] = &copyUser
	r.byMail[clean] = &copyUser
	return nil
}