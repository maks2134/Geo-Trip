package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	geov1 "github.com/makskozlov/atlasflow/proto/gen/go/atlasflow/geo/v1"
	"github.com/makskozlov/atlasflow/services/geo-places/adapter/geocode"
	"github.com/makskozlov/atlasflow/services/geo-places/adapter/handler"
	"github.com/makskozlov/atlasflow/services/geo-places/adapter/repo"
	"github.com/makskozlov/atlasflow/services/geo-places/domain"
	"github.com/makskozlov/atlasflow/services/geo-places/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer() (http.Handler, *handler.GeoHandler, *repo.MemoryPlaceRepository) {
	memRepo := repo.NewMemoryPlaceRepository()
	geo := geocode.NewMockGeocoder()
	geo.DefaultCity = "Paris"
	geo.DefaultCountry = "FR"

	uc := usecase.NewGeoUsecase(memRepo, geo, nil)
	h := handler.NewGeoHandler(uc)
	return handler.NewHTTPHandler(h), h, memRepo
}

func TestGeoHandler_SearchPlacesEndpoint(t *testing.T) {
	mux, _, r := setupTestServer()
	now := time.Now()

	_ = r.UpsertPlace(context.Background(), domain.Place{
		ID:          "p1",
		OwnerID:     "user1",
		Name:        "Louvre Museum",
		Category:    domain.CategorySight,
		Status:      domain.StatusVisited,
		CountryISO2: "FR",
		City:        "Paris",
		Location:    domain.LatLng{Latitude: 48.8606, Longitude: 2.3376},
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	body, _ := json.Marshal(geov1.SearchPlacesRequest{
		Query: "Louvre",
		BBox: &geov1.BoundingBox{
			MinLongitude: 2.0,
			MinLatitude:  48.0,
			MaxLongitude: 3.0,
			MaxLatitude:  49.0,
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/atlasflow.geo.v1.GeoService/SearchPlaces", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp geov1.SearchPlacesResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Len(t, resp.Places, 1)
	assert.Equal(t, "Louvre Museum", resp.Places[0].Name)
}

func TestGeoHandler_ReverseGeocodeEndpoint(t *testing.T) {
	mux, _, _ := setupTestServer()

	t.Run("valid coordinates", func(t *testing.T) {
		body, _ := json.Marshal(geov1.ReverseGeocodeRequest{Latitude: 48.8584, Longitude: 2.2945})
		req := httptest.NewRequest(http.MethodPost, "/atlasflow.geo.v1.GeoService/ReverseGeocode", bytes.NewReader(body))
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp geov1.ReverseGeocodeResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, "Paris", resp.Metadata.City)
		assert.Equal(t, "FR", resp.Metadata.CountryISO2)
	})

	t.Run("invalid coordinates bad request", func(t *testing.T) {
		body, _ := json.Marshal(geov1.ReverseGeocodeRequest{Latitude: 150.0, Longitude: 2.2945})
		req := httptest.NewRequest(http.MethodPost, "/atlasflow.geo.v1.GeoService/ReverseGeocode", bytes.NewReader(body))
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGeoHandler_MarkAndGetVisitedCountriesEndpoint(t *testing.T) {
	mux, _, _ := setupTestServer()

	// 1. Mark country
	markBody, _ := json.Marshal(geov1.MarkCountryVisitedRequest{
		UserID:      "u123",
		CountryCode: "IT",
		Visited:     true,
	})
	reqMark := httptest.NewRequest(http.MethodPost, "/atlasflow.geo.v1.GeoService/MarkCountryVisited", bytes.NewReader(markBody))
	wMark := httptest.NewRecorder()
	mux.ServeHTTP(wMark, reqMark)

	assert.Equal(t, http.StatusOK, wMark.Code)
	var markResp geov1.MarkCountryVisitedResponse
	require.NoError(t, json.NewDecoder(wMark.Body).Decode(&markResp))
	assert.Equal(t, int32(1), markResp.Stats.CountriesVisited)

	// 2. Get visited countries
	getBody, _ := json.Marshal(geov1.GetVisitedCountriesRequest{UserID: "u123"})
	reqGet := httptest.NewRequest(http.MethodPost, "/atlasflow.geo.v1.GeoService/GetVisitedCountries", bytes.NewReader(getBody))
	wGet := httptest.NewRecorder()
	mux.ServeHTTP(wGet, reqGet)

	assert.Equal(t, http.StatusOK, wGet.Code)
	var getResp geov1.GetVisitedCountriesResponse
	require.NoError(t, json.NewDecoder(wGet.Body).Decode(&getResp))
	assert.Equal(t, []string{"IT"}, getResp.CountryCodes)
	assert.Equal(t, int32(1), getResp.Stats.CountriesVisited)
}

func TestGeoHandler_DirectHandlerCalls(t *testing.T) {
	_, h, r := setupTestServer()
	ctx := context.Background()

	require.NoError(t, r.MarkCountryVisited(ctx, "user99", "ES", true))

	resp, err := h.GetVisitedCountries(ctx, &geov1.GetVisitedCountriesRequest{UserID: "user99"})
	require.NoError(t, err)
	assert.Contains(t, resp.CountryCodes, "ES")

	markResp, err := h.MarkCountryVisited(ctx, &geov1.MarkCountryVisitedRequest{
		UserID:      "user99",
		CountryCode: "PT",
		Visited:     true,
	})
	require.NoError(t, err)
	assert.Equal(t, int32(2), markResp.Stats.CountriesVisited)

	revResp, err := h.ReverseGeocode(ctx, &geov1.ReverseGeocodeRequest{Latitude: 40.4168, Longitude: -3.7038})
	require.NoError(t, err)
	assert.Equal(t, "FR", revResp.Metadata.CountryISO2)
}