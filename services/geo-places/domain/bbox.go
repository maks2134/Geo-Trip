package domain

type BoundingBox struct {
	MinLongitude float64
	MinLatitude  float64
	MaxLongitude float64
	MaxLatitude  float64
}

func (b BoundingBox) Validate() error {
	if b.MinLatitude < -90 || b.MinLatitude > 90 ||
		b.MaxLatitude < -90 || b.MaxLatitude > 90 ||
		b.MinLongitude < -180 || b.MinLongitude > 180 ||
		b.MaxLongitude < -180 || b.MaxLongitude > 180 {
		return ErrInvalidCoordinates
	}
	if b.MinLatitude > b.MaxLatitude || b.MinLongitude > b.MaxLongitude {
		return ErrInvalidBoundingBox
	}
	return nil
}

func (b BoundingBox) Contains(lat, lng float64) bool {
	return lat >= b.MinLatitude && lat <= b.MaxLatitude &&
		lng >= b.MinLongitude && lng <= b.MaxLongitude
}