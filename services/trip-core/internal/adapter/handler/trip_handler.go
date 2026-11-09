package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	tripv1 "github.com/makskozlov/atlasflow/proto/gen/go/atlasflow/trip/v1"
	"github.com/makskozlov/atlasflow/proto/gen/go/atlasflow/trip/v1/tripv1connect"
	"github.com/makskozlov/atlasflow/services/trip-core/internal/domain"
	"github.com/makskozlov/atlasflow/services/trip-core/internal/usecase"
)

type TripHandler struct {
	tripUsecase       *usecase.TripUsecase
	dayUsecase        *usecase.DayUsecase
	itineraryUsecase  *usecase.ItineraryUsecase
	budgetUsecase     *usecase.BudgetUsecase
	checklistUsecase  *usecase.ChecklistUsecase
	attachmentUsecase *usecase.AttachmentUsecase
	shareUsecase      *usecase.PublicShareUsecase
}

var _ tripv1connect.TripServiceHandler = (*TripHandler)(nil)

func NewTripHandler(
	tu *usecase.TripUsecase,
	du *usecase.DayUsecase,
	iu *usecase.ItineraryUsecase,
	bu *usecase.BudgetUsecase,
	cu *usecase.ChecklistUsecase,
	au *usecase.AttachmentUsecase,
	su *usecase.PublicShareUsecase,
) *TripHandler {
	return &TripHandler{
		tripUsecase:       tu,
		dayUsecase:        du,
		itineraryUsecase:  iu,
		budgetUsecase:     bu,
		checklistUsecase:  cu,
		attachmentUsecase: au,
		shareUsecase:      su,
	}
}

func mapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, domain.ErrNotFound) {
		return connect.NewError(connect.CodeNotFound, err)
	}
	if errors.Is(err, domain.ErrForbidden) {
		return connect.NewError(connect.CodePermissionDenied, err)
	}
	if errors.Is(err, domain.ErrValidation) || errors.Is(err, domain.ErrInvalidMimeType) || errors.Is(err, domain.ErrFileTooLarge) {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}
	if errors.Is(err, domain.ErrConflict) {
		return connect.NewError(connect.CodeAlreadyExists, err)
	}
	return connect.NewError(connect.CodeInternal, err)
}

