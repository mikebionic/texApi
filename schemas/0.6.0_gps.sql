CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS tbl_trip
(
    id            BIGSERIAL PRIMARY KEY,
    driver_id     INT       NOT NULL DEFAULT 0,
    vehicle_id    INT       NOT NULL DEFAULT 0,
    from_address  VARCHAR(800),
    to_address    VARCHAR(800),
    from_country  VARCHAR(500),
    to_country    VARCHAR(500),
    start_date    TIMESTAMP,
    end_date      TIMESTAMP,
    from_location GEOMETRY(POINT, 4326),
    to_location   GEOMETRY(POINT, 4326),
    distance_km   DECIMAL(10, 2), -- calculated total distance
    status        state_t   NOT NULL DEFAULT 'active',
    meta          TEXT      NOT NULL DEFAULT '',
    meta2         TEXT      NOT NULL DEFAULT '',
    meta3         TEXT      NOT NULL DEFAULT '',
    gps_logs      JSONB     NOT NULL DEFAULT '{}',
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted       INT       NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS tbl_gps_log
(
    id            BIGSERIAL PRIMARY KEY,
    company_id    INT REFERENCES tbl_company (id),
    vehicle_id    INT                         NOT NULL DEFAULT 0,
    driver_id     INT                         NOT NULL DEFAULT 0,
    offer_id      INT REFERENCES tbl_offer (id),        -- optional (to specify offer, use trip_id for combining them in one trip or filter by date)
    trip_id       INT REFERENCES tbl_trip (id),
    battery_level SMALLINT,                             -- device battery level (0-100)
    speed         DECIMAL(5, 2),                        -- Speed in km/h
    heading       DECIMAL(5, 2),                        -- degrees (0-359) North 0, 0° = North 90° = East 180° = South 270° = West
    accuracy      DECIMAL(7, 2),                        -- Accuracy of location in meters (optional)
    coordinates   GEOMETRY(POINT, 4326)       NOT NULL,
    status        state_t                     NOT NULL DEFAULT 'active',
    log_dt        TIMESTAMP WITHOUT TIME ZONE NOT NULL, -- When the location was recorded by device
    created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS tbl_offer_trip
(
    trip_id  INT     NOT NULL REFERENCES tbl_trip (id),
    offer_id INT     NOT NULL REFERENCES tbl_offer (id),
    is_main  BOOL    NOT NULL DEFAULT false,
    status   state_t NOT NULL DEFAULT 'active',
    deleted  INT     NOT NULL DEFAULT 0,
    UNIQUE (trip_id, offer_id, deleted)
);

CREATE INDEX idx_driver_log_dt ON tbl_gps_log(driver_id, log_dt);
CREATE INDEX idx_trip_log_dt ON tbl_gps_log(trip_id, log_dt);
CREATE INDEX idx_vehicle_log_dt ON tbl_gps_log(vehicle_id, log_dt);
CREATE INDEX idx_coordinates_gist ON tbl_gps_log USING GIST(coordinates);
CREATE INDEX idx_log_dt ON tbl_gps_log(log_dt);

CREATE INDEX idx_driver_dates ON tbl_trip(driver_id, start_date, end_date);
CREATE INDEX idx_vehicle_dates ON tbl_trip(vehicle_id, start_date, end_date);