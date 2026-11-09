---
name: public-share
description: >-
  Implements read-only public trip share links — slug generation, sanitized
  DTO, gateway routing /share/:slug, and guest map UI without private budget or
  documents. Use when building sharing, guest mode, or copy-itinerary CTA.
---

# Public read-only share

## Flow

```text
Owner → EnableShare(trip_id, { hide_budget: true, hide_private_docs: true })
     → trip-core creates Share { slug, trip_id, flags, revoked_at? }
Gateway GET /share/:slug
     → resolve slug → fetch PublicTripView from trip-core
Guest UI apps/web/app/share/[slug]
     → interactive map + day switcher + public notes
     → CTA "Copy route to my trips" (requires auth)
```

## Slug

- Short, URL-safe (`nanoid` 10–12 or base62). Unique index.
- Rotate/revoke: soft revoke sets `revoked_at`; gateway returns 404.

## PublicTripView (sanitized)

Include:

- Trip title, dates, cover, countries/cities labels
- Days with ordered public places (name, category, lat/lng, public notes, duration)
- Public attachments only

Exclude always (MVP):

- Expenses, budget limit/progress
- Private PDF/images
- Owner identity beyond display name if opted-in later
- Checklist personal packing state (optional: hide entirely in MVP)

## Gateway

- Public route **without** JWT.
- Rate-limit by IP.
- Do not accept mutations on public route except authenticated `ImportSharedTrip`.

## Frontend guest view

- Reuse `MapCanvas` + day dock; strip edit/DnD/upload controls.
- No calls to budget/checklist private endpoints.
- Mobile: bottom day sheet; desktop: left glass panel.

## Security checks

- [ ] Direct object ref on `trip_id` is impossible from guest — only slug
- [ ] Attachment URLs are either public objects or short-lived signed GET for `visibility=public`
- [ ] Revoked slug 404s
- [ ] Import creates a **copy** under the new owner (no live link to source trip data)
