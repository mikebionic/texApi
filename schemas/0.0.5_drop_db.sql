\c postgres;

DROP DATABASE IF EXISTS db_tex;
CREATE DATABASE db_tex;
\c db_tex;
CREATE EXTENSION IF NOT EXISTS pgcrypto;