

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
    'company_image',
    'company_background',
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
    context_uuid VARCHAR(100),-- References posts.id, messages.id, etc.
    formatting   VARCHAR(100), -- blur, mask, spoiler, 18+, ... (anything).

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

CREATE TYPE message_type_t AS ENUM ('text', 'media', 'system', 'link', 'voice', 'video', 'reply', 'forward', 'sticker', 'pin', 'reaction', 'edit', 'read', 'delete');
CREATE TYPE chat_type_t AS ENUM ('direct', 'group', 'channel');

CREATE TABLE tbl_conversation
(
    id                   SERIAL PRIMARY KEY,
    uuid                 UUID          NOT NULL DEFAULT gen_random_uuid(),
    chat_type            chat_type_t   NOT NULL DEFAULT 'direct',
    title                VARCHAR(200)  NOT NULL DEFAULT '',
    description          VARCHAR(1000),
    creator_id           INT           NOT NULL REFERENCES tbl_user (id),
    theme_color          VARCHAR(7)             DEFAULT '#FFFFFF',
    image_url            VARCHAR(200),
    background_image_url VARCHAR(500), -- also can be implemented with tbl_media
    background_blur      INT,  -- can add numbers to vary the blur intense
    last_message_id      INT,
    member_count         INT           NOT NULL DEFAULT 0,
    message_count        BIGINT        NOT NULL DEFAULT 0,
    auto_delete_duration INT           NOT NULL DEFAULT 0,  -- in minutes, should be beautified in Front-end
    invite_token         VARCHAR(200),
    is_public            BOOLEAN                DEFAULT false,
    public_url           VARCHAR(200),
    last_activity        TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at           TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active               INT           NOT NULL DEFAULT 1,
    deleted              INT           NOT NULL DEFAULT 0
);


