package domain

import "time"

type TripStatus string

const (
	TripStatusDraft     TripStatus = "draft"
	TripStatusUpcoming  TripStatus = "upcoming"
	TripStatusActive    TripStatus = "active"
	TripStatusCompleted TripStatus = "completed"
)

type Trip struct {
	ID          string
	OwnerID     string
	Title       string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	Status      TripStatus
	CoverURL    string
	Currency    string
	BudgetLimit float64
	ShareSlug   string
	IsPublic    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t Trip) AssertOwner(userID string) error {
	if t.OwnerID == "" || t.OwnerID != userID {
		return ErrForbidden
	}
	return nil
}