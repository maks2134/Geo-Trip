package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/adapter/repo"
	"github.com/makskozlov/atlasflow/services/trip-core/internal/domain"
	"github.com/makskozlov/atlasflow/services/trip-core/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

type simpleIDGen struct{ cnt int }

func (g *simpleIDGen) New() string {
	g.cnt++
	return "id-" + string(rune('0'+g.cnt))
}

type fakeStorage struct{}

func (f *fakeStorage) PresignPut(ctx context.Context, key, contentType string, exp time.Duration) (string, error) {
	return "https://storage.local/upload/" + key, nil
}
func (f *fakeStorage) PresignGet(ctx context.Context, key string, exp time.Duration) (string, error) {
	return "https://storage.local/download/" + key, nil
}

func setupTest(t *testing.T) (*repo.MemoryTripRepo, *TripSuite) {
	memRepo := repo.NewMemoryTripRepo()
	clock := fixedClock{t: time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)}
	idGen := &simpleIDGen{}
	storage := &fakeStorage{}

	return memRepo, &TripSuite{
		tripUsecase:       usecase.NewTripUsecase(memRepo, clock, idGen),
		dayUsecase:        usecase.NewDayUsecase(memRepo, clock, idGen),
		itineraryUsecase:  usecase.NewItineraryUsecase(memRepo, clock, idGen),
		budgetUsecase:     usecase.NewBudgetUsecase(memRepo, clock, idGen),
		checklistUsecase:  usecase.NewChecklistUsecase(memRepo, clock, idGen),
		attachmentUsecase: usecase.NewAttachmentUsecase(memRepo, storage, clock, idGen),
		shareUsecase:      usecase.NewPublicShareUsecase(memRepo, clock, idGen),
	}
}

type TripSuite struct {
	tripUsecase       *usecase.TripUsecase
	dayUsecase        *usecase.DayUsecase
	itineraryUsecase  *usecase.ItineraryUsecase
	budgetUsecase     *usecase.BudgetUsecase
	checklistUsecase  *usecase.ChecklistUsecase
	attachmentUsecase *usecase.AttachmentUsecase
	shareUsecase      *usecase.PublicShareUsecase
}

