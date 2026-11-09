package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/makskozlov/atlasflow/services/geo-places/adapter/repo"
	"github.com/makskozlov/atlasflow/services/geo-places/domain"
	"github.com/makskozlov/atlasflow/services/geo-places/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryPlaceRepository_BBoxAndStats(t *testing.T) {
	r := repo.NewMemoryPlaceRepository()
	ctx := context.Background()
	now := time.Now()

	p1 := domain.Place{
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
	}
	p2 := domain.Place{
		ID:          "p2",
		OwnerID:     "user1",
		Name:        "Colosseum",
		Category:    domain.CategorySight,
		Status:      domain.StatusVisited,
		CountryISO2: "IT",
		City:        "Rome",
		Location:    domain.LatLng{Latitude: 41.8902, Longitude: 12.4922},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	p3 := domain.Place{
		ID:          "p3",
		OwnerID:     "user1",
		Name:        "Tokyo Tower",
		Category:    domain.CategorySight,
		Status:      domain.StatusWishlist, // wishlist doesn't count for visited countries
		CountryISO2: "JP",
		City:        "Tokyo",
		Location:    domain.LatLng{Latitude: 35.6586, Longitude: 139.7454},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	require.NoError(t, r.UpsertPlace(ctx, p1))
	require.NoError(t, r.UpsertPlace(ctx, p2))
	require.NoError(t, r.UpsertPlace(ctx, p3))

	t.Run("search in Paris bbox", func(t *testing.T) {
		bbox := &domain.BoundingBox{
			MinLongitude: 2.0,
			MinLatitude:  48.0,
			MaxLongitude: 3.0,
			MaxLatitude:  49.0,
		}
		res, err := r.SearchPlaces(ctx, usecase.SearchFilter{OwnerID: "user1", BBox: bbox, Limit: 10})
		require.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "p1", res[0].ID)
	})

	t.Run("stats calculation", func(t *testing.T) {
		stats, err := r.GetTravelerStats(ctx, "user1")
		require.NoError(t, err)
		assert.Equal(t, int32(2), stats.CountriesVisited) // FR, IT (JP is wishlist)
		assert.Equal(t, int32(2), stats.CitiesVisited)    // Paris, Rome
		assert.Equal(t, int32(3), stats.PointsTotal)      // 3 points total
		assert.InDelta(t, (2.0/195.0)*100.0, stats.ExploredPercentage, 0.001)
	})

	t.Run("mark country visited explicitly", func(t *testing.T) {
		require.NoError(t, r.MarkCountryVisited(ctx, "user1", "ES", true))

		stats, err := r.GetTravelerStats(ctx, "user1")
		require.NoError(t, err)
		assert.Equal(t, int32(3), stats.CountriesVisited) // FR, IT, ES
		assert.Contains(t, stats.VisitedCountryCodes, "ES")

		// unmark
		require.NoError(t, r.MarkCountryVisited(ctx, "user1", "ES", false))
		statsAfter, err := r.GetTravelerStats(ctx, "user1")
		require.NoError(t, err)
		assert.Equal(t, int32(2), statsAfter.CountriesVisited)
	})
}