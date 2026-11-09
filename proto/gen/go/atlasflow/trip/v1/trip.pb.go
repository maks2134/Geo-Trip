package tripv1

import (
	"time"
)

type TripStatus int32
const (
	TripStatusUnspecified TripStatus = 0
	TripStatusDraft       TripStatus = 1
	TripStatusUpcoming    TripStatus = 2
	TripStatusActive      TripStatus = 3
	TripStatusCompleted   TripStatus = 4
)

type PlaceCategory int32
const (
	PlaceCategoryUnspecified PlaceCategory = 0
	PlaceCategoryFood        PlaceCategory = 1
	PlaceCategoryStay        PlaceCategory = 2
	PlaceCategorySight       PlaceCategory = 3
	PlaceCategoryNature      PlaceCategory = 4
	PlaceCategoryTransport   PlaceCategory = 5
)

type ExpenseCategory int32
const (
	ExpenseCategoryUnspecified ExpenseCategory = 0
	ExpenseCategoryStay        ExpenseCategory = 1
	ExpenseCategoryFood        ExpenseCategory = 2
	ExpenseCategoryTransport   ExpenseCategory = 3
	ExpenseCategoryTickets     ExpenseCategory = 4
	ExpenseCategoryOther       ExpenseCategory = 5
)

type ChecklistSection int32
const (
	ChecklistSectionUnspecified ChecklistSection = 0
	ChecklistSectionPreFlight   ChecklistSection = 1
	ChecklistSectionPacking     ChecklistSection = 2
	ChecklistSectionBuyOnSite   ChecklistSection = 3
)

type AttachmentVisibility int32
const (
	AttachmentVisibilityUnspecified AttachmentVisibility = 0
	AttachmentVisibilityPrivate     AttachmentVisibility = 1
	AttachmentVisibilityPublic      AttachmentVisibility = 2
)

type PlaceRef struct {
	PlaceID   string        `json:"place_id"`
	Name      string        `json:"name"`
	Latitude  float64       `json:"latitude"`
	Longitude float64       `json:"longitude"`
	Category  PlaceCategory `json:"category"`
	Address   string        `json:"address"`
}

