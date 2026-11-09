package usecase

import (
	"context"
	"fmt"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/domain"
)

type ItineraryUsecase struct {
	repo  TripRepository
	clock Clock
	ids   IDGen
}

func NewItineraryUsecase(repo TripRepository, clock Clock, ids IDGen) *ItineraryUsecase {
	return &ItineraryUsecase{repo: repo, clock: clock, ids: ids}
}

type AddPlaceInput struct {
	TripID          string
	UserID          string
	Place           domain.PlaceRef
	DayID           *string
	Position        string
	DurationMinutes int32
	Cost            *float64
	Notes           string
}

func (u *ItineraryUsecase) AddPlace(ctx context.Context, in AddPlaceInput) (domain.ItineraryPlaceItem, error) {
	t, err := u.repo.GetTrip(ctx, in.TripID)
	if err != nil {
		return domain.ItineraryPlaceItem{}, err
	}
	if err := t.AssertOwner(in.UserID); err != nil {
		return domain.ItineraryPlaceItem{}, err
	}

	pos := in.Position
	if pos == "" {
		pos, _ = domain.RankBetween("", "")
	}

	now := u.clock.Now()
	item := domain.ItineraryPlaceItem{
		ID:              u.ids.New(),
		TripID:          in.TripID,
		DayID:           in.DayID,
		Place:           in.Place,
		Position:        pos,
		DurationMinutes: in.DurationMinutes,
		Cost:            in.Cost,
		Notes:           in.Notes,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := u.repo.CreatePlaceItem(ctx, item); err != nil {
		return domain.ItineraryPlaceItem{}, fmt.Errorf("create place item: %w", err)
	}

	// If cost is specified, auto-link an expense
	if in.Cost != nil && *in.Cost > 0 {
		exp := domain.Expense{
			ID:          u.ids.New(),
			TripID:      in.TripID,
			DayID:       in.DayID,
			PlaceItemID: &item.ID,
			Category:    domain.ExpenseCategoryTickets, // default category for place visits
			Amount:      *in.Cost,
			Currency:    t.Currency,
			Title:       fmt.Sprintf("Place: %s", in.Place.Name),
			CreatedAt:   now,
		}
		_ = u.repo.CreateExpense(ctx, exp)
	}

	return item, nil
}

type ReorderPlaceInput struct {
	TripID       string
	UserID       string
	PlaceItemID  string
	TargetDayID  *string
	PrevPosition string
	NextPosition string
}

func (u *ItineraryUsecase) ReorderPlace(ctx context.Context, in ReorderPlaceInput) (domain.ItineraryPlaceItem, error) {
	t, err := u.repo.GetTrip(ctx, in.TripID)
	if err != nil {
		return domain.ItineraryPlaceItem{}, err
	}
	if err := t.AssertOwner(in.UserID); err != nil {
		return domain.ItineraryPlaceItem{}, err
	}

	item, err := u.repo.GetPlaceItem(ctx, in.PlaceItemID)
	if err != nil {
		return domain.ItineraryPlaceItem{}, err
	}
	if item.TripID != in.TripID {
		return domain.ItineraryPlaceItem{}, domain.ErrNotFound
	}

	newPos, err := domain.RankBetween(in.PrevPosition, in.NextPosition)
	if err != nil {
		return domain.ItineraryPlaceItem{}, fmt.Errorf("%w: rank calculation failed: %v", domain.ErrValidation, err)
	}

	item.Position = newPos
	item.DayID = in.TargetDayID
	item.UpdatedAt = u.clock.Now()

	// Single row update without full renumbering!
	if err := u.repo.UpdatePlaceItem(ctx, item); err != nil {
		return domain.ItineraryPlaceItem{}, fmt.Errorf("update place item position: %w", err)
	}
	return item, nil
}

func (u *ItineraryUsecase) RemovePlace(ctx context.Context, tripID, itemID, userID string) error {
	t, err := u.repo.GetTrip(ctx, tripID)
	if err != nil {
		return err
	}
	if err := t.AssertOwner(userID); err != nil {
		return err
	}
	return u.repo.DeletePlaceItem(ctx, tripID, itemID)
}

type OptimizeDayRouteInput struct {
	TripID    string
	DayID     string
	UserID    string
	LockStart bool
	LockEnd   bool
}

func (u *ItineraryUsecase) OptimizeDayRoute(ctx context.Context, in OptimizeDayRouteInput) (domain.TSPOptimizationResult, error) {
	t, err := u.repo.GetTrip(ctx, in.TripID)
	if err != nil {
		return domain.TSPOptimizationResult{}, err
	}
	if err := t.AssertOwner(in.UserID); err != nil {
		return domain.TSPOptimizationResult{}, err
	}

	allItems, err := u.repo.ListPlaceItems(ctx, in.TripID)
	if err != nil {
		return domain.TSPOptimizationResult{}, fmt.Errorf("list places for day: %w", err)
	}

	var places []domain.ItineraryPlaceItem
	for _, it := range allItems {
		if it.DayID != nil && *it.DayID == in.DayID {
			places = append(places, it)
		}
	}

	res := domain.OptimizeDayRoute2Opt(places, in.LockStart, in.LockEnd)

	// Update Lexorank positions sequentially
	var prevPos string
	now := u.clock.Now()
	for i := range res.OrderedPlaces {
		newPos, _ := domain.RankBetween(prevPos, "")
		res.OrderedPlaces[i].Position = newPos
		res.OrderedPlaces[i].UpdatedAt = now
		_ = u.repo.UpdatePlaceItem(ctx, res.OrderedPlaces[i])
		prevPos = newPos
	}

	return res, nil
}