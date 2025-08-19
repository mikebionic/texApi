CREATE TYPE entity_t AS ENUM ('individual', 'legal');
CREATE TYPE role_t AS ENUM ('system','admin','sender','carrier','driver','unknown');
CREATE TYPE plan_t AS ENUM ('start', 'standard', 'premium');
CREATE TYPE state_t AS ENUM ('active', 'enabled', 'disabled', 'deleted', 'pending', 'archived', 'working', 'completed');
CREATE TYPE visibility_t AS ENUM ('public', 'private', 'contacts');
CREATE TYPE notification_preference_t AS ENUM ('all', 'mentions', 'none');
CREATE TYPE sticker_type_t AS ENUM ('static', 'animated');
CREATE TYPE status_type_t AS ENUM ('pending', 'approved', 'declined');
CREATE TYPE status_t AS ENUM ('draft', 'active', 'archived', 'banned','deleted');


CREATE TABLE tbl_role (
    id SERIAL PRIMARY KEY,
    role  role_t NOT NULL DEFAULT 'unknown',
    uuid UUID DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT DEFAULT '',
    title VARCHAR(100) NOT NULL DEFAULT '',
    title_ru VARCHAR(100) NOT NULL DEFAULT '',
    subtitle VARCHAR(200) NOT NULL DEFAULT '',
    subtitle_ru VARCHAR(200) NOT NULL DEFAULT ''
);


CREATE TABLE tbl_user
(
    id            SERIAL PRIMARY KEY,
    uuid          UUID                  DEFAULT gen_random_uuid(),
    username      VARCHAR(200) NOT NULL DEFAULT '',
    password      VARCHAR(200) NOT NULL DEFAULT '',
    email         VARCHAR(100) NOT NULL DEFAULT '',
    phone         VARCHAR(100) NOT NULL DEFAULT '',
    role          role_t       NOT NULL DEFAULT 'unknown',
    role_id       INT REFERENCES tbl_role (id) ON DELETE CASCADE,
    company_id    INT          NOT NULL DEFAULT 0,
    driver_id     INT          NOT NULL DEFAULT 0,
    verified      INT          NOT NULL DEFAULT 0,
    meta          TEXT         NOT NULL DEFAULT '',
    meta2         TEXT         NOT NULL DEFAULT '',
    meta3         TEXT         NOT NULL DEFAULT '',
    refresh_token VARCHAR(500) NOT NULL DEFAULT '',
    otp_key       VARCHAR(20)  NOT NULL DEFAULT '',
    verify_time   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active        INT          NOT NULL DEFAULT 1,
    deleted       INT          NOT NULL DEFAULT 0
);


