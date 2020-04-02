CREATE USER useracl WITH PASSWORD 'p0stdb@';

create database acldev ENCODING "utf-8";
ALTER DATABASE acldev SET datestyle TO ISO, DMY;

\c acldev

GRANT CONNECT ON DATABASE acldev TO useracl;
GRANT USAGE ON SCHEMA public TO useracl;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO useracl;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO useracl;
grant ALL on database acldev to useracl;
