DROP DATABASE IF EXISTS db_tex;
create database db_tex;
create user mike with encrypted password 'pass123';
grant all privileges on database db_tex to mike;