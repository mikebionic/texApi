-- Main GPS location tracking table
CREATE TABLE tbl_gps_location
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID                  DEFAULT gen_random_uuid(),
    vehicle_id      INT          NOT NULL REFERENCES tbl_vehicle (id) ON DELETE CASCADE,
    driver_id       INT          NOT NULL REFERENCES tbl_driver (id) ON DELETE CASCADE,
    latitude        DECIMAL(10,7) NOT NULL,
    longitude       DECIMAL(10,7) NOT NULL,
    altitude        DECIMAL(10,2),           -- Optional altitude data in meters
    speed           DECIMAL(5,2),            -- Speed in km/h
    direction       DECIMAL(5,2),            -- Direction in degrees (0-359)
    accuracy        DECIMAL(7,2),            -- Accuracy of location in meters
    location_time   TIMESTAMP    NOT NULL,   -- When the location was recorded by device
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP, -- When record was added to DB
    meta            JSONB        NOT NULL DEFAULT '{}'::jsonb -- For additional data (battery level, etc.)
);

-- Create index for faster location queries
CREATE INDEX idx_gps_location_vehicle_time ON tbl_gps_location (vehicle_id, location_time);
CREATE INDEX idx_gps_location_driver_time ON tbl_gps_location (driver_id, location_time);
CREATE INDEX idx_gps_location_time ON tbl_gps_location (location_time);

-- Create spatial index for geo queries if using PostGIS extension
-- If you plan to use PostGIS (recommended), run:
-- CREATE EXTENSION postgis;
-- Then add:
-- SELECT AddGeometryColumn('tbl_gps_location', 'position', 4326, 'POINT', 2);
-- CREATE INDEX idx_gps_location_position ON tbl_gps_location USING GIST(position);

-- Device table for tracking which device is sending location data
CREATE TABLE tbl_gps_device
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID                  DEFAULT gen_random_uuid(),
    vehicle_id      INT REFERENCES tbl_vehicle (id) ON DELETE SET NULL,
    driver_id       INT REFERENCES tbl_driver (id) ON DELETE SET NULL,
    device_id       VARCHAR(100) NOT NULL, -- Device identifier (IMEI, etc.)
    device_type     VARCHAR(50)  NOT NULL, -- Mobile, dedicated GPS tracker, etc.
    last_ping       TIMESTAMP,             -- Last time device communicated with server
    is_active       BOOLEAN      NOT NULL DEFAULT true,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Trip tracking table for logical grouping of GPS points
CREATE TABLE tbl_gps_trip
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID                  DEFAULT gen_random_uuid(),
    vehicle_id      INT NOT NULL REFERENCES tbl_vehicle (id) ON DELETE CASCADE,
    driver_id       INT NOT NULL REFERENCES tbl_driver (id) ON DELETE CASCADE,
    start_time      TIMESTAMP NOT NULL,
    end_time        TIMESTAMP,            -- NULL if trip is ongoing
    start_location  JSONB,                -- Starting location details
    end_location    JSONB,                -- Ending location details
    distance        DECIMAL(10,2),        -- Total distance in kilometers
    avg_speed       DECIMAL(5,2),         -- Average speed in km/h
    max_speed       DECIMAL(5,2),         -- Maximum speed in km/h
    status          VARCHAR(20) NOT NULL DEFAULT 'active', -- active, completed, cancelled
    meta            JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table for geofencing - define areas for alerts/reporting
CREATE TABLE tbl_geofence
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID                  DEFAULT gen_random_uuid(),
    company_id      INT NOT NULL REFERENCES tbl_company (id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    fence_type      VARCHAR(20) NOT NULL, -- circle, polygon, rectangle
    coordinates     JSONB NOT NULL,       -- Format depends on fence_type
    radius          DECIMAL(10,2),        -- For circle type, in meters
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Geofence alerts table
CREATE TABLE tbl_geofence_event
(
    id              SERIAL PRIMARY KEY,
    geofence_id     INT NOT NULL REFERENCES tbl_geofence (id) ON DELETE CASCADE,
    vehicle_id      INT NOT NULL REFERENCES tbl_vehicle (id) ON DELETE CASCADE,
    driver_id       INT NOT NULL REFERENCES tbl_driver (id) ON DELETE CASCADE,
    event_type      VARCHAR(20) NOT NULL, -- enter, exit
    event_time      TIMESTAMP NOT NULL,
    location        JSONB NOT NULL,       -- Location at time of event
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);