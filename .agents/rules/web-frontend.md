---
trigger: glob
globs: "apps/web/**/*.{ts,tsx,css}, **/*.tsx"
description: Map-first Next.js UI rules for AtlasFlow glass panels and client state.
---

# Web frontend rules

1. Map occupies full viewport; chrome is floating glass (`backdrop-blur`, translucent border). No card grids in the hero/map canvas.
2. Server data → TanStack Query. Ephemeral UI/map (mode, selectedDayId, drawerPlaceId, filters) → Zustand.
3. Theme tokens: use CSS variables for Dark Midnight Cartography / Light Editorial Daylight; sync MapLibre style with theme.
4. Pin category colors: Food / Stay / Sight / Visited — use design tokens, not hard-coded one-offs in components.
5. Public `/share/[slug]` routes must not fetch budget or private attachment URLs.
6. Prefer skill `web-map-ui` for map/DnD/route rendering work; `public-share` for guest views.
