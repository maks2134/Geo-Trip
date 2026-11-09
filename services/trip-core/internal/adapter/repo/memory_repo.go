package repo

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/domain"
	"github.com/makskozlov/atlasflow/services/trip-core/internal/usecase"
)

var _ usecase.TripRepository = (*MemoryTripRepo)(nil)

type MemoryTripRepo struct {
	mu           sync.RWMutex
	trips        map[string]domain.Trip
	days         map[string]domain.TripDay
	items        map[string]domain.ItineraryPlaceItem
	expenses     map[string]domain.Expense
	checklist    map[string]domain.ChecklistItem
	attachments  map[string]domain.Attachment
	shares       map[string]domain.ShareRecord
	sharesByTrip map[string]string
}

func NewMemoryTripRepo() *MemoryTripRepo {
	return &MemoryTripRepo{
		trips:        make(map[string]domain.Trip),
		days:         make(map[string]domain.TripDay),
		items:        make(map[string]domain.ItineraryPlaceItem),
		expenses:     make(map[string]domain.Expense),
		checklist:    make(map[string]domain.ChecklistItem),
		attachments:  make(map[string]domain.Attachment),
		shares:       make(map[string]domain.ShareRecord),
		sharesByTrip: make(map[string]string),
	}
}

// Trips
func (r *MemoryTripRepo) CreateTrip(ctx context.Context, t domain.Trip) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.trips[t.ID] = t
	return nil
}

func (r *MemoryTripRepo) GetTrip(ctx context.Context, id string) (domain.Trip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.trips[id]
	if !ok {
		return domain.Trip{}, domain.ErrNotFound
	}
	return t, nil
}

func (r *MemoryTripRepo) ListTrips(ctx context.Context, ownerID string, status *domain.TripStatus, limit, offset int32) ([]domain.Trip, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.Trip
	for _, t := range r.trips {
		if t.OwnerID != ownerID {
			continue
		}
		if status != nil && t.Status != *status {
			continue
		}
		list = append(list, t)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.After(list[j].CreatedAt)
	})
	total := int32(len(list))
	start := int(offset)
	if start > len(list) {
		return []domain.Trip{}, total, nil
	}
	end := start + int(limit)
	if end > len(list) || limit == 0 {
		end = len(list)
	}
	return list[start:end], total, nil
}

func (r *MemoryTripRepo) UpdateTrip(ctx context.Context, t domain.Trip) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.trips[t.ID]; !ok {
		return domain.ErrNotFound
	}
	r.trips[t.ID] = t
	return nil
}

func (r *MemoryTripRepo) DeleteTrip(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.trips, id)
	return nil
}

// Days
func (r *MemoryTripRepo) CreateDay(ctx context.Context, d domain.TripDay) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.days[d.ID] = d
	return nil
}

func (r *MemoryTripRepo) GetDay(ctx context.Context, id string) (domain.TripDay, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.days[id]
	if !ok {
		return domain.TripDay{}, domain.ErrNotFound
	}
	return d, nil
}

func (r *MemoryTripRepo) ListDays(ctx context.Context, tripID string) ([]domain.TripDay, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []domain.TripDay
	for _, d := range r.days {
		if d.TripID == tripID {
			res = append(res, d)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].DayIndex < res[j].DayIndex
	})
	return res, nil
}

func (r *MemoryTripRepo) UpdateDay(ctx context.Context, d domain.TripDay) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.days[d.ID]; !ok {
		return domain.ErrNotFound
	}
	r.days[d.ID] = d
	return nil
}

func (r *MemoryTripRepo) DeleteDay(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.days, id)
	return nil
}

// Itinerary
func (r *MemoryTripRepo) CreatePlaceItem(ctx context.Context, it domain.ItineraryPlaceItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[it.ID] = it
	return nil
}

func (r *MemoryTripRepo) GetPlaceItem(ctx context.Context, id string) (domain.ItineraryPlaceItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	it, ok := r.items[id]
	if !ok {
		return domain.ItineraryPlaceItem{}, domain.ErrNotFound
	}
	return it, nil
}

