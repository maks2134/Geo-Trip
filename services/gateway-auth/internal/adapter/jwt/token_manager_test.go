package jwt_test

import (
	"testing"
	"time"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/jwt"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
)

type staticClock struct{ t time.Time }

func (s staticClock) Now() time.Time { return s.t }

func TestTokenManager_Lifecycle(t *testing.T) {
	clock := staticClock{t: time.Date(2026, 11, 9, 10, 0, 0, 0, time.UTC)}
	tm := jwt.NewTokenManager("my-ultra-secure-key-32-chars-long!", 10*time.Minute, 24*time.Hour, clock)

	user := domain.NewUser("usr-1", "test@atlasflow.local", "hash", "Tester", clock.Now())

	pair, record, err := tm.GenerateTokenPair(user)
	if err != nil {
		t.Fatalf("failed to generate token pair: %v", err)
	}

	// Validate access token
	claims, err := tm.ValidateAccessToken(pair.AccessToken)
	if err != nil {
		t.Fatalf("access token validation failed: %v", err)
	}
	if claims.UserID != "usr-1" || claims.Email != "test@atlasflow.local" {
		t.Errorf("unexpected claims: %+v", claims)
	}

	// Validate refresh token
	refreshClaims, err := tm.ValidateRefreshToken(pair.RefreshToken)
	if err != nil {
		t.Fatalf("refresh token validation failed: %v", err)
	}
	if refreshClaims.UserID != "usr-1" {
		t.Errorf("unexpected refresh claims: %+v", refreshClaims)
	}

	// Cross validation must fail (access token tested as refresh)
	_, err = tm.ValidateRefreshToken(pair.AccessToken)
	if err == nil {
		t.Errorf("expected error when validating access token as refresh token")
	}

	// Hash in record matches
	expectedHash := domain.HashToken(pair.RefreshToken)
	if record.TokenHash != expectedHash {
		t.Errorf("token hash mismatch: record=%s expected=%s", record.TokenHash, expectedHash)
	}
}