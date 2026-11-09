package usecase

import (
	"context"
	"time"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, record *domain.RefreshTokenRecord) error
	GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshTokenRecord, error)
	Revoke(ctx context.Context, tokenHash string, revokedAt time.Time) error
	RevokeAllForUser(ctx context.Context, userID string, revokedAt time.Time) error
}

type TokenManager interface {
	GenerateTokenPair(user *domain.User) (*domain.TokenPair, *domain.RefreshTokenRecord, error)
	ValidateAccessToken(tokenStr string) (*domain.AuthClaims, error)
	ValidateRefreshToken(tokenStr string) (*domain.AuthClaims, error)
}

type PublicTripClient interface {
	GetPublicTrip(ctx context.Context, slug string) (*domain.PublicTrip, error)
	ClonePublicTrip(ctx context.Context, slug, targetOwnerID string) (*domain.ClonedTripSummary, error)
}

type Clock interface {
	Now() time.Time
}

type IDGen interface {
	New() string
}
