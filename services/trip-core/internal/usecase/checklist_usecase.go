package usecase

import (
	"context"
	"fmt"
	"strings"

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

type TemplateItem struct {
	Section domain.ChecklistSection
	Title   string
}

// Starter template definitions (CAND-06)
var starterTemplates = map[string][]TemplateItem{
	"digital_nomad": {
		{domain.ChecklistSectionPreFlight, "Passport & Visa"},
		{domain.ChecklistSectionPreFlight, "Travel insurance policy"},
		{domain.ChecklistSectionPreFlight, "Offline maps & boarding pass downloaded"},
		{domain.ChecklistSectionPacking, "Laptop & charger"},
		{domain.ChecklistSectionPacking, "Universal power adapter"},
		{domain.ChecklistSectionPacking, "Power bank 20000mAh"},
		{domain.ChecklistSectionPacking, "Noise-canceling headphones"},
		{domain.ChecklistSectionPacking, "Backup credit / debit cards"},
		{domain.ChecklistSectionBuyOnSite, "Local SIM / eSIM data plan"},
		{domain.ChecklistSectionBuyOnSite, "Local currency / cash for emergency"},
	},
	"beach_sun": {
		{domain.ChecklistSectionPreFlight, "Hotel / Airbnb reservation confirmation"},
		{domain.ChecklistSectionPreFlight, "Airport transfer arranged"},
		{domain.ChecklistSectionPacking, "Swimwear & beach towel"},
		{domain.ChecklistSectionPacking, "Sunscreen SPF 50+ & lip balm"},
		{domain.ChecklistSectionPacking, "Sunglasses & sun hat"},
		{domain.ChecklistSectionPacking, "Waterproof phone pouch"},
		{domain.ChecklistSectionBuyOnSite, "After-sun aloe vera gel"},
		{domain.ChecklistSectionBuyOnSite, "Mosquito repellent"},
	},
	"mountain_trekking": {
		{domain.ChecklistSectionPreFlight, "National park permits & trail maps"},
		{domain.ChecklistSectionPreFlight, "Emergency rescue insurance"},
		{domain.ChecklistSectionPacking, "Trail hiking boots (broken in)"},
		{domain.ChecklistSectionPacking, "Waterproof shell jacket & pants"},
		{domain.ChecklistSectionPacking, "First aid kit with blister bandages"},
		{domain.ChecklistSectionPacking, "Headlamp with extra batteries"},
		{domain.ChecklistSectionPacking, "Hydration bladder or water filter"},
		{domain.ChecklistSectionBuyOnSite, "Trail energy bars & electrolytes"},
		{domain.ChecklistSectionBuyOnSite, "Camp stove gas canister"},
	},
	"winter_city": {
		{domain.ChecklistSectionPreFlight, "Museum & attraction pre-booked slots"},
		{domain.ChecklistSectionPacking, "Thermal base layers & wool socks"},
		{domain.ChecklistSectionPacking, "Insulated winter coat"},
		{domain.ChecklistSectionPacking, "Gloves, scarf & beanie"},
		{domain.ChecklistSectionPacking, "Compact storm umbrella"},
		{domain.ChecklistSectionBuyOnSite, "Public transit pass"},
	},
}

func (u *ChecklistUsecase) ApplyTemplate(ctx context.Context, tripID, userID, templateKey string) ([]domain.ChecklistItem, domain.ChecklistSummary, error) {
	t, err := u.repo.GetTrip(ctx, tripID)
	if err != nil {
		return nil, domain.ChecklistSummary{}, err
	}
	if err := t.AssertOwner(userID); err != nil {
		return nil, domain.ChecklistSummary{}, err
	}

	key := strings.ToLower(strings.TrimSpace(templateKey))
	tplItems, ok := starterTemplates[key]
	if !ok {
		return nil, domain.ChecklistSummary{}, fmt.Errorf("unknown checklist template: %s", templateKey)
	}

	existing, err := u.repo.ListChecklistItems(ctx, tripID, nil)
	if err != nil {
		return nil, domain.ChecklistSummary{}, err
	}

	existingTitles := make(map[string]bool)
	for _, it := range existing {
		existingTitles[strings.ToLower(strings.TrimSpace(it.Title))] = true
	}

	now := u.clock.Now()
	var lastPos string
	if len(existing) > 0 {
		lastPos = existing[len(existing)-1].Position
	}

	for _, ti := range tplItems {
		cleanTitle := strings.TrimSpace(ti.Title)
		if existingTitles[strings.ToLower(cleanTitle)] {
			continue // skip duplicate title
		}

		newPos, _ := domain.RankBetween(lastPos, "")
		lastPos = newPos

		item := domain.ChecklistItem{
			ID:        u.ids.New(),
			TripID:    tripID,
			Section:   ti.Section,
			Title:     cleanTitle,
			IsDone:    false,
			Position:  newPos,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := u.repo.CreateChecklistItem(ctx, item); err != nil {
			return nil, domain.ChecklistSummary{}, err
		}
		existingTitles[strings.ToLower(cleanTitle)] = true
	}

	all, err := u.repo.ListChecklistItems(ctx, tripID, nil)
	if err != nil {
		return nil, domain.ChecklistSummary{}, err
	}

	return all, domain.CalculateChecklistSummary(all), nil
}
