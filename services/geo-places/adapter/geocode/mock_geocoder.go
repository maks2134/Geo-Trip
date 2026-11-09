package geocode

import (
	"context"
	"strings"

	"github.com/makskozlov/atlasflow/services/geo-places/domain"
	"github.com/makskozlov/atlasflow/services/geo-places/usecase"
)

var _ usecase.Geocoder = (*MockGeocoder)(nil)

type MockGeocoder struct {
	DefaultCity    string
	DefaultCountry string
	Err            error
}

func NewMockGeocoder() *MockGeocoder {
	return &MockGeocoder{
		DefaultCity:    "Mock City",
		DefaultCountry: "MC",
	}
}

func (m *MockGeocoder) ReverseGeocode(ctx context.Context, loc domain.LatLng) (*domain.ReverseGeocodeMetadata, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if err := loc.Validate(); err != nil {
		return nil, err
	}
	return &domain.ReverseGeocodeMetadata{
		Name:        "Point of Interest",
		Address:     "123 Mockingbird Lane",
		City:        m.DefaultCity,
		CountryISO2: strings.ToUpper(m.DefaultCountry),
		CountryName: "Mockland",
		Location:    loc,
	}, nil
}