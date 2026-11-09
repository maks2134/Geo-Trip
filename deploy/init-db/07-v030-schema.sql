-- AtlasFlow v0.3.0 Database Schema Migration
-- Features: FEAT-03 (Multi-Currency Travel Ledger) & FEAT-06 (Scratch Milestones & UN Geoscheme)

-- ============================================================================
-- 1. Multi-Currency Travel Ledger Migration (FEAT-03 / CAND-01)
-- ============================================================================

ALTER TABLE expenses ADD COLUMN IF NOT EXISTS original_currency VARCHAR(3) DEFAULT 'USD';
ALTER TABLE expenses ADD COLUMN IF NOT EXISTS original_amount NUMERIC(12, 2);
ALTER TABLE expenses ADD COLUMN IF NOT EXISTS fx_rate NUMERIC(12, 6) DEFAULT 1.000000;
ALTER TABLE expenses ADD COLUMN IF NOT EXISTS base_amount NUMERIC(12, 2);
ALTER TABLE expenses ADD COLUMN IF NOT EXISTS is_custom_rate BOOLEAN DEFAULT false;

-- Backfill historical expense records so all financial aggregations remain accurate
UPDATE expenses
SET original_currency = currency,
    original_amount = amount,
    base_amount = amount,
    fx_rate = 1.000000
WHERE original_amount IS NULL;

-- ============================================================================
-- 2. UN Geoscheme Regional Taxonomy & Country Metadata (FEAT-06 / CAND-04)
-- ============================================================================

-- Add region and subregion columns to countries table
ALTER TABLE countries ADD COLUMN IF NOT EXISTS region TEXT NOT NULL DEFAULT '';
ALTER TABLE countries ADD COLUMN IF NOT EXISTS subregion TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_countries_region ON countries (region);
CREATE INDEX IF NOT EXISTS idx_countries_subregion ON countries (subregion);

