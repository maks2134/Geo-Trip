package handler

import (
	"encoding/json"
	"net/http"
	"strings"

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
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "forbidden") || strings.Contains(err.Error(), "unauthenticated") {
			status = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "too large") || strings.Contains(err.Error(), "unsupported") {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
	}

	withUser := func(r *http.Request) *http.Request {
		ctx := r.Context()
		if uid := r.Header.Get("X-User-Id"); uid != "" {
			ctx = ContextWithUserID(ctx, uid)
		}
		return r.WithContext(ctx)
	}

	// ConnectRPC procedure endpoints
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

	mux.HandleFunc("/atlasflow.trip.v1.TripService/ClonePublicTrip", func(w http.ResponseWriter, r *http.Request) {
		var req tripv1.ClonePublicTripRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp, err := h.ClonePublicTrip(withUser(r).Context(), &req)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	// REST / BFF Endpoints for gateway reverse proxy & frontend direct calls
	mux.HandleFunc("GET /api/v1/share/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		resp, err := h.GetPublicShareTrip(r.Context(), &tripv1.GetPublicShareTripRequest{Slug: slug})
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("POST /api/v1/share/{slug}/clone", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		uid := r.Header.Get("X-User-Id")
		resp, err := h.ClonePublicTrip(withUser(r).Context(), &tripv1.ClonePublicTripRequest{
			Slug:          slug,
			TargetOwnerID: uid,
		})
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("POST /api/v1/trips/clone-public/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		uid := r.Header.Get("X-User-Id")
		resp, err := h.ClonePublicTrip(withUser(r).Context(), &tripv1.ClonePublicTripRequest{
			Slug:          slug,
			TargetOwnerID: uid,
		})
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	// Attachments REST endpoints (MH-C)
	mux.HandleFunc("POST /api/v1/trips/{trip_id}/attachments/presign", func(w http.ResponseWriter, r *http.Request) {
		tripID := r.PathValue("trip_id")
		var req tripv1.CreateAttachmentPresignedURLRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, err)
			return
		}
		req.TripID = tripID
		resp, err := h.CreateAttachmentPresignedURL(withUser(r).Context(), &req)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("POST /api/v1/trips/{trip_id}/attachments/{id}/confirm", func(w http.ResponseWriter, r *http.Request) {
		tripID := r.PathValue("trip_id")
		attachmentID := r.PathValue("id")
		resp, err := h.ConfirmAttachmentUpload(withUser(r).Context(), &tripv1.ConfirmAttachmentUploadRequest{
			TripID:       tripID,
			AttachmentID: attachmentID,
		})
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("GET /api/v1/trips/{trip_id}/attachments", func(w http.ResponseWriter, r *http.Request) {
		tripID := r.PathValue("trip_id")
		placeID := r.URL.Query().Get("place_id")
		dayID := r.URL.Query().Get("day_id")
		var pID *string
		if placeID != "" {
			pID = &placeID
		}
		var dID *string
		if dayID != "" {
			dID = &dayID
		}
		resp, err := h.ListAttachments(withUser(r).Context(), &tripv1.ListAttachmentsRequest{
			TripID:  tripID,
			PlaceID: pID,
			DayID:   dID,
		})
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("DELETE /api/v1/trips/{trip_id}/attachments/{id}", func(w http.ResponseWriter, r *http.Request) {
		tripID := r.PathValue("trip_id")
		attachmentID := r.PathValue("id")
		resp, err := h.DeleteAttachment(withUser(r).Context(), &tripv1.DeleteAttachmentRequest{
			TripID:       tripID,
			AttachmentID: attachmentID,
		})
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	// Checklist starter templates (CAND-06)
	mux.HandleFunc("POST /api/v1/trips/{trip_id}/checklist/templates/{key}", func(w http.ResponseWriter, r *http.Request) {
		tripID := r.PathValue("trip_id")
		templateKey := r.PathValue("key")
		uid := r.Header.Get("X-User-Id")
		items, summary, err := h.checklistUsecase.ApplyTemplate(withUser(r).Context(), tripID, uid, templateKey)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"items":   items,
			"summary": summary,
		})
	})

	// Route optimization (FEAT-04 / CAND-08)
	mux.HandleFunc("/atlasflow.trip.v1.TripService/OptimizeDayRoute", func(w http.ResponseWriter, r *http.Request) {
		var req tripv1.OptimizeDayRouteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, err)
			return
		}
		resp, err := h.OptimizeDayRoute(withUser(r).Context(), &req)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("POST /api/v1/trips/{trip_id}/days/{day_id}/optimize", func(w http.ResponseWriter, r *http.Request) {
		tripID := r.PathValue("trip_id")
		dayID := r.PathValue("day_id")
		var body struct {
			LockStart bool `json:"lock_start"`
			LockEnd   bool `json:"lock_end"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		resp, err := h.OptimizeDayRoute(withUser(r).Context(), &tripv1.OptimizeDayRouteRequest{
			TripID:    tripID,
			DayID:     dayID,
			LockStart: body.LockStart,
			LockEnd:   body.LockEnd,
		})
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	return mux
}
