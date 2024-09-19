CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE languages (
                           id SERIAL PRIMARY KEY,
                           uuid UUID DEFAULT gen_random_uuid(),
                           code VARCHAR(5) UNIQUE NOT NULL,
                           name VARCHAR(100) NOT NULL
);

CREATE TABLE content_types (
                               id SERIAL PRIMARY KEY,
                               uuid UUID DEFAULT gen_random_uuid(),
                               type_name VARCHAR(50),
                               title TEXT,                            -- Title for the content type
                               description TEXT                       -- Description for the content type
);

CREATE TABLE content (
                         id SERIAL PRIMARY KEY,
                         uuid UUID DEFAULT gen_random_uuid(),
                         lang_id INT REFERENCES languages(id) ON DELETE SET NULL,
                         content_type_id INT REFERENCES content_types(id) ON DELETE SET NULL,
                         title TEXT,                  -- Title (unlimited)
                         subtitle TEXT,               -- Subtitle (unlimited)
                         description TEXT,            -- Description or content text
                         image_url TEXT,              -- Image URL (unlimited)
                         video_url TEXT,              -- Video URL (unlimited)
                         step INT,                    -- Display order
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         deleted INT DEFAULT 0        -- 0 if not deleted, 1 if deleted
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_content_updated_at
    BEFORE UPDATE ON content
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();