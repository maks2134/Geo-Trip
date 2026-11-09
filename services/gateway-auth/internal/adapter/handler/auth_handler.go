package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	authv1 "github.com/makskozlov/atlasflow/proto/gen/go/atlasflow/auth/v1"
	"github.com/makskozlov/atlasflow/proto/gen/go/atlasflow/auth/v1/authv1connect"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/middleware"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/usecase"
)

var _ authv1connect.AuthServiceHandler = (*AuthHandler)(nil)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// --- ConnectRPC AuthServiceHandler implementation ---

func (h *AuthHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	pair, user, err := h.authUsecase.Register(ctx, usecase.RegisterInput{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &authv1.RegisterResponse{
		User: &authv1.User{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.DisplayName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	pair, user, err := h.authUsecase.Login(ctx, usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{
		User: &authv1.User{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.DisplayName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	pair, err := h.authUsecase.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &authv1.RefreshTokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
	}, nil
}

func (h *AuthHandler) GetCurrentUser(ctx context.Context, _ *authv1.GetCurrentUserRequest) (*authv1.GetCurrentUserResponse, error) {
	identity, ok := middleware.UserFromContext(ctx)
	if !ok || identity.ID == "" {
		return nil, domain.ErrUnauthorized
	}

	user, err := h.authUsecase.GetCurrentUser(ctx, identity.ID)
	if err != nil {
		return nil, err
	}

	return &authv1.GetCurrentUserResponse{
		User: &authv1.User{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.DisplayName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

// --- HTTP REST Handlers ---

func (h *AuthHandler) RegisterHTTP(w http.ResponseWriter, r *http.Request) {
	var req authv1.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json request body")
		return
	}

	resp, err := h.Register(r.Context(), &req)
	if err != nil {
		MapDomainError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) LoginHTTP(w http.ResponseWriter, r *http.Request) {
	var req authv1.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json request body")
		return
	}

	resp, err := h.Login(r.Context(), &req)
	if err != nil {
		MapDomainError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) RefreshHTTP(w http.ResponseWriter, r *http.Request) {
	var req authv1.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json request body")
		return
	}

	resp, err := h.RefreshToken(r.Context(), &req)
	if err != nil {
		MapDomainError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) MeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := h.GetCurrentUser(r.Context(), &authv1.GetCurrentUserRequest{})
	if err != nil {
		MapDomainError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) LogoutHTTP(w http.ResponseWriter, r *http.Request) {
	var req authv1.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.RefreshToken != "" {
		_ = h.authUsecase.Logout(r.Context(), req.RefreshToken)
	}
	WriteJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// --- ConnectRPC HTTP Adapter Handler ---

func (h *AuthHandler) ServeConnectRPC(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case authv1connect.AuthServiceRegisterProcedure:
		h.RegisterHTTP(w, r)
	case authv1connect.AuthServiceLoginProcedure:
		h.LoginHTTP(w, r)
	case authv1connect.AuthServiceRefreshTokenProcedure:
		h.RefreshHTTP(w, r)
	case authv1connect.AuthServiceGetCurrentUserProcedure:
		h.MeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

func MapDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		WriteError(w, http.StatusConflict, "email already registered")
	case errors.Is(err, domain.ErrInvalidCredentials):
		WriteError(w, http.StatusUnauthorized, "invalid email or password")
	case errors.Is(err, domain.ErrInvalidToken), errors.Is(err, domain.ErrTokenExpired), errors.Is(err, domain.ErrTokenRevoked):
		WriteError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrTripNotFound):
		WriteError(w, http.StatusNotFound, "not found")
	case errors.Is(err, domain.ErrRateLimitExceeded):
		WriteError(w, http.StatusTooManyRequests, "rate limit exceeded")
	case errors.Is(err, domain.ErrUnauthorized):
		WriteError(w, http.StatusUnauthorized, "unauthorized")
	default:
		WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}