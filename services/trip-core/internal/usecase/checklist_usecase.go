package usecase

import (
	"context"
	"fmt"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/domain"
)

type ChecklistUsecase struct {
	repo  TripRepository
	clock Clock
	ids   IDGen
}

func NewChecklistUsecase(repo TripRepository, clock Clock, ids IDGen) *ChecklistUsecase {
	return &ChecklistUsecase{repo: repo, clock: clock, ids: ids}
}

type AddChecklistItemInput struct {
	TripID   string
	UserID   string
	Section  domain.ChecklistSection
	Title    string
	Position string
}

func (u *ChecklistUsecase) AddChecklistItem(ctx context.Context, in AddChecklistItemInput) (domain.ChecklistItem, domain.ChecklistSummary, error) {
	t, err := u.repo.GetTrip(ctx, in.TripID)
	if err != nil {
		return domain.ChecklistItem{}, domain.ChecklistSummary{}, err
	}
	if err := t.AssertOwner(in.UserID); err != nil {
		return domain.ChecklistItem{}, domain.ChecklistSummary{}, err
	}

	pos := in.Position
	if pos == "" {
		pos, _ = domain.RankBetween("", "")
	}

	now := u.clock.Now()
	item := domain.ChecklistItem{
		ID:        u.ids.New(),
		TripID:    in.TripID,
		Section:   in.Section,
		Title:     in.Title,
		IsDone:    false,
		Position:  pos,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.repo.CreateChecklistItem(ctx, item); err != nil {
		return domain.ChecklistItem{}, domain.ChecklistSummary{}, fmt.Errorf("create checklist item: %w", err)
	}

	items, err := u.repo.ListChecklistItems(ctx, in.TripID, nil)
	if err != nil {
		return domain.ChecklistItem{}, domain.ChecklistSummary{}, err
	}

	return item, domain.CalculateChecklistSummary(items), nil
}

func (u *ChecklistUsecase) ToggleItem(ctx context.Context, tripID, itemID, userID string, isDone bool) (domain.ChecklistItem, domain.ChecklistSummary, error) {
	t, err := u.repo.GetTrip(ctx, tripID)
	if err != nil {
		return domain.ChecklistItem{}, domain.ChecklistSummary{}, err
	}
	if err := t.AssertOwner(userID); err != nil {
		return domain.ChecklistItem{}, domain.ChecklistSummary{}, err
	}

	item, err := u.repo.GetChecklistItem(ctx, itemID)
	if err != nil {
		return domain.ChecklistItem{}, domain.ChecklistSummary{}, err
	}

	item.IsDone = isDone
	item.UpdatedAt = u.clock.Now()

	if err := u.repo.UpdateChecklistItem(ctx, item); err != nil {
		return domain.ChecklistItem{}, domain.ChecklistSummary{}, err
	}

	items, err := u.repo.ListChecklistItems(ctx, tripID, nil)
	if err != nil {
		return domain.ChecklistItem{}, domain.ChecklistSummary{}, err
	}

	return item, domain.CalculateChecklistSummary(items), nil
}

func (u *ChecklistUsecase) ListItems(ctx context.Context, tripID, userID string, section *domain.ChecklistSection) ([]domain.ChecklistItem, domain.ChecklistSummary, error) {
	t, err := u.repo.GetTrip(ctx, tripID)
	if err != nil {
		return nil, domain.ChecklistSummary{}, err
	}
	if err := t.AssertOwner(userID); err != nil {
		return nil, domain.ChecklistSummary{}, err
	}

	filtered, err := u.repo.ListChecklistItems(ctx, tripID, section)
	if err != nil {
		return nil, domain.ChecklistSummary{}, err
	}

	all, err := u.repo.ListChecklistItems(ctx, tripID, nil)
	if err != nil {
		return nil, domain.ChecklistSummary{}, err
	}

	return filtered, domain.CalculateChecklistSummary(all), nil
}