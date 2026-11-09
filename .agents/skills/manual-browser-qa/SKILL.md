---
name: manual-browser-qa
description: >-
  Manual browser QA checklist for AtlasFlow — boot local stack, drive the UI in
  a real browser against locked release ACs, capture pass/fail evidence, and
  gate release-ship. Use when running manual-qa agent or browser exploratory tests.
---

# Manual browser QA

## Preconditions

1. Read locked release: `docs/releases/vX.Y.Z.md`
2. Start stack with Taskfile (`task up`, `task dev:web`, or documented equivalents) — **not** Makefile
3. Confirm:
   - Web UI loads
   - Auth/gateway reachable
   - Map canvas renders (MapLibre)

## Smoke path (always)

| Step | Action | Pass if |
| --- | --- | --- |
| S1 | Open home / world map | Map full-viewport; glass panels float |
| S2 | Toggle World ↔ Trips | Mode switches without full reload crash |
| S3 | Open a trip workspace | Day timeline + map pins visible |
| S4 | Open place drawer | Metadata editable / visible |
| S5 | Open `/share/:slug` as guest | No budget/private docs; chrome OK |

## Feature matrix (run if in release)

| Area | Checks |
| --- | --- |
| i18n | Switch EN↔RU; plurals look right; author notes on share stay original language |
| Theme | `system` follows OS; manual light/dark persists; map style sync; layers survive style change |
| Satellite | Streets ↔ Satellite toggle; attribution visible; fallback on tile fail |
| MinIO photos | Upload cover/place image; preview; private voucher not on public share |
| Clone share | Copy to My Trips (auth); guest prompted to login; clone omits private data |
| Distance legs | ≥2 stops show km/time badges; day total updates after DnD |
| Checklist templates | Apply template; no dup titles; readiness % updates |
| Export | Download GPX and GeoJSON; files open/parse |

## Browser technique

1. Navigate → snapshot → act → snapshot (never click blind)
2. After locale/theme change, verify **no** hydration flash/wrong theme sticky bugs
3. Watch console for uncaught errors; note them in evidence
4. Mobile width: resize once (~390px) for share + day dock sanity

## Exit criteria

- All locked UI ACs: PASS (or PM-waived in release doc)
- Blockers: 0
- Write `docs/releases/vX.Y.Z-manual-qa.md` summary (optional but preferred) with PASS/FAIL table
- Signal orchestrator: `MANUAL_QA_PASSED=yes` before `release-ship`
