-- Script to create user accounts for existing employees
-- This migrates employees who don't have a user_id yet
-- Run this script manually in your PostgreSQL database

-- Enable pgcrypto extension for bcrypt password hashing
CREATE EXTENSION IF NOT EXISTS pgcrypto;

BEGIN;

-- Create users for employees that don't have user_id
INSERT INTO users (
    id,
    email,
    password_hash,
    first_name,
    last_name,
    role,
    is_active,
    created_at,
    updated_at
)
SELECT 
    gen_random_uuid(),
    e.email,
    -- Hash the DNI as password (you'll need to update this with actual bcrypt hash)
    -- Default password: their DNI in uppercase
    -- Users should change this on first login
    crypt(UPPER(e.dni), gen_salt('bf')),
    e.first_name,
    e.last_name,
    'employee'::user_role,
    e.is_active,
    NOW(),
    NOW()
FROM employees e
WHERE e.user_id IS NULL
  AND e.deleted_at IS NULL
  AND e.email IS NOT NULL
  AND e.dni IS NOT NULL
ON CONFLICT (email) DO NOTHING;

-- Update employees with their new user_id
UPDATE employees e
SET 
    user_id = u.id,
    updated_at = NOW()
FROM users u
WHERE e.email = u.email
  AND e.user_id IS NULL
  AND e.deleted_at IS NULL
  AND u.role = 'employee';

-- Verify the migration
SELECT 
    COUNT(*) as total_employees,
    COUNT(user_id) as employees_with_user,
    COUNT(*) - COUNT(user_id) as employees_without_user
FROM employees
WHERE deleted_at IS NULL;

-- Show employees with their user info
SELECT 
    e.id as employee_id,
    e.first_name,
    e.last_name,
    e.email,
    e.dni,
    e.user_id,
    u.id as user_id_from_users,
    u.role as user_role,
    u.is_active as user_is_active
FROM employees e
LEFT JOIN users u ON e.user_id = u.id
WHERE e.deleted_at IS NULL
ORDER BY e.created_at DESC;

COMMIT;

-- NOTES:
-- 1. This script uses PostgreSQL's crypt() function with bcrypt
-- 2. You may need to enable the pgcrypto extension first:
--    CREATE EXTENSION IF NOT EXISTS pgcrypto;
-- 3. Employees can login with:
--    - Username: their email
--    - Password: their DNI (uppercase)
-- 4. Consider forcing password change on first login
