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

// ClonePublicTrip creates a deep copy of a publicly shared trip for the authenticated target user.
// Strictly strips original owner's budget limits, expenses, checklists, and private attachments.
func (u *PublicShareUsecase) ClonePublicTrip(ctx context.Context, slug string, targetOwnerID string) (domain.Trip, error) {
	if strings.TrimSpace(targetOwnerID) == "" {
		return domain.Trip{}, domain.ErrForbidden
	}

	share, err := u.repo.GetShareBySlug(ctx, slug)
	if err != nil {
		return domain.Trip{}, domain.ErrNotFound
	}
	if !share.IsActive() {
		return domain.Trip{}, domain.ErrShareRevoked
	}

	origTrip, err := u.repo.GetTrip(ctx, share.TripID)
	if err != nil {
		return domain.Trip{}, domain.ErrNotFound
	}

	now := u.clock.Now()
	newTripID := u.ids.New()

	clonedTrip := domain.Trip{
		ID:          newTripID,
		OwnerID:     targetOwnerID,
		Title:       fmt.Sprintf("%s (Copy)", origTrip.Title),
		Description: origTrip.Description,
		StartDate:   origTrip.StartDate,
		EndDate:     origTrip.EndDate,
		Status:      domain.TripStatusDraft,
		CoverURL:    origTrip.CoverURL,
		Currency:    origTrip.Currency,
		BudgetLimit: 0, // Strict privacy & reset: no budget limit copied
		ShareSlug:   "",
		IsPublic:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.repo.CreateTrip(ctx, clonedTrip); err != nil {
		return domain.Trip{}, fmt.Errorf("create cloned trip: %w", err)
	}

	// Clone days
	origDays, err := u.repo.ListDays(ctx, origTrip.ID)
	if err != nil {
		return domain.Trip{}, err
	}

	dayIDMap := make(map[string]string) // oldDayID -> newDayID
	for _, od := range origDays {
		newDayID := u.ids.New()
		dayIDMap[od.ID] = newDayID
		clonedDay := domain.TripDay{
			ID:       newDayID,
			TripID:   newTripID,
			DayIndex: od.DayIndex,
			Date:     od.Date,
			Title:    od.Title,
			Notes:    od.Notes,
		}
		if err := u.repo.CreateDay(ctx, clonedDay); err != nil {
			return domain.Trip{}, fmt.Errorf("clone day: %w", err)
		}
	}

	// Clone places (itinerary items)
	origItems, err := u.repo.ListPlaceItems(ctx, origTrip.ID)
	if err != nil {
		return domain.Trip{}, err
	}

	for _, oi := range origItems {
		var newDayID *string
		if oi.DayID != nil && *oi.DayID != "" {
			if mapped, ok := dayIDMap[*oi.DayID]; ok {
				newDayID = &mapped
			}
		}

		clonedItem := domain.ItineraryPlaceItem{
			ID:              u.ids.New(),
			TripID:          newTripID,
			DayID:           newDayID,
			Place:           oi.Place,
			Position:        oi.Position,
			DurationMinutes: oi.DurationMinutes,
			Cost:            nil, // Strict privacy: do not copy place cost
			Notes:           oi.Notes,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		if err := u.repo.CreatePlaceItem(ctx, clonedItem); err != nil {
			return domain.Trip{}, fmt.Errorf("clone place item: %w", err)
		}
	}

	// Clone public-only attachments
	origAttachments, _ := u.repo.ListAttachments(ctx, origTrip.ID)
	for _, oa := range origAttachments {
		if oa.Visibility == domain.AttachmentVisibilityPublic && oa.IsReady {
			var newPlaceID *string // places got new IDs, so attachments are cloned as trip attachments
			clonedAtt := domain.Attachment{
				ID:          u.ids.New(),
				TripID:      newTripID,
				PlaceID:     newPlaceID,
				DayID:       nil,
				Filename:    oa.Filename,
				ContentType: oa.ContentType,
				FileSize:    oa.FileSize,
				StorageKey:  oa.StorageKey,
				DownloadURL: oa.DownloadURL,
				Visibility:  domain.AttachmentVisibilityPublic,
				IsReady:     true,
				OwnerType:   oa.OwnerType,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			_ = u.repo.CreateAttachment(ctx, clonedAtt)
		}
	}

	return clonedTrip, nil
}
