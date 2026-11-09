# Agent Operating Rules — AtlasFlow

## Default behavior

1. Before coding a feature, identify the owning service (`trip-core` | `geo-places` | `gateway-auth` | `apps/web`) and keep the change inside that boundary.
2. Prefer extending existing packages/patterns over inventing new frameworks.
3. Ship small, compilable increments: domain → usecase → repo/handler → tests → UI wiring.
4. When a skill matches the task, read and follow it from `.agents/skills/`.
5. Product discovery / release scope → `business-analyst` + `product-manager`. A release feature set closes only with both approvals; then ship via `release-ship` (GitHub + semver tag/Release).
6. Before ship: `verifier` (automated) then `manual-qa` (real browser clicks against locked ACs). Do not ship without `MANUAL_QA_PASSED`.

## Backend (Go)

- Clean Architecture layers: `domain/` (entities, pure rules) → `usecase/` (orchestration, ports as interfaces) → `adapter/` (HTTP/Connect, Postgres/PostGIS, S3) → `cmd/` (wiring only).
- Dependency rule: inward only. Domain never imports adapters, frameworks, or SQL drivers.
- Errors: wrap with `%w`; map to typed domain errors at usecase boundary; handlers translate to API codes.
- Persistence: parameterized SQL only; PostGIS geography in `geo-places` only.
- Ordering (itinerary/backlog): lexorank / fractional `position` strings — never dense integer arrays that require full renumber.
- Auth: gateway validates JWT; services trust gateway identity headers / internal mTLS — never re-implement login in trip/geo.

## Frontend (Next.js)

- Map is the canvas; panels float (glass). Do not build dashboard-first layouts for core screens.
- Data: TanStack Query for server state; Zustand for ephemeral map/UI (selected day, drawer, filters).
- Keep place/trip mutations optimistic only when rollback is defined.
- Public share view must not request private budget/docs endpoints.

## Quality bar

- No god packages (`common`, `utils` dumping grounds).
- Every usecase method has at least one table-driven unit test for happy path + key error.
- Do not invent product scope beyond MVP modules listed in `GEMINI.md`.
- Use Taskfile (`task`), never add Makefile targets.
