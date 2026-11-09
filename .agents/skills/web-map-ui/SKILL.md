---
name: web-map-ui
description: >-
  Builds AtlasFlow map-first Next.js UI — MapLibre canvas, glass panels, world
  vs trip modes, day timeline DnD, route polylines, place drawer, Zustand +
  TanStack Query patterns. Use when working in apps/web or on Figma-aligned UI.
---

# Web map-first UI

## Architecture

```text
apps/web/
  app/                    # Next.js app router
  src/
    features/
      world-map/
      trip-workspace/
      place-drawer/
      budget/
      checklists/
      share/
    entities/             # trip, place types + api mappers
    shared/
      ui/                 # GlassPanel, Dock, tokens
      map/                # MapCanvas, layers, flyTo helpers
      lib/query/          # QueryClient provider
```

Feature folders own UI + hooks; no cross-feature deep imports upward.

## State split

| Concern | Store |
| --- | --- |
| Trips, places, budget, checklists | TanStack Query |
| UI mode `world|trips`, selectedTripId, selectedDayId, drawerPlaceId, filters | Zustand map UI store |
| Theme dark/light | Zustand or next-themes → drives CSS vars + map style URL |

## Map canvas rules

1. Single `MapCanvas` mounts once; feature layers toggle by mode.
2. World mode: country fill GeoJSON + filtered pins.
3. Trip mode: numbered pins for selected day + LineString route; dim other days.
4. On `selectedDayId` change → `fitBounds` / `flyTo` bbox of that day's coordinates (padding for left glass panel).

## Glass UI tokens (CSS vars)

```css
:root[data-theme="dark"] {
  --map-base: #0B0F17;
  --glass: rgba(17, 23, 34, 0.78);
  --glass-border: rgba(255, 255, 255, 0.09);
  --text: #F8FAFC;
  --text-muted: #94A3B8;
  --accent: #38BDF8;
  --pin-food: #F97316;
  --pin-stay: #A855F7;
  --pin-sight: #10B981;
  --pin-visited: #FACC15;
}
:root[data-theme="light"] {
  --map-base: #F4F1EA;
  --glass: rgba(255, 255, 255, 0.82);
  --glass-border: rgba(15, 23, 42, 0.08);
  --text: #0F172A;
  --text-muted: #64748B;
  --accent: #2563EB;
  --pin-food: #EA580C;
  --pin-stay: #9333EA;
  --pin-sight: #059669;
  --pin-visited: #CA8A04;
}
```

`GlassPanel`: `background: var(--glass); backdrop-filter: blur(20px); border: 1px solid var(--glass-border);`

## DnD timeline

- Library: `@dnd-kit/core` (preferred) or equivalent.
- Droppable: backlog + each day.
- On drag end → call `reorderItem` mutation with neighbors; invalidate trip query.
- Optimistic update allowed if rollback restores previous order.

## Place drawer

Floating over map (not a route change). Sections: photo/title, category+status, rating, duration, cost, notes, links, attachments, actions.

Geocode search: debounce 200–300ms → geo-places suggest → pin preview on map before confirm.

## Query keys

```ts
['trips']
['trip', tripId]
['places', { bbox, status, category }]
['geo-stats']
['share', slug]
```

## Do / Don't

- Do keep map interactions at 60fps — avoid re-creating Map instance on panel toggles.
- Don't put budget charts in the first viewport of world map.
- Don't fetch all world places without bbox.
