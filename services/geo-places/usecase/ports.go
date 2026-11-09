package usecase

import (
	"context"
	"time"

	"github.com/makskozlov/atlasflow/services/geo-places/domain"
)

type SearchFilter struct {
	OwnerID  string
	Query    string
	BBox     *domain.BoundingBox
	Status   *domain.PlaceStatus
	Category *domain.PlaceCategory
	Limit    int
}

type PlaceRepository interface {
	UpsertPlace(ctx context.Context, place domain.Place) error
	GetPlace(ctx context.Context, id, ownerID string) (*domain.Place, error)
	SearchPlaces(ctx context.Context, filter SearchFilter) ([]domain.Place, error)
	GetTravelerStats(ctx context.Context, ownerID string) (*domain.TravelerStats, error)
	MarkCountryVisited(ctx context.Context, ownerID, countryCode string, visited bool) error
	GetVisitedCountries(ctx context.Context, ownerID string) ([]string, error)
	GetVisitedCountryPolygons(ctx context.Context, countryCodes []string) ([]domain.CountryFeature, error)
}

type Geocoder interface {
	ReverseGeocode(ctx context.Context, loc domain.LatLng) (*domain.ReverseGeocodeMetadata, error)
}

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}