CREATE TABLE tbl_company
(
    id                       SERIAL PRIMARY KEY,
    uuid                     UUID                  DEFAULT gen_random_uuid(),
    user_id                  INT          NOT NULL REFERENCES tbl_user (id) ON DELETE CASCADE,
    role                     role_t        NOT NULL DEFAULT 'unknown',
    role_id                  INT          NOT NULL REFERENCES tbl_role (id) ON DELETE CASCADE,
    plan                     plan_t        NOT NULL DEFAULT 'standard',
    plan_active              INT           NOT NULL DEFAULT 0,
    company_name             VARCHAR(100) NOT NULL DEFAULT '',
    first_name               VARCHAR(100) NOT NULL DEFAULT '',
    last_name                VARCHAR(100) NOT NULL DEFAULT '',
    patronymic_name          VARCHAR(100) NOT NULL DEFAULT '',
    about                    VARCHAR(2000) NOT NULL DEFAULT '',
    phone                    VARCHAR(100) NOT NULL DEFAULT '',
    phone2                   VARCHAR(100) NOT NULL DEFAULT '',
    phone3                   VARCHAR(100) NOT NULL DEFAULT '',
    email                    VARCHAR(100) NOT NULL DEFAULT '',
    email2                   VARCHAR(100) NOT NULL DEFAULT '',
    email3                   VARCHAR(100) NOT NULL DEFAULT '',
    meta                     TEXT         NOT NULL DEFAULT '',
    meta2                    TEXT         NOT NULL DEFAULT '',
    meta3                    TEXT         NOT NULL DEFAULT '',
    address                  VARCHAR(200) NOT NULL DEFAULT '',
    country                  VARCHAR(200) NOT NULL DEFAULT '',
    country_id               INT          NOT NULL DEFAULT 0,
    city_id                  INT          NOT NULL DEFAULT 0,
    image_url                VARCHAR(200) NOT NULL DEFAULT '',
    verified                 INT           NOT NULL DEFAULT 0, -- VERIFIED BADGE
    entity                   entity_t     NOT NULL DEFAULT 'individual',
    featured                 INT          NOT NULL DEFAULT 0,
    rating                   INT          NOT NULL DEFAULT 0,
    partner                  INT          NOT NULL DEFAULT 0,
    successful_ops           INT          NOT NULL DEFAULT 0,
    view_count               INT          NOT NULL DEFAULT 0,

    self_destruct_duration   INT          NOT NULL DEFAULT 0,    -- duration in minutes
    passkey                  VARCHAR(100) NOT NULL DEFAULT '',
    blacklist                TEXT[]                DEFAULT '{}',
    login_devices            TEXT[]                DEFAULT '{}',
    show_avatar              visibility_t NOT NULL DEFAULT 'public',
    show_bio                 visibility_t NOT NULL DEFAULT 'public',
    show_last_seen           visibility_t NOT NULL DEFAULT 'public',
    show_phone_number        visibility_t NOT NULL DEFAULT 'public',
    receive_calls            visibility_t NOT NULL DEFAULT 'public',
    invite_group             visibility_t NOT NULL DEFAULT 'public',
    notifications_chat       INT                   DEFAULT 0,
    notifications_group      INT                   DEFAULT 0,
    notifications_story      INT                   DEFAULT 0,
    notifications_reactions  INT                   DEFAULT 0,

    avatar_exceptions        TEXT[]                DEFAULT '{}', -- ID of company
    bio_exceptions           TEXT[]                DEFAULT '{}',
    last_seen_exceptions     TEXT[]                DEFAULT '{}',
    phone_number_exceptions  TEXT[]                DEFAULT '{}',
    receive_calls_exceptions TEXT[]                DEFAULT '{}',
    invite_group_exceptions  TEXT[]                DEFAULT '{}',

    last_active              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at               TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active                   INT          NOT NULL DEFAULT 1,
    deleted                  INT          NOT NULL DEFAULT 0
        CONSTRAINT rating_range CHECK (rating >= 0 AND rating <= 5)
);

