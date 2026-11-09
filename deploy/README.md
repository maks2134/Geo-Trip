# AtlasFlow / NomadGlass — Local Infrastructure

This directory contains the Docker Compose environment for local AtlasFlow development:
- **PostgreSQL 17** with **PostGIS 3.5** spatial engine.
- **MinIO** object storage (S3 compatible) with automatic bucket bootstrapping.
- **MinIO MC** init container to configure the `atlasflow-media` bucket and public policy.
- **Redis 7** for ephemeral caching (bbox queries, session states).

---

## Quickstart

### 1. Copy Environment Configuration
```bash
cp deploy/.env.example deploy/.env
```

### 2. Start Services
From within `deploy/`:
```bash
docker compose up -d
```
Or from the project root:
```bash
docker compose -f deploy/docker-compose.yml --env-file deploy/.env up -d
```

### 3. Check Service Health
```bash
docker compose -f deploy/docker-compose.yml ps
```
All services (`postgres`, `minio`, `redis`) should transition to `healthy`. `minio-init` will run once to create buckets and exit with status 0.
