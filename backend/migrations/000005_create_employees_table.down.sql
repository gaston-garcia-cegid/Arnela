-- Drop employees table and indexes
DROP TRIGGER IF EXISTS update_employees_updated_at ON employees;
DROP INDEX IF EXISTS idx_employees_deleted_at;
DROP INDEX IF EXISTS idx_employees_position;
DROP INDEX IF EXISTS idx_employees_is_active;
DROP INDEX IF EXISTS idx_employees_dni;
DROP INDEX IF EXISTS idx_employees_email;
DROP INDEX IF EXISTS idx_employees_user_id;
DROP TABLE IF EXISTS employees CASCADE;
