---
name: atlasflow-mvp
description: >-
  AtlasFlow/NomadGlass MVP product map — module boundaries, ownership, data
  flow, and screen-to-service mapping. Use when planning features, scaffolding
  the monorepo, deciding where code belongs, or clarifying MVP scope.
---

# AtlasFlow MVP map

## Ownership

| Module | Owner | Key entities |
| --- | --- | --- |
| Auth + public share entry | `gateway-auth` | User, Session, ShareSlug routing |
| Trips, days, backlog, budget, checklists, file meta | `trip-core` | Trip, Day, PlaceRef, Expense, ChecklistItem, AttachmentMeta |
| Points, bbox, geocode, country stats | `geo-places` | Place, CountryVisit, GeoStats |
| Map UI + glass panels | `apps/web` | Map store, query hooks, drawers |

## Cross-service rules

- `trip-core` stores `place_id` + denormalized lat/lng snapshot for itinerary rendering offline of geo service.
- Authoritative geometry and country aggregation live in `geo-places`.
- `gateway-auth` terminates TLS/JWT and exposes BFF/public share; forward `X-User-Id` (or internal JWT) to services.
- Files: client gets presigned URL from `trip-core`; binary never streams through trip-core.

## Screen → API (MVP)

| Screen | Primary reads/writes |
| --- | --- |
| World map + stats | `geo-places` stats + places by bbox/filter; trip list from `trip-core` |
| Active trip workspace | trip + days + places + reorder; map route from ordered places |
| Place drawer | place CRUD; attachments presign |
| Budget / checklists tabs | expenses + checklist endpoints on `trip-core` |
| Public share | gateway `/share/:slug` → sanitized trip DTO (no budget/private docs) |

## Status enums (keep stable)

- Trip: `draft` | `upcoming` | `active` | `completed`
- Place pin: `visited` | `planned` | `wishlist`
- Place category: `food` | `stay` | `sight` | `nature` | `transport`

## Anti-scope (do not build in MVP)

- Multi-user collaborative editing, live cursors
- Social feed / followers
- Offline-first sync engine
- Custom map tile hosting
- AI itinerary generation

## When scaffolding

1. Create `go.work` + modules per service under `services/*`.
2. Shared contracts in `proto/`; generate into each service or `pkg/gen`.
3. `deploy/docker-compose.yml` for PostGIS + MinIO before feature work that needs them.

See also: `go-clean-arch`, `trip-core`, `geo-postgis`, `web-map-ui`, `public-share`.
