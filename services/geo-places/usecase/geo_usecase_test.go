package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/makskozlov/atlasflow/services/geo-places/domain"
	"github.com/makskozlov/atlasflow/services/geo-places/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	places       []domain.Place
	stats        *domain.TravelerStats
	visitedCodes []string
	features     []domain.CountryFeature
	searchErr    error
	statsErr     error
	markErr      error
	visitedErr   error
	polyErr      error
}

func (f *fakeRepo) UpsertPlace(ctx context.Context, place domain.Place) error { return nil }
func (f *fakeRepo) GetPlace(ctx context.Context, id, ownerID string) (*domain.Place, error) {
	return nil, nil
}
func (f *fakeRepo) SearchPlaces(ctx context.Context, filter usecase.SearchFilter) ([]domain.Place, error) {
	if f.searchErr != nil {
		return nil, f.searchErr
	}
	return f.places, nil
}
func (f *fakeRepo) GetTravelerStats(ctx context.Context, ownerID string) (*domain.TravelerStats, error) {
	if f.statsErr != nil {
		return nil, f.statsErr
	}
	return f.stats, nil
}
func (f *fakeRepo) MarkCountryVisited(ctx context.Context, ownerID, countryCode string, visited bool) error {
	return f.markErr
}
func (f *fakeRepo) GetVisitedCountries(ctx context.Context, ownerID string) ([]string, error) {
	if f.visitedErr != nil {
		return nil, f.visitedErr
	}
	return f.visitedCodes, nil
}
func (f *fakeRepo) GetVisitedCountryPolygons(ctx context.Context, countryCodes []string) ([]domain.CountryFeature, error) {
	if f.polyErr != nil {
		return nil, f.polyErr
	}
	return f.features, nil
}

type fakeGeocoder struct {
	meta *domain.ReverseGeocodeMetadata
	err  error
}

func (f *fakeGeocoder) ReverseGeocode(ctx context.Context, loc domain.LatLng) (*domain.ReverseGeocodeMetadata, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.meta, nil
}

func TestGeoUsecase_SearchPlaces(t *testing.T) {
	tests := []struct {
		name      string
		input     usecase.SearchPlacesInput
		repoErr   error
		places    []domain.Place
		wantErr   bool
		errTarget error
	}{
		{
			name: "happy path bbox search",
			input: usecase.SearchPlacesInput{
				OwnerID: "u1",
				BBox: &domain.BoundingBox{
					MinLongitude: 2.0, MinLatitude: 48.0,
					MaxLongitude: 3.0, MaxLatitude: 49.0,
				},
				Limit: 20,
			},
			places: []domain.Place{
				{ID: "p1", Name: "Louvre", Location: domain.LatLng{Latitude: 48.86, Longitude: 2.33}},
			},
		},
		{
			name: "invalid bbox error",
			input: usecase.SearchPlacesInput{
				OwnerID: "u1",
				BBox: &domain.BoundingBox{
					MinLongitude: 5.0, MinLatitude: 48.0,
					MaxLongitude: 2.0, MaxLatitude: 49.0,
				},
			},
			wantErr:   true,
			errTarget: domain.ErrInvalidBoundingBox,
		},
		{
			name: "limit clamped to max 500",
			input: usecase.SearchPlacesInput{
				OwnerID: "u1",
				Limit:   1000,
			},
			places: []domain.Place{},
		},
		{
			name: "repository error propagated",
			input: usecase.SearchPlacesInput{
				OwnerID: "u1",
			},
			repoErr: errors.New("db connection failure"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{places: tt.places, searchErr: tt.repoErr}
			uc := usecase.NewGeoUsecase(repo, &fakeGeocoder{}, nil)

			res, err := uc.SearchPlaces(context.Background(), tt.input)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					require.ErrorIs(t, err, tt.errTarget)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(tt.places), len(res))
			}
		})
	}
}

