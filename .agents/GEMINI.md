# AtlasFlow / NomadGlass — Project Context

Personal travel scratch-map + day-by-day trip planner MVP.

## Stack

- Monorepo: `/Users/makskozlov/pet`
- `apps/web` — Next.js (TypeScript, Tailwind), MapLibre/Mapbox, Zustand, TanStack Query
- `services/gateway-auth` — HTTP/gRPC gateway, JWT + refresh, public `/share/:slug`
- `services/trip-core` — trips, days, lexorank reorder, budget, checklists, file metadata (S3/MinIO presigned)
- `services/geo-places` — PostGIS `GEOGRAPHY(Point, 4326)`, bbox search, geocoding, country/city stats
- `pkg/`, `proto/` — shared Go 1.26 libs + ConnectRPC/protobuf
- `deploy/` — docker-compose: PostgreSQL+PostGIS, MinIO, optional Redis

## Product modules (MVP)

1. World Scratch Map — visited countries polygons, traveler stats, pin filters (Visited / Planned / Wishlist + categories)
2. Trip Planner — trip card, backlog, day timeline, DnD, route polyline + numbered pins, day flyTo
3. Place Drawer — geocode/click pin, metadata, notes, links, PDF/image attachments
4. Trip Budget — place costs + day/trip expenses, limit, category/day breakdown
5. Checklists — pre-flight / packing / buy-on-site, readiness %
6. Public share — `/share/<slug>` read-only; hide budget + private docs

## UI

Adaptive Glass UI, canvas-first map (100% viewport). Floating panels: backdrop-blur glass, 1px translucent border, dark/light tokens synced with map style.

## Non-negotiables

- Go services: Clean Architecture (domain → usecase → adapters). No business logic in handlers/DB layers.
- Frontend: thin UI; server truth via TanStack Query; map state in Zustand; no DB/API secrets in client.
- Prefer ConnectRPC contracts in `proto/`; no ad-hoc cross-service DB access.
