# Prompt: BA + PM — next features & release

Before paste: new chat → model **Argon / Claude thinking + High** → permissions **Auto**.  
In `/agents` you may select `product-manager` or `business-analyst` as primary; prefer starting as PM and let it invoke BA.

---

## Primary prompt (paste as-is)

```text
Work in multi-agent mode as Product Manager (`product-manager`) and Business Analyst (`business-analyst`).
PM is the product-phase orchestrator: invoke BA, decide scope / priority / version.
BA owns requirements, user stories, acceptance criteria, and risks. Do not write application code.

Repository: /Users/makskozlov/pet
Context: .agents/GEMINI.md, .agents/AGENTS.md, skills atlasflow-mvp + release-ship.
Product: AtlasFlow / NomadGlass — world scratch map + day-by-day trip planner (glass, map-first).

### Session goals
1) Design the feature set for the next release.
2) MUST-HAVE — include and fully specify:
   A) Multilingual UI (i18n; at least RU + EN; architecture that allows adding languages)
   B) System theme: auto from OS (`prefers-color-scheme`) + manual light/dark override; keep map style in sync
   C) MinIO photo uploads (trip cover, place photos) via presigned URLs; MIME/size limits; privacy rules
   D) Satellite map basemap (if MapLibre/Mapbox/tiles allow) — Streets ↔ Satellite toggle; document fallback if not feasible
3) Propose 6–12 additional features scored with ICE; pull some into this release, park the rest in backlog.
4) Jointly (BA + PM) close the release feature set: semver version, in/out scope, acceptance criteria.
5) After `BA_APPROVED` + `PM_APPROVED` + `RELEASE_LOCKED` — prepare eng handoff. When implementation + verifier are green — follow skill `release-ship`: push to GitHub and create the version (annotated tag + GitHub Release). In this session, if code is not ready yet, stop at the locked release doc + eng brief; run ship commands only when code is ready or I say "ship".

### Output format
Create/update:
docs/releases/vNEXT.md  (rename to vX.Y.Z after PM picks the version)

File structure:
- Theme / Goal
- Version (PM proposal)
- MUST-HAVE (A–D) — full stories + Given/When/Then from BA
- Selected additional features for this release
- Backlog (parked)
- Out of scope
- Dependencies / risks (map tile keys, MinIO buckets, i18n tooling)
- Implementation tracks (eng agents: web-map, trip-core, gateway, deploy-infra, …)
- Approvals: BA_APPROVED / PM_APPROVED / status
- Changelog draft
- Ship checklist (git push, tag vX.Y.Z, gh release create)

### Rules
- Do not expand into realtime multiplayer / social network.
- Do not break public-share privacy.
- Wire MinIO photos into the existing trip-core attachment / presign pattern.
- System theme default = system; persist user override.
- i18n: no hard-coded UI strings in new components; keys + ICU/plurals where needed.
- Satellite: explicitly state feasible / not feasible and why; list API keys under risks if required.
- Dev/ops task runner: use **Taskfile** (`Taskfile.yml` + [go-task](https://taskfile.dev)), **not Makefile**. Any release eng brief, scripts, or docs must reference `task <name>` commands; if a Makefile exists, plan migrating commands into Taskfile and do not add new Make targets.
- After the release feature set is closed, when code is ready it MUST be pushed to GitHub with a version — product rule, not optional.

Ask one short clarifying question only if you cannot choose a version or cut MUST-HAVEs without it. Otherwise use sensible defaults and produce the release draft immediately, then deepen AC with BA.
```

---

## Short follow-up after lock (when code is done)

```text
Release docs/releases/vX.Y.Z.md is already BA_APPROVED + PM_APPROVED + RELEASE_LOCKED, verifier is green.
Execute release-ship: commit (no secrets), push to GitHub, annotated tag vX.Y.Z, gh release create with notes from docs/releases/vX.Y.Z.md.
```
