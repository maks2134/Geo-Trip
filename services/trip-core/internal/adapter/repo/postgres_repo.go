package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/domain"
	"github.com/makskozlov/atlasflow/services/trip-core/internal/usecase"
)

var _ usecase.TripRepository = (*PostgresTripRepo)(nil)

type PostgresTripRepo struct {
	db *sql.DB
}

func NewPostgresTripRepo(db *sql.DB) *PostgresTripRepo {
	return &PostgresTripRepo{db: db}
}

// Trips
func (r *PostgresTripRepo) CreateTrip(ctx context.Context, t domain.Trip) error {
	const q = `
		INSERT INTO trips (id, owner_id, title, description, start_date, end_date, status, cover_url, currency, budget_limit, share_slug, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	_, err := r.db.ExecContext(ctx, q,
		t.ID, t.OwnerID, t.Title, t.Description, t.StartDate, t.EndDate, string(t.Status),
		t.CoverURL, t.Currency, t.BudgetLimit, t.ShareSlug, t.IsPublic, t.CreatedAt, t.UpdatedAt,
	)
	return err
}

func (r *PostgresTripRepo) GetTrip(ctx context.Context, id string) (domain.Trip, error) {
	const q = `
		SELECT id, owner_id, title, description, start_date, end_date, status, cover_url, currency, budget_limit, share_slug, is_public, created_at, updated_at
		FROM trips WHERE id = $1`
	var t domain.Trip
	var status string
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&t.ID, &t.OwnerID, &t.Title, &t.Description, &t.StartDate, &t.EndDate, &status,
		&t.CoverURL, &t.Currency, &t.BudgetLimit, &t.ShareSlug, &t.IsPublic, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Trip{}, domain.ErrNotFound
	}
	t.Status = domain.TripStatus(status)
	return t, err
}

func (r *PostgresTripRepo) ListTrips(ctx context.Context, ownerID string, status *domain.TripStatus, limit, offset int32) ([]domain.Trip, int32, error) {
	var count int32
	countQ := `SELECT COUNT(*) FROM trips WHERE owner_id = $1`
	args := []any{ownerID}
	if status != nil {
		countQ += ` AND status = $2`
		args = append(args, string(*status))
	}
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&count); err != nil {
		return nil, 0, err
	}

	q := `SELECT id, owner_id, title, description, start_date, end_date, status, cover_url, currency, budget_limit, share_slug, is_public, created_at, updated_at
	      FROM trips WHERE owner_id = $1`
	queryArgs := []any{ownerID}
	argIdx := 2
	if status != nil {
		q += fmt.Sprintf(" AND status = $%d", argIdx)
		queryArgs = append(queryArgs, string(*status))
		argIdx++
	}
	q += " ORDER BY created_at DESC"
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
		queryArgs = append(queryArgs, limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, q, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var trips []domain.Trip
	for rows.Next() {
		var t domain.Trip
		var st string
		if err := rows.Scan(
			&t.ID, &t.OwnerID, &t.Title, &t.Description, &t.StartDate, &t.EndDate, &st,
			&t.CoverURL, &t.Currency, &t.BudgetLimit, &t.ShareSlug, &t.IsPublic, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		t.Status = domain.TripStatus(st)
		trips = append(trips, t)
	}
	return trips, count, rows.Err()
}

func (r *PostgresTripRepo) UpdateTrip(ctx context.Context, t domain.Trip) error {
	const q = `
		UPDATE trips
		SET title = $1, description = $2, start_date = $3, end_date = $4, status = $5, cover_url = $6,
		    currency = $7, budget_limit = $8, share_slug = $9, is_public = $10, updated_at = $11
		WHERE id = $12`
	res, err := r.db.ExecContext(ctx, q,
		t.Title, t.Description, t.StartDate, t.EndDate, string(t.Status), t.CoverURL,
		t.Currency, t.BudgetLimit, t.ShareSlug, t.IsPublic, t.UpdatedAt, t.ID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PostgresTripRepo) DeleteTrip(ctx context.Context, id string) error {
	const q = `DELETE FROM trips WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

// Days
func (r *PostgresTripRepo) CreateDay(ctx context.Context, d domain.TripDay) error {
	const q = `
		INSERT INTO trip_days (id, trip_id, day_index, date, title, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())`
	_, err := r.db.ExecContext(ctx, q, d.ID, d.TripID, d.DayIndex, d.Date, d.Title, d.Notes)
	return err
}

func (r *PostgresTripRepo) GetDay(ctx context.Context, id string) (domain.TripDay, error) {
	const q = `SELECT id, trip_id, day_index, date, title, notes FROM trip_days WHERE id = $1`
	var d domain.TripDay
	err := r.db.QueryRowContext(ctx, q, id).Scan(&d.ID, &d.TripID, &d.DayIndex, &d.Date, &d.Title, &d.Notes)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.TripDay{}, domain.ErrNotFound
	}
	return d, err
}