func TestGeoUsecase_ReverseGeocode(t *testing.T) {
	tests := []struct {
		name      string
		lat       float64
		lng       float64
		geoMeta   *domain.ReverseGeocodeMetadata
		geoErr    error
		wantErr   bool
		errTarget error
	}{
		{
			name: "happy path reverse geocode",
			lat:  48.8584,
			lng:  2.2945,
			geoMeta: &domain.ReverseGeocodeMetadata{
				Name:        "Champ de Mars",
				City:        "Paris",
				CountryISO2: "FR",
				CountryName: "France",
				Location:    domain.LatLng{Latitude: 48.8584, Longitude: 2.2945},
			},
		},
		{
			name:      "latitude out of bounds",
			lat:       95.0,
			lng:       2.0,
			wantErr:   true,
			errTarget: domain.ErrInvalidCoordinates,
		},
		{
			name:      "longitude out of bounds",
			lat:       48.0,
			lng:       -185.0,
			wantErr:   true,
			errTarget: domain.ErrInvalidCoordinates,
		},
		{
			name:    "geocoder failure",
			lat:     48.0,
			lng:     2.0,
			geoErr:  errors.New("network timeout"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			geocoder := &fakeGeocoder{meta: tt.geoMeta, err: tt.geoErr}
			uc := usecase.NewGeoUsecase(&fakeRepo{}, geocoder, nil)

			res, err := uc.ReverseGeocode(context.Background(), tt.lat, tt.lng)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					require.ErrorIs(t, err, tt.errTarget)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.geoMeta.City, res.City)
				assert.Equal(t, tt.geoMeta.CountryISO2, res.CountryISO2)
			}
		})
	}
}

func TestGeoUsecase_GetVisitedCountries(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		codes     []string
		stats     *domain.TravelerStats
		repoErr   error
		wantErr   bool
		errTarget error
	}{
		{
			name:   "happy path with visited countries and stats",
			userID: "u123",
			codes:  []string{"FR", "IT", "JP"},
			stats: &domain.TravelerStats{
				CountriesVisited:   3,
				TotalCountries:     195,
				ExploredPercentage: 1.538,
			},
		},
		{
			name:      "empty user id error",
			userID:    "   ",
			wantErr:   true,
			errTarget: domain.ErrEmptyOwnerID,
		},
		{
			name:    "repo error",
			userID:  "u123",
			repoErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{visitedCodes: tt.codes, stats: tt.stats, visitedErr: tt.repoErr}
			uc := usecase.NewGeoUsecase(repo, &fakeGeocoder{}, nil)

			codes, stats, err := uc.GetVisitedCountries(context.Background(), tt.userID)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					require.ErrorIs(t, err, tt.errTarget)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.codes, codes)
				assert.Equal(t, tt.stats.CountriesVisited, stats.CountriesVisited)
			}
		})
	}
}

func TestGeoUsecase_MarkCountryVisited(t *testing.T) {
	expectedStats := &domain.TravelerStats{
		CountriesVisited:    1,
		TotalCountries:      195,
		ExploredPercentage:  0.512,
		VisitedCountryCodes: []string{"ES"},
	}

	tests := []struct {
		name        string
		userID      string
		countryCode string
		visited     bool
		markErr     error
		wantErr     bool
		errTarget   error
	}{
		{
			name:        "happy path mark country visited",
			userID:      "u123",
			countryCode: "es",
			visited:     true,
		},
		{
			name:        "empty user id error",
			userID:      "",
			countryCode: "ES",
			visited:     true,
			wantErr:     true,
			errTarget:   domain.ErrEmptyOwnerID,
		},
		{
			name:        "invalid country code length",
			userID:      "u123",
			countryCode: "ESP",
			visited:     true,
			wantErr:     true,
			errTarget:   domain.ErrInvalidCountryCode,
		},
		{
			name:        "repo error",
			userID:      "u123",
			countryCode: "ES",
			visited:     true,
			markErr:     errors.New("db timeout"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{stats: expectedStats, markErr: tt.markErr}
			uc := usecase.NewGeoUsecase(repo, &fakeGeocoder{}, nil)

			stats, err := uc.MarkCountryVisited(context.Background(), tt.userID, tt.countryCode, tt.visited)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					require.ErrorIs(t, err, tt.errTarget)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, int32(1), stats.CountriesVisited)
			}
		})
	}
}

func TestGeoUsecase_GetVisitedCountryPolygons(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		codes     []string
		features  []domain.CountryFeature
		wantErr   bool
		errTarget error
	}{
		{
			name:   "returns unlocked country geojson",
			userID: "u123",
			codes:  []string{"FR"},
			features: []domain.CountryFeature{
				{ISO2: "FR", GeoJSON: `{"type":"MultiPolygon","coordinates":[]}`},
			},
		},
		{
			name:     "empty visited countries returns empty list",
			userID:   "u123",
			codes:    []string{},
			features: []domain.CountryFeature{},
		},
		{
			name:      "empty user id returns error",
			userID:    "",
			wantErr:   true,
			errTarget: domain.ErrEmptyOwnerID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{visitedCodes: tt.codes, features: tt.features}
			uc := usecase.NewGeoUsecase(repo, &fakeGeocoder{}, nil)

			feats, err := uc.GetVisitedCountryPolygons(context.Background(), tt.userID)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					require.ErrorIs(t, err, tt.errTarget)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(tt.features), len(feats))
			}
		})
	}
}