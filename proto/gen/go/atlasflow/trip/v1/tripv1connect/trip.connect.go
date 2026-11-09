package tripv1connect

import (
	"context"

	tripv1 "github.com/makskozlov/atlasflow/proto/gen/go/atlasflow/trip/v1"
)

type TripServiceHandler interface {
	CreateTrip(context.Context, *tripv1.CreateTripRequest) (*tripv1.CreateTripResponse, error)
	GetTrip(context.Context, *tripv1.GetTripRequest) (*tripv1.GetTripResponse, error)
	ListTrips(context.Context, *tripv1.ListTripsRequest) (*tripv1.ListTripsResponse, error)
	UpdateTrip(context.Context, *tripv1.UpdateTripRequest) (*tripv1.UpdateTripResponse, error)
	DeleteTrip(context.Context, *tripv1.DeleteTripRequest) (*tripv1.DeleteTripResponse, error)

	AddDay(context.Context, *tripv1.AddDayRequest) (*tripv1.AddDayResponse, error)
	UpdateDay(context.Context, *tripv1.UpdateDayRequest) (*tripv1.UpdateDayResponse, error)
	DeleteDay(context.Context, *tripv1.DeleteDayRequest) (*tripv1.DeleteDayResponse, error)

	AddPlaceToTrip(context.Context, *tripv1.AddPlaceToTripRequest) (*tripv1.AddPlaceToTripResponse, error)
	ReorderPlace(context.Context, *tripv1.ReorderPlaceRequest) (*tripv1.ReorderPlaceResponse, error)
	RemovePlaceFromTrip(context.Context, *tripv1.RemovePlaceFromTripRequest) (*tripv1.RemovePlaceFromTripResponse, error)

	AddExpense(context.Context, *tripv1.AddExpenseRequest) (*tripv1.AddExpenseResponse, error)
	ListExpenses(context.Context, *tripv1.ListExpensesRequest) (*tripv1.ListExpensesResponse, error)
	DeleteExpense(context.Context, *tripv1.DeleteExpenseRequest) (*tripv1.DeleteExpenseResponse, error)

	AddChecklistItem(context.Context, *tripv1.AddChecklistItemRequest) (*tripv1.AddChecklistItemResponse, error)
	ToggleChecklistItem(context.Context, *tripv1.ToggleChecklistItemRequest) (*tripv1.ToggleChecklistItemResponse, error)
	ListChecklistItems(context.Context, *tripv1.ListChecklistItemsRequest) (*tripv1.ListChecklistItemsResponse, error)

	CreateAttachmentPresignedURL(context.Context, *tripv1.CreateAttachmentPresignedURLRequest) (*tripv1.CreateAttachmentPresignedURLResponse, error)
	ConfirmAttachmentUpload(context.Context, *tripv1.ConfirmAttachmentUploadRequest) (*tripv1.ConfirmAttachmentUploadResponse, error)
	ListAttachments(context.Context, *tripv1.ListAttachmentsRequest) (*tripv1.ListAttachmentsResponse, error)
	DeleteAttachment(context.Context, *tripv1.DeleteAttachmentRequest) (*tripv1.DeleteAttachmentResponse, error)

	GetPublicShareTrip(context.Context, *tripv1.GetPublicShareTripRequest) (*tripv1.GetPublicShareTripResponse, error)
	ClonePublicTrip(context.Context, *tripv1.ClonePublicTripRequest) (*tripv1.ClonePublicTripResponse, error)
}
