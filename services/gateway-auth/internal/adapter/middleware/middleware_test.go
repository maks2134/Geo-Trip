package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/jwt"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/middleware"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
)

type staticClock struct{ t time.Time }

func (s staticClock) Now() time.Time { return s.t }

func TestAuthMiddleware(t *testing.T) {
	clock := staticClock{t: time.Now()}
	tm := jwt.NewTokenManager("secret-key-at-least-32-bytes-long-ok", 15*time.Minute, 24*time.Hour, clock)
	mw := middleware.NewAuthMiddleware(tm)

	user := domain.NewUser("uid-99", "user@atlasflow.local", "hash", "Tester", time.Now())
	pair, _, _ := tm.GenerateTokenPair(user)

	tests := []struct {
		name           string
		authHeader     string
		spoofedUserID  string
		expectedStatus int
		verifyHeader   bool
	}{
		{
			name:           "valid token forwards X-User-Id",
			authHeader:     "Bearer " + pair.AccessToken,
			spoofedUserID:  "attacker-id",
			expectedStatus: http.StatusOK,
			verifyHeader:   true,
		},
		{
			name:           "missing auth header rejected",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "malformed token rejected",
			authHeader:     "Bearer invalid.jwt.token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := mw.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.verifyHeader {
					forwardedID := r.Header.Get(middleware.HeaderUserID)
					if forwardedID != "uid-99" {
						t.Errorf("expected %s to be uid-99, got %q", middleware.HeaderUserID, forwardedID)
					}
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			if tc.spoofedUserID != "" {
				req.Header.Set(middleware.HeaderUserID, tc.spoofedUserID)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Fatalf("expected status %d, got %d", tc.expectedStatus, rec.Code)
			}
		})
	}
}

func TestIPRateLimiter(t *testing.T) {
	limiter := middleware.NewIPRateLimiter(10, 2)
	handler := limiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/share/slug-1", nil)
	req.RemoteAddr = "192.168.1.1:1234"

	// 1st request ok
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req)
	if rec1.Code != http.StatusOK {
		t.Errorf("req 1 failed: %d", rec1.Code)
	}

	// 2nd request ok (burst = 2)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req)
	if rec2.Code != http.StatusOK {
		t.Errorf("req 2 failed: %d", rec2.Code)
	}

	// 3rd immediate request must be 429
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req)
	if rec3.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429 Too Many Requests, got %d", rec3.Code)
	}
}