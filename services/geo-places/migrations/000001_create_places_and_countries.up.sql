CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS places (
    id           UUID PRIMARY KEY,
    owner_id     TEXT NOT NULL,
    name         TEXT NOT NULL,
    category     TEXT NOT NULL, -- food|stay|sight|nature|transport
    status       TEXT NOT NULL, -- visited|planned|wishlist
    country_iso2 CHAR(2),
    city         TEXT,
    address      TEXT,
    geom         geography(Point, 4326) NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS places_geom_gix ON places USING GIST (geom);
CREATE INDEX IF NOT EXISTS places_owner_status_idx ON places (owner_id, status);
CREATE INDEX IF NOT EXISTS places_owner_idx ON places (owner_id);

CREATE TABLE IF NOT EXISTS visited_countries (
    user_id      TEXT NOT NULL,
    country_iso2 CHAR(2) NOT NULL,
    visited_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, country_iso2)
);

CREATE TABLE IF NOT EXISTS countries (
    iso2 CHAR(2) PRIMARY KEY,
    name TEXT NOT NULL,
    geom geography(MultiPolygon, 4326) NOT NULL
);

CREATE INDEX IF NOT EXISTS countries_geom_gix ON countries USING GIST (geom);