func (h *TripHandler) CreateTrip(ctx context.Context, req *tripv1.CreateTripRequest) (*tripv1.CreateTripResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	t, err := h.tripUsecase.CreateTrip(ctx, usecase.CreateTripInput{
		OwnerID:     uid,
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		CoverURL:    req.CoverURL,
		Currency:    req.Currency,
		BudgetLimit: req.BudgetLimit,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &tripv1.CreateTripResponse{
		Trip: &tripv1.Trip{
			ID:          t.ID,
			OwnerID:     t.OwnerID,
			Title:       t.Title,
			Description: t.Description,
			StartDate:   t.StartDate,
			EndDate:     t.EndDate,
			CoverURL:    t.CoverURL,
			Currency:    t.Currency,
			BudgetLimit: t.BudgetLimit,
			ShareSlug:   t.ShareSlug,
			IsPublic:    t.IsPublic,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		},
	}, nil
}

func (h *TripHandler) GetTrip(ctx context.Context, req *tripv1.GetTripRequest) (*tripv1.GetTripResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	t, err := h.tripUsecase.GetTrip(ctx, req.TripID, uid)
	if err != nil {
		return nil, mapError(err)
	}

	return &tripv1.GetTripResponse{
		Trip: &tripv1.Trip{
			ID:          t.ID,
			OwnerID:     t.OwnerID,
			Title:       t.Title,
			Description: t.Description,
			StartDate:   t.StartDate,
			EndDate:     t.EndDate,
			CoverURL:    t.CoverURL,
			Currency:    t.Currency,
			BudgetLimit: t.BudgetLimit,
			ShareSlug:   t.ShareSlug,
			IsPublic:    t.IsPublic,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		},
	}, nil
}

func (h *TripHandler) ListTrips(ctx context.Context, req *tripv1.ListTripsRequest) (*tripv1.ListTripsResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	trips, total, err := h.tripUsecase.ListTrips(ctx, uid, nil, req.Limit, req.Offset)
	if err != nil {
		return nil, mapError(err)
	}

	var protoTrips []tripv1.Trip
	for _, t := range trips {
		protoTrips = append(protoTrips, tripv1.Trip{
			ID:          t.ID,
			OwnerID:     t.OwnerID,
			Title:       t.Title,
			Description: t.Description,
			StartDate:   t.StartDate,
			EndDate:     t.EndDate,
			CoverURL:    t.CoverURL,
			Currency:    t.Currency,
			BudgetLimit: t.BudgetLimit,
			ShareSlug:   t.ShareSlug,
			IsPublic:    t.IsPublic,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}

	return &tripv1.ListTripsResponse{
		Trips:      protoTrips,
		TotalCount: total,
	}, nil
}

func (h *TripHandler) UpdateTrip(ctx context.Context, req *tripv1.UpdateTripRequest) (*tripv1.UpdateTripResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	t, err := h.tripUsecase.UpdateTrip(ctx, usecase.UpdateTripInput{
		TripID:      req.TripID,
		UserID:      uid,
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		CoverURL:    req.CoverURL,
		Currency:    req.Currency,
		BudgetLimit: req.BudgetLimit,
		IsPublic:    req.IsPublic,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &tripv1.UpdateTripResponse{
		Trip: &tripv1.Trip{
			ID:          t.ID,
			OwnerID:     t.OwnerID,
			Title:       t.Title,
			Description: t.Description,
			StartDate:   t.StartDate,
			EndDate:     t.EndDate,
			CoverURL:    t.CoverURL,
			Currency:    t.Currency,
			BudgetLimit: t.BudgetLimit,
			ShareSlug:   t.ShareSlug,
			IsPublic:    t.IsPublic,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		},
	}, nil
}

func (h *TripHandler) DeleteTrip(ctx context.Context, req *tripv1.DeleteTripRequest) (*tripv1.DeleteTripResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	if err := h.tripUsecase.DeleteTrip(ctx, req.TripID, uid); err != nil {
		return nil, mapError(err)
	}
	return &tripv1.DeleteTripResponse{Success: true}, nil
}

func (h *TripHandler) AddDay(ctx context.Context, req *tripv1.AddDayRequest) (*tripv1.AddDayResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	d, err := h.dayUsecase.AddDay(ctx, usecase.AddDayInput{
		TripID: req.TripID,
		UserID: uid,
		Date:   req.Date,
		Title:  req.Title,
		Notes:  req.Notes,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &tripv1.AddDayResponse{
		Day: &tripv1.TripDay{
			ID:       d.ID,
			TripID:   d.TripID,
			DayIndex: d.DayIndex,
			Date:     d.Date,
			Title:    d.Title,
			Notes:    d.Notes,
		},
	}, nil
}

func (h *TripHandler) UpdateDay(ctx context.Context, req *tripv1.UpdateDayRequest) (*tripv1.UpdateDayResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	d, err := h.dayUsecase.UpdateDay(ctx, usecase.UpdateDayInput{
		TripID: req.TripID,
		DayID:  req.DayID,
		UserID: uid,
		Date:   req.Date,
		Title:  req.Title,
		Notes:  req.Notes,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &tripv1.UpdateDayResponse{
		Day: &tripv1.TripDay{
			ID:       d.ID,
			TripID:   d.TripID,
			DayIndex: d.DayIndex,
			Date:     d.Date,
			Title:    d.Title,
			Notes:    d.Notes,
		},
	}, nil
}

func (h *TripHandler) DeleteDay(ctx context.Context, req *tripv1.DeleteDayRequest) (*tripv1.DeleteDayResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	if err := h.dayUsecase.DeleteDay(ctx, req.TripID, req.DayID, uid); err != nil {
		return nil, mapError(err)
	}
	return &tripv1.DeleteDayResponse{Success: true}, nil
}

func (h *TripHandler) AddPlaceToTrip(ctx context.Context, req *tripv1.AddPlaceToTripRequest) (*tripv1.AddPlaceToTripResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	it, err := h.itineraryUsecase.AddPlace(ctx, usecase.AddPlaceInput{
		TripID: req.TripID,
		UserID: uid,
		Place: domain.PlaceRef{
			PlaceID:   req.Place.PlaceID,
			Name:      req.Place.Name,
			Latitude:  req.Place.Latitude,
			Longitude: req.Place.Longitude,
			Category:  domain.PlaceCategory(req.Place.Category),
			Address:   req.Place.Address,
		},
		DayID:           req.DayID,
		Position:        req.Position,
		DurationMinutes: req.DurationMinutes,
		Cost:            req.Cost,
		Notes:           req.Notes,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &tripv1.AddPlaceToTripResponse{
		Item: &tripv1.ItineraryPlaceItem{
			ID:     it.ID,
			TripID: it.TripID,
			DayID:  it.DayID,
			Place: tripv1.PlaceRef{
				PlaceID:   it.Place.PlaceID,
				Name:      it.Place.Name,
				Latitude:  it.Place.Latitude,
				Longitude: it.Place.Longitude,
				Category:  tripv1.PlaceCategory(it.Place.Category),
				Address:   it.Place.Address,
			},
			Position:        it.Position,
			DurationMinutes: it.DurationMinutes,
			Cost:            it.Cost,
			Notes:           it.Notes,
			CreatedAt:       it.CreatedAt,
			UpdatedAt:       it.UpdatedAt,
		},
	}, nil
}

func (h *TripHandler) ReorderPlace(ctx context.Context, req *tripv1.ReorderPlaceRequest) (*tripv1.ReorderPlaceResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	it, err := h.itineraryUsecase.ReorderPlace(ctx, usecase.ReorderPlaceInput{
		TripID:       req.TripID,
		UserID:       uid,
		PlaceItemID:  req.PlaceItemID,
		TargetDayID:  req.TargetDayID,
		PrevPosition: req.PrevPosition,
		NextPosition: req.NextPosition,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &tripv1.ReorderPlaceResponse{
		Item: &tripv1.ItineraryPlaceItem{
			ID:     it.ID,
			TripID: it.TripID,
			DayID:  it.DayID,
			Place: tripv1.PlaceRef{
				PlaceID:   it.Place.PlaceID,
				Name:      it.Place.Name,
				Latitude:  it.Place.Latitude,
				Longitude: it.Place.Longitude,
				Category:  tripv1.PlaceCategory(it.Place.Category),
				Address:   it.Place.Address,
			},
			Position:        it.Position,
			DurationMinutes: it.DurationMinutes,
			Cost:            it.Cost,
			Notes:           it.Notes,
			CreatedAt:       it.CreatedAt,
			UpdatedAt:       it.UpdatedAt,
		},
	}, nil
}

func (h *TripHandler) RemovePlaceFromTrip(ctx context.Context, req *tripv1.RemovePlaceFromTripRequest) (*tripv1.RemovePlaceFromTripResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	if err := h.itineraryUsecase.RemovePlace(ctx, req.TripID, req.PlaceItemID, uid); err != nil {
		return nil, mapError(err)
	}
	return &tripv1.RemovePlaceFromTripResponse{Success: true}, nil
}

func (h *TripHandler) AddExpense(ctx context.Context, req *tripv1.AddExpenseRequest) (*tripv1.AddExpenseResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	exp, summary, err := h.budgetUsecase.AddExpense(ctx, usecase.AddExpenseInput{
		TripID:      req.TripID,
		UserID:      uid,
		DayID:       req.DayID,
		PlaceItemID: req.PlaceItemID,
		Category:    domain.ExpenseCategory(req.Category),
		Amount:      req.Amount,
		Currency:    req.Currency,
		Title:       req.Title,
		Notes:       req.Notes,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &tripv1.AddExpenseResponse{
		Expense: &tripv1.Expense{
			ID:          exp.ID,
			TripID:      exp.TripID,
			DayID:       exp.DayID,
			PlaceItemID: exp.PlaceItemID,
			Category:    tripv1.ExpenseCategory(exp.Category),
			Amount:      exp.Amount,
			Currency:    exp.Currency,
			Title:       exp.Title,
			Notes:       exp.Notes,
			CreatedAt:   exp.CreatedAt,
		},
		BudgetSummary: &tripv1.BudgetSummary{
			BudgetLimit:     summary.BudgetLimit,
			TotalSpent:      summary.TotalSpent,
			Currency:        summary.Currency,
			SpentByCategory: summary.SpentByCategory,
			SpentByDay:      summary.SpentByDay,
		},
	}, nil
}

func (h *TripHandler) ListExpenses(ctx context.Context, req *tripv1.ListExpensesRequest) (*tripv1.ListExpensesResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	var cat *domain.ExpenseCategory
	if req.Category != nil {
		c := domain.ExpenseCategory(*req.Category)
		cat = &c
	}

	expenses, summary, err := h.budgetUsecase.ListExpenses(ctx, req.TripID, uid, req.DayID, cat)
	if err != nil {
		return nil, mapError(err)
	}

	var protoExpenses []tripv1.Expense
	for _, e := range expenses {
		protoExpenses = append(protoExpenses, tripv1.Expense{
			ID:          e.ID,
			TripID:      e.TripID,
			DayID:       e.DayID,
			PlaceItemID: e.PlaceItemID,
			Category:    tripv1.ExpenseCategory(e.Category),
			Amount:      e.Amount,
			Currency:    e.Currency,
			Title:       e.Title,
			Notes:       e.Notes,
			CreatedAt:   e.CreatedAt,
		})
	}

	return &tripv1.ListExpensesResponse{
		Expenses: protoExpenses,
		BudgetSummary: &tripv1.BudgetSummary{
			BudgetLimit:     summary.BudgetLimit,
			TotalSpent:      summary.TotalSpent,
			Currency:        summary.Currency,
			SpentByCategory: summary.SpentByCategory,
			SpentByDay:      summary.SpentByDay,
		},
	}, nil
}

func (h *TripHandler) DeleteExpense(ctx context.Context, req *tripv1.DeleteExpenseRequest) (*tripv1.DeleteExpenseResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	summary, err := h.budgetUsecase.DeleteExpense(ctx, req.TripID, req.ExpenseID, uid)
	if err != nil {
		return nil, mapError(err)
	}

	return &tripv1.DeleteExpenseResponse{
		Success: true,
		BudgetSummary: &tripv1.BudgetSummary{
			BudgetLimit:     summary.BudgetLimit,
			TotalSpent:      summary.TotalSpent,
			Currency:        summary.Currency,
			SpentByCategory: summary.SpentByCategory,
			SpentByDay:      summary.SpentByDay,
		},
	}, nil
}

func (h *TripHandler) AddChecklistItem(ctx context.Context, req *tripv1.AddChecklistItemRequest) (*tripv1.AddChecklistItemResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	it, summary, err := h.checklistUsecase.AddChecklistItem(ctx, usecase.AddChecklistItemInput{
		TripID:   req.TripID,
		UserID:   uid,
		Section:  domain.ChecklistSection(req.Section),
		Title:    req.Title,
		Position: req.Position,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &tripv1.AddChecklistItemResponse{
		Item: &tripv1.ChecklistItem{
			ID:        it.ID,
			TripID:    it.TripID,
			Section:   tripv1.ChecklistSection(it.Section),
			Title:     it.Title,
			IsDone:    it.IsDone,
			Position:  it.Position,
			CreatedAt: it.CreatedAt,
			UpdatedAt: it.UpdatedAt,
		},
		Summary: &tripv1.ChecklistSummary{
			TotalItems:          summary.TotalItems,
			CompletedItems:      summary.CompletedItems,
			ReadinessPercentage: summary.ReadinessPercentage,
		},
	}, nil
}

func (h *TripHandler) ToggleChecklistItem(ctx context.Context, req *tripv1.ToggleChecklistItemRequest) (*tripv1.ToggleChecklistItemResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	it, summary, err := h.checklistUsecase.ToggleItem(ctx, req.TripID, req.ItemID, uid, req.IsDone)
	if err != nil {
		return nil, mapError(err)
	}

	return &tripv1.ToggleChecklistItemResponse{
		Item: &tripv1.ChecklistItem{
			ID:        it.ID,
			TripID:    it.TripID,
			Section:   tripv1.ChecklistSection(it.Section),
			Title:     it.Title,
			IsDone:    it.IsDone,
			Position:  it.Position,
			CreatedAt: it.CreatedAt,
			UpdatedAt: it.UpdatedAt,
		},
		Summary: &tripv1.ChecklistSummary{
			TotalItems:          summary.TotalItems,
			CompletedItems:      summary.CompletedItems,
			ReadinessPercentage: summary.ReadinessPercentage,
		},
	}, nil
}

func (h *TripHandler) ListChecklistItems(ctx context.Context, req *tripv1.ListChecklistItemsRequest) (*tripv1.ListChecklistItemsResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	var sec *domain.ChecklistSection
	if req.Section != nil {
		s := domain.ChecklistSection(*req.Section)
		sec = &s
	}

	items, summary, err := h.checklistUsecase.ListItems(ctx, req.TripID, uid, sec)
	if err != nil {
		return nil, mapError(err)
	}

	var protoItems []tripv1.ChecklistItem
	for _, it := range items {
		protoItems = append(protoItems, tripv1.ChecklistItem{
			ID:        it.ID,
			TripID:    it.TripID,
			Section:   tripv1.ChecklistSection(it.Section),
			Title:     it.Title,
			IsDone:    it.IsDone,
			Position:  it.Position,
			CreatedAt: it.CreatedAt,
			UpdatedAt: it.UpdatedAt,
		})
	}

	return &tripv1.ListChecklistItemsResponse{
		Items: protoItems,
		Summary: &tripv1.ChecklistSummary{
			TotalItems:          summary.TotalItems,
			CompletedItems:      summary.CompletedItems,
			ReadinessPercentage: summary.ReadinessPercentage,
		},
	}, nil
}

func (h *TripHandler) CreateAttachmentPresignedURL(ctx context.Context, req *tripv1.CreateAttachmentPresignedURLRequest) (*tripv1.CreateAttachmentPresignedURLResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	out, err := h.attachmentUsecase.CreateUpload(ctx, usecase.CreateUploadInput{
		TripID:      req.TripID,
		UserID:      uid,
		PlaceID:     req.PlaceID,
		DayID:       req.DayID,
		Filename:    req.Filename,
		ContentType: req.ContentType,
		FileSize:    req.FileSize,
		Visibility:  domain.AttachmentVisibility(req.Visibility),
		OwnerType:   req.OwnerType,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &tripv1.CreateAttachmentPresignedURLResponse{
		UploadURL:    out.UploadURL,
		DownloadURL:  out.DownloadURL,
		FileKey:      out.FileKey,
		AttachmentID: out.AttachmentID,
	}, nil
}

func (h *TripHandler) ConfirmAttachmentUpload(ctx context.Context, req *tripv1.ConfirmAttachmentUploadRequest) (*tripv1.ConfirmAttachmentUploadResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	att, err := h.attachmentUsecase.ConfirmUpload(ctx, req.TripID, req.AttachmentID, uid)
	if err != nil {
		return nil, mapError(err)
	}

	return &tripv1.ConfirmAttachmentUploadResponse{
		Attachment: &tripv1.AttachmentMeta{
			ID:          att.ID,
			TripID:      att.TripID,
			PlaceID:     att.PlaceID,
			DayID:       att.DayID,
			Filename:    att.Filename,
			ContentType: att.ContentType,
			FileSize:    att.FileSize,
			StorageKey:  att.StorageKey,
			DownloadURL: att.DownloadURL,
			Visibility:  tripv1.AttachmentVisibility(att.Visibility),
			CreatedAt:   att.CreatedAt,
			IsReady:     att.IsReady,
			OwnerType:   att.OwnerType,
		},
	}, nil
}

func (h *TripHandler) ListAttachments(ctx context.Context, req *tripv1.ListAttachmentsRequest) (*tripv1.ListAttachmentsResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	var vis *domain.AttachmentVisibility
	if req.Visibility != nil {
		v := domain.AttachmentVisibility(*req.Visibility)
		vis = &v
	}

	list, err := h.attachmentUsecase.ListAttachments(ctx, req.TripID, req.PlaceID, req.DayID, vis, uid)
	if err != nil {
		return nil, mapError(err)
	}

	var res []tripv1.AttachmentMeta
	for _, a := range list {
		res = append(res, tripv1.AttachmentMeta{
			ID:          a.ID,
			TripID:      a.TripID,
			PlaceID:     a.PlaceID,
			DayID:       a.DayID,
			Filename:    a.Filename,
			ContentType: a.ContentType,
			FileSize:    a.FileSize,
			StorageKey:  a.StorageKey,
			DownloadURL: a.DownloadURL,
			Visibility:  tripv1.AttachmentVisibility(a.Visibility),
			CreatedAt:   a.CreatedAt,
			IsReady:     a.IsReady,
			OwnerType:   a.OwnerType,
		})
	}

	return &tripv1.ListAttachmentsResponse{Attachments: res}, nil
}

func (h *TripHandler) DeleteAttachment(ctx context.Context, req *tripv1.DeleteAttachmentRequest) (*tripv1.DeleteAttachmentResponse, error) {
	uid, err := GetUserID(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	if err := h.attachmentUsecase.DeleteAttachment(ctx, req.TripID, req.AttachmentID, uid); err != nil {
		return nil, mapError(err)
	}

	return &tripv1.DeleteAttachmentResponse{Success: true}, nil
}

func mapTripStatus(s domain.TripStatus) tripv1.TripStatus {
	switch s {
	case domain.TripStatusDraft:
		return tripv1.TripStatusDraft
	case domain.TripStatusUpcoming:
		return tripv1.TripStatusUpcoming
	case domain.TripStatusActive:
		return tripv1.TripStatusActive
	case domain.TripStatusCompleted:
		return tripv1.TripStatusCompleted
	default:
		return tripv1.TripStatusUnspecified
	}
}

func (h *TripHandler) ClonePublicTrip(ctx context.Context, req *tripv1.ClonePublicTripRequest) (*tripv1.ClonePublicTripResponse, error) {
	targetOwnerID := req.TargetOwnerID
	if targetOwnerID == "" {
		if uid, err := GetUserID(ctx); err == nil && uid != "" {
			targetOwnerID = uid
		}
	}
	if targetOwnerID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required to clone trip"))
	}

	cloned, err := h.shareUsecase.ClonePublicTrip(ctx, req.Slug, targetOwnerID)
	if err != nil {
		return nil, mapError(err)
	}

	protoTrip := &tripv1.Trip{
		ID:          cloned.ID,
		OwnerID:     cloned.OwnerID,
		Title:       cloned.Title,
		Description: cloned.Description,
		StartDate:   cloned.StartDate,
		EndDate:     cloned.EndDate,
		Status:      mapTripStatus(cloned.Status),
		CoverURL:    cloned.CoverURL,
		Currency:    cloned.Currency,
		BudgetLimit: cloned.BudgetLimit,
		CreatedAt:   cloned.CreatedAt,
		UpdatedAt:   cloned.UpdatedAt,
	}

	return &tripv1.ClonePublicTripResponse{Trip: protoTrip}, nil
}

func (h *TripHandler) GetPublicShareTrip(ctx context.Context, req *tripv1.GetPublicShareTripRequest) (*tripv1.GetPublicShareTripResponse, error) {
	pubView, err := h.shareUsecase.GetPublicTripBySlug(ctx, req.Slug)
	if err != nil {
		return nil, mapError(err)
	}

	var pubDays []tripv1.PublicTripDay
	for _, d := range pubView.Days {
		var places []tripv1.PublicPlaceItem
		for _, p := range d.Places {
			places = append(places, tripv1.PublicPlaceItem{
				ID: p.ID,
				Place: tripv1.PlaceRef{
					PlaceID:   p.Place.PlaceID,
					Name:      p.Place.Name,
					Latitude:  p.Place.Latitude,
					Longitude: p.Place.Longitude,
					Category:  tripv1.PlaceCategory(p.Place.Category),
					Address:   p.Place.Address,
				},
				Position:        p.Position,
				DurationMinutes: p.DurationMinutes,
				Notes:           p.Notes,
			})
		}
		pubDays = append(pubDays, tripv1.PublicTripDay{
			ID:       d.ID,
			DayIndex: d.DayIndex,
			Date:     d.Date,
			Title:    d.Title,
			Notes:    d.Notes,
			Places:   places,
		})
	}

	var pubAttachments []tripv1.PublicAttachmentMeta
	for _, a := range pubView.Attachments {
		pubAttachments = append(pubAttachments, tripv1.PublicAttachmentMeta{
			ID:          a.ID,
			Filename:    a.Filename,
			ContentType: a.ContentType,
			DownloadURL: a.DownloadURL,
			PlaceID:     a.PlaceID,
			DayID:       a.DayID,
		})
	}

	return &tripv1.GetPublicShareTripResponse{
		Trip: &tripv1.PublicTrip{
			ID:               pubView.ID,
			Title:            pubView.Title,
			Description:      pubView.Description,
			StartDate:        pubView.StartDate,
			EndDate:          pubView.EndDate,
			CoverURL:         pubView.CoverURL,
			Days:             pubDays,
			Attachments:      pubAttachments,
			OwnerDisplayName: pubView.OwnerDisplayName,
		},
	}, nil
}