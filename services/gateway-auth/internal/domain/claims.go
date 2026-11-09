package domain

import "github.com/golang-jwt/jwt/v5"

type AuthClaims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}