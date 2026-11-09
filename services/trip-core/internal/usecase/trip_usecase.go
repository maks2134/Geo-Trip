package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/domain"
)

type TripUsecase struct {
	repo  TripRepository
	clock Clock
	ids   IDGen
}

func NewTripUsecase(repo TripRepository, clock Clock, ids IDGen) *TripUsecase {
	return &TripUsecase{repo: repo, clock: clock, ids: ids}
}

type CreateTripInput struct {
	OwnerID     string
	Title       string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	CoverURL    string
	Currency    string
	BudgetLimit float64
}

func (u *TripUsecase) CreateTrip(ctx context.Context, in CreateTripInput) (domain.Trip, error) {
	if in.OwnerID == "" {
		return domain.Trip{}, domain.ErrForbidden
	}
	if in.Title == "" {
		return domain.Trip{}, fmt.Errorf("%w: title is required", domain.ErrValidation)
	}
	now := u.clock.Now()
	currency := in.Currency
	if currency == "" {
		currency = "USD"
	}

	t := domain.Trip{
		ID:          u.ids.New(),
		OwnerID:     in.OwnerID,
		Title:       in.Title,
		Description: in.Description,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
		Status:      domain.TripStatusDraft,
		CoverURL:    in.CoverURL,
		Currency:    currency,
		BudgetLimit: in.BudgetLimit,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.repo.CreateTrip(ctx, t); err != nil {
		return domain.Trip{}, fmt.Errorf("create trip: %w", err)
	}
	return t, nil
}

func (u *TripUsecase) GetTrip(ctx context.Context, tripID, userID string) (domain.Trip, error) {
	t, err := u.repo.GetTrip(ctx, tripID)
	if err != nil {
		return domain.Trip{}, err
	}
	if err := t.AssertOwner(userID); err != nil {
		return domain.Trip{}, err
	}
	return t, nil
}

func (u *TripUsecase) ListTrips(ctx context.Context, ownerID string, status *domain.TripStatus, limit, offset int32) ([]domain.Trip, int32, error) {
	if ownerID == "" {
		return nil, 0, domain.ErrForbidden
	}
	return u.repo.ListTrips(ctx, ownerID, status, limit, offset)
}

type UpdateTripInput struct {
	TripID      string
	UserID      string
	Title       *string
	Description *string
	StartDate   *time.Time
	EndDate     *time.Time
	Status      *domain.TripStatus
	CoverURL    *string
	Currency    *string
	BudgetLimit *float64
	IsPublic    *bool
}

func (u *TripUsecase) UpdateTrip(ctx context.Context, in UpdateTripInput) (domain.Trip, error) {
	t, err := u.repo.GetTrip(ctx, in.TripID)
	if err != nil {
		return domain.Trip{}, err
	}
	if err := t.AssertOwner(in.UserID); err != nil {
		return domain.Trip{}, err
	}

	if in.Title != nil {
		t.Title = *in.Title
	}
	if in.Description != nil {
		t.Description = *in.Description
	}
	if in.StartDate != nil {
		t.StartDate = *in.StartDate
	}
	if in.EndDate != nil {
		t.EndDate = *in.EndDate
	}
	if in.Status != nil {
		t.Status = *in.Status
	}
	if in.CoverURL != nil {
		t.CoverURL = *in.CoverURL
	}
	if in.Currency != nil {
		t.Currency = *in.Currency
	}
	if in.BudgetLimit != nil {
		t.BudgetLimit = *in.BudgetLimit
	}
	if in.IsPublic != nil {
		t.IsPublic = *in.IsPublic
	}
	t.UpdatedAt = u.clock.Now()

	if err := u.repo.UpdateTrip(ctx, t); err != nil {
		return domain.Trip{}, fmt.Errorf("update trip: %w", err)
	}
	return t, nil
}

func (u *TripUsecase) DeleteTrip(ctx context.Context, tripID, userID string) error {
	t, err := u.repo.GetTrip(ctx, tripID)
	if err != nil {
		return err
	}
	if err := t.AssertOwner(userID); err != nil {
		return err
	}
	return u.repo.DeleteTrip(ctx, tripID)
}