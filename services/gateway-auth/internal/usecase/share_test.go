package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/usecase"
)

type fakeTripClient struct {
	trips map[string]*domain.PublicTrip
}

func (f *fakeTripClient) GetPublicTrip(ctx context.Context, slug string) (*domain.PublicTrip, error) {
	trip, ok := f.trips[slug]
	if !ok {
		return nil, errors.New("not found in trip-core")
	}
	return trip, nil
}

func (f *fakeTripClient) ClonePublicTrip(ctx context.Context, slug, targetOwnerID string) (*domain.ClonedTripSummary, error) {
	trip, ok := f.trips[slug]
	if !ok {
		return nil, domain.ErrTripNotFound
	}
	return &domain.ClonedTripSummary{
		ID:          "cloned-1",
		OwnerID:     targetOwnerID,
		Title:       trip.Title + " (Copy)",
		Description: trip.Description,
		StartDate:   trip.StartDate,
		EndDate:     trip.EndDate,
		Status:      "draft",
		CoverURL:    trip.CoverURL,
		BudgetLimit: 0,
	}, nil
}

func TestPublicShareUsecase_GetPublicTrip(t *testing.T) {
	client := &fakeTripClient{
		trips: map[string]*domain.PublicTrip{
			"tokyo-drift-2026": {
				ID:               "trip-1",
				Slug:             "tokyo-drift-2026",
				Title:            "Tokyo Explorer",
				Description:      "7 days in Tokyo",
				StartDate:        time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
				CoverURL:         "https://example.com/tokyo.jpg",
				OwnerDisplayName: "Alex",
				Days: []domain.PublicTripDay{
					{
						ID:       "day-1",
						DayIndex: 1,
						Title:    "Arrival & Shibuya",
						Places: []domain.PublicPlaceItem{
							{
								ID:              "p-1",
								Name:            "Shibuya Crossing",
								Category:        "SIGHT",
								Latitude:        35.6595,
								Longitude:       139.7004,
								DurationMinutes: 60,
							},
						},
					},
				},
			},
		},
	}

	uc := usecase.NewPublicShareUsecase(client)

	tests := []struct {
		name    string
		slug    string
		wantErr bool
	}{
		{
			name:    "valid active slug",
			slug:    "tokyo-drift-2026",
			wantErr: false,
		},
		{
			name:    "empty slug",
			slug:    "   ",
			wantErr: true,
		},
		{
			name:    "nonexistent or revoked slug",
			slug:    "unknown-trip-slug",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trip, err := uc.GetPublicTrip(context.Background(), tc.slug)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for slug %q, got nil", tc.slug)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if trip.Slug != tc.slug || trip.Title != "Tokyo Explorer" {
				t.Errorf("unexpected trip payload: %+v", trip)
			}
		})
	}

	t.Run("clone public trip success", func(t *testing.T) {
		cloned, err := uc.ClonePublicTrip(context.Background(), "tokyo-drift-2026", "user-99")
		if err != nil {
			t.Fatalf("unexpected error cloning trip: %v", err)
		}
		if cloned.OwnerID != "user-99" {
			t.Errorf("expected owner user-99, got %s", cloned.OwnerID)
		}
		if cloned.BudgetLimit != 0 {
			t.Errorf("expected 0 budget limit on cloned trip, got %f", cloned.BudgetLimit)
		}
	})

	t.Run("clone public trip unauthenticated error", func(t *testing.T) {
		_, err := uc.ClonePublicTrip(context.Background(), "tokyo-drift-2026", "")
		if !errors.Is(err, domain.ErrUnauthorized) {
			t.Errorf("expected ErrUnauthorized, got %v", err)
		}
	})
}