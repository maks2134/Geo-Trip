CREATE TABLE IF NOT EXISTS trips (
    id VARCHAR(64) PRIMARY KEY,
    owner_id VARCHAR(64) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    cover_url TEXT NOT NULL DEFAULT '',
    currency VARCHAR(8) NOT NULL DEFAULT 'USD',
    budget_limit NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    share_slug VARCHAR(64) NOT NULL DEFAULT '',
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_trips_owner ON trips (owner_id);
CREATE INDEX IF NOT EXISTS idx_trips_share_slug ON trips (share_slug) WHERE share_slug != '';

CREATE TABLE IF NOT EXISTS trip_days (
    id VARCHAR(64) PRIMARY KEY,
    trip_id VARCHAR(64) NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    day_index INT NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    title VARCHAR(255) NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_trip_days_trip ON trip_days (trip_id, day_index);

CREATE TABLE IF NOT EXISTS itinerary_items (
    id VARCHAR(64) PRIMARY KEY,
    trip_id VARCHAR(64) NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    day_id VARCHAR(64) REFERENCES trip_days(id) ON DELETE SET NULL,
    place_id VARCHAR(128) NOT NULL,
    place_name VARCHAR(255) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    category INT NOT NULL DEFAULT 0,
    address TEXT NOT NULL DEFAULT '',
    position VARCHAR(64) NOT NULL,
    duration_minutes INT NOT NULL DEFAULT 60,
    cost NUMERIC(12, 2),
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_itinerary_trip_day_pos ON itinerary_items (trip_id, day_id, position);

CREATE TABLE IF NOT EXISTS expenses (
    id VARCHAR(64) PRIMARY KEY,
    trip_id VARCHAR(64) NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    day_id VARCHAR(64) REFERENCES trip_days(id) ON DELETE SET NULL,
    place_item_id VARCHAR(64) REFERENCES itinerary_items(id) ON DELETE SET NULL,
    category INT NOT NULL DEFAULT 0,
    amount NUMERIC(12, 2) NOT NULL,
    currency VARCHAR(8) NOT NULL DEFAULT 'USD',
    title VARCHAR(255) NOT NULL,
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_expenses_trip ON expenses (trip_id);

CREATE TABLE IF NOT EXISTS checklist_items (
    id VARCHAR(64) PRIMARY KEY,
    trip_id VARCHAR(64) NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    section INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    is_done BOOLEAN NOT NULL DEFAULT FALSE,
    position VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_checklist_trip_pos ON checklist_items (trip_id, section, position);

CREATE TABLE IF NOT EXISTS attachments (
    id VARCHAR(64) PRIMARY KEY,
    trip_id VARCHAR(64) NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    place_id VARCHAR(64),
    day_id VARCHAR(64),
    filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(128) NOT NULL,
    file_size BIGINT NOT NULL DEFAULT 0,
    storage_key TEXT NOT NULL,
    download_url TEXT NOT NULL DEFAULT '',
    visibility INT NOT NULL DEFAULT 1,
    is_ready BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_attachments_trip ON attachments (trip_id);

CREATE TABLE IF NOT EXISTS trip_shares (
    id VARCHAR(64) PRIMARY KEY,
    trip_id VARCHAR(64) NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    slug VARCHAR(64) UNIQUE NOT NULL,
    hide_budget BOOLEAN NOT NULL DEFAULT TRUE,
    hide_private_docs BOOLEAN NOT NULL DEFAULT TRUE,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_trip_shares_slug ON trip_shares (slug);
