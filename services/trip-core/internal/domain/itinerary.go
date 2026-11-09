package domain

import "time"

type PlaceCategory int32

const (
	PlaceCategoryUnspecified PlaceCategory = 0
	PlaceCategoryFood        PlaceCategory = 1
	PlaceCategoryStay        PlaceCategory = 2
	PlaceCategorySight       PlaceCategory = 3
	PlaceCategoryNature      PlaceCategory = 4
	PlaceCategoryTransport   PlaceCategory = 5
)

type PlaceRef struct {
	PlaceID   string
	Name      string
	Latitude  float64
	Longitude float64
	Category  PlaceCategory
	Address   string
}

type ItineraryPlaceItem struct {
	ID              string
	TripID          string
	DayID           *string // nil means in Backlog
	Place           PlaceRef
	Position        string
	DurationMinutes int32
	Cost            *float64
	Notes           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}