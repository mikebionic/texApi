-- Мои заявки
CREATE TABLE tbl_offer
(
    id             SERIAL PRIMARY KEY,
    uuid           UUID                                                          DEFAULT gen_random_uuid(),
    user_id        INT            REFERENCES tbl_user (id) ON DELETE SET NULL    DEFAULT 0,
    company_id     INT            REFERENCES tbl_company (id) ON DELETE SET NULL DEFAULT 0,
    driver_id      INT            REFERENCES tbl_driver (id) ON DELETE SET NULL  DEFAULT 0,
    vehicle_id   INT            REFERENCES tbl_vehicle (id) ON DELETE SET NULL DEFAULT 0,
    cost_per_km    DECIMAL(10, 2) NOT NULL                                       DEFAULT 0.0,
    from_country   VARCHAR(100)   NOT NULL                                       DEFAULT '',
    from_region    VARCHAR(100)   NOT NULL                                       DEFAULT '',
    to_country     VARCHAR(100)   NOT NULL                                       DEFAULT '',
    to_region      VARCHAR(100)   NOT NULL                                       DEFAULT '',
    view_count     INT            NOT NULL                                       DEFAULT 0,
    validity_start DATE           NOT NULL                                       DEFAULT CURRENT_TIMESTAMP,
    validity_end   DATE           NOT NULL                                       DEFAULT CURRENT_TIMESTAMP,
    note           TEXT,
    created_at     TIMESTAMP                                                     DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP                                                     DEFAULT CURRENT_TIMESTAMP,
    deleted        INT                                                           DEFAULT 0
);

-- Мои отклики
CREATE TYPE response_state_t AS ENUM ('pending', 'accepted', 'declined');
CREATE TABLE tbl_response
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID                                                         DEFAULT gen_random_uuid(),
    company_id     INT            REFERENCES tbl_company (id) ON DELETE SET NULL DEFAULT 0,
    bid_user_id    INT              REFERENCES tbl_user (id) ON DELETE SET NULL DEFAULT 0,
    state      response_state_t NOT NULL                                    DEFAULT 'pending',
    created_at TIMESTAMP                                                    DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP                                                    DEFAULT CURRENT_TIMESTAMP,
    deleted    INT
);



INSERT INTO tbl_offer (user_id, company_id, driver_id, vehicle_id, cost_per_km, from_country, from_region, to_country, to_region, view_count, validity_start, validity_end, note)
VALUES (1, 1, 1, 1, 400.25, 'Germany', 'Berlin', 'Italy', 'Rome', 10, '2024-11-01', '2024-11-10', 'Urgent transport needed from Berlin to Rome.'),
    (2, 2, 2, 2, 350.50, 'France', 'Paris', 'Spain', 'Barcelona', 15, '2024-11-05', '2024-11-12', 'Looking for reliable transport for goods from Paris to Barcelona.'),
    (3, 3, 3, 3, 299.75, 'UK', 'London', 'Netherlands', 'Amsterdam', 20, '2024-11-07', '2024-11-14', 'Request for timely delivery from London to Amsterdam.'),
    (1, 1, 4, 4, 100.00, 'Poland', 'Warsaw', 'Hungary', 'Budapest', 5, '2024-11-10', '2024-11-17', 'Need cargo transport from Warsaw to Budapest.'),
    (2, 2, 5, 5, 99.99, 'Belgium', 'Brussels', 'Austria', 'Vienna', 8, '2024-11-12', '2024-11-20', 'Looking for a driver to transport goods from Brussels to Vienna.');

INSERT INTO tbl_response (company_id, bid_user_id, state)
VALUES    (2, 2, 'declined'),
    (3, 3, 'accepted'),
    (1, 1, 'pending'),
    (2, 2, 'declined'),
    (3, 3, 'declined'),
    (1, 1, 'accepted'),
    (2, 2, 'declined'),
    (1, 1, 'pending'),
    (2, 2, 'pending'),
    (1, 1, 'accepted'),
    (2, 2, 'declined'),
    (1, 1, 'declined'),
    (2, 2, 'declined'),
    (3, 3, 'pending'),
    (2, 2, 'accepted'),
    (1, 1, 'accepted'),
    (2, 2, 'pending'),
    (3, 3, 'declined'),
    (1, 1, 'declined');