CREATE TABLE tbl_sessions (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL REFERENCES tbl_user(id) ON DELETE CASCADE,
  company_id INT NOT NULL DEFAULT 0,
  refresh_token VARCHAR(500) NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  device_name VARCHAR(200) NOT NULL DEFAULT '',
  device_model VARCHAR(200) NOT NULL DEFAULT '',
  device_firmware VARCHAR(200) NOT NULL DEFAULT '',
  app_name VARCHAR(100) NOT NULL DEFAULT '',
  app_version VARCHAR(50) NOT NULL DEFAULT '',
  user_agent TEXT NOT NULL DEFAULT '',
  ip_address VARCHAR(50) NOT NULL DEFAULT '',
  login_method VARCHAR(20) NOT NULL DEFAULT 'password', -- password, oauth, otp
  last_used_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_sessions_user_id ON tbl_sessions(user_id);
CREATE INDEX idx_sessions_refresh_token ON tbl_sessions(refresh_token);
CREATE INDEX idx_sessions_expires_at ON tbl_sessions(expires_at);

CREATE TABLE tbl_firebase_token (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES tbl_user(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    device_type TEXT, -- Optional: 'android', 'ios', 'web', etc.
    meta text,
    meta2 text,
    meta3 text,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	active INT NOT NULL DEFAULT 1,
	deleted INT NOT NULL DEFAULT 0
);


CREATE TABLE tbl_plan_moves(
   id SERIAL PRIMARY KEY,
   user_id    INT           NOT NULL DEFAULT 0,
   company_id INT           NOT NULL DEFAULT 0,
   status     status_type_t NOT NULL DEFAULT 'pending',
   valid_until TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '1 month'), -- month TODO: reconfigure in backend
   created_at TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
   deleted    INT           NOT NULL DEFAULT 0
);

CREATE TABLE tbl_verify_request
(
    id         SERIAL PRIMARY KEY,
    user_id    INT           NOT NULL DEFAULT 0,
    company_id INT           NOT NULL DEFAULT 0,
    status     status_type_t NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted    INT           NOT NULL DEFAULT 0
);

CREATE TABLE tbl_user_log
(
    id         SERIAL PRIMARY KEY,
    user_id    INT          NOT NULL DEFAULT 0,
    company_id INT          NOT NULL DEFAULT 0,
    role       role_t       NOT NULL DEFAULT 'unknown',
    action     VARCHAR(200) NOT NULL DEFAULT '',
    details     VARCHAR(500) NOT NULL DEFAULT '',
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted    INT          NOT NULL DEFAULT 0
);


CREATE TABLE tbl_driver
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID                  DEFAULT gen_random_uuid(),
    company_id      INT          NOT NULL REFERENCES tbl_company (id) ON DELETE CASCADE,
    first_name      VARCHAR(100) NOT NULL DEFAULT '',
    last_name       VARCHAR(100) NOT NULL DEFAULT '',
    patronymic_name VARCHAR(100) NOT NULL DEFAULT '',
    phone           VARCHAR(100) NOT NULL DEFAULT '',
    email           VARCHAR(100) NOT NULL DEFAULT '',
    featured        INT          NOT NULL DEFAULT 0,
    rating          INT          NOT NULL DEFAULT 0,
    partner         INT          NOT NULL DEFAULT 0,
    successful_ops  INT          NOT NULL DEFAULT 0,
    image_url       VARCHAR(200) NOT NULL DEFAULT '',
    meta            TEXT         NOT NULL DEFAULT '',
    meta2           TEXT         NOT NULL DEFAULT '',
    meta3           TEXT         NOT NULL DEFAULT '',
    available       INT          NOT NULL DEFAULT 1,
    view_count      INT          NOT NULL DEFAULT 0,
    block_reason    VARCHAR(500) NOT NULL DEFAULT '',
    created_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    active          INT          NOT NULL DEFAULT 1,
    deleted         INT          NOT NULL DEFAULT 0
);

CREATE TABLE tbl_version
(
    id                          SERIAL PRIMARY KEY,
    uuid                        UUID          NOT NULL DEFAULT gen_random_uuid(),
    version_number              VARCHAR(50)   NOT NULL DEFAULT '1.0.0', -- e.g., "1.2.3", "2.0.0-beta.1"
    version_code                INT           NOT NULL DEFAULT 0, -- incrementing integer for comparison
    title                       VARCHAR(200)  NOT NULL DEFAULT '',
    description                 VARCHAR(1000),
    platform                    VARCHAR(20)   NOT NULL DEFAULT 'unknown', -- 'ios', 'android', 'web', 'desktop'
    minimal_platform_version    VARCHAR(50),  -- minimum OS version required (e.g., "iOS 14.0", "Android 8.0")
    download_url                VARCHAR(500), -- URL for downloading the app
    file_size                   BIGINT,       -- in bytes
    checksum                    VARCHAR(128), -- integrity verification
    changelog                   TEXT,
    release_notes               TEXT,
    is_critical_update          BOOLEAN       NOT NULL DEFAULT FALSE, -- forces update
    is_beta                     BOOLEAN       NOT NULL DEFAULT FALSE, -- flag
    auto_update_enabled         BOOLEAN       NOT NULL DEFAULT TRUE,  -- allow auto-updates
    rollout_percentage          INT           NOT NULL DEFAULT 100,   -- gradual rollout (0-100)
    active_at                   TIMESTAMP     DEFAULT CURRENT_TIMESTAMP,
    deprecated_at               TIMESTAMP,    -- Extra, when this version becomes deprecated
    end_of_life_at              TIMESTAMP,    -- Extra, when this version is no longer supported
    created_at                  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active                      INT           NOT NULL DEFAULT 1,
    deleted                     INT           NOT NULL DEFAULT 0,

    -- Constraints
    CONSTRAINT unique_version_platform UNIQUE (version_number, platform, deleted),
    CONSTRAINT valid_rollout_percentage CHECK (rollout_percentage >= 0 AND rollout_percentage <= 100)
);

