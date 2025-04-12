
CREATE TYPE media_type AS ENUM (
    'image',
    'video',
    'audio',
    'application',
    'text',
    'document',
    'voice',
    'unknown'
    );

CREATE TYPE media_context AS ENUM (
    'post',
    'story',
    'voice',
    'message',
    'document',
    'profile_image',
    'profile_background',
    'chat_image',
    'chat_background',
    'unknown'
    );

CREATE TABLE tbl_media
(
    id           SERIAL PRIMARY KEY,
    uuid         UUID          NOT NULL DEFAULT gen_random_uuid(),
    user_id      INT           NOT NULL REFERENCES tbl_user (id), -- by whom uploaded
    company_id   INT           NOT NULL DEFAULT 0,                -- by whom uploaded

-- Content metadata
    media_type   media_type    NOT NULL DEFAULT 'unknown',
    context      media_context NOT NULL DEFAULT 'unknown',
    context_id   INT                    DEFAULT 0,                -- References posts.id, messages.id, etc.
    context_uuid VARCHAR(100)           DEFAULT '',-- References posts.id, messages.id, etc.

-- File information
    filename     VARCHAR(255)  NOT NULL DEFAULT '',
    file_path    VARCHAR(900)           DEFAULT '',               -- without ENV.UPLOAD_PATH
    thumb_path   VARCHAR(900)           DEFAULT '',
    thumb_fn     VARCHAR(255)  NOT NULL DEFAULT '',
    original_fn  VARCHAR(255)  NOT NULL DEFAULT '',

-- Media metadata
    mime_type    VARCHAR(20)            DEFAULT '',
    file_size    BIGINT                 DEFAULT 0,-- in bytes
    duration     INT                    DEFAULT 0,-- For videos and audio
    width        INT                    DEFAULT 0,-- For images and videos
    height       INT                    DEFAULT 0,-- For images and videos

    meta         TEXT                   DEFAULT '',
    meta2        TEXT                   DEFAULT '',
    meta3        TEXT                   DEFAULT '',
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active       INT           NOT NULL DEFAULT 1,
    deleted      INT           NOT NULL DEFAULT 0
);

CREATE TYPE chat_type_t AS ENUM ('direct', 'group', 'channel');

CREATE TABLE tbl_conversation
(
    id                   SERIAL PRIMARY KEY,
    uuid                 UUID          NOT NULL DEFAULT gen_random_uuid(),
    chat_type            chat_type_t   NOT NULL DEFAULT 'direct',
    title                VARCHAR(200)  NOT NULL DEFAULT '',
    description          VARCHAR(1000) NOT NULL DEFAULT '',
    creator_id           INT           NOT NULL REFERENCES tbl_user (id),
    theme_color          VARCHAR(7)             DEFAULT '#FFFFFF',
    image_url            VARCHAR(200)           DEFAULT '',
    background_image_url VARCHAR(500)  NOT NULL DEFAULT '', -- also can be implemented with tbl_media
    background_blur      INT                    DEFAULT 0,  -- can add numbers to vary the blur intense
    last_message_id      INT                    DEFAULT 0,
    member_count         INT           NOT NULL DEFAULT 0,
    message_count        BIGINT        NOT NULL DEFAULT 0,
    auto_delete_duration INT           NOT NULL DEFAULT 0,  -- in minutes, should be beautified in Front-end
    last_activity        TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at           TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active               INT           NOT NULL DEFAULT 1,
    deleted              INT           NOT NULL DEFAULT 0
);

CREATE TYPE message_type_t AS ENUM ('text', 'media', 'system', 'reply', 'forward', 'sticker', 'pin', 'reaction', 'edit', 'delete');

CREATE TABLE tbl_message
(
    id                SERIAL PRIMARY KEY,
    uuid              UUID           NOT NULL DEFAULT gen_random_uuid(),
    conversation_id   INT            NOT NULL REFERENCES tbl_conversation (id) ON DELETE CASCADE,
    sender_id         INT            NOT NULL REFERENCES tbl_user (id),
    message_type      message_type_t NOT NULL DEFAULT 'text',
    content           TEXT                    DEFAULT '',
    reply_to_id       INT                     DEFAULT 0,
    forwarded_from_id INT                     DEFAULT 0, -- TODO: Is it okay to reference as a Foreign Key?
    media_id          INT                     DEFAULT 0,
    sticker_id        INT                     DEFAULT 0,
    is_edited         BOOLEAN        NOT NULL DEFAULT false,
    is_pinned         BOOLEAN        NOT NULL DEFAULT false,
    is_read           BOOLEAN        NOT NULL DEFAULT false,
    is_delivered      BOOLEAN        NOT NULL DEFAULT false,
    is_silent         BOOLEAN                 DEFAULT false,
    edited_at         TIMESTAMP               DEFAULT CURRENT_TIMESTAMP,
    created_at        TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active            INT            NOT NULL DEFAULT 1,
    deleted           INT            NOT NULL DEFAULT 0
);

CREATE TABLE tbl_message_media
(
    id         SERIAL PRIMARY KEY,
    message_id INTEGER   NOT NULL REFERENCES tbl_message (id),
    media_id   INTEGER   NOT NULL REFERENCES tbl_media (id),
    is_primary BOOLEAN   NOT NULL DEFAULT FALSE,
    sort_order INTEGER   NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (message_id, media_id)
);


CREATE TABLE tbl_message_reaction
(
    id         SERIAL PRIMARY KEY,
    message_id INT         NOT NULL REFERENCES tbl_message (id) ON DELETE CASCADE,
    user_id    INT         NOT NULL REFERENCES tbl_user (id) ON DELETE CASCADE,
    company_id INT         NOT NULL REFERENCES tbl_company (id) ON DELETE CASCADE,
    emoji      VARCHAR(50) NOT NULL,
    reaction   VARCHAR(50) NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (message_id, company_id, reaction)
);

CREATE TYPE notification_preference_t AS ENUM ('all', 'mentions', 'none');

CREATE TABLE tbl_conversation_member
(
    id                      SERIAL PRIMARY KEY,
    conversation_id         INT       NOT NULL REFERENCES tbl_conversation (id) ON DELETE CASCADE,
    user_id                 INT       NOT NULL REFERENCES tbl_user (id),
    is_admin                BOOLEAN   NOT NULL        DEFAULT false,
    nickname                VARCHAR(100),
    last_read_message_id    INT,
    unread_count            INT       NOT NULL        DEFAULT 0,
    notification_preference notification_preference_t DEFAULT 'all',
    muted_until             TIMESTAMP,
    joined_at               TIMESTAMP NOT NULL        DEFAULT CURRENT_TIMESTAMP,
    left_at                 TIMESTAMP,
    created_at              TIMESTAMP NOT NULL        DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP NOT NULL        DEFAULT CURRENT_TIMESTAMP,
    active               INT           NOT NULL DEFAULT 1,
    deleted              INT           NOT NULL DEFAULT 0,
    UNIQUE (conversation_id, user_id)
);
