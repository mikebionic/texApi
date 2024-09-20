DROP DATABASE IF EXISTS db_tex;
create database db_tex;
create user mike with encrypted password 'PASSWORD_PLACEHOLDER';
grant all privileges on database db_tex to mike;