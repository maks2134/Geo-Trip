package authv1connect

import (
	"context"

	authv1 "github.com/makskozlov/atlasflow/proto/gen/go/atlasflow/auth/v1"
)

const (
	AuthServiceName                 = "atlasflow.auth.v1.AuthService"
	AuthServiceRegisterProcedure       = "/atlasflow.auth.v1.AuthService/Register"
	AuthServiceLoginProcedure          = "/atlasflow.auth.v1.AuthService/Login"
	AuthServiceRefreshTokenProcedure   = "/atlasflow.auth.v1.AuthService/RefreshToken"
	AuthServiceGetCurrentUserProcedure = "/atlasflow.auth.v1.AuthService/GetCurrentUser"
)

type AuthServiceHandler interface {
	Register(context.Context, *authv1.RegisterRequest) (*authv1.RegisterResponse, error)
	Login(context.Context, *authv1.LoginRequest) (*authv1.LoginResponse, error)
	RefreshToken(context.Context, *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error)
	GetCurrentUser(context.Context, *authv1.GetCurrentUserRequest) (*authv1.GetCurrentUserResponse, error)
}
