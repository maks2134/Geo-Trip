---
name: trip-core
description: >-
  Implements trip planner domain in trip-core — trips, days, backlog, lexorank
  DnD reorder, budget tracker, checklists, and attachment metadata with
  presigned uploads. Use when building itinerary, expenses, todos, or file
  vouchers.
---

# Trip Core skill

## Aggregate sketch

```text
Trip
 ├── Days[] (date, index)
 │    └── Items[] (placeRef, position, duration, cost?, notes?)
 ├── Backlog[] (placeRef, position)   // unassigned
 ├── Expenses[] (optional day_id, category, amount, currency)
 ├── ChecklistSections[] → Items[]
 └── Attachments[] (owner: trip|day|place, visibility, storage_key)
```

`placeRef` = `{ place_id, name, lat, lng, category }` snapshot so itinerary map works if geo-places is slow.

## Lexorank / position

Never use `1,2,3…` that forces renumbering whole day.

```go
// position is a lexicographic string between neighbors
func RankBetween(left, right string) (string, error)
```

On DnD:

1. Client sends `{ item_id, target_day_id|null, before_id?, after_id? }`
2. Usecase loads neighbors, computes `RankBetween`, saves one row
3. If ranks collide / exhausted → rebalance that day only

Backlog uses `day_id IS NULL` + same position scheme.

## Day route

Server may return ordered items; **client** draws polyline and asks routing API optionally. For MVP distance estimate:

```text
sum ST_Distance(geography) between consecutive points (via geo-places helper or local haversine on snapshots)
```

Keep routing provider behind a port if added later.

## Budget

- Expense categories: `stay` | `food` | `transport` | `tickets` | `other`
- Place card cost creates/links an expense with `place_item_id`
- Summary: `limit`, `spent`, `by_category[]`, `by_day[]`
- Currency: single trip currency in MVP (`EUR`/`USD`/…) — no FX engine

## Checklists

Sections fixed for MVP:

1. `pre_flight` — visa, insurance, tickets
2. `packing`
3. `buy_on_site`

Readiness % = `done / total` across sections; expose on trip card DTO.

## Attachments

```text
1. Client → CreateUpload(filename, content_type, owner) 
2. trip-core validates MIME (pdf, jpeg, png, webp), size; writes AttachmentMeta(private by default)
3. Returns presigned PUT URL (MinIO/S3)
4. Client uploads directly
5. ConfirmUpload marks ready
```

Public share includes only attachments with `visibility=public`.

## Usecase checklist (authz)

Every mutating method:

1. Load trip
2. `trip.AssertOwner(userID)`
3. Apply change
4. Save

## Tests to always add

- Reorder between two items updates only `position`
- Moving backlog → day clears backlog and sets day_id
- Budget summary math with place costs + free expenses
- Checklist readiness rounding
- Forbidden when `owner_id` mismatch
