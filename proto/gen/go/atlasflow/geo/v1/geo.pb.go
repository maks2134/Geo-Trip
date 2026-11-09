package geov1

import (
	"time"
)

type PlaceStatus int32
const (
	PlaceStatusUnspecified PlaceStatus = 0
	PlaceStatusVisited     PlaceStatus = 1
	PlaceStatusPlanned     PlaceStatus = 2
	PlaceStatusWishlist    PlaceStatus = 3
)

type PlaceCategory int32
const (
	PlaceCategoryUnspecified PlaceCategory = 0
	PlaceCategoryFood        PlaceCategory = 1
	PlaceCategoryStay        PlaceCategory = 2
	PlaceCategorySight       PlaceCategory = 3
	PlaceCategoryNature      PlaceCategory = 4
	PlaceCategoryTransport   PlaceCategory = 5
)

type BoundingBox struct {
	MinLongitude float64 `json:"min_longitude"`
	MinLatitude  float64 `json:"min_latitude"`
	MaxLongitude float64 `json:"max_longitude"`
	MaxLatitude  float64 `json:"max_latitude"`
}

type Place struct {
	ID          string        `json:"id"`
	OwnerID     string        `json:"owner_id"`
	Name        string        `json:"name"`
	Category    PlaceCategory `json:"category"`
	Status      PlaceStatus   `json:"status"`
	CountryISO2 string        `json:"country_iso2"`
	City        string        `json:"city"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	Address     string        `json:"address"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type ReverseGeocodeMetadata struct {
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	City        string  `json:"city"`
	CountryISO2 string  `json:"country_iso2"`
	CountryName string  `json:"country_name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type TravelerStats struct {
	CountriesVisited     int32    `json:"countries_visited"`
	TotalCountries       int32    `json:"total_countries"`
	ExploredPercentage   float64  `json:"explored_percentage"`
	CitiesVisited        int32    `json:"cities_visited"`
	PointsTotal          int32    `json:"points_total"`
	DistanceKM           float64  `json:"distance_km"`
	VisitedCountryCodes  []string `json:"visited_country_codes"`
}

type SearchPlacesRequest struct {
	Query    string         `json:"query"`
	BBox     *BoundingBox   `json:"bbox,omitempty"`
	Limit    int32          `json:"limit"`
	Status   *PlaceStatus   `json:"status,omitempty"`
	Category *PlaceCategory `json:"category,omitempty"`
}

type SearchPlacesResponse struct {
	Places []Place `json:"places"`
}

type ReverseGeocodeRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Language  string  `json:"language"`
}

type ReverseGeocodeResponse struct {
	Name             string                 `json:"name"`
	City             string                 `json:"city"`
	Country          string                 `json:"country"`
	CountryCode      string                 `json:"country_code"`
	FormattedAddress string                 `json:"formatted_address"`
	CategoryHint     string                 `json:"category_hint"`
	Metadata         ReverseGeocodeMetadata `json:"metadata"`
}

type RegionProgress struct {
	RegionKey            string  `json:"region_key"`
	RegionName           string  `json:"region_name"`
	VisitedCountries     int32   `json:"visited_countries"`
	TotalCountries       int32   `json:"total_countries"`
	CompletionPercentage float64 `json:"completion_percentage"`
}

type MilestoneBadge struct {
	BadgeKey    string `json:"badge_key"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	IsUnlocked  bool   `json:"is_unlocked"`
	UnlockedAt  string `json:"unlocked_at"`
}

type GetRegionalMilestonesRequest struct {
	UserID string `json:"user_id"`
}

type GetRegionalMilestonesResponse struct {
	Regions []RegionProgress `json:"regions"`
	Badges  []MilestoneBadge `json:"badges"`
}

type GetVisitedCountriesRequest struct {
	UserID string `json:"user_id"`
}

type GetVisitedCountriesResponse struct {
	CountryCodes []string      `json:"country_codes"`
	Stats        TravelerStats `json:"stats"`
}

type MarkCountryVisitedRequest struct {
	UserID      string `json:"user_id"`
	CountryCode string `json:"country_code"`
	Visited     bool   `json:"visited"`
}

type MarkCountryVisitedResponse struct {
	Stats TravelerStats `json:"stats"`
}
