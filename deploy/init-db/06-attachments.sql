-- AtlasFlow v0.2.0: MH-C Media Vault & Attachments Migration
-- Idempotent schema definition ensuring all attachment metadata & trip cover fields exist.

-- 1. Ensure trips table has cover_url
ALTER TABLE trips ADD COLUMN IF NOT EXISTS cover_url TEXT NOT NULL DEFAULT '';

-- 2. Ensure attachments table exists
CREATE TABLE IF NOT EXISTS attachments (
    id VARCHAR(64) PRIMARY KEY,
    trip_id VARCHAR(64) NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    place_id VARCHAR(64),
    day_id VARCHAR(64),
    filename VARCHAR(255) NOT NULL DEFAULT '',
    content_type VARCHAR(128) NOT NULL DEFAULT '',
    file_size BIGINT NOT NULL DEFAULT 0,
    storage_key TEXT NOT NULL DEFAULT '',
    download_url TEXT NOT NULL DEFAULT '',
    visibility INT NOT NULL DEFAULT 1,
    is_ready BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Idempotently add any missing columns if attachments was created by an earlier script
ALTER TABLE attachments ADD COLUMN IF NOT EXISTS trip_id VARCHAR(64) REFERENCES trips(id) ON DELETE CASCADE;
ALTER TABLE attachments ADD COLUMN IF NOT EXISTS place_id VARCHAR(64);
ALTER TABLE attachments ADD COLUMN IF NOT EXISTS day_id VARCHAR(64);
ALTER TABLE attachments ADD COLUMN IF NOT EXISTS filename VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE attachments ADD COLUMN IF NOT EXISTS content_type VARCHAR(128) NOT NULL DEFAULT '';
ALTER TABLE attachments ADD COLUMN IF NOT EXISTS file_size BIGINT NOT NULL DEFAULT 0;
ALTER TABLE attachments ADD COLUMN IF NOT EXISTS storage_key TEXT NOT NULL DEFAULT '';
ALTER TABLE attachments ADD COLUMN IF NOT EXISTS download_url TEXT NOT NULL DEFAULT '';
ALTER TABLE attachments ADD COLUMN IF NOT EXISTS visibility INT NOT NULL DEFAULT 1;
ALTER TABLE attachments ADD COLUMN IF NOT EXISTS is_ready BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE attachments ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE attachments ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- 4. Create performance and lookup indexes
CREATE INDEX IF NOT EXISTS idx_attachments_trip ON attachments (trip_id);
CREATE INDEX IF NOT EXISTS idx_attachments_place ON attachments (place_id) WHERE place_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_attachments_day ON attachments (day_id) WHERE day_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_attachments_visibility ON attachments (trip_id, visibility, is_ready);