func (r *PostgresTripRepo) ListDays(ctx context.Context, tripID string) ([]domain.TripDay, error) {
	const q = `SELECT id, trip_id, day_index, date, title, notes FROM trip_days WHERE trip_id = $1 ORDER BY day_index ASC`
	rows, err := r.db.QueryContext(ctx, q, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []domain.TripDay
	for rows.Next() {
		var d domain.TripDay
		if err := rows.Scan(&d.ID, &d.TripID, &d.DayIndex, &d.Date, &d.Title, &d.Notes); err != nil {
			return nil, err
		}
		days = append(days, d)
	}
	return days, rows.Err()
}

func (r *PostgresTripRepo) UpdateDay(ctx context.Context, d domain.TripDay) error {
	const q = `UPDATE trip_days SET date = $1, title = $2, notes = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, q, d.Date, d.Title, d.Notes, d.ID)
	return err
}

func (r *PostgresTripRepo) DeleteDay(ctx context.Context, id string) error {
	const q = `DELETE FROM trip_days WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

// Itinerary Items
func (r *PostgresTripRepo) CreatePlaceItem(ctx context.Context, it domain.ItineraryPlaceItem) error {
	const q = `
		INSERT INTO itinerary_items (id, trip_id, day_id, place_id, place_name, latitude, longitude, category, address, position, duration_minutes, cost, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err := r.db.ExecContext(ctx, q,
		it.ID, it.TripID, it.DayID, it.Place.PlaceID, it.Place.Name, it.Place.Latitude, it.Place.Longitude,
		int(it.Place.Category), it.Place.Address, it.Position, it.DurationMinutes, it.Cost, it.Notes,
		it.CreatedAt, it.UpdatedAt,
	)
	return err
}

func (r *PostgresTripRepo) GetPlaceItem(ctx context.Context, id string) (domain.ItineraryPlaceItem, error) {
	const q = `
		SELECT id, trip_id, day_id, place_id, place_name, latitude, longitude, category, address, position, duration_minutes, cost, notes, created_at, updated_at
		FROM itinerary_items WHERE id = $1`
	var it domain.ItineraryPlaceItem
	var cat int
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&it.ID, &it.TripID, &it.DayID, &it.Place.PlaceID, &it.Place.Name, &it.Place.Latitude, &it.Place.Longitude,
		&cat, &it.Place.Address, &it.Position, &it.DurationMinutes, &it.Cost, &it.Notes,
		&it.CreatedAt, &it.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ItineraryPlaceItem{}, domain.ErrNotFound
	}
	it.Place.Category = domain.PlaceCategory(cat)
	return it, err
}

func (r *PostgresTripRepo) ListPlaceItems(ctx context.Context, tripID string) ([]domain.ItineraryPlaceItem, error) {
	const q = `
		SELECT id, trip_id, day_id, place_id, place_name, latitude, longitude, category, address, position, duration_minutes, cost, notes, created_at, updated_at
		FROM itinerary_items WHERE trip_id = $1 ORDER BY position ASC`
	rows, err := r.db.QueryContext(ctx, q, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.ItineraryPlaceItem
	for rows.Next() {
		var it domain.ItineraryPlaceItem
		var cat int
		if err := rows.Scan(
			&it.ID, &it.TripID, &it.DayID, &it.Place.PlaceID, &it.Place.Name, &it.Place.Latitude, &it.Place.Longitude,
			&cat, &it.Place.Address, &it.Position, &it.DurationMinutes, &it.Cost, &it.Notes,
			&it.CreatedAt, &it.UpdatedAt,
		); err != nil {
			return nil, err
		}
		it.Place.Category = domain.PlaceCategory(cat)
		items = append(items, it)
	}
	return items, rows.Err()
}

func (r *PostgresTripRepo) UpdatePlaceItem(ctx context.Context, it domain.ItineraryPlaceItem) error {
	const q = `
		UPDATE itinerary_items
		SET day_id = $1, position = $2, duration_minutes = $3, cost = $4, notes = $5, updated_at = $6
		WHERE id = $7`
	_, err := r.db.ExecContext(ctx, q, it.DayID, it.Position, it.DurationMinutes, it.Cost, it.Notes, it.UpdatedAt, it.ID)
	return err
}

func (r *PostgresTripRepo) DeletePlaceItem(ctx context.Context, tripID, id string) error {
	const q = `DELETE FROM itinerary_items WHERE trip_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, q, tripID, id)
	return err
}

// Expenses
func (r *PostgresTripRepo) CreateExpense(ctx context.Context, e domain.Expense) error {
	const q = `
		INSERT INTO expenses (id, trip_id, day_id, place_item_id, category, amount, currency, title, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, q,
		e.ID, e.TripID, e.DayID, e.PlaceItemID, int(e.Category), e.Amount, e.Currency, e.Title, e.Notes, e.CreatedAt,
	)
	return err
}

func (r *PostgresTripRepo) ListExpenses(ctx context.Context, tripID string, dayID *string, category *domain.ExpenseCategory) ([]domain.Expense, error) {
	q := `SELECT id, trip_id, day_id, place_item_id, category, amount, currency, title, notes, created_at FROM expenses WHERE trip_id = $1`
	args := []any{tripID}
	argIdx := 2
	if dayID != nil {
		q += fmt.Sprintf(" AND day_id = $%d", argIdx)
		args = append(args, *dayID)
		argIdx++
	}
	if category != nil {
		q += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, int(*category))
		argIdx++
	}
	q += " ORDER BY created_at ASC"

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Expense
	for rows.Next() {
		var e domain.Expense
		var cat int
		if err := rows.Scan(
			&e.ID, &e.TripID, &e.DayID, &e.PlaceItemID, &cat, &e.Amount, &e.Currency, &e.Title, &e.Notes, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		e.Category = domain.ExpenseCategory(cat)
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *PostgresTripRepo) DeleteExpense(ctx context.Context, tripID, id string) error {
	const q = `DELETE FROM expenses WHERE trip_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, q, tripID, id)
	return err
}

