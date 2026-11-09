package domain

import "time"

type ShareRecord struct {
	ID              string
	TripID          string
	Slug            string
	HideBudget      bool
	HidePrivateDocs bool
	RevokedAt       *time.Time
	CreatedAt       time.Time
}

func (s ShareRecord) IsActive() bool {
	return s.RevokedAt == nil
}

type PublicTripView struct {
	ID               string
	Title            string
	Description      string
	StartDate        time.Time
	EndDate          time.Time
	CoverURL         string
	Days             []PublicDayView
	Backlog          []PublicPlaceView
	Attachments      []PublicAttachmentView
	OwnerDisplayName string
}

type PublicDayView struct {
	ID       string
	DayIndex int32
	Date     time.Time
	Title    string
	Notes    string
	Places   []PublicPlaceView
}

type PublicPlaceView struct {
	ID              string
	Place           PlaceRef
	Position        string
	DurationMinutes int32
	Notes           string
}

type PublicAttachmentView struct {
	ID          string
	Filename    string
	ContentType string
	DownloadURL string
	PlaceID     *string
	DayID       *string
}

// SanitizeForPublic strips budget details and private attachments strictly.
func SanitizeForPublic(
	trip Trip,
	days []TripDay,
	items []ItineraryPlaceItem,
	attachments []Attachment,
	ownerName string,
) PublicTripView {
	dayMap := make(map[string][]PublicPlaceView)
	var backlog []PublicPlaceView

	for _, it := range items {
		pv := PublicPlaceView{
			ID:              it.ID,
			Place:           it.Place,
			Position:        it.Position,
			DurationMinutes: it.DurationMinutes,
			Notes:           it.Notes, // public notes
		}
		if it.DayID == nil || *it.DayID == "" {
			backlog = append(backlog, pv)
		} else {
			dayMap[*it.DayID] = append(dayMap[*it.DayID], pv)
		}
	}

	var pubDays []PublicDayView
	for _, d := range days {
		pubDays = append(pubDays, PublicDayView{
			ID:       d.ID,
			DayIndex: d.DayIndex,
			Date:     d.Date,
			Title:    d.Title,
			Notes:    d.Notes,
			Places:   dayMap[d.ID],
		})
	}

	var pubAttachments []PublicAttachmentView
	for _, a := range attachments {
		if a.Visibility == AttachmentVisibilityPublic && a.IsReady {
			pubAttachments = append(pubAttachments, PublicAttachmentView{
				ID:          a.ID,
				Filename:    a.Filename,
				ContentType: a.ContentType,
				DownloadURL: a.DownloadURL,
				PlaceID:     a.PlaceID,
				DayID:       a.DayID,
			})
		}
	}

	return PublicTripView{
		ID:               trip.ID,
		Title:            trip.Title,
		Description:      trip.Description,
		StartDate:        trip.StartDate,
		EndDate:          trip.EndDate,
		CoverURL:         trip.CoverURL,
		Days:             pubDays,
		Backlog:          backlog,
		Attachments:      pubAttachments,
		OwnerDisplayName: ownerName,
	}
}