func (r *MemoryTripRepo) ListPlaceItems(ctx context.Context, tripID string) ([]domain.ItineraryPlaceItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []domain.ItineraryPlaceItem
	for _, it := range r.items {
		if it.TripID == tripID {
			res = append(res, it)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Position < res[j].Position
	})
	return res, nil
}

func (r *MemoryTripRepo) UpdatePlaceItem(ctx context.Context, it domain.ItineraryPlaceItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[it.ID]; !ok {
		return domain.ErrNotFound
	}
	r.items[it.ID] = it
	return nil
}

func (r *MemoryTripRepo) DeletePlaceItem(ctx context.Context, tripID, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}

// Expenses
func (r *MemoryTripRepo) CreateExpense(ctx context.Context, e domain.Expense) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.expenses[e.ID] = e
	return nil
}

func (r *MemoryTripRepo) ListExpenses(ctx context.Context, tripID string, dayID *string, category *domain.ExpenseCategory) ([]domain.Expense, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []domain.Expense
	for _, e := range r.expenses {
		if e.TripID != tripID {
			continue
		}
		if dayID != nil && (e.DayID == nil || *e.DayID != *dayID) {
			continue
		}
		if category != nil && e.Category != *category {
			continue
		}
		res = append(res, e)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].CreatedAt.Before(res[j].CreatedAt)
	})
	return res, nil
}

func (r *MemoryTripRepo) DeleteExpense(ctx context.Context, tripID, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.expenses, id)
	return nil
}

// Checklist
func (r *MemoryTripRepo) CreateChecklistItem(ctx context.Context, it domain.ChecklistItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checklist[it.ID] = it
	return nil
}

func (r *MemoryTripRepo) GetChecklistItem(ctx context.Context, id string) (domain.ChecklistItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	it, ok := r.checklist[id]
	if !ok {
		return domain.ChecklistItem{}, domain.ErrNotFound
	}
	return it, nil
}

func (r *MemoryTripRepo) ListChecklistItems(ctx context.Context, tripID string, section *domain.ChecklistSection) ([]domain.ChecklistItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []domain.ChecklistItem
	for _, it := range r.checklist {
		if it.TripID != tripID {
			continue
		}
		if section != nil && it.Section != *section {
			continue
		}
		res = append(res, it)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Position < res[j].Position
	})
	return res, nil
}

func (r *MemoryTripRepo) UpdateChecklistItem(ctx context.Context, it domain.ChecklistItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.checklist[it.ID]; !ok {
		return domain.ErrNotFound
	}
	r.checklist[it.ID] = it
	return nil
}

// Attachments
func (r *MemoryTripRepo) CreateAttachment(ctx context.Context, a domain.Attachment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.attachments[a.ID] = a
	return nil
}

func (r *MemoryTripRepo) GetAttachment(ctx context.Context, id string) (domain.Attachment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.attachments[id]
	if !ok {
		return domain.Attachment{}, domain.ErrNotFound
	}
	return a, nil
}

func (r *MemoryTripRepo) ListAttachments(ctx context.Context, tripID string) ([]domain.Attachment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []domain.Attachment
	for _, a := range r.attachments {
		if a.TripID == tripID {
			res = append(res, a)
		}
	}
	return res, nil
}

func (r *MemoryTripRepo) ConfirmAttachmentReady(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.attachments[id]
	if !ok {
		return domain.ErrNotFound
	}
	a.IsReady = true
	r.attachments[id] = a
	return nil
}

// Shares
func (r *MemoryTripRepo) CreateShare(ctx context.Context, s domain.ShareRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shares[s.Slug] = s
	r.sharesByTrip[s.TripID] = s.Slug
	return nil
}

func (r *MemoryTripRepo) GetShareBySlug(ctx context.Context, slug string) (domain.ShareRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.shares[slug]
	if !ok {
		return domain.ShareRecord{}, domain.ErrNotFound
	}
	return s, nil
}

func (r *MemoryTripRepo) RevokeShare(ctx context.Context, tripID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	slug, ok := r.sharesByTrip[tripID]
	if !ok {
		return nil
	}
	s := r.shares[slug]
	now := time.Now()
	s.RevokedAt = &now
	r.shares[slug] = s
	return nil
}