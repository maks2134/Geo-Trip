---
name: release-ship
description: >-
  Closes an AtlasFlow release after BA+PM approval, then ships to GitHub with
  semver tag and GitHub Release. Use when locking a feature set, versioning,
  tagging, pushing, or publishing a release.
---

# Release lock & GitHub ship

## Who decides

| Step | Owner |
| --- | --- |
| Requirements + AC | `business-analyst` |
| Priority + release cut + version | `product-manager` |
| Implementation | orchestrator + eng agents |
| Verify | `verifier` |
| GitHub push / tag / Release | PM-directed ship (eng executes) |

Release is locked only when the release doc has:

- `BA_APPROVED: yes`
- `PM_APPROVED: yes`
- status `RELEASE_LOCKED`

## Versioning (semver)

- `MAJOR` — breaking API / data migrations users must handle
- `MINOR` — new features, backward compatible (typical feature release)
- `PATCH` — fixes / polish only

Tag format: `vX.Y.Z` (example `v0.2.0`).

## Artifact

Create or update:

`docs/releases/vX.Y.Z.md`

Must include: theme, in/out scope, feature list with AC, approvals, changelog bullets.

## Ship sequence (after code is green)

Run only when user asked to ship **or** PM explicitly ordered ship after lock+verify:

```bash
git status
git checkout -B release/vX.Y.Z   # if not already on release branch
# stage relevant files — never .env secrets
git add -A
git commit -m "$(cat <<'EOF'
release: vX.Y.Z — <theme>

EOF
)"
git push -u origin HEAD
git tag -a "vX.Y.Z" -m "AtlasFlow vX.Y.Z — <theme>"
git push origin "vX.Y.Z"
gh release create "vX.Y.Z" --title "vX.Y.Z — <theme>" --notes-file docs/releases/vX.Y.Z.md
```

If remote/repo missing: create via `gh repo create` only when user/PM approved hosting.

## Never

- Ship without BA+PM lock
- Force-push main/master
- Commit secrets (`.env`, keys, MinIO creds)
- Skip verifier when features touched code
