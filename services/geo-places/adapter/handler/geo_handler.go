package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	geov1 "github.com/makskozlov/atlasflow/proto/gen/go/atlasflow/geo/v1"
	"github.com/makskozlov/atlasflow/proto/gen/go/atlasflow/geo/v1/geov1connect"
	"github.com/makskozlov/atlasflow/services/geo-places/domain"
	"github.com/makskozlov/atlasflow/services/geo-places/usecase"
)

type GeoHandler struct {
	uc *usecase.GeoUsecase
}

var _ geov1connect.GeoServiceHandler = (*GeoHandler)(nil)

func NewGeoHandler(uc *usecase.GeoUsecase) *GeoHandler {
	return &GeoHandler{uc: uc}
}

type contextKey string

const UserIDContextKey contextKey = "user_id"

func extractUserID(ctx context.Context, fallback string) string {
	if val, ok := ctx.Value(UserIDContextKey).(string); ok && val != "" {
		return val
	}
	return fallback
}

func (h *GeoHandler) SearchPlaces(ctx context.Context, req *geov1.SearchPlacesRequest) (*geov1.SearchPlacesResponse, error) {
	if req == nil {
		req = &geov1.SearchPlacesRequest{}
	}

	var bbox *domain.BoundingBox
	if req.BBox != nil {
		bbox = &domain.BoundingBox{
			MinLongitude: req.BBox.MinLongitude,
			MinLatitude:  req.BBox.MinLatitude,
			MaxLongitude: req.BBox.MaxLongitude,
			MaxLatitude:  req.BBox.MaxLatitude,
		}
	}

	var status *domain.PlaceStatus
	if req.Status != nil {
		s := protoToDomainStatus(*req.Status)
		status = &s
	}

	var category *domain.PlaceCategory
	if req.Category != nil {
		c := protoToDomainCategory(*req.Category)
		category = &c
	}

	in := usecase.SearchPlacesInput{
		OwnerID:  extractUserID(ctx, ""),
		Query:    req.Query,
		BBox:     bbox,
		Status:   status,
		Category: category,
		Limit:    int(req.Limit),
	}

	places, err := h.uc.SearchPlaces(ctx, in)
	if err != nil {
		return nil, mapDomainError(err)
	}

	protoPlaces := make([]geov1.Place, len(places))
	for i, p := range places {
		protoPlaces[i] = domainToProtoPlace(p)
	}

	return &geov1.SearchPlacesResponse{Places: protoPlaces}, nil
}

func (h *GeoHandler) ReverseGeocode(ctx context.Context, req *geov1.ReverseGeocodeRequest) (*geov1.ReverseGeocodeResponse, error) {
	if req == nil {
		return nil, domain.ErrInvalidCoordinates
	}

	meta, err := h.uc.ReverseGeocode(ctx, req.Latitude, req.Longitude)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &geov1.ReverseGeocodeResponse{
		Metadata: geov1.ReverseGeocodeMetadata{
			Name:        meta.Name,
			Address:     meta.Address,
			City:        meta.City,
			CountryISO2: meta.CountryISO2,
			CountryName: meta.CountryName,
			Latitude:    meta.Location.Latitude,
			Longitude:   meta.Location.Longitude,
		},
	}, nil
}

func (h *GeoHandler) GetVisitedCountries(ctx context.Context, req *geov1.GetVisitedCountriesRequest) (*geov1.GetVisitedCountriesResponse, error) {
	userID := ""
	if req != nil {
		userID = req.UserID
	}
	userID = extractUserID(ctx, userID)

	codes, stats, err := h.uc.GetVisitedCountries(ctx, userID)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &geov1.GetVisitedCountriesResponse{
		CountryCodes: codes,
		Stats:        domainToProtoStats(stats),
	}, nil
}

func (h *GeoHandler) MarkCountryVisited(ctx context.Context, req *geov1.MarkCountryVisitedRequest) (*geov1.MarkCountryVisitedResponse, error) {
	if req == nil {
		return nil, domain.ErrEmptyOwnerID
	}

	userID := extractUserID(ctx, req.UserID)
	stats, err := h.uc.MarkCountryVisited(ctx, userID, req.CountryCode, req.Visited)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &geov1.MarkCountryVisitedResponse{
		Stats: domainToProtoStats(stats),
	}, nil
}

func protoToDomainStatus(s geov1.PlaceStatus) domain.PlaceStatus {
	switch s {
	case geov1.PlaceStatusVisited:
		return domain.StatusVisited
	case geov1.PlaceStatusPlanned:
		return domain.StatusPlanned
	case geov1.PlaceStatusWishlist:
		return domain.StatusWishlist
	default:
		return domain.StatusUnspecified
	}
}

