-- ============================================
-- Arnela Database Setup Script
-- PostgreSQL 16
-- ============================================

-- Connect to postgres database first to create user and database
-- Run this script as postgres superuser:
-- psql -U postgres -f setup_database.sql

-- 1. Create user (if not exists)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_user WHERE usename = 'arnela_user') THEN
        CREATE USER arnela_user WITH PASSWORD 'arnela_secure_pass_2024';
        RAISE NOTICE 'User arnela_user created successfully';
    ELSE
        RAISE NOTICE 'User arnela_user already exists';
    END IF;
END
$$;

-- 2. Drop database if exists (CAUTION: This will delete all data!)
-- Uncomment the next line only if you want to recreate the database from scratch
-- DROP DATABASE IF EXISTS arnela_db;

-- 3. Create database
SELECT 'CREATE DATABASE arnela_db OWNER arnela_user'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'arnela_db')\gexec

-- 4. Grant privileges
GRANT ALL PRIVILEGES ON DATABASE arnela_db TO arnela_user;

-- 5. Connect to the new database
\c arnela_db

-- 6. Grant schema privileges
GRANT ALL ON SCHEMA public TO arnela_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO arnela_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO arnela_user;

-- 7. Set default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO arnela_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO arnela_user;

-- Success message
\echo '✅ Database setup completed successfully!'
\echo ''
\echo 'Database: arnela_db'
\echo 'User: arnela_user'
\echo 'Password: arnela_secure_pass_2024'
\echo ''
\echo 'You can now run the backend application.'
\echo 'Migrations will run automatically on startup.'
