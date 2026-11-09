# Service template snippets

## Domain entity

```go
package domain

type TripID string

type TripStatus string

const (
    TripDraft     TripStatus = "draft"
    TripUpcoming  TripStatus = "upcoming"
    TripActive    TripStatus = "active"
    TripCompleted TripStatus = "completed"
)

type Trip struct {
    ID        TripID
    OwnerID   string
    Title     string
    StartDate time.Time
    EndDate   time.Time
    Status    TripStatus
    CoverURL  string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewTrip(id, owner, title string, start, end, now time.Time) Trip {
    return Trip{
        ID: TripID(id), OwnerID: owner, Title: title,
        StartDate: start, EndDate: end, Status: TripDraft,
        CreatedAt: now, UpdatedAt: now,
    }
}

func (t Trip) AssertOwner(userID string) error {
    if t.OwnerID != userID {
        return ErrForbidden
    }
    return nil
}
```

## Postgres adapter sketch

```go
type TripRepo struct{ pool *pgxpool.Pool }

func (r *TripRepo) Save(ctx context.Context, t domain.Trip) error {
    const q = `
        INSERT INTO trips (id, owner_id, title, start_date, end_date, status, cover_url, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
        ON CONFLICT (id) DO UPDATE SET
          title=EXCLUDED.title, start_date=EXCLUDED.start_date, end_date=EXCLUDED.end_date,
          status=EXCLUDED.status, cover_url=EXCLUDED.cover_url, updated_at=EXCLUDED.updated_at`
    _, err := r.pool.Exec(ctx, q,
        t.ID, t.OwnerID, t.Title, t.StartDate, t.EndDate, t.Status, t.CoverURL, t.CreatedAt, t.UpdatedAt,
    )
    return err
}
```

## Connect handler sketch

```go
func (h *TripHandler) CreateTrip(ctx context.Context, req *connect.Request[tripv1.CreateTripRequest]) (*connect.Response[tripv1.Trip], error) {
    userID, err := auth.UserID(ctx)
    if err != nil {
        return nil, connect.NewError(connect.CodeUnauthenticated, err)
    }
    out, err := h.svc.Create(ctx, usecase.CreateTripInput{
        OwnerID: userID,
        Title:   req.Msg.GetTitle(),
        Start:   req.Msg.GetStartDate().AsTime(),
        End:     req.Msg.GetEndDate().AsTime(),
    })
    if err != nil {
        return nil, mapErr(err)
    }
    return connect.NewResponse(toProtoTrip(out)), nil
}
```