-- Indexes for better performance
CREATE INDEX idx_version_platform_active ON tbl_version(platform, active, deleted);
CREATE INDEX idx_version_active_at ON tbl_version(active_at);
CREATE INDEX idx_version_code ON tbl_version(version_code);
CREATE INDEX idx_version_uuid ON tbl_version(uuid);

CREATE TYPE plan_level_t AS ENUM ('none', 'basic', 'advanced', 'full');

CREATE TABLE tbl_plan
(
    id                   SERIAL PRIMARY KEY,
    uuid                 UUID           NOT NULL DEFAULT gen_random_uuid(),
    name            VARCHAR(100)   NOT NULL,                   -- e.g., "Roolz Basic", "DAT Enhanced"
    code            VARCHAR(50)    NOT NULL,                   -- e.g., "ROOLZ_BASIC", "DAT_ENHANCED"
    provider             VARCHAR(50)    NOT NULL,                   -- e.g., "Roolz", "DAT", "TEX"
    region               VARCHAR(50)    NOT NULL,                   -- e.g., "CIS", "USA", "Turkmenistan"

    price_usd            DECIMAL(10, 2) NOT NULL,
    price_local          DECIMAL(10, 2),                            -- local currency amount
    local_currency       VARCHAR(10),                               -- e.g., "TMT"
    billing_cycle        VARCHAR(20)    NOT NULL DEFAULT 'monthly', -- monthly, yearly, etc.

    load_posts_limit     INT,                                       -- NULL for unlimited, number for limited
    load_posts_unlimited BOOLEAN        NOT NULL DEFAULT FALSE,

    gps_tracking_level   plan_level_t    NOT NULL DEFAULT 'none',    -- 'none', 'basic', 'advanced', 'full'
    gps_has_eta          BOOLEAN        NOT NULL DEFAULT FALSE,

    rate_tools_level     plan_level_t    NOT NULL DEFAULT 'none',    -- 'none', 'basic', 'advanced'
    rate_tools_features  TEXT[],                                    -- array of features like ['LaneMakers', 'RateView']

    edocs_available      BOOLEAN        NOT NULL DEFAULT FALSE,
    edocs_limit          INT,                                       -- NULL for unlimited
    edocs_has_archiving  BOOLEAN        NOT NULL DEFAULT FALSE,

    support_level        VARCHAR(20)    NOT NULL DEFAULT 'none',    -- 'none', 'email', 'phone', 'priority', 'dedicated'

    payment_guarantee    BOOLEAN        NOT NULL DEFAULT FALSE,
    api_access           BOOLEAN        NOT NULL DEFAULT FALSE,

    display_order        INT            NOT NULL DEFAULT 0,
    is_popular           BOOLEAN        NOT NULL DEFAULT FALSE,
    is_recommended       BOOLEAN        NOT NULL DEFAULT FALSE,
    description          VARCHAR(1000),
    features_summary     TEXT,                                      -- human-readable feature summary

    available_from       TIMESTAMP               DEFAULT CURRENT_TIMESTAMP,
    available_until      TIMESTAMP,                                 -- for limited-time plans

    meta                 TEXT,
    meta2                TEXT,
    meta3                TEXT,

    created_at           TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active               INT            NOT NULL DEFAULT 1,
    deleted              INT            NOT NULL DEFAULT 0,

    CONSTRAINT unique_plan_code UNIQUE (code, deleted),
    CONSTRAINT unique_plan_provider_region UNIQUE (provider, name, region, deleted),
    CONSTRAINT valid_price CHECK (price_usd >= 0),
    CONSTRAINT valid_load_posts_limit CHECK (load_posts_limit IS NULL OR load_posts_limit >= 0),
    CONSTRAINT valid_edocs_limit CHECK (edocs_limit IS NULL OR edocs_limit >= 0)
);

CREATE INDEX idx_plan_provider ON tbl_plan(provider) WHERE deleted = 0;
CREATE INDEX idx_plan_region ON tbl_plan(region) WHERE deleted = 0;
CREATE INDEX idx_plan_active ON tbl_plan(active, deleted);
CREATE INDEX idx_plan_price ON tbl_plan(price_usd) WHERE deleted = 0;

INSERT INTO tbl_plan (name, code, provider, region, price_usd, price_local, local_currency,
                      load_posts_limit, load_posts_unlimited, gps_tracking_level, rate_tools_level,
                      edocs_available, support_level, payment_guarantee, api_access, display_order)
