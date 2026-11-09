package domain_test

import (
	"testing"

	"github.com/makskozlov/atlasflow/services/geo-places/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoundingBox_Validate(t *testing.T) {
	tests := []struct {
		name    string
		bbox    domain.BoundingBox
		wantErr error
	}{
		{
			name: "valid bbox",
			bbox: domain.BoundingBox{
				MinLongitude: 2.0,
				MinLatitude:  48.0,
				MaxLongitude: 3.0,
				MaxLatitude:  49.0,
			},
		},
		{
			name: "invalid latitude bounds",
			bbox: domain.BoundingBox{
				MinLongitude: 0.0,
				MinLatitude:  -95.0,
				MaxLongitude: 10.0,
				MaxLatitude:  10.0,
			},
			wantErr: domain.ErrInvalidCoordinates,
		},
		{
			name: "min lat greater than max lat",
			bbox: domain.BoundingBox{
				MinLongitude: 0.0,
				MinLatitude:  50.0,
				MaxLongitude: 10.0,
				MaxLatitude:  40.0,
			},
			wantErr: domain.ErrInvalidBoundingBox,
		},
		{
			name: "min lng greater than max lng",
			bbox: domain.BoundingBox{
				MinLongitude: 50.0,
				MinLatitude:  10.0,
				MaxLongitude: 20.0,
				MaxLatitude:  20.0,
			},
			wantErr: domain.ErrInvalidBoundingBox,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bbox.Validate()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBoundingBox_Contains(t *testing.T) {
	bbox := domain.BoundingBox{
		MinLongitude: 2.0,
		MinLatitude:  48.0,
		MaxLongitude: 3.0,
		MaxLatitude:  49.0,
	}

	assert.True(t, bbox.Contains(48.5, 2.5))
	assert.True(t, bbox.Contains(48.0, 2.0)) // on boundary
	assert.False(t, bbox.Contains(49.5, 2.5)) // north outside
	assert.False(t, bbox.Contains(48.5, 1.5)) // west outside
}