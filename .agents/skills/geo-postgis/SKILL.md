---
name: geo-postgis
description: >-
  Builds geo-places service features with PostGIS — points, bbox queries,
  geocoding ports, visited country polygons, and traveler stats. Use when
  working on map pins, spatial SQL, country unlock stats, or geocoding.
---

# Geo / PostGIS (geo-places)

## Schema essentials

```sql
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE places (
  id           uuid PRIMARY KEY,
  owner_id     text NOT NULL,
  name         text NOT NULL,
  category     text NOT NULL, -- food|stay|sight|nature|transport
  status       text NOT NULL, -- visited|planned|wishlist
  country_iso2 char(2),
  city         text,
  geom         geography(Point, 4326) NOT NULL,
  created_at   timestamptz NOT NULL,
  updated_at   timestamptz NOT NULL
);

CREATE INDEX places_geom_gix ON places USING GIST (geom);
CREATE INDEX places_owner_status_idx ON places (owner_id, status);
```

Country polygons: store Natural Earth / GADM polygons in `countries` (`iso2`, `geom geography(MultiPolygon,4326)`). Do not hand-draw MVP polygons in app code.

## Bbox query (viewport)

```sql
SELECT id, name, category, status, country_iso2, city,
       ST_Y(geom::geometry) AS lat, ST_X(geom::geometry) AS lng
FROM places
WHERE owner_id = $1
  AND ($2::text IS NULL OR status = $2)
  AND ST_Intersects(
        geom,
        ST_MakeEnvelope($3, $4, $5, $6, 4326)::geography
      )
LIMIT $7;
```

Args: west, south, east, north as lng/lat. Always clamp `LIMIT` (e.g. 500).

## Stats aggregation

```sql
-- visited countries count
SELECT COUNT(DISTINCT country_iso2)
FROM places
WHERE owner_id = $1 AND status = 'visited' AND country_iso2 IS NOT NULL;

-- optional: unlocked polygons for map fill
SELECT c.iso2, ST_AsGeoJSON(c.geom) AS geojson
FROM countries c
WHERE c.iso2 IN (
  SELECT DISTINCT country_iso2 FROM places
  WHERE owner_id = $1 AND status = 'visited' AND country_iso2 IS NOT NULL
);
```

Traveler widget fields:

- `countries_visited` / `195`
- `cities_visited`
- `explored_pct` = countries_visited / 195
- `points_total`
- `distance_km` — computed in `trip-core` from day route legs; geo-places may expose `ST_Distance` helper for two points

## Clean arch ports

```go
type PlaceRepository interface {
    Upsert(ctx context.Context, p Place) error
    ListInBBox(ctx context.Context, q BBoxQuery) ([]Place, error)
    Stats(ctx context.Context, ownerID string) (TravelerStats, error)
    VisitedCountryGeoJSON(ctx context.Context, ownerID string) ([]CountryFeature, error)
}

type Geocoder interface {
    Suggest(ctx context.Context, query string, proximity *LatLng) ([]GeocodeHit, error)
}
```

Geocoder adapter wraps Mapbox/Nominatim/etc. Domain must not import the HTTP client.

## Domain rules

- Coordinates validated: lat ∈ [-90,90], lng ∈ [-180,180]
- Creating a pin with `status=visited` may set `country_iso2` from geocode reverse or polygon contain (`ST_Covers`)
- Wishlist/planned pins count toward points but **not** opened-country fill

## Anti-patterns

- Storing lat/lng as float columns without geography (bbox and distance become wrong)
- Returning full world GeoJSON every request — filter by visited ISO codes only
- Calling external geocoder from the handler without timeouts/cancellation
