package domain

import "errors"

var (
	ErrInvalidCoordinates = errors.New("invalid coordinates: latitude must be in [-90, 90], longitude in [-180, 180]")
	ErrInvalidBoundingBox = errors.New("invalid bounding box: min coordinates must be <= max coordinates")
	ErrEmptyOwnerID       = errors.New("owner id cannot be empty")
	ErrEmptyPlaceName     = errors.New("place name cannot be empty")
	ErrInvalidCountryCode = errors.New("invalid country code: must be a 2-letter ISO code")
	ErrInvalidStatus      = errors.New("invalid place status")
	ErrInvalidCategory    = errors.New("invalid place category")
	ErrNotFound           = errors.New("place not found")
	ErrGeocodingFailed    = errors.New("geocoding service unavailable or failed")
)