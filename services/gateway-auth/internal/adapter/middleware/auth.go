package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/usecase"
)

type contextKey string

const (
	UserContextKey  contextKey = "user_identity"
	HeaderUserID    string     = "X-User-Id"
	HeaderUserEmail string     = "X-User-Email"
)

type UserIdentity struct {
	ID    string
	Email string
}

func UserFromContext(ctx context.Context) (UserIdentity, bool) {
	val, ok := ctx.Value(UserContextKey).(UserIdentity)
	return val, ok
}

type AuthMiddleware struct {
	tokenManager usecase.TokenManager
}

func NewAuthMiddleware(tokenManager usecase.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenManager: tokenManager}
}

// StripUntrustedHeaders removes any caller-supplied X-User-Id / X-User-Email
// headers to eliminate spoofing attacks on downstream microservices.
func StripUntrustedHeaders(r *http.Request) {
	r.Header.Del(HeaderUserID)
	r.Header.Del(HeaderUserEmail)
}

// RequireAuth enforces a valid Bearer JWT. Injects identity into context
// and sets X-User-Id on the request for downstream forwarding.
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		StripUntrustedHeaders(r)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, `{"error":"invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		claims, err := m.tokenManager.ValidateAccessToken(parts[1])
		if err != nil {
			http.Error(w, `{"error":"invalid or expired access token"}`, http.StatusUnauthorized)
			return
		}

		identity := UserIdentity{
			ID:    claims.UserID,
			Email: claims.Email,
		}

		// Inject into context
		ctx := context.WithValue(r.Context(), UserContextKey, identity)
		r = r.WithContext(ctx)

		// Forward identity headers to internal services
		r.Header.Set(HeaderUserID, identity.ID)
		r.Header.Set(HeaderUserEmail, identity.Email)

		next.ServeHTTP(w, r)
	})
}

// OptionalAuth parses Bearer JWT if present, but allows guest access if absent.
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		StripUntrustedHeaders(r)

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				if claims, err := m.tokenManager.ValidateAccessToken(parts[1]); err == nil {
					identity := UserIdentity{
						ID:    claims.UserID,
						Email: claims.Email,
					}
					ctx := context.WithValue(r.Context(), UserContextKey, identity)
					r = r.WithContext(ctx)
					r.Header.Set(HeaderUserID, identity.ID)
					r.Header.Set(HeaderUserEmail, identity.Email)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}