package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/domain"
)

type PublicShareUsecase struct {
	repo  TripRepository
	clock Clock
	ids   IDGen
}

func NewPublicShareUsecase(repo TripRepository, clock Clock, ids IDGen) *PublicShareUsecase {
	return &PublicShareUsecase{repo: repo, clock: clock, ids: ids}
}

func generateSlug() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	s := base64.RawURLEncoding.EncodeToString(b)
	return strings.ToLower(strings.ReplaceAll(s, "-", ""))[:10]
}

func (u *PublicShareUsecase) EnableShare(ctx context.Context, tripID, userID string) (string, error) {
	t, err := u.repo.GetTrip(ctx, tripID)
	if err != nil {
		return "", err
	}
	if err := t.AssertOwner(userID); err != nil {
		return "", err
	}

	slug := generateSlug()
	rec := domain.ShareRecord{
		ID:              u.ids.New(),
		TripID:          tripID,
		Slug:            slug,
		HideBudget:      true,
		HidePrivateDocs: true,
		CreatedAt:       u.clock.Now(),
	}

	if err := u.repo.CreateShare(ctx, rec); err != nil {
		return "", fmt.Errorf("create share: %w", err)
	}

	t.ShareSlug = slug
	t.IsPublic = true
	t.UpdatedAt = u.clock.Now()
	_ = u.repo.UpdateTrip(ctx, t)

	return slug, nil
}

func (u *PublicShareUsecase) RevokeShare(ctx context.Context, tripID, userID string) error {
	t, err := u.repo.GetTrip(ctx, tripID)
	if err != nil {
		return err
	}
	if err := t.AssertOwner(userID); err != nil {
		return err
	}

	t.ShareSlug = ""
	t.IsPublic = false
	t.UpdatedAt = u.clock.Now()
	_ = u.repo.UpdateTrip(ctx, t)

	return u.repo.RevokeShare(ctx, tripID)
}

func (u *PublicShareUsecase) GetPublicTripBySlug(ctx context.Context, slug string) (domain.PublicTripView, error) {
	share, err := u.repo.GetShareBySlug(ctx, slug)
	if err != nil {
		return domain.PublicTripView{}, domain.ErrNotFound
	}
	if !share.IsActive() {
		return domain.PublicTripView{}, domain.ErrShareRevoked
	}

	t, err := u.repo.GetTrip(ctx, share.TripID)
	if err != nil {
		return domain.PublicTripView{}, domain.ErrNotFound
	}

	days, err := u.repo.ListDays(ctx, t.ID)
	if err != nil {
		return domain.PublicTripView{}, err
	}

	items, err := u.repo.ListPlaceItems(ctx, t.ID)
	if err != nil {
		return domain.PublicTripView{}, err
	}

	attachments, err := u.repo.ListAttachments(ctx, t.ID)
	if err != nil {
		return domain.PublicTripView{}, err
	}

	// Strictly sanitizes trip: strips budget, expenses, and private attachments
	return domain.SanitizeForPublic(t, days, items, attachments, "Traveler"), nil
}