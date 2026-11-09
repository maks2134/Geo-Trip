package usecase

import (
	"context"
	"time"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/domain"
)

type TripRepository interface {
	// Trips
	CreateTrip(ctx context.Context, t domain.Trip) error
	GetTrip(ctx context.Context, id string) (domain.Trip, error)
	ListTrips(ctx context.Context, ownerID string, status *domain.TripStatus, limit, offset int32) ([]domain.Trip, int32, error)
	UpdateTrip(ctx context.Context, t domain.Trip) error
	DeleteTrip(ctx context.Context, id string) error

	// Days
	CreateDay(ctx context.Context, d domain.TripDay) error
	GetDay(ctx context.Context, id string) (domain.TripDay, error)
	ListDays(ctx context.Context, tripID string) ([]domain.TripDay, error)
	UpdateDay(ctx context.Context, d domain.TripDay) error
	DeleteDay(ctx context.Context, id string) error

	// Itinerary
	CreatePlaceItem(ctx context.Context, it domain.ItineraryPlaceItem) error
	GetPlaceItem(ctx context.Context, id string) (domain.ItineraryPlaceItem, error)
	ListPlaceItems(ctx context.Context, tripID string) ([]domain.ItineraryPlaceItem, error)
	UpdatePlaceItem(ctx context.Context, it domain.ItineraryPlaceItem) error
	DeletePlaceItem(ctx context.Context, tripID, id string) error

	// Expenses
	CreateExpense(ctx context.Context, e domain.Expense) error
	ListExpenses(ctx context.Context, tripID string, dayID *string, category *domain.ExpenseCategory) ([]domain.Expense, error)
	DeleteExpense(ctx context.Context, tripID, id string) error

	// Checklist
	CreateChecklistItem(ctx context.Context, it domain.ChecklistItem) error
	GetChecklistItem(ctx context.Context, id string) (domain.ChecklistItem, error)
	ListChecklistItems(ctx context.Context, tripID string, section *domain.ChecklistSection) ([]domain.ChecklistItem, error)
	UpdateChecklistItem(ctx context.Context, it domain.ChecklistItem) error

	// Attachments
	CreateAttachment(ctx context.Context, a domain.Attachment) error
	GetAttachment(ctx context.Context, id string) (domain.Attachment, error)
	ListAttachments(ctx context.Context, tripID string) ([]domain.Attachment, error)
	ListAttachmentsFiltered(ctx context.Context, tripID string, placeID, dayID *string, visibility *domain.AttachmentVisibility) ([]domain.Attachment, error)
	ConfirmAttachmentReady(ctx context.Context, id string) error
	DeleteAttachment(ctx context.Context, tripID, id string) error

	// Share
	CreateShare(ctx context.Context, s domain.ShareRecord) error
	GetShareBySlug(ctx context.Context, slug string) (domain.ShareRecord, error)
	RevokeShare(ctx context.Context, tripID string) error
}

type StoragePort interface {
	PresignPut(ctx context.Context, key, contentType string, expires time.Duration) (string, error)
	PresignGet(ctx context.Context, key string, expires time.Duration) (string, error)
	DeleteObject(ctx context.Context, key string) error
}

type Clock interface {
	Now() time.Time
}

type IDGen interface {
	New() string
}