// Checklist
func (r *PostgresTripRepo) CreateChecklistItem(ctx context.Context, it domain.ChecklistItem) error {
	const q = `
		INSERT INTO checklist_items (id, trip_id, section, title, is_done, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, q, it.ID, it.TripID, int(it.Section), it.Title, it.IsDone, it.Position, it.CreatedAt, it.UpdatedAt)
	return err
}

func (r *PostgresTripRepo) GetChecklistItem(ctx context.Context, id string) (domain.ChecklistItem, error) {
	const q = `SELECT id, trip_id, section, title, is_done, position, created_at, updated_at FROM checklist_items WHERE id = $1`
	var it domain.ChecklistItem
	var sec int
	err := r.db.QueryRowContext(ctx, q, id).Scan(&it.ID, &it.TripID, &sec, &it.Title, &it.IsDone, &it.Position, &it.CreatedAt, &it.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ChecklistItem{}, domain.ErrNotFound
	}
	it.Section = domain.ChecklistSection(sec)
	return it, err
}

func (r *PostgresTripRepo) ListChecklistItems(ctx context.Context, tripID string, section *domain.ChecklistSection) ([]domain.ChecklistItem, error) {
	q := `SELECT id, trip_id, section, title, is_done, position, created_at, updated_at FROM checklist_items WHERE trip_id = $1`
	args := []any{tripID}
	if section != nil {
		q += " AND section = $2"
		args = append(args, int(*section))
	}
	q += " ORDER BY position ASC"

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.ChecklistItem
	for rows.Next() {
		var it domain.ChecklistItem
		var sec int
		if err := rows.Scan(&it.ID, &it.TripID, &sec, &it.Title, &it.IsDone, &it.Position, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, err
		}
		it.Section = domain.ChecklistSection(sec)
		list = append(list, it)
	}
	return list, rows.Err()
}

func (r *PostgresTripRepo) UpdateChecklistItem(ctx context.Context, it domain.ChecklistItem) error {
	const q = `UPDATE checklist_items SET is_done = $1, title = $2, position = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, q, it.IsDone, it.Title, it.Position, it.UpdatedAt, it.ID)
	return err
}

