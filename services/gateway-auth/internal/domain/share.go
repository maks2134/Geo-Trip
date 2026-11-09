package domain

import "time"

type PublicPlaceItem struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Address         string  `json:"address"`
	Position        string  `json:"position"`
	DurationMinutes int32   `json:"duration_minutes"`
	Notes           string  `json:"notes"`
}

type PublicTripDay struct {
	ID       string            `json:"id"`
	DayIndex int32             `json:"day_index"`
	Date     time.Time         `json:"date"`
	Title    string            `json:"title"`
	Notes    string            `json:"notes"`
	Places   []PublicPlaceItem `json:"places"`
}

type PublicAttachmentMeta struct {
	ID          string  `json:"id"`
	Filename    string  `json:"filename"`
	ContentType string  `json:"content_type"`
	DownloadURL string  `json:"download_url"`
	PlaceID     *string `json:"place_id,omitempty"`
	DayID       *string `json:"day_id,omitempty"`
}

// PublicTrip contains strictly sanitized, public trip view data.
// It explicitly omits budget, expenses, and private checklists/documents.
type PublicTrip struct {
	ID               string                 `json:"id"`
	Slug             string                 `json:"slug"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	StartDate        time.Time              `json:"start_date"`
	EndDate          time.Time              `json:"end_date"`
	CoverURL         string                 `json:"cover_url"`
	Days             []PublicTripDay        `json:"days"`
	Backlog          []PublicPlaceItem      `json:"backlog"`
	Attachments      []PublicAttachmentMeta `json:"attachments"`
	OwnerDisplayName string                 `json:"owner_display_name"`
}