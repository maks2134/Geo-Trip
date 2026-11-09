package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
)

type PublicShareUsecase struct {
	tripClient PublicTripClient
}

func NewPublicShareUsecase(tripClient PublicTripClient) *PublicShareUsecase {
	return &PublicShareUsecase{tripClient: tripClient}
}

// GetPublicTrip resolves the public trip view for guest visitors.
// Guarantees slug validation and ensures no private details leak.
func (u *PublicShareUsecase) GetPublicTrip(ctx context.Context, slug string) (*domain.PublicTrip, error) {
	cleanSlug := strings.TrimSpace(slug)
	if cleanSlug == "" || len(cleanSlug) < 3 || len(cleanSlug) > 128 {
		return nil, domain.ErrTripNotFound
	}

	trip, err := u.tripClient.GetPublicTrip(ctx, cleanSlug)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrTripNotFound, err)
	}
	if trip == nil {
		return nil, domain.ErrTripNotFound
	}

	return trip, nil
}

// ClonePublicTrip clones an active public trip into the authenticated user's workspace.
func (u *PublicShareUsecase) ClonePublicTrip(ctx context.Context, slug, targetOwnerID string) (*domain.ClonedTripSummary, error) {
	cleanSlug := strings.TrimSpace(slug)
	if cleanSlug == "" || len(cleanSlug) < 3 || len(cleanSlug) > 128 {
		return nil, domain.ErrTripNotFound
	}
	if strings.TrimSpace(targetOwnerID) == "" {
		return nil, domain.ErrUnauthorized
	}

	trip, err := u.tripClient.ClonePublicTrip(ctx, cleanSlug, targetOwnerID)
	if err != nil {
		return nil, fmt.Errorf("clone trip: %w", err)
	}
	if trip == nil {
		return nil, domain.ErrTripNotFound
	}

	return trip, nil
}
