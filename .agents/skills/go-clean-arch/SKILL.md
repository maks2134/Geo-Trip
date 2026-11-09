---
name: go-clean-arch
description: >-
  Implements AtlasFlow Go services with Clean Architecture — domain, usecase
  ports, adapters, cmd wiring, ConnectRPC handlers, and tests. Use when creating
  or changing gateway-auth, trip-core, geo-places, pkg, or proto Go code.
---

# Go Clean Architecture (AtlasFlow)

## Target layout (per service)

```text
services/trip-core/
  cmd/trip-core/main.go
  internal/
    domain/           # entities, value objects, domain errors (no infra imports)
    usecase/          # application services + port interfaces
    adapter/
      http/           # Connect/HTTP handlers, DTO mappers
      postgres/       # repository implementations
      s3/             # presign client
    app/              # composition root helpers (optional)
  migrations/
```

Same pattern for `geo-places` and `gateway-auth` (auth may add `adapter/jwt`).

## Dependency rule

```text
cmd → adapter → usecase → domain
         ↘──────────────↗
```

- `domain` imports **only** stdlib (and maybe tiny shared `pkg/ids` if pure).
- `usecase` imports `domain` + defines interfaces for repos/clocks/ID gens.
- `adapter` implements interfaces; maps storage/proto ↔ domain.
- `main` wires concrete adapters into usecases.

## Port style

```go
// usecase/ports.go
type TripRepository interface {
    Get(ctx context.Context, id TripID) (Trip, error)
    Save(ctx context.Context, t Trip) error
}

type Clock interface{ Now() time.Time }
type IDGen interface{ New() string }
```

## Usecase skeleton

```go
type TripService struct {
    trips TripRepository
    ids   IDGen
    clock Clock
}

func (s *TripService) Create(ctx context.Context, in CreateTripInput) (Trip, error) {
    if err := in.Validate(); err != nil {
        return Trip{}, err
    }
    t := NewTrip(s.ids.New(), in.OwnerID, in.Title, in.Start, in.End, s.clock.Now())
    if err := s.trips.Save(ctx, t); err != nil {
        return Trip{}, fmt.Errorf("save trip: %w", err)
    }
    return t, nil
}
```

## Domain errors

```go
var (
    ErrNotFound   = errors.New("not found")
    ErrForbidden  = errors.New("forbidden")
    ErrValidation = errors.New("validation")
    ErrConflict   = errors.New("conflict")
)
```

Handlers map via `errors.Is` → Connect/HTTP codes. Do not leak SQL errors.

## Handler rules

1. Auth identity from context (set by gateway middleware) — never trust body `user_id`.
2. Decode request → call usecase → encode response.
3. No SQL, S3, or map math in handlers.

## Testing

| Layer | What to test | How |
| --- | --- | --- |
| domain | invariants, lexorank helpers | pure unit |
| usecase | orchestration, authz checks | fake repos |
| adapter/postgres | SQL + scans | testcontainers / integration |
| handler | mapping + status codes | lightly, with fake usecase |

Prefer table-driven tests. One file `*_test.go` beside the code under test.

## Wiring (`main`)

Keep `main` boring: load config → logger → DB pool → repos → usecases → handler mux → serve + graceful shutdown.

## Shared `pkg/`

Only cross-cutting pure infra: config, slog setup, tracing, auth context helpers. **No** trip/geo business types in `pkg/`.

## Checklist before finishing a backend change

- [ ] New logic is in domain/usecase, not handler/repo
- [ ] Ports updated if new IO needed
- [ ] Domain error + handler mapping
- [ ] Table-driven usecase test
- [ ] Proto/API contract updated if external shape changed

Detailed templates: [references/service-template.md](references/service-template.md)
