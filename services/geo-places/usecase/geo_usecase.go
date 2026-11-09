package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/makskozlov/atlasflow/services/geo-places/domain"
)

const (
	DefaultLimit = 50
	MaxLimit     = 500
)

type SearchPlacesInput struct {
	OwnerID  string
	Query    string
	BBox     *domain.BoundingBox
	Status   *domain.PlaceStatus
	Category *domain.PlaceCategory
	Limit    int
}

type GeoUsecase struct {
	repo     PlaceRepository
	geocoder Geocoder
	clock    Clock
}

func NewGeoUsecase(repo PlaceRepository, geocoder Geocoder, clock Clock) *GeoUsecase {
	if clock == nil {
		clock = RealClock{}
	}
	return &GeoUsecase{
		repo:     repo,
		geocoder: geocoder,
		clock:    clock,
	}
}

func (u *GeoUsecase) SearchPlaces(ctx context.Context, in SearchPlacesInput) ([]domain.Place, error) {
	if in.BBox != nil {
		if err := in.BBox.Validate(); err != nil {
			return nil, fmt.Errorf("validate bbox: %w", err)
		}
	}

	limit := in.Limit
	if limit <= 0 {
		limit = DefaultLimit
	} else if limit > MaxLimit {
		limit = MaxLimit
	}

	filter := SearchFilter{
		OwnerID:  in.OwnerID,
		Query:    strings.TrimSpace(in.Query),
		BBox:     in.BBox,
		Status:   in.Status,
		Category: in.Category,
		Limit:    limit,
	}

	places, err := u.repo.SearchPlaces(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("search places: %w", err)
	}
	return places, nil
}

func (u *GeoUsecase) ReverseGeocode(ctx context.Context, lat, lng float64) (*domain.ReverseGeocodeMetadata, error) {
	loc := domain.LatLng{Latitude: lat, Longitude: lng}
	if err := loc.Validate(); err != nil {
		return nil, err
	}

	meta, err := u.geocoder.ReverseGeocode(ctx, loc)
	if err != nil {
		return nil, fmt.Errorf("reverse geocode: %w", err)
	}
	return meta, nil
}

func (u *GeoUsecase) GetVisitedCountries(ctx context.Context, userID string) ([]string, *domain.TravelerStats, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, nil, domain.ErrEmptyOwnerID
	}

	codes, err := u.repo.GetVisitedCountries(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("get visited countries: %w", err)
	}

	stats, err := u.repo.GetTravelerStats(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("get traveler stats: %w", err)
	}

	return codes, stats, nil
}

func (u *GeoUsecase) MarkCountryVisited(ctx context.Context, userID, countryCode string, visited bool) (*domain.TravelerStats, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrEmptyOwnerID
	}

	cleanCode := strings.ToUpper(strings.TrimSpace(countryCode))
	if len(cleanCode) != 2 {
		return nil, domain.ErrInvalidCountryCode
	}

	if err := u.repo.MarkCountryVisited(ctx, userID, cleanCode, visited); err != nil {
		return nil, fmt.Errorf("mark country visited: %w", err)
	}

	stats, err := u.repo.GetTravelerStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get updated stats: %w", err)
	}

	return stats, nil
}

func (u *GeoUsecase) GetVisitedCountryPolygons(ctx context.Context, userID string) ([]domain.CountryFeature, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrEmptyOwnerID
	}

	codes, err := u.repo.GetVisitedCountries(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get visited countries for polygons: %w", err)
	}

	if len(codes) == 0 {
		return []domain.CountryFeature{}, nil
	}

	features, err := u.repo.GetVisitedCountryPolygons(ctx, codes)
	if err != nil {
		return nil, fmt.Errorf("get visited polygons: %w", err)
	}

	return features, nil
}

func (u *GeoUsecase) GetRegionalMilestones(ctx context.Context, userID string) ([]domain.RegionProgress, []domain.MilestoneBadge, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, nil, domain.ErrEmptyOwnerID
	}

	codes, err := u.repo.GetVisitedCountries(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("get visited countries for milestones: %w", err)
	}

	regions := domain.EvaluateRegionalProgress(codes)
	badges := domain.EvaluateMilestoneBadges(codes)

	return regions, badges, nil
}