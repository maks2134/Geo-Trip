package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenRevoked       = errors.New("token revoked")
	ErrInvalidInput       = errors.New("invalid input")
	ErrRateLimitExceeded  = errors.New("rate limit exceeded")
	ErrTripNotFound       = errors.New("trip not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInternal           = errors.New("internal server error")
)