-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_clients_updated_at ON clients;
DROP FUNCTION IF EXISTS update_clients_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_clients_email;
DROP INDEX IF EXISTS idx_clients_dni;
DROP INDEX IF EXISTS idx_clients_phone;
DROP INDEX IF EXISTS idx_clients_user_id;
DROP INDEX IF EXISTS idx_clients_last_name;
DROP INDEX IF EXISTS idx_clients_is_active;
DROP INDEX IF EXISTS idx_clients_deleted_at;
DROP INDEX IF EXISTS idx_clients_email_unique;
DROP INDEX IF EXISTS idx_clients_dni_unique;

-- Drop table
DROP TABLE IF EXISTS clients;
