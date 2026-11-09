package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
)

type RegisterInput struct {
	Email       string
	Password    string
	DisplayName string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthUsecase struct {
	users  UserRepository
	tokens RefreshTokenRepository
	jwt    TokenManager
	clock  Clock
	ids    IDGen
}

func NewAuthUsecase(
	users UserRepository,
	tokens RefreshTokenRepository,
	jwt TokenManager,
	clock Clock,
	ids IDGen,
) *AuthUsecase {
	return &AuthUsecase{
		users:  users,
		tokens: tokens,
		jwt:    jwt,
		clock:  clock,
		ids:    ids,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, in RegisterInput) (*domain.TokenPair, *domain.User, error) {
	if err := domain.ValidateEmail(in.Email); err != nil {
		return nil, nil, err
	}
	if err := domain.ValidatePassword(in.Password); err != nil {
		return nil, nil, err
	}

	cleanEmail := strings.ToLower(strings.TrimSpace(in.Email))
	existing, err := u.users.GetByEmail(ctx, cleanEmail)
	if err == nil && existing != nil {
		return nil, nil, domain.ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, nil, fmt.Errorf("check email: %w", err)
	}

	hash, err := domain.HashPassword(in.Password)
	if err != nil {
		return nil, nil, err
	}

	now := u.clock.Now()
	displayName := strings.TrimSpace(in.DisplayName)
	if displayName == "" {
		displayName = strings.Split(cleanEmail, "@")[0]
	}

	user := domain.NewUser(u.ids.New(), cleanEmail, hash, displayName, now)
	if err := u.users.Create(ctx, user); err != nil {
		return nil, nil, fmt.Errorf("create user: %w", err)
	}

	pair, record, err := u.jwt.GenerateTokenPair(user)
	if err != nil {
		return nil, nil, fmt.Errorf("generate tokens: %w", err)
	}

	if err := u.tokens.Save(ctx, record); err != nil {
		return nil, nil, fmt.Errorf("save refresh token: %w", err)
	}

	return pair, user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, in LoginInput) (*domain.TokenPair, *domain.User, error) {
	cleanEmail := strings.ToLower(strings.TrimSpace(in.Email))
	if cleanEmail == "" || in.Password == "" {
		return nil, nil, domain.ErrInvalidCredentials
	}

	user, err := u.users.GetByEmail(ctx, cleanEmail)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, nil, domain.ErrInvalidCredentials
		}
		return nil, nil, fmt.Errorf("get user: %w", err)
	}

	if !user.CheckPassword(in.Password) {
		return nil, nil, domain.ErrInvalidCredentials
	}

	pair, record, err := u.jwt.GenerateTokenPair(user)
	if err != nil {
		return nil, nil, fmt.Errorf("generate tokens: %w", err)
	}

	if err := u.tokens.Save(ctx, record); err != nil {
		return nil, nil, fmt.Errorf("save refresh token: %w", err)
	}

	return pair, user, nil
}

// RefreshToken implements atomic refresh token rotation.
// The presented token is invalidated and a fresh pair is minted.
func (u *AuthUsecase) RefreshToken(ctx context.Context, refreshTokenStr string) (*domain.TokenPair, error) {
	claims, err := u.jwt.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidToken, err)
	}

	tokenHash := domain.HashToken(refreshTokenStr)
	record, err := u.tokens.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	now := u.clock.Now()
	if record.IsRevoked() {
		return nil, domain.ErrTokenRevoked
	}
	if record.IsExpired(now) {
		return nil, domain.ErrTokenExpired
	}

	// Revoke old token (rotation)
	if err := u.tokens.Revoke(ctx, tokenHash, now); err != nil {
		return nil, fmt.Errorf("revoke old refresh token: %w", err)
	}

	user, err := u.users.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user for refresh: %w", err)
	}

	pair, newRecord, err := u.jwt.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("generate token pair: %w", err)
	}

	if err := u.tokens.Save(ctx, newRecord); err != nil {
		return nil, fmt.Errorf("save new refresh token: %w", err)
	}

	return pair, nil
}

func (u *AuthUsecase) GetCurrentUser(ctx context.Context, userID string) (*domain.User, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrUnauthorized
	}
	user, err := u.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func (u *AuthUsecase) Logout(ctx context.Context, refreshTokenStr string) error {
	tokenHash := domain.HashToken(refreshTokenStr)
	return u.tokens.Revoke(ctx, tokenHash, u.clock.Now())
}