func domainToProtoStatus(s domain.PlaceStatus) geov1.PlaceStatus {
	switch s {
	case domain.StatusVisited:
		return geov1.PlaceStatusVisited
	case domain.StatusPlanned:
		return geov1.PlaceStatusPlanned
	case domain.StatusWishlist:
		return geov1.PlaceStatusWishlist
	default:
		return geov1.PlaceStatusUnspecified
	}
}

func protoToDomainCategory(c geov1.PlaceCategory) domain.PlaceCategory {
	switch c {
	case geov1.PlaceCategoryFood:
		return domain.CategoryFood
	case geov1.PlaceCategoryStay:
		return domain.CategoryStay
	case geov1.PlaceCategorySight:
		return domain.CategorySight
	case geov1.PlaceCategoryNature:
		return domain.CategoryNature
	case geov1.PlaceCategoryTransport:
		return domain.CategoryTransport
	default:
		return domain.CategoryUnspecified
	}
}

func domainToProtoCategory(c domain.PlaceCategory) geov1.PlaceCategory {
	switch c {
	case domain.CategoryFood:
		return geov1.PlaceCategoryFood
	case domain.CategoryStay:
		return geov1.PlaceCategoryStay
	case domain.CategorySight:
		return geov1.PlaceCategorySight
	case domain.CategoryNature:
		return geov1.PlaceCategoryNature
	case domain.CategoryTransport:
		return geov1.PlaceCategoryTransport
	default:
		return geov1.PlaceCategoryUnspecified
	}
}

func domainToProtoPlace(p domain.Place) geov1.Place {
	return geov1.Place{
		ID:          p.ID,
		OwnerID:     p.OwnerID,
		Name:        p.Name,
		Category:    domainToProtoCategory(p.Category),
		Status:      domainToProtoStatus(p.Status),
		CountryISO2: p.CountryISO2,
		City:        p.City,
		Latitude:    p.Location.Latitude,
		Longitude:   p.Location.Longitude,
		Address:     p.Address,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func domainToProtoStats(s *domain.TravelerStats) geov1.TravelerStats {
	if s == nil {
		return geov1.TravelerStats{TotalCountries: domain.TotalRecognizedCountries}
	}
	return geov1.TravelerStats{
		CountriesVisited:    s.CountriesVisited,
		TotalCountries:      s.TotalCountries,
		ExploredPercentage:  s.ExploredPercentage,
		CitiesVisited:       s.CitiesVisited,
		PointsTotal:         s.PointsTotal,
		DistanceKM:          s.DistanceKM,
		VisitedCountryCodes: s.VisitedCountryCodes,
	}
}

func mapDomainError(err error) error {
	return err
}

func NewHTTPHandler(handler geov1connect.GeoServiceHandler) http.Handler {
	mux := http.NewServeMux()

	writeJSON := func(w http.ResponseWriter, status int, v any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(v)
	}

	writeError := func(w http.ResponseWriter, err error) {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrInvalidCoordinates) ||
			errors.Is(err, domain.ErrInvalidBoundingBox) ||
			errors.Is(err, domain.ErrEmptyOwnerID) ||
			errors.Is(err, domain.ErrEmptyPlaceName) ||
			errors.Is(err, domain.ErrInvalidCountryCode) ||
			errors.Is(err, domain.ErrInvalidStatus) ||
			errors.Is(err, domain.ErrInvalidCategory) {
			status = http.StatusBadRequest
		} else if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, domain.ErrGeocodingFailed) {
			status = http.StatusBadGateway
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
	}

	withAuth := func(r *http.Request) context.Context {
		ctx := r.Context()
		if uid := r.Header.Get("X-User-Id"); uid != "" {
			ctx = context.WithValue(ctx, UserIDContextKey, uid)
		}
		return ctx
	}

	mux.HandleFunc("/atlasflow.geo.v1.GeoService/SearchPlaces", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req geov1.SearchPlacesRequest
		if r.Body != nil && r.ContentLength > 0 {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeError(w, domain.ErrInvalidCoordinates)
				return
			}
		}
		res, err := handler.SearchPlaces(withAuth(r), &req)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})

	mux.HandleFunc("/atlasflow.geo.v1.GeoService/ReverseGeocode", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req geov1.ReverseGeocodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, domain.ErrInvalidCoordinates)
			return
		}
		res, err := handler.ReverseGeocode(withAuth(r), &req)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})

	mux.HandleFunc("/atlasflow.geo.v1.GeoService/GetVisitedCountries", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req geov1.GetVisitedCountriesRequest
		if r.Body != nil && r.ContentLength > 0 {
			_ = json.NewDecoder(r.Body).Decode(&req)
		}
		res, err := handler.GetVisitedCountries(withAuth(r), &req)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})

	mux.HandleFunc("/atlasflow.geo.v1.GeoService/MarkCountryVisited", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req geov1.MarkCountryVisitedRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, domain.ErrInvalidCountryCode)
			return
		}
		res, err := handler.MarkCountryVisited(withAuth(r), &req)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return mux
}