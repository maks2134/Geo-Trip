package handler

import (
	"net/http"
	"strings"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/middleware"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/usecase"
)

type ShareHandler struct {
	shareUsecase *usecase.PublicShareUsecase
}

func NewShareHandler(shareUsecase *usecase.PublicShareUsecase) *ShareHandler {
	return &ShareHandler{
		shareUsecase: shareUsecase,
	}
}

type PublicTripResponse struct {
	Trip *domain.PublicTrip `json:"trip"`
}

type ClonedTripResponse struct {
	Trip *domain.ClonedTripSummary `json:"trip"`
}

// GetShareTrip handles public, unauthenticated trip view requests by slug.
// Rate limiting is applied and any incoming private headers are completely stripped.
// Output strictly contains sanitized fields (no budget, no private attachments).
func (h *ShareHandler) GetShareTrip(w http.ResponseWriter, r *http.Request) {
	middleware.StripUntrustedHeaders(r)

	slug := r.PathValue("slug")
	if slug == "" {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) > 0 {
			slug = parts[len(parts)-1]
		}
	}

	trip, err := h.shareUsecase.GetPublicTrip(r.Context(), slug)
	if err != nil {
		MapDomainError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, PublicTripResponse{Trip: trip})
}

// CloneSharedTrip clones an active public trip into the authenticated user's workspace (CAND-10).
func (h *ShareHandler) CloneSharedTrip(w http.ResponseWriter, r *http.Request) {
	identity, ok := middleware.UserFromContext(r.Context())
	if !ok || identity.ID == "" {
		WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "authentication required to clone trip"})
		return
	}

	slug := r.PathValue("slug")
	if slug == "" {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) > 0 {
			slug = parts[len(parts)-1]
		}
	}

	cloned, err := h.shareUsecase.ClonePublicTrip(r.Context(), slug, identity.ID)
	if err != nil {
		MapDomainError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, ClonedTripResponse{Trip: cloned})
}
