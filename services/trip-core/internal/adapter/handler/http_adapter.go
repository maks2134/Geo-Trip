package handler

import (
	"encoding/json"
	"net/http"

	tripv1 "github.com/makskozlov/atlasflow/proto/gen/go/atlasflow/trip/v1"
)

func NewHTTPHandler(h *TripHandler) http.Handler {
	mux := http.NewServeMux()

	writeJSON := func(w http.ResponseWriter, status int, v any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(v)
	}

	writeErr := func(w http.ResponseWriter, err error) {
		status := http.StatusInternalServerError
		writeJSON(w, status, map[string]string{"error": err.Error()})
	}

	withUser := func(r *http.Request) *http.Request {
		ctx := r.Context()
		if uid := r.Header.Get("X-User-Id"); uid != "" {
			ctx = ContextWithUserID(ctx, uid)
		}
		return r.WithContext(ctx)
	}

	mux.HandleFunc("/atlasflow.trip.v1.TripService/CreateTrip", func(w http.ResponseWriter, r *http.Request) {
		var req tripv1.CreateTripRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, err)
			return
		}
		resp, err := h.CreateTrip(withUser(r).Context(), &req)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("/atlasflow.trip.v1.TripService/GetTrip", func(w http.ResponseWriter, r *http.Request) {
		var req tripv1.GetTripRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err := h.GetTrip(withUser(r).Context(), &req)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("/atlasflow.trip.v1.TripService/ListTrips", func(w http.ResponseWriter, r *http.Request) {
		var req tripv1.ListTripsRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err := h.ListTrips(withUser(r).Context(), &req)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("/atlasflow.trip.v1.TripService/GetPublicShareTrip", func(w http.ResponseWriter, r *http.Request) {
		var req tripv1.GetPublicShareTripRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err := h.GetPublicShareTrip(r.Context(), &req)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	return mux
}