// Attachments
func (r *PostgresTripRepo) CreateAttachment(ctx context.Context, a domain.Attachment) error {
	const q = `
		INSERT INTO attachments (id, trip_id, place_id, day_id, filename, content_type, file_size, storage_key, download_url, visibility, is_ready, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, q,
		a.ID, a.TripID, a.PlaceID, a.DayID, a.Filename, a.ContentType, a.FileSize,
		a.StorageKey, a.DownloadURL, int(a.Visibility), a.IsReady, a.CreatedAt,
	)
	return err
}

func (r *PostgresTripRepo) GetAttachment(ctx context.Context, id string) (domain.Attachment, error) {
	const q = `
		SELECT id, trip_id, place_id, day_id, filename, content_type, file_size, storage_key, download_url, visibility, is_ready, created_at
		FROM attachments WHERE id = $1`
	var a domain.Attachment
	var vis int
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&a.ID, &a.TripID, &a.PlaceID, &a.DayID, &a.Filename, &a.ContentType, &a.FileSize,
		&a.StorageKey, &a.DownloadURL, &vis, &a.IsReady, &a.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Attachment{}, domain.ErrNotFound
	}
	a.Visibility = domain.AttachmentVisibility(vis)
	return a, err
}

func (r *PostgresTripRepo) ListAttachments(ctx context.Context, tripID string) ([]domain.Attachment, error) {
	return r.ListAttachmentsFiltered(ctx, tripID, nil, nil, nil)
}

func (r *PostgresTripRepo) ListAttachmentsFiltered(
	ctx context.Context,
	tripID string,
	placeID, dayID *string,
	visibility *domain.AttachmentVisibility,
) ([]domain.Attachment, error) {
	q := `
		SELECT id, trip_id, place_id, day_id, filename, content_type, file_size, storage_key, download_url, visibility, is_ready, created_at
		FROM attachments WHERE trip_id = $1`
	args := []interface{}{tripID}
	argIdx := 2

	if placeID != nil {
		q += fmt.Sprintf(" AND place_id = $%d", argIdx)
		args = append(args, *placeID)
		argIdx++
	}
	if dayID != nil {
		q += fmt.Sprintf(" AND day_id = $%d", argIdx)
		args = append(args, *dayID)
		argIdx++
	}
	if visibility != nil {
		q += fmt.Sprintf(" AND visibility = $%d", argIdx)
		args = append(args, int(*visibility))
		argIdx++
	}

	q += " ORDER BY created_at ASC"

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Attachment
	for rows.Next() {
		var a domain.Attachment
		var vis int
		if err := rows.Scan(
			&a.ID, &a.TripID, &a.PlaceID, &a.DayID, &a.Filename, &a.ContentType, &a.FileSize,
			&a.StorageKey, &a.DownloadURL, &vis, &a.IsReady, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		a.Visibility = domain.AttachmentVisibility(vis)
		list = append(list, a)
	}
	return list, rows.Err()
}

func (r *PostgresTripRepo) ConfirmAttachmentReady(ctx context.Context, id string) error {
	const q = `UPDATE attachments SET is_ready = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

func (r *PostgresTripRepo) DeleteAttachment(ctx context.Context, tripID, id string) error {
	const q = `DELETE FROM attachments WHERE id = $1 AND trip_id = $2`
	_, err := r.db.ExecContext(ctx, q, id, tripID)
	return err
}

// Shares
func (r *PostgresTripRepo) CreateShare(ctx context.Context, s domain.ShareRecord) error {
	const q = `
		INSERT INTO trip_shares (id, trip_id, slug, hide_budget, hide_private_docs, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, q, s.ID, s.TripID, s.Slug, s.HideBudget, s.HidePrivateDocs, s.CreatedAt)
	return err
}

func (r *PostgresTripRepo) GetShareBySlug(ctx context.Context, slug string) (domain.ShareRecord, error) {
	const q = `SELECT id, trip_id, slug, hide_budget, hide_private_docs, revoked_at, created_at FROM trip_shares WHERE slug = $1`
	var s domain.ShareRecord
	err := r.db.QueryRowContext(ctx, q, slug).Scan(&s.ID, &s.TripID, &s.Slug, &s.HideBudget, &s.HidePrivateDocs, &s.RevokedAt, &s.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ShareRecord{}, domain.ErrNotFound
	}
	return s, err
}

func (r *PostgresTripRepo) RevokeShare(ctx context.Context, tripID string) error {
	const q = `UPDATE trip_shares SET revoked_at = NOW() WHERE trip_id = $1 AND revoked_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, tripID)
	return err
}