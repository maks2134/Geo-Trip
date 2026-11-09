package domain

import (
	"math"
	"strings"
	"time"
)

const TotalRecognizedCountries = 195

type TravelerStats struct {
	CountriesVisited     int32
	TotalCountries       int32
	ExploredPercentage   float64
	CitiesVisited        int32
	PointsTotal          int32
	DistanceKM           float64
	VisitedCountryCodes  []string
}

func NewTravelerStats(visitedCodes []string, citiesCount int32, pointsTotal int32, distanceKM float64) TravelerStats {
	countriesVisited := int32(len(visitedCodes))
	var exploredPct float64
	if TotalRecognizedCountries > 0 {
		exploredPct = (float64(countriesVisited) / float64(TotalRecognizedCountries)) * 100.0
	}
	return TravelerStats{
		CountriesVisited:    countriesVisited,
		TotalCountries:      TotalRecognizedCountries,
		ExploredPercentage:  exploredPct,
		CitiesVisited:       citiesCount,
		PointsTotal:         pointsTotal,
		DistanceKM:          distanceKM,
		VisitedCountryCodes: visitedCodes,
	}
}

type CountryFeature struct {
	ISO2    string
	GeoJSON string
}

type ReverseGeocodeMetadata struct {
	Name        string
	Address     string
	City        string
	CountryISO2 string
	CountryName string
	Location    LatLng
}

// v0.3.0 UN Geoscheme regional progress & milestones (FEAT-06)
type RegionProgress struct {
	RegionKey            string
	RegionName           string
	VisitedCountries     int32
	TotalCountries       int32
	CompletionPercentage float64
}

type MilestoneBadge struct {
	BadgeKey    string
	Title       string
	Description string
	Icon        string
	IsUnlocked  bool
	UnlockedAt  string
}

type UNRegionDef struct {
	Key       string
	Name      string
	Continent string
	Countries []string
}

var UNRegions = []UNRegionDef{
	{Key: "europe_western", Name: "Western Europe", Continent: "Europe", Countries: []string{"FR", "DE", "NL", "BE", "AT", "CH", "LU"}},
	{Key: "europe_northern", Name: "Northern Europe", Continent: "Europe", Countries: []string{"GB", "IE", "SE", "NO", "DK", "FI", "IS", "EE", "LV", "LT"}},
	{Key: "europe_southern", Name: "Southern Europe", Continent: "Europe", Countries: []string{"IT", "ES", "PT", "GR", "HR", "SI", "RS", "ME", "BA", "MK", "AL", "CY", "MT"}},
	{Key: "europe_eastern", Name: "Eastern Europe", Continent: "Europe", Countries: []string{"PL", "CZ", "SK", "HU", "RO", "BG", "UA", "BY", "MD", "RU"}},
	{Key: "asia_eastern", Name: "Eastern Asia", Continent: "Asia", Countries: []string{"JP", "KR", "CN", "TW", "MN", "HK"}},
	{Key: "asia_southeastern", Name: "South-Eastern Asia", Continent: "Asia", Countries: []string{"TH", "VN", "ID", "MY", "SG", "PH", "KH", "LA", "MM"}},
	{Key: "asia_southern", Name: "Southern Asia", Continent: "Asia", Countries: []string{"IN", "NP", "LK", "MV"}},
	{Key: "asia_western", Name: "Western Asia", Continent: "Asia", Countries: []string{"TR", "GE", "AM", "AZ", "AE", "IL", "JO", "SA", "QA"}},
	{Key: "americas_northern", Name: "Northern America", Continent: "Americas", Countries: []string{"US", "CA"}},
	{Key: "americas_latin", Name: "Latin America", Continent: "Americas", Countries: []string{"MX", "CR", "PA", "BR", "AR", "CL", "CO", "PE", "UY"}},
	{Key: "africa", Name: "Africa", Continent: "Africa", Countries: []string{"EG", "MA", "TN", "ZA", "KE", "TZ", "MU", "SC"}},
	{Key: "oceania", Name: "Oceania", Continent: "Oceania", Countries: []string{"AU", "NZ", "FJ"}},
}