-- Dedicated UN Geoscheme lookup table
CREATE TABLE IF NOT EXISTS regions_un (
    country_iso2 CHAR(2) PRIMARY KEY,
    country_name TEXT NOT NULL,
    region TEXT NOT NULL,
    subregion TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_regions_un_region ON regions_un (region);
CREATE INDEX IF NOT EXISTS idx_regions_un_subregion ON regions_un (subregion);

-- Populate UN Geoscheme country taxonomy
INSERT INTO regions_un (country_iso2, country_name, region, subregion) VALUES
    -- Europe
    ('FR', 'France', 'Europe', 'Western Europe'),
    ('DE', 'Germany', 'Europe', 'Western Europe'),
    ('NL', 'Netherlands', 'Europe', 'Western Europe'),
    ('BE', 'Belgium', 'Europe', 'Western Europe'),
    ('AT', 'Austria', 'Europe', 'Western Europe'),
    ('CH', 'Switzerland', 'Europe', 'Western Europe'),
    ('LU', 'Luxembourg', 'Europe', 'Western Europe'),
    ('GB', 'United Kingdom', 'Europe', 'Northern Europe'),
    ('IE', 'Ireland', 'Europe', 'Northern Europe'),
    ('SE', 'Sweden', 'Europe', 'Northern Europe'),
    ('NO', 'Norway', 'Europe', 'Northern Europe'),
    ('DK', 'Denmark', 'Europe', 'Northern Europe'),
    ('FI', 'Finland', 'Europe', 'Northern Europe'),
    ('IS', 'Iceland', 'Europe', 'Northern Europe'),
    ('EE', 'Estonia', 'Europe', 'Northern Europe'),
    ('LV', 'Latvia', 'Europe', 'Northern Europe'),
    ('LT', 'Lithuania', 'Europe', 'Northern Europe'),
    ('IT', 'Italy', 'Europe', 'Southern Europe'),
    ('ES', 'Spain', 'Europe', 'Southern Europe'),
    ('PT', 'Portugal', 'Europe', 'Southern Europe'),
    ('GR', 'Greece', 'Europe', 'Southern Europe'),
    ('HR', 'Croatia', 'Europe', 'Southern Europe'),
    ('SI', 'Slovenia', 'Europe', 'Southern Europe'),
    ('RS', 'Serbia', 'Europe', 'Southern Europe'),
    ('ME', 'Montenegro', 'Europe', 'Southern Europe'),
    ('BA', 'Bosnia and Herzegovina', 'Europe', 'Southern Europe'),
    ('MK', 'North Macedonia', 'Europe', 'Southern Europe'),
    ('AL', 'Albania', 'Europe', 'Southern Europe'),
    ('CY', 'Cyprus', 'Europe', 'Southern Europe'),
    ('MT', 'Malta', 'Europe', 'Southern Europe'),
    ('PL', 'Poland', 'Europe', 'Eastern Europe'),
    ('CZ', 'Czechia', 'Europe', 'Eastern Europe'),
    ('SK', 'Slovakia', 'Europe', 'Eastern Europe'),
    ('HU', 'Hungary', 'Europe', 'Eastern Europe'),
    ('RO', 'Romania', 'Europe', 'Eastern Europe'),
    ('BG', 'Bulgaria', 'Europe', 'Eastern Europe'),
    ('UA', 'Ukraine', 'Europe', 'Eastern Europe'),
    ('BY', 'Belarus', 'Europe', 'Eastern Europe'),
    ('MD', 'Moldova', 'Europe', 'Eastern Europe'),
    ('RU', 'Russian Federation', 'Europe', 'Eastern Europe'),

    -- Asia
    ('JP', 'Japan', 'Asia', 'Eastern Asia'),
    ('KR', 'South Korea', 'Asia', 'Eastern Asia'),
    ('CN', 'China', 'Asia', 'Eastern Asia'),
    ('TW', 'Taiwan', 'Asia', 'Eastern Asia'),
    ('MN', 'Mongolia', 'Asia', 'Eastern Asia'),
    ('HK', 'Hong Kong', 'Asia', 'Eastern Asia'),
    ('TH', 'Thailand', 'Asia', 'South-Eastern Asia'),
    ('VN', 'Vietnam', 'Asia', 'South-Eastern Asia'),
    ('ID', 'Indonesia', 'Asia', 'South-Eastern Asia'),
    ('MY', 'Malaysia', 'Asia', 'South-Eastern Asia'),
    ('SG', 'Singapore', 'Asia', 'South-Eastern Asia'),
    ('PH', 'Philippines', 'Asia', 'South-Eastern Asia'),
    ('KH', 'Cambodia', 'Asia', 'South-Eastern Asia'),
    ('LA', 'Laos', 'Asia', 'South-Eastern Asia'),
    ('MM', 'Myanmar', 'Asia', 'South-Eastern Asia'),
    ('IN', 'India', 'Asia', 'Southern Asia'),
    ('NP', 'Nepal', 'Asia', 'Southern Asia'),
    ('LK', 'Sri Lanka', 'Asia', 'Southern Asia'),
    ('MV', 'Maldives', 'Asia', 'Southern Asia'),
    ('TR', 'Turkey', 'Asia', 'Western Asia'),
    ('GE', 'Georgia', 'Asia', 'Western Asia'),
    ('AM', 'Armenia', 'Asia', 'Western Asia'),
    ('AZ', 'Azerbaijan', 'Asia', 'Western Asia'),
    ('AE', 'United Arab Emirates', 'Asia', 'Western Asia'),
    ('IL', 'Israel', 'Asia', 'Western Asia'),
    ('JO', 'Jordan', 'Asia', 'Western Asia'),
    ('SA', 'Saudi Arabia', 'Asia', 'Western Asia'),
    ('QA', 'Qatar', 'Asia', 'Western Asia'),

    -- Americas
    ('US', 'United States of America', 'Americas', 'Northern America'),
    ('CA', 'Canada', 'Americas', 'Northern America'),
    ('MX', 'Mexico', 'Americas', 'Central America'),
    ('CR', 'Costa Rica', 'Americas', 'Central America'),
    ('PA', 'Panama', 'Americas', 'Central America'),
    ('BR', 'Brazil', 'Americas', 'South America'),
    ('AR', 'Argentina', 'Americas', 'South America'),
    ('CL', 'Chile', 'Americas', 'South America'),
    ('CO', 'Colombia', 'Americas', 'South America'),
    ('PE', 'Peru', 'Americas', 'South America'),
    ('UY', 'Uruguay', 'Americas', 'South America'),

    -- Africa
    ('EG', 'Egypt', 'Africa', 'Northern Africa'),
    ('MA', 'Morocco', 'Africa', 'Northern Africa'),
    ('TN', 'Tunisia', 'Africa', 'Northern Africa'),
    ('ZA', 'South Africa', 'Africa', 'Sub-Saharan Africa'),
    ('KE', 'Kenya', 'Africa', 'Sub-Saharan Africa'),
    ('TZ', 'Tanzania', 'Africa', 'Sub-Saharan Africa'),
    ('MU', 'Mauritius', 'Africa', 'Sub-Saharan Africa'),
    ('SC', 'Seychelles', 'Africa', 'Sub-Saharan Africa'),

    -- Oceania
    ('AU', 'Australia', 'Oceania', 'Australia and New Zealand'),
    ('NZ', 'New Zealand', 'Oceania', 'Australia and New Zealand'),
    ('FJ', 'Fiji', 'Oceania', 'Melanesia')
ON CONFLICT (country_iso2) DO NOTHING;

-- Synchronize region taxonomy to existing countries table records
UPDATE countries c
SET region = r.region,
    subregion = r.subregion
FROM regions_un r
WHERE c.iso2 = r.country_iso2 AND (c.region = '' OR c.subregion = '');

-- ============================================================================
-- 3. Scratch Milestones Table (FEAT-06 / CAND-04)
-- ============================================================================

CREATE TABLE IF NOT EXISTS user_milestones (
    user_id UUID NOT NULL,
    milestone_key VARCHAR(64) NOT NULL,
    unlocked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    PRIMARY KEY (user_id, milestone_key)
);

CREATE INDEX IF NOT EXISTS idx_user_milestones_unlocked_at ON user_milestones (unlocked_at);
