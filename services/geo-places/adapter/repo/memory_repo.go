package repo

import (
	"context"
	"sort"
	"strings"
	"sync"

	"github.com/makskozlov/atlasflow/services/geo-places/domain"
	"github.com/makskozlov/atlasflow/services/geo-places/usecase"
)

type MemoryPlaceRepository struct {
	mu               sync.RWMutex
	places           map[string]domain.Place
	visitedCountries map[string]map[string]bool // userID -> countryISO2 -> true
	countryPolygons  map[string]string          // countryISO2 -> GeoJSON string
}

func NewMemoryPlaceRepository() *MemoryPlaceRepository {
	return &MemoryPlaceRepository{
		places:           make(map[string]domain.Place),
		visitedCountries: make(map[string]map[string]bool),
		countryPolygons:  make(map[string]string),
	}
}

func (m *MemoryPlaceRepository) SetCountryPolygon(iso2, geojson string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.countryPolygons[strings.ToUpper(iso2)] = geojson
}

func (m *MemoryPlaceRepository) UpsertPlace(ctx context.Context, p domain.Place) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.places[p.ID] = p
	return nil
}

func (m *MemoryPlaceRepository) GetPlace(ctx context.Context, id, ownerID string) (*domain.Place, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.places[id]
	if !ok || (ownerID != "" && p.OwnerID != ownerID) {
		return nil, domain.ErrNotFound
	}
	return &p, nil
}

func (m *MemoryPlaceRepository) SearchPlaces(ctx context.Context, filter usecase.SearchFilter) ([]domain.Place, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []domain.Place
	q := strings.ToLower(filter.Query)

	for _, p := range m.places {
		if filter.OwnerID != "" && p.OwnerID != filter.OwnerID {
			continue
		}
		if filter.Status != nil && *filter.Status != domain.StatusUnspecified && p.Status != *filter.Status {
			continue
		}
		if filter.Category != nil && *filter.Category != domain.CategoryUnspecified && p.Category != *filter.Category {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(p.Name), q) && !strings.Contains(strings.ToLower(p.City), q) {
			continue
		}
		if filter.BBox != nil && !filter.BBox.Contains(p.Location.Latitude, p.Location.Longitude) {
			continue
		}

		results = append(results, p)
		if filter.Limit > 0 && len(results) >= filter.Limit {
			break
		}
	}

	return results, nil
}

func (m *MemoryPlaceRepository) MarkCountryVisited(ctx context.Context, ownerID, countryCode string, visited bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cleanCode := strings.ToUpper(countryCode)
	if m.visitedCountries[ownerID] == nil {
		m.visitedCountries[ownerID] = make(map[string]bool)
	}

	if visited {
		m.visitedCountries[ownerID][cleanCode] = true
	} else {
		delete(m.visitedCountries[ownerID], cleanCode)
	}
	return nil
}

func (m *MemoryPlaceRepository) GetVisitedCountries(ctx context.Context, ownerID string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	countrySet := make(map[string]bool)

	// Collect from explicit mark
	if marked, ok := m.visitedCountries[ownerID]; ok {
		for code := range marked {
			countrySet[code] = true
		}
	}

	// Collect from visited pins
	for _, p := range m.places {
		if p.OwnerID == ownerID && p.Status == domain.StatusVisited && len(p.CountryISO2) == 2 {
			countrySet[strings.ToUpper(p.CountryISO2)] = true
		}
	}

	codes := make([]string, 0, len(countrySet))
	for code := range countrySet {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes, nil
}

func (m *MemoryPlaceRepository) GetTravelerStats(ctx context.Context, ownerID string) (*domain.TravelerStats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	countrySet := make(map[string]bool)
	citySet := make(map[string]bool)
	pointsTotal := int32(0)

	if marked, ok := m.visitedCountries[ownerID]; ok {
		for code := range marked {
			countrySet[code] = true
		}
	}

	for _, p := range m.places {
		if p.OwnerID != ownerID {
			continue
		}
		pointsTotal++
		if p.Status == domain.StatusVisited {
			if len(p.CountryISO2) == 2 {
				countrySet[strings.ToUpper(p.CountryISO2)] = true
			}
			if strings.TrimSpace(p.City) != "" {
				citySet[strings.ToLower(strings.TrimSpace(p.City))] = true
			}
		}
	}

	codes := make([]string, 0, len(countrySet))
	for code := range countrySet {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	stats := domain.NewTravelerStats(codes, int32(len(citySet)), pointsTotal, 0.0)
	return &stats, nil
}

func (m *MemoryPlaceRepository) GetVisitedCountryPolygons(ctx context.Context, countryCodes []string) ([]domain.CountryFeature, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var features []domain.CountryFeature
	for _, code := range countryCodes {
		clean := strings.ToUpper(code)
		if geojson, ok := m.countryPolygons[clean]; ok {
			features = append(features, domain.CountryFeature{
				ISO2:    clean,
				GeoJSON: geojson,
			})
		}
	}
	return features, nil
}