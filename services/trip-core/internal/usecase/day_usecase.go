package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/domain"
)

type DayUsecase struct {
	repo  TripRepository
	clock Clock
	ids   IDGen
}

func NewDayUsecase(repo TripRepository, clock Clock, ids IDGen) *DayUsecase {
	return &DayUsecase{repo: repo, clock: clock, ids: ids}
}

type AddDayInput struct {
	TripID string
	UserID string
	Date   time.Time
	Title  string
	Notes  string
}

func (u *DayUsecase) AddDay(ctx context.Context, in AddDayInput) (domain.TripDay, error) {
	t, err := u.repo.GetTrip(ctx, in.TripID)
	if err != nil {
		return domain.TripDay{}, err
	}
	if err := t.AssertOwner(in.UserID); err != nil {
		return domain.TripDay{}, err
	}

	days, err := u.repo.ListDays(ctx, in.TripID)
	if err != nil {
		return domain.TripDay{}, err
	}

	d := domain.TripDay{
		ID:       u.ids.New(),
		TripID:   in.TripID,
		DayIndex: int32(len(days) + 1),
		Date:     in.Date,
		Title:    in.Title,
		Notes:    in.Notes,
	}

	if err := u.repo.CreateDay(ctx, d); err != nil {
		return domain.TripDay{}, fmt.Errorf("create day: %w", err)
	}
	return d, nil
}

type UpdateDayInput struct {
	TripID string
	DayID  string
	UserID string
	Date   *time.Time
	Title  *string
	Notes  *string
}

func (u *DayUsecase) UpdateDay(ctx context.Context, in UpdateDayInput) (domain.TripDay, error) {
	t, err := u.repo.GetTrip(ctx, in.TripID)
	if err != nil {
		return domain.TripDay{}, err
	}
	if err := t.AssertOwner(in.UserID); err != nil {
		return domain.TripDay{}, err
	}

	d, err := u.repo.GetDay(ctx, in.DayID)
	if err != nil {
		return domain.TripDay{}, err
	}

	if in.Date != nil {
		d.Date = *in.Date
	}
	if in.Title != nil {
		d.Title = *in.Title
	}
	if in.Notes != nil {
		d.Notes = *in.Notes
	}

	if err := u.repo.UpdateDay(ctx, d); err != nil {
		return domain.TripDay{}, fmt.Errorf("update day: %w", err)
	}
	return d, nil
}

func (u *DayUsecase) DeleteDay(ctx context.Context, tripID, dayID, userID string) error {
	t, err := u.repo.GetTrip(ctx, tripID)
	if err != nil {
		return err
	}
	if err := t.AssertOwner(userID); err != nil {
		return err
	}
	return u.repo.DeleteDay(ctx, dayID)
}