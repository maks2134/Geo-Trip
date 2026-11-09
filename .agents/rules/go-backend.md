---
trigger: glob
globs: "**/*.go, **/go.mod, **/*.sql, **/proto/**/*.proto"
description: Clean Architecture and Go conventions for AtlasFlow backend services.
---

# Go backend rules

1. Place new domain types in `internal/<bounded-context>/domain` (or `domain/`). Never put SQL tags or protobuf tags on domain entities.
2. Usecases depend on port interfaces defined next to the usecase (or in `domain`). Adapters implement ports.
3. Handlers/adapters: parse → call usecase → map DTO. No business branching beyond transport concerns.
4. Use `context.Context` as first arg on all I/O methods. Honor cancellation.
5. SQL: `$1` placeholders only. Geography columns only in `geo-places`.
6. Prefer table-driven tests with `testify`. Mock ports, not concrete Postgres, in usecase tests.
7. Follow skills: `go-clean-arch`, `trip-core`, `geo-postgis` when implementing related features.
