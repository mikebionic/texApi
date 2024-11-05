-- Мои заявки
CREATE TABLE tbl_my_request
(
    id             SERIAL PRIMARY KEY,
    uuid           UUID                                                          DEFAULT gen_random_uuid(),
    company_id     INT            REFERENCES tbl_company (id) ON DELETE SET NULL DEFAULT 0,
    driver_id      INT            REFERENCES tbl_driver (id) ON DELETE SET NULL  DEFAULT 0,
    transport_id   INT            REFERENCES tbl_vehicle (id) ON DELETE SET NULL DEFAULT 0,
    cost_per_km    DECIMAL(10, 2) NOT NULL                                       DEFAULT 0.0,
    from_country   VARCHAR(100)   NOT NULL                                       DEFAULT '',
    from_region    VARCHAR(100)   NOT NULL                                       DEFAULT '',
    to_country     VARCHAR(100)   NOT NULL                                       DEFAULT '',
    to_region      VARCHAR(100)   NOT NULL                                       DEFAULT '',
    validity_start DATE           NOT NULL                                       DEFAULT CURRENT_TIMESTAMP,
    validity_end   DATE           NOT NULL                                       DEFAULT CURRENT_TIMESTAMP,
    note           TEXT,
    created_at     TIMESTAMP                                                     DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP                                                     DEFAULT CURRENT_TIMESTAMP,
    deleted        INT                                                           DEFAULT 0
);

-- Мои отклики
CREATE TYPE response_state_t AS ENUM ('accepted', 'declined');
CREATE TABLE tbl_response
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID                                                         DEFAULT gen_random_uuid(),
    user_id    INT              REFERENCES tbl_user (id) ON DELETE SET NULL DEFAULT 0,
    created_at TIMESTAMP                                                    DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP                                                    DEFAULT CURRENT_TIMESTAMP,
    state      response_state_t NOT NULL                                    DEFAULT 0,
    deleted    INT
);

-- Мои ставки
CREATE TABLE tbl_my_request
(
    id             SERIAL PRIMARY KEY,
    uuid           UUID                                                          DEFAULT gen_random_uuid(),
    company_id     INT            REFERENCES tbl_company (id) ON DELETE SET NULL DEFAULT 0,
    driver_id      INT            REFERENCES tbl_driver (id) ON DELETE SET NULL  DEFAULT 0,
    transport_id   INT            REFERENCES tbl_vehicle (id) ON DELETE SET NULL DEFAULT 0,
    cost_per_km    DECIMAL(10, 2) NOT NULL                                       DEFAULT 0.0,
    from_country   VARCHAR(100)   NOT NULL                                       DEFAULT '',
    from_region    VARCHAR(100)   NOT NULL                                       DEFAULT '',
    to_country     VARCHAR(100)   NOT NULL                                       DEFAULT '',
    to_region      VARCHAR(100)   NOT NULL                                       DEFAULT '',
    validity_start DATE           NOT NULL                                       DEFAULT CURRENT_TIMESTAMP,
    validity_end   DATE           NOT NULL                                       DEFAULT CURRENT_TIMESTAMP,
    note           TEXT,
    created_at     TIMESTAMP                                                     DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP                                                     DEFAULT CURRENT_TIMESTAMP,
    deleted        INT                                                           DEFAULT 0
);
