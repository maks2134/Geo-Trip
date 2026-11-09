-- AtlasFlow Initial Database Setup
-- Executed automatically on first container run by PostgreSQL entrypoint

-- 1. Spatial database extension (PostGIS)
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;

-- 2. UUID & cryptographic generation extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 3. Schema permissions
GRANT ALL ON SCHEMA public TO CURRENT_USER;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO CURRENT_USER;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO CURRENT_USER;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO CURRENT_USER;

-- 4. Verification notice in server log
DO $$
BEGIN
    RAISE NOTICE 'AtlasFlow PostgreSQL initialized with PostGIS: %', PostGIS_Full_Version();
END $$;
