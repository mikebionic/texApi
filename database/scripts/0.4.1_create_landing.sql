CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE languages (
       id SERIAL PRIMARY KEY,
       uuid UUID DEFAULT gen_random_uuid (),
       code VARCHAR(5) UNIQUE NOT NULL,
       name VARCHAR(100) NOT NULL
);

CREATE TABLE content_types (
           id SERIAL PRIMARY KEY,
           uuid UUID DEFAULT gen_random_uuid (),
           type_name VARCHAR(50) DEFAULT '',
           title TEXT DEFAULT '',
           description TEXT DEFAULT ''
);

CREATE TABLE content (
     id SERIAL PRIMARY KEY,
     uuid UUID DEFAULT gen_random_uuid (),
     lang_id INT REFERENCES languages(id) ON DELETE SET NULL DEFAULT 0,
     content_type_id INT REFERENCES content_types(id) ON DELETE SET NULL DEFAULT 0,
     title TEXT DEFAULT '',
     subtitle TEXT DEFAULT '',
     description TEXT DEFAULT '',
     image_url TEXT DEFAULT '',
     video_url TEXT DEFAULT '',
     step INT DEFAULT 0,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     deleted INT DEFAULT 0
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