func EvaluateRegionalProgress(visitedCodes []string) []RegionProgress {
	visitedSet := make(map[string]bool, len(visitedCodes))
	for _, c := range visitedCodes {
		visitedSet[strings.ToUpper(strings.TrimSpace(c))] = true
	}

	var progress []RegionProgress
	for _, reg := range UNRegions {
		var count int32
		for _, iso := range reg.Countries {
			if visitedSet[iso] {
				count++
			}
		}
		total := int32(len(reg.Countries))
		pct := 0.0
		if total > 0 {
			pct = math.Round((float64(count)/float64(total))*1000) / 10
		}
		progress = append(progress, RegionProgress{
			RegionKey:            reg.Key,
			RegionName:           reg.Name,
			VisitedCountries:     count,
			TotalCountries:       total,
			CompletionPercentage: pct,
		})
	}
	return progress
}

func EvaluateMilestoneBadges(visitedCodes []string) []MilestoneBadge {
	visitedSet := make(map[string]bool, len(visitedCodes))
	continentsSet := make(map[string]bool)
	nowStr := time.Now().UTC().Format(time.RFC3339)

	for _, c := range visitedCodes {
		code := strings.ToUpper(strings.TrimSpace(c))
		visitedSet[code] = true
		for _, reg := range UNRegions {
			for _, iso := range reg.Countries {
				if iso == code {
					continentsSet[reg.Continent] = true
				}
			}
		}
	}

	count := len(visitedSet)
	hasAll := func(codes ...string) bool {
		for _, c := range codes {
			if !visitedSet[c] {
				return false
			}
		}
		return true
	}
	countIn := func(codes ...string) int {
		var cnt int
		for _, c := range codes {
			if visitedSet[c] {
				cnt++
			}
		}
		return cnt
	}

	badges := []MilestoneBadge{
		{
			BadgeKey:    "first_step",
			Title:       "First Step",
			Description: "Unlocked your first country on the world scratch map",
			Icon:        "footprints",
			IsUnlocked:  count >= 1,
			UnlockedAt:  nowStr,
		},
		{
			BadgeKey:    "wanderer_5",
			Title:       "Wanderer",
			Description: "Visited 5 countries across the globe",
			Icon:        "compass",
			IsUnlocked:  count >= 5,
			UnlockedAt:  nowStr,
		},
		{
			BadgeKey:    "explorer_10",
			Title:       "Explorer",
			Description: "Reached 10 countries milestone",
			Icon:        "map",
			IsUnlocked:  count >= 10,
			UnlockedAt:  nowStr,
		},
		{
			BadgeKey:    "globetrotter_25",
			Title:       "Globetrotter",
			Description: "Explored 25 countries worldwide",
			Icon:        "globe",
			IsUnlocked:  count >= 25,
			UnlockedAt:  nowStr,
		},
		{
			BadgeKey:    "centurion_100",
			Title:       "Centurion",
			Description: "Visited 100+ countries — legendary nomad tier",
			Icon:        "crown",
			IsUnlocked:  count >= 100,
			UnlockedAt:  nowStr,
		},
		{
			BadgeKey:    "nordic_cross",
			Title:       "Nordic Cross",
			Description: "Visited Sweden, Norway, Denmark, Finland, and Iceland",
			Icon:        "mountain-snow",
			IsUnlocked:  hasAll("SE", "NO", "DK", "FI", "IS"),
			UnlockedAt:  nowStr,
		},
		{
			BadgeKey:    "balkan_explorer",
			Title:       "Balkan Explorer",
			Description: "Explored 5+ countries in the Balkan region",
			Icon:        "compass",
			IsUnlocked:  countIn("HR", "SI", "RS", "ME", "BA", "MK", "AL", "GR", "BG") >= 5,
			UnlockedAt:  nowStr,
		},
		{
			BadgeKey:    "southeast_asia",
			Title:       "Southeast Asian Trail",
			Description: "Explored 4+ countries in Southeast Asia",
			Icon:        "palmtree",
			IsUnlocked:  countIn("TH", "VN", "ID", "MY", "SG", "PH", "KH", "LA") >= 4,
			UnlockedAt:  nowStr,
		},
		{
			BadgeKey:    "continents_3",
			Title:       "Three Continents",
			Description: "Stepped foot on destinations across 3 continents",
			Icon:        "earth",
			IsUnlocked:  len(continentsSet) >= 3,
			UnlockedAt:  nowStr,
		},
		{
			BadgeKey:    "continents_5",
			Title:       "Five Continents Club",
			Description: "Traveled across 5 continents",
			Icon:        "sparkles",
			IsUnlocked:  len(continentsSet) >= 5,
			UnlockedAt:  nowStr,
		},
	}

	return badges
}
