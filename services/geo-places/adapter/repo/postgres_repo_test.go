package repo_test

import (
	"testing"

	"github.com/makskozlov/atlasflow/services/geo-places/adapter/repo"
	"github.com/makskozlov/atlasflow/services/geo-places/domain"
	"github.com/makskozlov/atlasflow/services/geo-places/usecase"
	"github.com/stretchr/testify/assert"
)

func TestPostgresPlaceRepository_BuildSearchSQL(t *testing.T) {
	r := repo.NewPostgresPlaceRepository(nil)
	status := domain.StatusVisited
	cat := domain.CategoryFood

	bbox := &domain.BoundingBox{
		MinLongitude: 2.0,
		MinLatitude:  48.0,
		MaxLongitude: 3.0,
		MaxLatitude:  49.0,
	}

	filter := usecase.SearchFilter{
		OwnerID:  "u123",
		Query:    "bistro",
		BBox:     bbox,
		Status:   &status,
		Category: &cat,
		Limit:    25,
	}

	query, args := r.BuildSearchSQL(filter)

	assert.Contains(t, query, "owner_id = $1")
	assert.Contains(t, query, "status = $2")
	assert.Contains(t, query, "category = $3")
	assert.Contains(t, query, "name ILIKE $4")
	assert.Contains(t, query, "ST_Intersects")
	assert.Contains(t, query, "LIMIT $9")

	assert.Equal(t, "u123", args[0])
	assert.Equal(t, "visited", args[1])
	assert.Equal(t, "food", args[2])
	assert.Equal(t, "%bistro%", args[3])
	assert.Equal(t, 2.0, args[4])
	assert.Equal(t, 48.0, args[5])
	assert.Equal(t, 3.0, args[6])
	assert.Equal(t, 49.0, args[7])
	assert.Equal(t, 25, args[8])
}