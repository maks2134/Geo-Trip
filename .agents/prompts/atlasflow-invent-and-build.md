# Prompt: invent features → lock → implement via every agent → ship

Before paste: new chat → **Argon / Claude thinking + High** → permissions **Auto**.  
One session: product invents, then engineering builds, each specialist agent fully owns its track.

---

## Primary prompt (paste as-is)

```text
/teamwork

Run a full AtlasFlow product+engineering cycle in multi-agent mode in /Users/makskozlov/pet.
You (orchestrator) coordinate EVERY specialist agent below end-to-end. Do not skip agents. Do not collapse all work into one generic coder. Each agent must receive a full brief, own exclusive paths, finish its deliverables, and report AC evidence before the next dependent track starts (parallel only when paths do not collide).

Remote: https://github.com/maks2134/Geo-Trip
Follow: .agents/GEMINI.md, .agents/AGENTS.md
Skills: atlasflow-mvp, go-clean-arch, trip-core, geo-postgis, web-map-ui, public-share, release-ship, manual-browser-qa

============================================================
PHASE 0 — PRODUCT (invent, do not code yet)
============================================================
Invoke in order:
1) `business-analyst` — invent and deeply specify a strong next-release feature set for AtlasFlow/NomadGlass (scratch map + day-by-day planner, glass map-first).
2) `product-manager` — prioritize with ICE, cut release scope, choose semver, lock the release.

Product rules for invention:
- Stay inside product vision (travel map + itinerary). No realtime multiplayer / social network / AI itinerary generator.
- Invent creatively: propose 8–14 candidates, then select a coherent release theme.
- Prefer high-leverage traveler value + technical fit with existing services (gateway-auth, trip-core, geo-places, apps/web, MinIO, PostGIS, MapLibre).
- You MAY include/extend: i18n, system theme, MinIO media, satellite map, share/clone, exports, checklists, budget, geo stats — but you are free to invent better ideas if they fit.
- Every selected feature needs BA-quality Given/When/Then ACs, scope in/out, service impact, risks.
- Write/lock: docs/releases/vX.Y.Z.md with BA_APPROVED=yes, PM_APPROVED=yes, status=RELEASE_LOCKED, Implementation Tracks naming the eng agents below.
- Task runner standard: Taskfile.yml (`task`), NEVER Makefile (migrate/remove Make).

Stop Phase 0 only when the release doc is locked. Then immediately continue to Phase 1 without waiting for me unless a product decision is truly blocked.

============================================================
PHASE 1 — ENGINEERING (implement locked doc only)
============================================================
Source of truth: the locked docs/releases/vX.Y.Z.md from Phase 0.
Implement ALL in-scope features. Do NOT implement Backlog/Out of scope.

Invoke / brief EACH agent thoroughly (full ownership brief = paths + must-not-touch + AC IDs + skills + done criteria):

1) `deploy-infra`
   Own: deploy/, docker-compose.yml, Taskfile.yml
   Deliver infra needed by the release (MinIO/PostGIS/env samples/task recipes). Delete/avoid Makefile.

2) `proto-contracts`
   Own: proto/ (+ codegen)
   Deliver every new/changed RPC/message required by locked ACs.

3) `go-trip`
   Own: services/trip-core/
   Clean Architecture. Domain/usecase/adapters/tests for trip-side features.

4) `go-gateway`
   Own: services/gateway-auth/
   Auth, BFF/public share routes, sanitization, any new HTTP surfaces.

5) `go-geo`
   Own: services/geo-places/
   PostGIS/geocode/stats/distance features touched by the release.

6) `web-map`
   Own: apps/web/
   Map-first Adaptive Glass UI for all user-facing ACs.

7) `verifier`
   Own: verification only (no feature edits)
   Run task test / task build / web build; Clean Arch import check; privacy checks for /share; AC checklist against the release doc. Fail back to the owning author agent by name until green.

8) `manual-qa`
   Own: real-browser manual testing only (no feature code)
   After verifier is green: boot/use the running app, open a real browser, and personally click through every user-facing locked AC (skill `manual-browser-qa`). Cover map modes, i18n, theme, satellite, uploads, share privacy, and any selected release features. File bugs to owning eng agents with repro+evidence; re-test until `MANUAL_QA_PASSED=yes`. Prefer writing `docs/releases/vX.Y.Z-manual-qa.md` PASS/FAIL table.

Parallelize only non-overlapping tracks after deploy+proto foundations. If a worker is blocked, re-brief or replace that worker — do not silently absorb its files into another agent.

============================================================
PHASE 2 — SHIP
============================================================
When verifier is green AND `manual-qa` reports `MANUAL_QA_PASSED=yes` against locked UI ACs, execute release-ship:
- commit (no secrets, no .env)
- push to GitHub (maks2134/Geo-Trip)
- annotated tag vX.Y.Z
- gh release create with notes from docs/releases/vX.Y.Z.md

Do NOT ship on automated tests alone — browser manual QA is mandatory.

============================================================
WORKING STYLE
============================================================
- Integrity mode: development
- Prefer real runnable code over docs after Phase 0
- Preserve public-share privacy always
- Report to me as milestone summaries only (Phase 0 lock → each track done → verifier → manual-qa → ship)
- Ask me only if GitHub auth/remote is missing, browser cannot start, or a locked AC is self-contradictory

Start Phase 0 now: invent features with BA, lock with PM, then build with every eng agent, run manual browser QA, then ship.
```

---

## Swarm fallback (no /teamwork)

```text
Run a native swarm. You are the coordinator — do not write feature code yourself. Work every agent role thoroughly.

PHASE 0: Add @ba and @pm.
@ba invents 8–14 AtlasFlow features with full ACs (vision: scratch map + day planner; no multiplayer/social/AI itinerary).
@pm ICE-prioritizes, cuts release, writes docs/releases/vX.Y.Z.md, sets BA_APPROVED+PM_APPROVED+RELEASE_LOCKED, Taskfile not Makefile.
After lock, kill product workers if idle.

PHASE 1: Add eng workers with exclusive paths and full AC briefs from the locked doc:
@deploy deploy/ + Taskfile.yml
@proto proto/
@trip services/trip-core/
@gateway services/gateway-auth/
@geo services/geo-places/
@web apps/web/
@verifier tests/AC/privacy only
@qa manual browser tester (skill manual-browser-qa) — click real UI; no code edits

Order: deploy+proto → parallel trip/gateway/geo → web → verifier loop → @qa browser pass (@ authors on failures) until MANUAL_QA_PASSED.

PHASE 2: When automated verifier + manual QA are green, release-ship to https://github.com/maks2134/Geo-Trip (push, tag vX.Y.Z, gh release create). Do not ship without browser QA.

Keep assigning on the board until invent→build→manual-qa→ship is complete. Milestone reports only.
```
