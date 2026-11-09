# Copy-paste into Antigravity

Use in a **new chat**. Prefer `/teamwork` (paid). If teamwork unavailable, use the Swarm fallback.

## Model setup (before paste)

1. In the chat model picker, set the **parent/orchestrator** to a strong model:
   - Best: **Gemini 4 Argon** (Pro/Ultra) or **Claude Opus/Sonnet (thinking)** with effort **High**
   - Budget: **Gemini 3.8 Flash** + effort Medium (weaker planning)
2. Custom agents already pin tiers in frontmatter (`pro` / `flash`) — the parent should invoke those agents so routing is automatic by role.
3. For Swarm, add this line to your prompt if you want mixed models:
   `Pin worker models: deploy+proto=flash; gateway/trip/geo/web/verifier=your strong model; switch on quota exhaustion only.`


---

## Primary: /teamwork

```text
/teamwork

Build the full AtlasFlow / NomadGlass MVP in this workspace: /Users/makskozlov/pet

You are the orchestrator. Work in multi-agent mode end-to-end without waiting for step-by-step instructions after I approve the plan. Spawn and delegate to the project custom agents (and/or equivalent Teamwork workers) with strict non-overlapping ownership:

- proto-contracts → proto/ (+ codegen)
- go-gateway → services/gateway-auth/
- go-trip → services/trip-core/
- go-geo → services/geo-places/
- web-map → apps/web/
- deploy-infra → deploy/
- verifier → review + run tests only (no feature edits)

Always load and follow:
- .agents/GEMINI.md
- .agents/AGENTS.md
- skills: atlasflow-mvp, go-clean-arch, trip-core, geo-postgis, web-map-ui, public-share

### Product (MVP only)
Hybrid personal world scratch-map + day-by-day trip planner:
1) World map + traveler stats + pin filters (visited/planned/wishlist + categories)
2) Trip planner: backlog, day timeline, lexorank DnD, numbered route polyline, day flyTo
3) Place drawer: geocode/click pin, metadata, notes, links, PDF/image attachments via presigned S3/MinIO
4) Budget tracker (limit, categories, by day)
5) Checklists (pre_flight / packing / buy_on_site) + readiness %
6) Public read-only /share/<slug> (hide budget + private docs)

### Stack
- Go 1.26 Clean Architecture services + ConnectRPC/protobuf in proto/
- apps/web: Next.js + TypeScript + Tailwind + MapLibre + Zustand + TanStack Query
- Adaptive Glass UI, canvas-first map (map = 100% viewport)
- deploy/: docker-compose with PostgreSQL+PostGIS, MinIO, optional Redis
- go.work monorepo; shared pkg/ for cross-cutting infra only (no business types)

### Execution order
1. Scaffold monorepo + deploy infra
2. Proto contracts
3. Parallel: gateway-auth, trip-core, geo-places (Clean Arch)
4. apps/web core screens wired to APIs (mock only if a service is blocked; replace with real calls)
5. Public share path through gateway
6. Verifier: go test, builds, architecture boundary check, share privacy check

### Constraints
- Work IN this repo path (do not relocate to ~/teamwork_projects unless unavoidable; if you must, sync results back here)
- Integrity mode: development
- Exclusive file ownership per worker — never two writers on the same files
- No scope beyond MVP
- Prefer real runnable code over docs
- Task runner: Taskfile.yml (`task`), NOT Makefile — migrate any Make usage into Taskfile
- Acceptance: docker compose up deps; services boot; core APIs work; web shows world+trip modes; share link is read-only without private data

Start Phase 1 interview only for decisions that block coding (module path naming, map provider key env name). Default everything else sensibly and proceed after my single approval.
```

---

## Fallback: Swarm (no /teamwork)

```text
Run a native swarm. You are the coordinator — do not write feature code yourself.

Goal: implement AtlasFlow / NomadGlass MVP in /Users/makskozlov/pet per .agents/GEMINI.md and skills atlasflow-mvp, go-clean-arch, trip-core, geo-postgis, web-map-ui, public-share.

Add workers (and narrator/monitor if those skills exist) with exclusive paths and workspace share/branch for authors:

1) @deploy — own deploy/ ; PostGIS+MinIO compose
2) @proto — own proto/ ; MVP ConnectRPC contracts
3) @gateway — own services/gateway-auth/ ; JWT + /share/:slug
4) @trip — own services/trip-core/ ; trips/days/lexorank/budget/checklists/attachments
5) @geo — own services/geo-places/ ; PostGIS bbox+stats+geocode port
6) @web — own apps/web/ ; map-first glass UI + share page
7) @verifier — tests/review only; @ authors with failures

Order: deploy+proto first, then gateway/trip/geo in parallel, then web, then verifier loop until green.

Brief each worker with owned paths, must-not-touch, acceptance, and which skill to load. Keep assigning on the board until MVP acceptance is met. Report only milestone summaries to me.
```
