DROP DATABASE IF EXISTS db_tex;

CREATE DATABASE db_tex;

\c db_tex;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

DROP TABLE IF EXISTS tbl_content CASCADE;

DROP TABLE IF EXISTS tbl_content_type CASCADE;

DROP TABLE IF EXISTS tbl_language CASCADE;

DROP TRIGGER IF EXISTS update_content_updated_at ON tbl_content;

DROP FUNCTION IF EXISTS update_updated_at_column;

CREATE TABLE
    tbl_language (
        id SERIAL PRIMARY KEY,
        uuid UUID DEFAULT gen_random_uuid (),
        code VARCHAR(5) UNIQUE NOT NULL,
        name VARCHAR(100) NOT NULL
    );

CREATE TABLE
    tbl_content_type (
        id SERIAL PRIMARY KEY,
        uuid UUID DEFAULT gen_random_uuid (),
        name VARCHAR(50) DEFAULT '',
        title TEXT DEFAULT '',
        title_ru TEXT DEFAULT '',
        description TEXT DEFAULT '',
        parent_id INT DEFAULT 0,
        parent_name TEXT DEFAULT ''
    );

CREATE TABLE
    tbl_content (
        id SERIAL PRIMARY KEY,
        uuid UUID DEFAULT gen_random_uuid (),
        lang_id INT REFERENCES tbl_language (id) ON DELETE SET NULL DEFAULT 0,
        content_type_id INT REFERENCES tbl_content_type (id) ON DELETE SET NULL DEFAULT 0,
        title TEXT DEFAULT '',
        slogan TEXT DEFAULT '',
        subtitle TEXT DEFAULT '',
        description TEXT DEFAULT '',
        count INT DEFAULT 0,
        count_type TEXT DEFAULT '',
        image_url TEXT DEFAULT '',
        video_url TEXT DEFAULT '',
        step INT DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        active INT DEFAULT 1,
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
    BEFORE UPDATE ON tbl_content
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();