func TestTripAuthz(t *testing.T) {
	_, suite := setupTest(t)
	ctx := context.Background()

	trip, err := suite.tripUsecase.CreateTrip(ctx, usecase.CreateTripInput{
		OwnerID: "owner-1",
		Title:   "Tokyo Trip",
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		userID  string
		wantErr error
	}{
		{name: "Owner access allowed", userID: "owner-1", wantErr: nil},
		{name: "Non-owner forbidden", userID: "stranger-2", wantErr: domain.ErrForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := suite.tripUsecase.GetTrip(ctx, trip.ID, tt.userID)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestItineraryLexorankReorder(t *testing.T) {
	_, suite := setupTest(t)
	ctx := context.Background()

	trip, err := suite.tripUsecase.CreateTrip(ctx, usecase.CreateTripInput{
		OwnerID: "owner-1",
		Title:   "Paris Trip",
	})
	require.NoError(t, err)

	day, err := suite.dayUsecase.AddDay(ctx, usecase.AddDayInput{
		TripID: trip.ID,
		UserID: "owner-1",
		Title:  "Day 1",
	})
	require.NoError(t, err)

	p1, err := suite.itineraryUsecase.AddPlace(ctx, usecase.AddPlaceInput{
		TripID:   trip.ID,
		UserID:   "owner-1",
		DayID:    &day.ID,
		Position: "a0",
		Place:    domain.PlaceRef{Name: "Eiffel Tower"},
	})
	require.NoError(t, err)

	p2, err := suite.itineraryUsecase.AddPlace(ctx, usecase.AddPlaceInput{
		TripID:   trip.ID,
		UserID:   "owner-1",
		DayID:    &day.ID,
		Position: "a5",
		Place:    domain.PlaceRef{Name: "Louvre"},
	})
	require.NoError(t, err)

	p3, err := suite.itineraryUsecase.AddPlace(ctx, usecase.AddPlaceInput{
		TripID:   trip.ID,
		UserID:   "owner-1",
		DayID:    &day.ID,
		Position: "b0",
		Place:    domain.PlaceRef{Name: "Montmartre"},
	})
	require.NoError(t, err)

	reordered, err := suite.itineraryUsecase.ReorderPlace(ctx, usecase.ReorderPlaceInput{
		TripID:       trip.ID,
		UserID:       "owner-1",
		PlaceItemID:  p3.ID,
		TargetDayID:  &day.ID,
		PrevPosition: p1.Position,
		NextPosition: p2.Position,
	})
	require.NoError(t, err)

	assert.True(t, reordered.Position > p1.Position, "new pos must be greater than prev")
	assert.True(t, reordered.Position < p2.Position, "new pos must be smaller than next")

	backlogItem, err := suite.itineraryUsecase.AddPlace(ctx, usecase.AddPlaceInput{
		TripID:   trip.ID,
		UserID:   "owner-1",
		DayID:    nil,
		Position: "c0",
		Place:    domain.PlaceRef{Name: "Notre Dame"},
	})
	require.NoError(t, err)
	assert.Nil(t, backlogItem.DayID)

	moved, err := suite.itineraryUsecase.ReorderPlace(ctx, usecase.ReorderPlaceInput{
		TripID:       trip.ID,
		UserID:       "owner-1",
		PlaceItemID:  backlogItem.ID,
		TargetDayID:  &day.ID,
		PrevPosition: reordered.Position,
		NextPosition: p2.Position,
	})
	require.NoError(t, err)
	assert.NotNil(t, moved.DayID)
	assert.Equal(t, day.ID, *moved.DayID)
}

func TestBudgetMath(t *testing.T) {
	_, suite := setupTest(t)
	ctx := context.Background()

	trip, err := suite.tripUsecase.CreateTrip(ctx, usecase.CreateTripInput{
		OwnerID:     "owner-1",
		Title:       "Budget Trip",
		BudgetLimit: 1000.00,
		Currency:    "EUR",
	})
	require.NoError(t, err)

	_, summary, err := suite.budgetUsecase.AddExpense(ctx, usecase.AddExpenseInput{
		TripID:   trip.ID,
		UserID:   "owner-1",
		Category: domain.ExpenseCategoryFood,
		Amount:   50.25,
		Title:    "Dinner",
	})
	require.NoError(t, err)
	assert.Equal(t, 50.25, summary.TotalSpent)
	assert.Equal(t, 50.25, summary.SpentByCategory["2"])

	_, summary2, err := suite.budgetUsecase.AddExpense(ctx, usecase.AddExpenseInput{
		TripID:   trip.ID,
		UserID:   "owner-1",
		Category: domain.ExpenseCategoryStay,
		Amount:   200.50,
		Title:    "Hotel",
	})
	require.NoError(t, err)
	assert.Equal(t, 250.75, summary2.TotalSpent)
	assert.Equal(t, 200.50, summary2.SpentByCategory["1"])
	assert.Equal(t, 1000.00, summary2.BudgetLimit)
}

func TestChecklistReadiness(t *testing.T) {
	_, suite := setupTest(t)
	ctx := context.Background()

	trip, err := suite.tripUsecase.CreateTrip(ctx, usecase.CreateTripInput{
		OwnerID: "owner-1",
		Title:   "Checklist Trip",
	})
	require.NoError(t, err)

	it1, s1, err := suite.checklistUsecase.AddChecklistItem(ctx, usecase.AddChecklistItemInput{
		TripID:  trip.ID,
		UserID:  "owner-1",
		Section: domain.ChecklistSectionPreFlight,
		Title:   "Passport",
	})
	require.NoError(t, err)
	assert.Equal(t, int32(1), s1.TotalItems)
	assert.Equal(t, 0.0, s1.ReadinessPercentage)

	_, _, _ = suite.checklistUsecase.AddChecklistItem(ctx, usecase.AddChecklistItemInput{
		TripID:  trip.ID,
		UserID:  "owner-1",
		Section: domain.ChecklistSectionPacking,
		Title:   "Toothbrush",
	})

	_, _, _ = suite.checklistUsecase.AddChecklistItem(ctx, usecase.AddChecklistItemInput{
		TripID:  trip.ID,
		UserID:  "owner-1",
		Section: domain.ChecklistSectionBuyOnSite,
		Title:   "SIM card",
	})

	_, s2, err := suite.checklistUsecase.ToggleItem(ctx, trip.ID, it1.ID, "owner-1", true)
	require.NoError(t, err)
	assert.Equal(t, int32(3), s2.TotalItems)
	assert.Equal(t, int32(1), s2.CompletedItems)
	assert.Equal(t, 33.3, s2.ReadinessPercentage)
}

func TestPublicShareSanitization(t *testing.T) {
	_, suite := setupTest(t)
	ctx := context.Background()

	trip, err := suite.tripUsecase.CreateTrip(ctx, usecase.CreateTripInput{
		OwnerID:     "owner-1",
		Title:       "Secret Trip",
		BudgetLimit: 5000,
	})
	require.NoError(t, err)

	_, _, err = suite.budgetUsecase.AddExpense(ctx, usecase.AddExpenseInput{
		TripID:   trip.ID,
		UserID:   "owner-1",
		Category: domain.ExpenseCategoryFood,
		Amount:   120,
		Title:    "Secret dinner",
	})
	require.NoError(t, err)

	pubMeta, err := suite.attachmentUsecase.CreateUpload(ctx, usecase.CreateUploadInput{
		TripID:      trip.ID,
		UserID:      "owner-1",
		Filename:    "ticket.pdf",
		ContentType: "application/pdf",
		FileSize:    1024,
		Visibility:  domain.AttachmentVisibilityPublic,
	})
	require.NoError(t, err)
	err = suite.attachmentUsecase.ConfirmUpload(ctx, trip.ID, pubMeta.AttachmentID, "owner-1")
	require.NoError(t, err)

	privMeta, err := suite.attachmentUsecase.CreateUpload(ctx, usecase.CreateUploadInput{
		TripID:      trip.ID,
		UserID:      "owner-1",
		Filename:    "passport.jpg",
		ContentType: "image/jpeg",
		FileSize:    2048,
		Visibility:  domain.AttachmentVisibilityPrivate,
	})
	require.NoError(t, err)
	_ = suite.attachmentUsecase.ConfirmUpload(ctx, trip.ID, privMeta.AttachmentID, "owner-1")

	slug, err := suite.shareUsecase.EnableShare(ctx, trip.ID, "owner-1")
	require.NoError(t, err)

	pubView, err := suite.shareUsecase.GetPublicTripBySlug(ctx, slug)
	require.NoError(t, err)

	assert.Equal(t, trip.Title, pubView.Title)
	assert.Len(t, pubView.Attachments, 1)
	assert.Equal(t, "ticket.pdf", pubView.Attachments[0].Filename)

	err = suite.shareUsecase.RevokeShare(ctx, trip.ID, "owner-1")
	require.NoError(t, err)

	_, err = suite.shareUsecase.GetPublicTripBySlug(ctx, slug)
	assert.ErrorIs(t, err, domain.ErrShareRevoked)
}
