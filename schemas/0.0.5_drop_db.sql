-- 1. Connect to the default 'postgres' database (or any existing database)
\c postgres;

-- 2. Drop the target database if it exists
DROP DATABASE IF EXISTS db_tex;

-- 3. Create a new empty database
CREATE DATABASE db_tex;

-- 4. Connect to the newly created database
\c db_tex;

CREATE EXTENSION IF NOT EXISTS pgcrypto;