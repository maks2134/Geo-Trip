package geov1connect

import (
	"context"

	geov1 "github.com/makskozlov/atlasflow/proto/gen/go/atlasflow/geo/v1"
)

type GeoServiceHandler interface {
	SearchPlaces(context.Context, *geov1.SearchPlacesRequest) (*geov1.SearchPlacesResponse, error)
	ReverseGeocode(context.Context, *geov1.ReverseGeocodeRequest) (*geov1.ReverseGeocodeResponse, error)
	GetVisitedCountries(context.Context, *geov1.GetVisitedCountriesRequest) (*geov1.GetVisitedCountriesResponse, error)
	MarkCountryVisited(context.Context, *geov1.MarkCountryVisitedRequest) (*geov1.MarkCountryVisitedResponse, error)
	GetRegionalMilestones(context.Context, *geov1.GetRegionalMilestonesRequest) (*geov1.GetRegionalMilestonesResponse, error)
}
