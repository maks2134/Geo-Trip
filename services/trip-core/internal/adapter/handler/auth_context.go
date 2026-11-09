package handler

import (
	"context"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/domain"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

// ContextWithUserID returns a new context carrying the user ID.
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}

// GetUserID extracts the authenticated user ID from context.
func GetUserID(ctx context.Context) (string, error) {
	if val, ok := ctx.Value(UserIDContextKey).(string); ok && val != "" {
		return val, nil
	}
	// Also support plain string key or fallback for local development/internal calls
	return "", domain.ErrForbidden
}
