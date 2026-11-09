package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"` // in seconds
	ExpiresAt    time.Time `json:"expires_at"`
}

type RefreshTokenRecord struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (r *RefreshTokenRecord) IsRevoked() bool {
	return r.RevokedAt != nil
}

func (r *RefreshTokenRecord) IsExpired(now time.Time) bool {
	return now.After(r.ExpiresAt)
}

func (r *RefreshTokenRecord) Revoke(now time.Time) {
	r.RevokedAt = &now
}