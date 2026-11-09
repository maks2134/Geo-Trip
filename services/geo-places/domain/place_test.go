package domain_test

import (
	"testing"
	"time"

	"github.com/makskozlov/atlasflow/services/geo-places/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLatLng_Validate(t *testing.T) {
	tests := []struct {
		name    string
		latLng  domain.LatLng
		wantErr error
	}{
		{
			name:   "valid coordinates",
			latLng: domain.LatLng{Latitude: 48.8566, Longitude: 2.3522},
		},
		{
			name:   "valid boundary coordinates",
			latLng: domain.LatLng{Latitude: -90.0, Longitude: 180.0},
		},
		{
			name:    "latitude too high",
			latLng:  domain.LatLng{Latitude: 90.0001, Longitude: 0.0},
			wantErr: domain.ErrInvalidCoordinates,
		},
		{
			name:    "latitude too low",
			latLng:  domain.LatLng{Latitude: -90.1, Longitude: 0.0},
			wantErr: domain.ErrInvalidCoordinates,
		},
		{
			name:    "longitude too high",
			latLng:  domain.LatLng{Latitude: 0.0, Longitude: 180.01},
			wantErr: domain.ErrInvalidCoordinates,
		},
		{
			name:    "longitude too low",
			latLng:  domain.LatLng{Latitude: 0.0, Longitude: -180.05},
			wantErr: domain.ErrInvalidCoordinates,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.latLng.Validate()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPlace_Validate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		place   domain.Place
		wantErr error
	}{
		{
			name: "valid place",
			place: domain.Place{
				ID:          "p1",
				OwnerID:     "user1",
				Name:        "Eiffel Tower",
				Category:    domain.CategorySight,
				Status:      domain.StatusVisited,
				CountryISO2: "FR",
				City:        "Paris",
				Location:    domain.LatLng{Latitude: 48.8584, Longitude: 2.2945},
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		{
			name: "empty owner id",
			place: domain.Place{
				Name:     "Eiffel Tower",
				Location: domain.LatLng{Latitude: 48.8584, Longitude: 2.2945},
			},
			wantErr: domain.ErrEmptyOwnerID,
		},
		{
			name: "empty name",
			place: domain.Place{
				OwnerID:  "user1",
				Name:     "  ",
				Location: domain.LatLng{Latitude: 48.8584, Longitude: 2.2945},
			},
			wantErr: domain.ErrEmptyPlaceName,
		},
		{
			name: "invalid location",
			place: domain.Place{
				OwnerID:  "user1",
				Name:     "Point",
				Location: domain.LatLng{Latitude: 100, Longitude: 0},
			},
			wantErr: domain.ErrInvalidCoordinates,
		},
		{
			name: "invalid country iso2 code length",
			place: domain.Place{
				OwnerID:     "user1",
				Name:        "Point",
				CountryISO2: "FRA",
				Location:    domain.LatLng{Latitude: 48.8584, Longitude: 2.2945},
			},
			wantErr: domain.ErrInvalidCountryCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.place.Validate()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParsePlaceStatus(t *testing.T) {
	tests := []struct {
		in      string
		want    domain.PlaceStatus
		wantErr bool
	}{
		{"visited", domain.StatusVisited, false},
		{"VISITED", domain.StatusVisited, false},
		{"planned", domain.StatusPlanned, false},
		{"wishlist", domain.StatusWishlist, false},
		{"", domain.StatusUnspecified, false},
		{"invalid", domain.StatusUnspecified, true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := domain.ParsePlaceStatus(tt.in)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}