type ItineraryPlaceItem struct {
	ID              string    `json:"id"`
	TripID          string    `json:"trip_id"`
	DayID           *string   `json:"day_id,omitempty"`
	Place           PlaceRef  `json:"place"`
	Position        string    `json:"position"`
	DurationMinutes int32     `json:"duration_minutes"`
	Cost            *float64  `json:"cost,omitempty"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type TripDay struct {
	ID        string               `json:"id"`
	TripID    string               `json:"trip_id"`
	DayIndex  int32                `json:"day_index"`
	Date      time.Time            `json:"date"`
	Title     string               `json:"title"`
	Notes     string               `json:"notes"`
	Places    []ItineraryPlaceItem `json:"places"`
}

type Expense struct {
	ID          string          `json:"id"`
	TripID      string          `json:"trip_id"`
	DayID       *string         `json:"day_id,omitempty"`
	PlaceItemID *string         `json:"place_item_id,omitempty"`
	Category    ExpenseCategory `json:"category"`
	Amount      float64         `json:"amount"`
	Currency    string          `json:"currency"`
	Title       string          `json:"title"`
	Notes       string          `json:"notes"`
	CreatedAt   time.Time       `json:"created_at"`
}

type BudgetSummary struct {
	BudgetLimit      float64            `json:"budget_limit"`
	TotalSpent       float64            `json:"total_spent"`
	Currency         string             `json:"currency"`
	SpentByCategory  map[string]float64 `json:"spent_by_category"`
	SpentByDay       map[string]float64 `json:"spent_by_day"`
}

type ChecklistItem struct {
	ID        string           `json:"id"`
	TripID    string           `json:"trip_id"`
	Section   ChecklistSection `json:"section"`
	Title     string           `json:"title"`
	IsDone    bool             `json:"is_done"`
	Position  string           `json:"position"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type ChecklistSummary struct {
	TotalItems          int32   `json:"total_items"`
	CompletedItems      int32   `json:"completed_items"`
	ReadinessPercentage float64 `json:"readiness_percentage"`
}

type AttachmentMeta struct {
	ID          string               `json:"id"`
	TripID      string               `json:"trip_id"`
	PlaceID     *string              `json:"place_id,omitempty"`
	DayID       *string              `json:"day_id,omitempty"`
	Filename    string               `json:"filename"`
	ContentType string               `json:"content_type"`
	FileSize    int64                `json:"file_size"`
	StorageKey  string               `json:"storage_key"`
	DownloadURL string               `json:"download_url"`
	Visibility  AttachmentVisibility `json:"visibility"`
	CreatedAt   time.Time            `json:"created_at"`
}

type Trip struct {
	ID               string               `json:"id"`
	OwnerID          string               `json:"owner_id"`
	Title            string               `json:"title"`
	Description      string               `json:"description"`
	StartDate        time.Time            `json:"start_date"`
	EndDate          time.Time            `json:"end_date"`
	Status           TripStatus           `json:"status"`
	CoverURL         string               `json:"cover_url"`
	Currency         string               `json:"currency"`
	BudgetLimit      float64              `json:"budget_limit"`
	Days             []TripDay            `json:"days"`
	Backlog          []ItineraryPlaceItem `json:"backlog"`
	BudgetSummary    BudgetSummary        `json:"budget_summary"`
	ChecklistSummary ChecklistSummary     `json:"checklist_summary"`
	Attachments      []AttachmentMeta     `json:"attachments"`
	ShareSlug        string               `json:"share_slug"`
	IsPublic         bool                 `json:"is_public"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
}

type PublicPlaceItem struct {
	ID              string   `json:"id"`
	Place           PlaceRef `json:"place"`
	Position        string   `json:"position"`
	DurationMinutes int32    `json:"duration_minutes"`
	Notes           string   `json:"notes"`
}

type PublicTripDay struct {
	ID        string            `json:"id"`
	DayIndex  int32             `json:"day_index"`
	Date      time.Time         `json:"date"`
	Title     string            `json:"title"`
	Notes     string            `json:"notes"`
	Places    []PublicPlaceItem `json:"places"`
}

type PublicAttachmentMeta struct {
	ID          string  `json:"id"`
	Filename    string  `json:"filename"`
	ContentType string  `json:"content_type"`
	DownloadURL string  `json:"download_url"`
	PlaceID     *string `json:"place_id,omitempty"`
	DayID       *string `json:"day_id,omitempty"`
}

type PublicTrip struct {
	ID               string                 `json:"id"`
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

type CreateTripRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CoverURL    string    `json:"cover_url"`
	Currency    string    `json:"currency"`
	BudgetLimit float64   `json:"budget_limit"`
}

type CreateTripResponse struct {
	Trip *Trip `json:"trip"`
}

type GetTripRequest struct {
	TripID string `json:"trip_id"`
}

type GetTripResponse struct {
	Trip *Trip `json:"trip"`
}

type ListTripsRequest struct {
	Status *TripStatus `json:"status,omitempty"`
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
}

type ListTripsResponse struct {
	Trips      []Trip `json:"trips"`
	TotalCount int32  `json:"total_count"`
}

type UpdateTripRequest struct {
	TripID      string      `json:"trip_id"`
	Title       *string     `json:"title,omitempty"`
	Description *string     `json:"description,omitempty"`
	StartDate   *time.Time  `json:"start_date,omitempty"`
	EndDate     *time.Time  `json:"end_date,omitempty"`
	Status      *TripStatus `json:"status,omitempty"`
	CoverURL    *string     `json:"cover_url,omitempty"`
	Currency    *string     `json:"currency,omitempty"`
	BudgetLimit *float64    `json:"budget_limit,omitempty"`
	IsPublic    *bool       `json:"is_public,omitempty"`
}

type UpdateTripResponse struct {
	Trip *Trip `json:"trip"`
}

type DeleteTripRequest struct {
	TripID string `json:"trip_id"`
}

type DeleteTripResponse struct {
	Success bool `json:"success"`
}

type AddDayRequest struct {
	TripID string    `json:"trip_id"`
	Date   time.Time `json:"date"`
	Title  string    `json:"title"`
	Notes  string    `json:"notes"`
}

type AddDayResponse struct {
	Day *TripDay `json:"day"`
}

type UpdateDayRequest struct {
	TripID string     `json:"trip_id"`
	DayID  string     `json:"day_id"`
	Date   *time.Time `json:"date,omitempty"`
	Title  *string    `json:"title,omitempty"`
	Notes  *string    `json:"notes,omitempty"`
}

type UpdateDayResponse struct {
	Day *TripDay `json:"day"`
}

type DeleteDayRequest struct {
	TripID string `json:"trip_id"`
	DayID  string `json:"day_id"`
}

type DeleteDayResponse struct {
	Success bool `json:"success"`
}

type AddPlaceToTripRequest struct {
	TripID          string   `json:"trip_id"`
	Place           PlaceRef `json:"place"`
	DayID           *string  `json:"day_id,omitempty"`
	Position        string   `json:"position"`
	DurationMinutes int32    `json:"duration_minutes"`
	Cost            *float64 `json:"cost,omitempty"`
	Notes           string   `json:"notes"`
}

type AddPlaceToTripResponse struct {
	Item *ItineraryPlaceItem `json:"item"`
}

type ReorderPlaceRequest struct {
	TripID       string  `json:"trip_id"`
	PlaceItemID  string  `json:"place_item_id"`
	TargetDayID  *string `json:"target_day_id,omitempty"`
	PrevPosition string  `json:"prev_position"`
	NextPosition string  `json:"next_position"`
}

type ReorderPlaceResponse struct {
	Item *ItineraryPlaceItem `json:"item"`
}

type RemovePlaceFromTripRequest struct {
	TripID      string `json:"trip_id"`
	PlaceItemID string `json:"place_item_id"`
}

type RemovePlaceFromTripResponse struct {
	Success bool `json:"success"`
}

type AddExpenseRequest struct {
	TripID      string          `json:"trip_id"`
	DayID       *string         `json:"day_id,omitempty"`
	PlaceItemID *string         `json:"place_item_id,omitempty"`
	Category    ExpenseCategory `json:"category"`
	Amount      float64         `json:"amount"`
	Currency    string          `json:"currency"`
	Title       string          `json:"title"`
	Notes       string          `json:"notes"`
}

type AddExpenseResponse struct {
	Expense       *Expense       `json:"expense"`
	BudgetSummary *BudgetSummary `json:"budget_summary"`
}

type ListExpensesRequest struct {
	TripID   string           `json:"trip_id"`
	DayID    *string          `json:"day_id,omitempty"`
	Category *ExpenseCategory `json:"category,omitempty"`
}

type ListExpensesResponse struct {
	Expenses      []Expense      `json:"expenses"`
	BudgetSummary *BudgetSummary `json:"budget_summary"`
}

type DeleteExpenseRequest struct {
	TripID    string `json:"trip_id"`
	ExpenseID string `json:"expense_id"`
}

type DeleteExpenseResponse struct {
	Success       bool           `json:"success"`
	BudgetSummary *BudgetSummary `json:"budget_summary"`
}

type AddChecklistItemRequest struct {
	TripID   string           `json:"trip_id"`
	Section  ChecklistSection `json:"section"`
	Title    string           `json:"title"`
	Position string           `json:"position"`
}

type AddChecklistItemResponse struct {
	Item    *ChecklistItem    `json:"item"`
	Summary *ChecklistSummary `json:"summary"`
}

type ToggleChecklistItemRequest struct {
	TripID string `json:"trip_id"`
	ItemID string `json:"item_id"`
	IsDone bool   `json:"is_done"`
}

type ToggleChecklistItemResponse struct {
	Item    *ChecklistItem    `json:"item"`
	Summary *ChecklistSummary `json:"summary"`
}

type ListChecklistItemsRequest struct {
	TripID  string            `json:"trip_id"`
	Section *ChecklistSection `json:"section,omitempty"`
}

type ListChecklistItemsResponse struct {
	Items   []ChecklistItem   `json:"items"`
	Summary *ChecklistSummary `json:"summary"`
}

type CreateAttachmentPresignedURLRequest struct {
	TripID      string               `json:"trip_id"`
	PlaceID     *string              `json:"place_id,omitempty"`
	DayID       *string              `json:"day_id,omitempty"`
	Filename    string               `json:"filename"`
	ContentType string               `json:"content_type"`
	Visibility  AttachmentVisibility `json:"visibility"`
}

type CreateAttachmentPresignedURLResponse struct {
	UploadURL    string `json:"upload_url"`
	DownloadURL  string `json:"download_url"`
	FileKey      string `json:"file_key"`
	AttachmentID string `json:"attachment_id"`
}

type GetPublicShareTripRequest struct {
	Slug string `json:"slug"`
}

type GetPublicShareTripResponse struct {
	Trip *PublicTrip `json:"trip"`
}
