package domain

import (
	"math"
)

type TSPOptimizationResult struct {
	OrderedPlaceIDs     []string
	OriginalDistanceKM  float64
	OptimizedDistanceKM float64
	DistanceSavedKM     float64
	SavingsPercent      float64
	OrderedPlaces       []ItineraryPlaceItem
}

// HaversineKM calculates the great-circle distance between two points on Earth in kilometers.
func HaversineKM(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKM = 6371.0

	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	rLat1 := lat1 * math.Pi / 180.0
	rLat2 := lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rLat1)*math.Cos(rLat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKM * c
}

// CalculateRouteTotalDistanceKM computes the total consecutive travel distance in kilometers.
func CalculateRouteTotalDistanceKM(places []ItineraryPlaceItem) float64 {
	if len(places) < 2 {
		return 0.0
	}
	var total float64
	for i := 0; i < len(places)-1; i++ {
		p1 := places[i].Place
		p2 := places[i+1].Place
		total += HaversineKM(p1.Latitude, p1.Longitude, p2.Latitude, p2.Longitude)
	}
	return math.Round(total*100) / 100
}

// OptimizeDayRoute2Opt performs a 2-opt heuristic optimization on the sequential stops of a day.
// If lockStart is true, the first stop (e.g. hotel) remains fixed at index 0.
// If lockEnd is true, the last stop (e.g. airport/dinner) remains fixed at index len-1.
func OptimizeDayRoute2Opt(places []ItineraryPlaceItem, lockStart bool, lockEnd bool) TSPOptimizationResult {
	n := len(places)
	if n <= 2 {
		total := CalculateRouteTotalDistanceKM(places)
		ids := make([]string, n)
		for i, p := range places {
			ids[i] = p.ID
		}
		return TSPOptimizationResult{
			OrderedPlaceIDs:     ids,
			OriginalDistanceKM:  total,
			OptimizedDistanceKM: total,
			DistanceSavedKM:     0,
			SavingsPercent:      0,
			OrderedPlaces:       places,
		}
	}

	origDistance := CalculateRouteTotalDistanceKM(places)

	// Work on a copy of places
	current := make([]ItineraryPlaceItem, n)
	copy(current, places)

	startIdx := 0
	if lockStart {
		startIdx = 1
	}

	endIdx := n - 1
	if lockEnd {
		endIdx = n - 2
	}

	improved := true
	maxIterations := 100
	iteration := 0

	for improved && iteration < maxIterations {
		improved = false
		iteration++

		for i := startIdx; i < endIdx; i++ {
			for k := i + 1; k <= endIdx; k++ {
				// Reverse segment between i and k
				reversed := make([]ItineraryPlaceItem, n)
				copy(reversed, current)

				// Reverse subarray reversed[i:k+1]
				for l, r := i, k; l < r; l, r = l+1, r-1 {
					reversed[l], reversed[r] = reversed[r], reversed[l]
				}

				newDist := CalculateRouteTotalDistanceKM(reversed)
				if newDist < CalculateRouteTotalDistanceKM(current)-0.01 {
					current = reversed
					improved = true
				}
			}
		}
	}

	optDistance := CalculateRouteTotalDistanceKM(current)
	saved := math.Max(0, math.Round((origDistance-optDistance)*100)/100)
	var pct float64
	if origDistance > 0 {
		pct = math.Round((saved/origDistance)*1000) / 10
	}

	ids := make([]string, n)
	for i, p := range current {
		ids[i] = p.ID
	}

	return TSPOptimizationResult{
		OrderedPlaceIDs:     ids,
		OriginalDistanceKM:  origDistance,
		OptimizedDistanceKM: optDistance,
		DistanceSavedKM:     saved,
		SavingsPercent:      pct,
		OrderedPlaces:       current,
	}
}
