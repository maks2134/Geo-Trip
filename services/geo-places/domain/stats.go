package domain

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