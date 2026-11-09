# Prompt: implement locked release `v0.2.0`

Before paste: new chat → **Argon / Claude thinking + High** → permissions **Auto**.  
Primary: `orchestrator`. Source of truth: `docs/releases/v0.2.0.md` (already `RELEASE_LOCKED`).

---

## Primary prompt (paste as-is)

```text
/teamwork

Implement AtlasFlow release v0.2.0 end-to-end in multi-agent mode. You are the engineering orchestrator (`orchestrator`).

Repository: /Users/makskozlov/pet
Remote target: https://github.com/maks2134/Geo-Trip

### Source of truth (do not invent scope)
Read and implement ONLY:
- docs/releases/v0.2.0.md
Confirmed lock in that file: BA_APPROVED=yes, PM_APPROVED=yes, status=RELEASE_LOCKED, version=v0.2.0.

Also follow:
- .agents/GEMINI.md
- .agents/AGENTS.md
- skills: atlasflow-mvp, go-clean-arch, trip-core, geo-postgis, web-map-ui, public-share, release-ship

Meet every Given/When/Then AC in the release doc. Do NOT implement Backlog (v0.3+) or Out of Scope items.

### In-scope summary (from locked doc)
MUST-HAVE:
- MH-A i18n ru+en (ICU plurals, zero hardcoded UI strings, share chrome localized / author content untouched)
- MH-B system/light/dark theme + MapLibre style sync without wiping GeoJSON layers
- MH-C MinIO presigned photo/PDF uploads (cover, place gallery, private vouchers; MIME/size limits; share privacy)
- MH-D Streets ↔ Satellite toggle + fallback
- OPS-1 migrate Makefile → Taskfile.yml (`task …` only)

SELECTED extras:
- CAND-10 Copy to My Trips (ClonePublicTrip; no budget/private attachments)
- CAND-01 day leg distance/time matrix
- CAND-06 packing checklist starter templates
- CAND-08 GPX / GeoJSON export

### Agents / exclusive ownership (from §8 Implementation Tracks)
1) deploy-infra → Taskfile.yml, delete Makefile usage, deploy/, docker-compose.yml (MinIO CORS, S3_PUBLIC_ENDPOINT, bucket atlasflow-media, attachment/cover migrations as specified)
2) proto-contracts → proto/ (attachment presign/confirm/delete + ClonePublicTrip RPCs + codegen)
3) go-trip → services/trip-core/ (AttachmentMeta, presign adapter, ClonePublicTrip, templates, public filter tests)
4) go-gateway → services/gateway-auth/ (presign BFF, POST /api/v1/share/{slug}/clone, Accept-Language, share sanitization)
5) go-geo → services/geo-places/ (locale geocode hint, distance helper)
6) web-map → apps/web/ (MH-A/B/C/D + CAND-10/01/06/08)
7) verifier → read-only: task test, task build, web build, Clean Arch check, hardcoded-string audit, /share privacy audit

Brief each worker with owned paths, must-not-touch, AC IDs (MH-*/CAND-*), and skills. Parallelize after deploy+proto foundations. Never two writers on the same paths.

### Engineering rules
- Clean Architecture Go; map-first Adaptive Glass UI
- Taskfile only — no new Makefile targets
- Honor risks in §7 (presign public vs internal MinIO host, style.load layer re-sync, satellite attribution, hydration, clone privacy)
- No secrets in git
- Integrity mode: development
- Work IN this repo path

### Execution order
1. Confirm release doc lock
2. Track 1 Taskfile + MinIO/deploy wiring
3. Track 2 proto
4. Parallel Tracks 3–5
5. Track 6 web
6. Track 7 verifier until green against ACs
7. release-ship for v0.2.0:
   - commit without secrets
   - push to GitHub (maks2134/Geo-Trip)
   - annotated tag v0.2.0
   - gh release create "v0.2.0" --title "v0.2.0 — Global Nomad Experience" --notes-file docs/releases/v0.2.0.md

### Acceptance
- All locked ACs pass
- task test / task build / web build green
- Share/clone never leak budget or private attachments
- GitHub has commits + tag v0.2.0 + Release

Ask only if a locked AC is ambiguous or GitHub auth/remote blocks ship. Otherwise start immediately from the release doc.
```

---

## Swarm fallback (no /teamwork)

```text
Run a native swarm. You are the coordinator — do not write feature code yourself.

Implement ONLY docs/releases/v0.2.0.md (BA_APPROVED + PM_APPROVED + RELEASE_LOCKED) in /Users/makskozlov/pet.
Follow .agents/GEMINI.md and skills atlasflow-mvp, go-clean-arch, trip-core, geo-postgis, web-map-ui, public-share, release-ship.

Workers (exclusive paths; Taskfile not Makefile):
@deploy Taskfile.yml + deploy/ + MinIO CORS
@proto proto/ attachment + ClonePublicTrip
@trip services/trip-core/
@gateway services/gateway-auth/
@geo services/geo-places/
@web apps/web/ (MH-A/B/C/D + CAND-10/01/06/08)
@verifier AC/tests only

Order: deploy+proto → parallel trip/gateway/geo → web → verifier → release-ship (push GitHub, tag v0.2.0, gh release create with docs/releases/v0.2.0.md).
Brief from AC IDs. No backlog/out-of-scope. Report milestones only.
```