VALUES
    ('Roolz Basic', 'ROOLZ_BASIC', 'Roolz', 'CIS', 30.00, NULL, NULL, 50, FALSE, 'basic', 'basic', TRUE, 'email', FALSE, FALSE, 1),
    ('Roolz Pro', 'ROOLZ_PRO', 'Roolz', 'CIS', 100.00, NULL, NULL, NULL, TRUE, 'advanced', 'advanced', TRUE, 'priority', FALSE, FALSE, 2),

    ('DAT Standard', 'DAT_STANDARD', 'DAT', 'USA', 45.00, NULL, NULL, NULL, FALSE, 'basic', 'none', TRUE, 'email', FALSE, FALSE, 3),
    ('DAT Enhanced', 'DAT_ENHANCED', 'DAT', 'USA', 85.00, NULL, NULL, NULL, FALSE, 'advanced', 'basic', TRUE, 'phone', TRUE, FALSE, 4),
    ('DAT Select', 'DAT_SELECT', 'DAT', 'USA', 135.00, NULL, NULL, NULL, FALSE, 'advanced', 'advanced', FALSE, 'dedicated', TRUE, TRUE, 5),

    ('TEX Start', 'TEX_START', 'TEX', 'Turkmenistan', 10.00, 350.00, 'TMT', 30, FALSE, 'basic', 'none', TRUE, 'email', FALSE, FALSE, 6),
    ('TEX Pro', 'TEX_PRO', 'TEX', 'Turkmenistan', 30.00, 900.00, 'TMT', NULL, TRUE, 'full', 'basic', TRUE, 'priority', TRUE, FALSE, 7),
    ('TEX Enterprise', 'TEX_ENTERPRISE', 'TEX', 'Turkmenistan', 80.00, 1800.00, 'TMT', NULL, TRUE, 'advanced', 'advanced', TRUE, 'dedicated', TRUE, TRUE, 8);

CREATE VIEW v_plan_summary AS
SELECT
    id,
    uuid,
    name,
    provider,
    region,
    CASE
        WHEN price_local IS NOT NULL THEN CONCAT('$', price_usd, ' (~', price_local, ' ', local_currency, ')')
        ELSE CONCAT('$', price_usd)
        END as price_display,
    CASE
        WHEN load_posts_unlimited THEN 'Unlimited'
        WHEN load_posts_limit IS NOT NULL THEN CONCAT(load_posts_limit, '/month')
        ELSE 'Limited'
        END as load_posts_display,
    gps_tracking_level,
    rate_tools_level,
    edocs_available,
    support_level,
    payment_guarantee,
    api_access,
    active
FROM tbl_plan
WHERE deleted = 0
ORDER BY display_order;


CREATE TABLE tbl_organization
(
    id               SERIAL PRIMARY KEY,
    uuid             UUID         NOT NULL DEFAULT gen_random_uuid(),
    name             VARCHAR(200) NOT NULL,
    description_en   VARCHAR(200),
    description_ru   VARCHAR(200),
    description_tk   VARCHAR(200),
    email            VARCHAR(200),
    image_url        VARCHAR(500),
    logo_url         VARCHAR(500),
    icon_url         VARCHAR(500),
    banner_url       VARCHAR(500),
    website_url      VARCHAR(500),
    about_text       TEXT,
    refund_text      TEXT,
    delivery_text    TEXT,
    contact_text     TEXT,
    terms_conditions TEXT,
    privacy_policy   TEXT,
    address1         VARCHAR(200),
    address2         VARCHAR(200),
    address3         VARCHAR(200),
    address4         VARCHAR(200),
    address_title1   VARCHAR(200),
    address_title2   VARCHAR(200),
    address_title3   VARCHAR(200),
    address_title4   VARCHAR(200),
    contact_phone1   VARCHAR(30),
    contact_phone2   VARCHAR(30),
    contact_phone3   VARCHAR(30),
    contact_phone4   VARCHAR(30),
    contact_title1   VARCHAR(200),
    contact_title2   VARCHAR(200),
    contact_title3   VARCHAR(200),
    contact_title4   VARCHAR(200),
    meta             TEXT,
    meta2            TEXT,
    meta3            TEXT,
    created_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active           INT          NOT NULL DEFAULT 1,
    deleted          INT          NOT NULL DEFAULT 0
);

