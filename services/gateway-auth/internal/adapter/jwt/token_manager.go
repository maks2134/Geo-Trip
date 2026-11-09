package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/usecase"
)

type TokenManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	clock      usecase.Clock
}

func NewTokenManager(
	secret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	clock usecase.Clock,
) *TokenManager {
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}
	if refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}
	return &TokenManager{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		clock:      clock,
	}
}

func (m *TokenManager) GenerateTokenPair(user *domain.User) (*domain.TokenPair, *domain.RefreshTokenRecord, error) {
	now := m.clock.Now()

	// 1. Access Token
	accessExpiry := now.Add(m.accessTTL)
	accessClaims := domain.AuthClaims{
		UserID:    user.ID,
		Email:     user.Email,
		TokenType: domain.TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			ID:        uuid.NewString(),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(m.secret)
	if err != nil {
		return nil, nil, fmt.Errorf("sign access token: %w", err)
	}

	// 2. Refresh Token
	refreshExpiry := now.Add(m.refreshTTL)
	refreshClaims := domain.AuthClaims{
		UserID:    user.ID,
		Email:     user.Email,
		TokenType: domain.TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			ID:        uuid.NewString(),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(m.secret)
	if err != nil {
		return nil, nil, fmt.Errorf("sign refresh token: %w", err)
	}

	// Compute SHA-256 hash for safe DB persistence
	tokenHash := domain.HashToken(refreshToken)
	record := &domain.RefreshTokenRecord{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: refreshExpiry,
		CreatedAt: now,
	}

	pair := &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(m.accessTTL.Seconds()),
		ExpiresAt:    accessExpiry,
	}

	return pair, record, nil
}

func (m *TokenManager) ValidateAccessToken(tokenStr string) (*domain.AuthClaims, error) {
	return m.validateToken(tokenStr, domain.TokenTypeAccess)
}

func (m *TokenManager) ValidateRefreshToken(tokenStr string) (*domain.AuthClaims, error) {
	return m.validateToken(tokenStr, domain.TokenTypeRefresh)
}

func (m *TokenManager) validateToken(tokenStr string, expectedType domain.TokenType) (*domain.AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	}, jwt.WithTimeFunc(m.clock.Now))

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*domain.AuthClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	if claims.TokenType != expectedType {
		return nil, fmt.Errorf("invalid token type: expected %s, got %s", expectedType, claims.TokenType)
	}

	return claims, nil
}