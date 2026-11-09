package domain

import (
	"fmt"
	"strings"
	"time"
)

type PlaceStatus string

const (
	StatusUnspecified PlaceStatus = "unspecified"
	StatusVisited     PlaceStatus = "visited"
	StatusPlanned     PlaceStatus = "planned"
	StatusWishlist    PlaceStatus = "wishlist"
)

func ParsePlaceStatus(s string) (PlaceStatus, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "visited":
		return StatusVisited, nil
	case "planned":
		return StatusPlanned, nil
	case "wishlist":
		return StatusWishlist, nil
	case "", "unspecified":
		return StatusUnspecified, nil
	default:
		return StatusUnspecified, fmt.Errorf("%w: unknown status %q", ErrInvalidStatus, s)
	}
}

type PlaceCategory string

const (
	CategoryUnspecified PlaceCategory = "unspecified"
	CategoryFood        PlaceCategory = "food"
	CategoryStay        PlaceCategory = "stay"
	CategorySight       PlaceCategory = "sight"
	CategoryNature      PlaceCategory = "nature"
	CategoryTransport   PlaceCategory = "transport"
)

func ParsePlaceCategory(s string) (PlaceCategory, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "food":
		return CategoryFood, nil
	case "stay":
		return CategoryStay, nil
	case "sight":
		return CategorySight, nil
	case "nature":
		return CategoryNature, nil
	case "transport":
		return CategoryTransport, nil
	case "", "unspecified":
		return CategoryUnspecified, nil
	default:
		return CategoryUnspecified, fmt.Errorf("%w: unknown category %q", ErrInvalidCategory, s)
	}
}

type LatLng struct {
	Latitude  float64
	Longitude float64
}

func (l LatLng) Validate() error {
	if l.Latitude < -90 || l.Latitude > 90 || l.Longitude < -180 || l.Longitude > 180 {
		return ErrInvalidCoordinates
	}
	return nil
}

type Place struct {
	ID          string
	OwnerID     string
	Name        string
	Category    PlaceCategory
	Status      PlaceStatus
	CountryISO2 string
	City        string
	Address     string
	Location    LatLng
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p Place) Validate() error {
	if strings.TrimSpace(p.OwnerID) == "" {
		return ErrEmptyOwnerID
	}
	if strings.TrimSpace(p.Name) == "" {
		return ErrEmptyPlaceName
	}
	if err := p.Location.Validate(); err != nil {
		return err
	}
	if p.CountryISO2 != "" && len(strings.TrimSpace(p.CountryISO2)) != 2 {
		return ErrInvalidCountryCode
	}
	return nil
}