CREATE TABLE tbl_packaging_type
(
    id             SERIAL PRIMARY KEY,
    name_ru        VARCHAR(255)   NOT NULL DEFAULT '',
    name_en        VARCHAR(255)   NOT NULL DEFAULT '',
    name_tk        VARCHAR(255)   NOT NULL DEFAULT '',
    category_ru    VARCHAR(255)   NOT NULL DEFAULT '',
    category_en    VARCHAR(255)   NOT NULL DEFAULT '',
    category_tk    VARCHAR(255)   NOT NULL DEFAULT '',
    material       VARCHAR(255)   NOT NULL DEFAULT '',
    dimensions     VARCHAR(255)   NOT NULL DEFAULT '',
    weight         DECIMAL(10, 2) NOT NULL DEFAULT 0.0,
    description_ru TEXT           NOT NULL DEFAULT '',
    description_en TEXT           NOT NULL DEFAULT '',
    description_tk TEXT           NOT NULL DEFAULT '',
    active         INT            NOT NULL DEFAULT 0,
    deleted        INT            NOT NULL DEFAULT 0
);

CREATE TABLE tbl_vehicle
(
    id                  SERIAL PRIMARY KEY,
    uuid                UUID                                                    DEFAULT gen_random_uuid(),
    company_id          INT          NOT NULL REFERENCES tbl_company (id) ON DELETE CASCADE,
    vehicle_type_id     INT REFERENCES tbl_vehicle_type (id) ON DELETE CASCADE  DEFAULT 1,
    vehicle_brand_id    INT REFERENCES tbl_vehicle_brand (id) ON DELETE CASCADE DEFAULT 1,
    vehicle_model_id    INT REFERENCES tbl_vehicle_model (id) ON DELETE CASCADE DEFAULT 1,
    year_of_issue       VARCHAR(10)  NOT NULL                                   DEFAULT '',
    mileage             INT          NOT NULL                                   DEFAULT 0,
    numberplate         VARCHAR(20)  NOT NULL                                   DEFAULT '',
    trailer_numberplate VARCHAR(20)  NOT NULL                                   DEFAULT '',
    gps                 INT          NOT NULL                                   DEFAULT 0,
    photo1_url          VARCHAR(200) NOT NULL                                   DEFAULT '',
    photo2_url          VARCHAR(200) NOT NULL                                   DEFAULT '',
    photo3_url          VARCHAR(200) NOT NULL                                   DEFAULT '',
    docs1_url           VARCHAR(200) NOT NULL                                   DEFAULT '',
    docs2_url           VARCHAR(200) NOT NULL                                   DEFAULT '',
    docs3_url           VARCHAR(200) NOT NULL                                   DEFAULT '',
    view_count          INT          NOT NULL                                   DEFAULT 0,
    meta                TEXT         NOT NULL                                   DEFAULT '',
    meta2               TEXT         NOT NULL                                   DEFAULT '',
    meta3               TEXT         NOT NULL                                   DEFAULT '',
    available           INT          NOT NULL                                   DEFAULT 1,
    created_at          TIMESTAMP    NOT NULL                                   DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP    NOT NULL                                   DEFAULT CURRENT_TIMESTAMP,
    active              INT          NOT NULL                                   DEFAULT 1,
    deleted             INT          NOT NULL                                   DEFAULT 0
);


INSERT INTO tbl_role (role, name, description, title, subtitle, title_ru, subtitle_ru) VALUES
   ('system', 'system', 'System level access','','','',''),
   ('admin', 'admin', 'Has full access to manage the system','','','',''),
   ('sender', 'sender', 'Can place orders and track deliveries', 'Sender', 'I am looking for transport','Отправитель', 'Я ищу транспорт'),
   ('carrier', 'carrier_personal', 'Responsible for delivering orders using their personal vehicle', 'Carrier', 'Personal vehicle','Перевозчик','Личный автотранспорт'),
   ('carrier', 'carrier_owner', 'Responsible for delivering orders with a fleet of vehicles they own', 'Carrier', 'Fleet owner','Перевозчик','Владелец парка автотранспорта'),
   ('carrier', 'carrier_company', 'Responsible for delivering orders through a logistics company', 'Carrier', 'Logistics company','Перевозчик','Логистическая кампания');