CREATE TABLE tbl_message
(
    id                SERIAL PRIMARY KEY,
    uuid              UUID           NOT NULL DEFAULT gen_random_uuid(),
    conversation_id   INT            NOT NULL REFERENCES tbl_conversation (id) ON DELETE CASCADE,
    sender_id         INT            NOT NULL REFERENCES tbl_user (id),
    message_type      message_type_t NOT NULL DEFAULT 'text',
    content           VARCHAR(800)           NOT NULL DEFAULT '',
    reply_to_id       INT                     DEFAULT 0,
    forwarded_from_id INT                     DEFAULT 0,    -- TODO: Is it okay to reference as a Foreign Key?
    media_id          INT                     DEFAULT 0,
    sticker_id        INT                     DEFAULT 0,
    is_edited         BOOLEAN        NOT NULL DEFAULT false,
    is_pinned         BOOLEAN        NOT NULL DEFAULT false,
    is_delivered      BOOLEAN,
    is_silent         BOOLEAN,
    edited_at         TIMESTAMP,                            -- check if I send edited_at info
    read_at           TIMESTAMP,
    read_by           JSONB                   DEFAULT '[]',
    deleted_for       JSONB                   DEFAULT '[]', -- with time [{'user_id': 1, 'dt':'2023-02-11 12:00:00'}, {'user_id': 2, 'dt':'2023-02-11 12:00:00'}]
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


-- AUTO-DELETE TODO: check UTC rules
CREATE OR REPLACE FUNCTION check_message_auto_delete()
    RETURNS TRIGGER AS $$
DECLARE
    conv_auto_delete_duration INT;
    company_self_destruct_duration INT;
BEGIN
    -- Get conversation's auto-delete duration
    SELECT auto_delete_duration INTO conv_auto_delete_duration
    FROM tbl_conversation WHERE id = NEW.conversation_id;

-- Get sender's company self-destruct duration
    SELECT self_destruct_duration INTO company_self_destruct_duration
    FROM tbl_company p
             JOIN tbl_user u ON u.company_id = p.id
    WHERE u.id = NEW.sender_id;

-- AUTO DELETE WITH CHECKING PROFILE SETTINGS:
    IF company_self_destruct_duration > 0 AND
       (conv_auto_delete_duration = 0 OR company_self_destruct_duration < conv_auto_delete_duration) THEN
        NEW.deleted = 1 AT TIME ZONE 'UTC' + (company_self_destruct_duration || ' minutes')::INTERVAL;
    ELSIF conv_auto_delete_duration > 0 THEN
        NEW.deleted = 1 AT TIME ZONE 'UTC' + (conv_auto_delete_duration || ' minutes')::INTERVAL;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER message_auto_delete_trigger
    BEFORE INSERT ON tbl_message
    FOR EACH ROW
EXECUTE FUNCTION check_message_auto_delete();


CREATE TABLE tbl_message_reaction
(
    id         SERIAL PRIMARY KEY,
    message_id INT         NOT NULL REFERENCES tbl_message (id) ON DELETE CASCADE,
    user_id    INT         NOT NULL REFERENCES tbl_user (id) ON DELETE CASCADE,
    company_id INT         NOT NULL REFERENCES tbl_company (id) ON DELETE CASCADE,
    emoji      VARCHAR(50) NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (message_id, company_id, emoji)
);

CREATE TABLE tbl_conversation_member
(
    id                      SERIAL PRIMARY KEY,
    conversation_id         INT       NOT NULL REFERENCES tbl_conversation (id) ON DELETE CASCADE,
    user_id                 INT       NOT NULL REFERENCES tbl_user (id),
    is_admin                BOOLEAN   NOT NULL        DEFAULT false, -- set by admin
    nickname                VARCHAR(100),
    privileges              TEXT[],                                  -- set by tbl_conversation.creator_id
    last_read_message_id    INT,                                     -- this approach also good but not pragmatic
    unread_count            INT       NOT NULL        DEFAULT 0,
    notification_preference notification_preference_t DEFAULT 'all',
    muted_until             TIMESTAMP,
    joined_at               TIMESTAMP NOT NULL        DEFAULT CURRENT_TIMESTAMP,
    left_at                 TIMESTAMP,
    created_at              TIMESTAMP NOT NULL        DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP NOT NULL        DEFAULT CURRENT_TIMESTAMP,
    active                  INT       NOT NULL        DEFAULT 1,
    deleted                 INT       NOT NULL        DEFAULT 0
);



CREATE TABLE tbl_company_contact
(
    id                SERIAL PRIMARY KEY,
    company_id        INT       NOT NULL REFERENCES tbl_company (id) ON DELETE CASCADE, -- my company
    system_company_id INT       NOT NULL REFERENCES tbl_company (id),                   -- other company (contact)
    first_name        VARCHAR(100)       DEFAULT '',
    last_name         VARCHAR(100)       DEFAULT '',
    nickname          VARCHAR(100)       DEFAULT '',
    phone_number      VARCHAR(100)       DEFAULT '',
    email             VARCHAR(100)       DEFAULT '',
    link              VARCHAR(500)       DEFAULT '',
    type              VARCHAR(50)        DEFAULT 'personal',                            -- e.g., 'personal', 'work', etc.
    custom_label      VARCHAR(100),
    is_favorite       BOOLEAN            DEFAULT false,
    blocked           BOOLEAN            DEFAULT false,
    blocked_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reason            TEXT,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (company_id, system_company_id)
);


CREATE TABLE tbl_call_room
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID         NOT NULL DEFAULT gen_random_uuid(),
    conversation_id INT          NOT NULL DEFAULT 0,
    max_user        INT          NOT NULL DEFAULT 2,
    user_ids        TEXT[]       NOT NULL DEFAULT '{}',
    profile_ids     TEXT[]       NOT NULL DEFAULT '{}',
    title           VARCHAR(200) NOT NULL DEFAULT '',
    hex             VARCHAR(60)  NOT NULL DEFAULT '',
    duration        INT          NOT NULL DEFAULT 60, -- 1h
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active          INT          NOT NULL DEFAULT 1,
    deleted         INT          NOT NULL DEFAULT 0
);