package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/jwt"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/repo"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/usecase"
)

type fixedClock struct{ t time.Time }

func (f fixedClock) Now() time.Time { return f.t }

type simpleIDGen struct{ id string }

func (s simpleIDGen) New() string { return s.id }

func setupAuthUsecase() (*usecase.AuthUsecase, *repo.InMemoryUserRepository, *repo.InMemoryTokenRepository, *jwt.TokenManager) {
	clock := fixedClock{t: time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)}
	idGen := simpleIDGen{id: "user-test-123"}
	userRepo := repo.NewInMemoryUserRepository()
	tokenRepo := repo.NewInMemoryTokenRepository()
	tm := jwt.NewTokenManager("test-jwt-secret-minimum-32-characters-length-ok", 15*time.Minute, 7*24*time.Hour, clock)
	uc := usecase.NewAuthUsecase(userRepo, tokenRepo, tm, clock, idGen)
	return uc, userRepo, tokenRepo, tm
}

func TestAuthUsecase_Register(t *testing.T) {
	tests := []struct {
		name        string
		input       usecase.RegisterInput
		prepare     func(r *repo.InMemoryUserRepository)
		wantErr     error
		checkResult func(t *testing.T, pair *domain.TokenPair, u *domain.User)
	}{
		{
			name: "successful registration",
			input: usecase.RegisterInput{
				Email:       "traveler@example.com",
				Password:    "supersecret123",
				DisplayName: "Globe Trotter",
			},
			wantErr: nil,
			checkResult: func(t *testing.T, pair *domain.TokenPair, u *domain.User) {
				if pair == nil || pair.AccessToken == "" || pair.RefreshToken == "" {
					t.Fatalf("expected non-empty tokens, got %+v", pair)
				}
				if u.Email != "traveler@example.com" {
					t.Errorf("expected email traveler@example.com, got %s", u.Email)
				}
				if u.DisplayName != "Globe Trotter" {
					t.Errorf("expected display name Globe Trotter, got %s", u.DisplayName)
				}
			},
		},
		{
			name: "duplicate email",
			input: usecase.RegisterInput{
				Email:       "existing@example.com",
				Password:    "password123",
				DisplayName: "Existing",
			},
			prepare: func(r *repo.InMemoryUserRepository) {
				_ = r.Create(context.Background(), domain.NewUser("u-1", "existing@example.com", "hash", "Existing", time.Now()))
			},
			wantErr: domain.ErrEmailAlreadyExists,
		},
		{
			name: "invalid email format",
			input: usecase.RegisterInput{
				Email:    "not-an-email",
				Password: "password123",
			},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "short password",
			input: usecase.RegisterInput{
				Email:    "valid@example.com",
				Password: "short",
			},
			wantErr: domain.ErrInvalidInput,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uc, userRepo, _, _ := setupAuthUsecase()
			if tc.prepare != nil {
				tc.prepare(userRepo)
			}

			pair, u, err := uc.Register(context.Background(), tc.input)
			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tc.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.checkResult != nil {
				tc.checkResult(t, pair, u)
			}
		})
	}
}

func TestAuthUsecase_Login(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase()
	ctx := context.Background()

	hash, _ := domain.HashPassword("correctpassword123")
	existingUser := domain.NewUser("u-login", "login@example.com", hash, "Login User", time.Now())
	_ = userRepo.Create(ctx, existingUser)

	tests := []struct {
		name    string
		input   usecase.LoginInput
		wantErr error
	}{
		{
			name: "successful login",
			input: usecase.LoginInput{
				Email:    "login@example.com",
				Password: "correctpassword123",
			},
			wantErr: nil,
		},
		{
			name: "wrong password",
			input: usecase.LoginInput{
				Email:    "login@example.com",
				Password: "wrongpassword",
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name: "user not found",
			input: usecase.LoginInput{
				Email:    "nonexistent@example.com",
				Password: "correctpassword123",
			},
			wantErr: domain.ErrInvalidCredentials,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pair, user, err := uc.Login(ctx, tc.input)
			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if pair.AccessToken == "" || user.ID != "u-login" {
				t.Errorf("login response invalid: user=%+v pair=%+v", user, pair)
			}
		})
	}
}

func TestAuthUsecase_RefreshTokenRotation(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase()
	ctx := context.Background()

	hash, _ := domain.HashPassword("secret12345")
	u := domain.NewUser("u-rot", "rot@example.com", hash, "Rotation", time.Now())
	_ = userRepo.Create(ctx, u)

	pair, _, err := uc.Login(ctx, usecase.LoginInput{Email: "rot@example.com", Password: "secret12345"})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	// 1. First refresh should succeed
	newPair, err := uc.RefreshToken(ctx, pair.RefreshToken)
	if err != nil {
		t.Fatalf("expected refresh to succeed, got %v", err)
	}
	if newPair.AccessToken == "" || newPair.RefreshToken == pair.RefreshToken {
		t.Fatalf("expected new refreshed tokens")
	}

	// 2. Reusing old refresh token must be rejected (revoked by rotation)
	_, err = uc.RefreshToken(ctx, pair.RefreshToken)
	if err == nil {
		t.Fatalf("expected old refresh token to be rejected after rotation")
	}
}