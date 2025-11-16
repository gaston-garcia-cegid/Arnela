-- Drop NIF indexes
DROP INDEX IF EXISTS idx_clients_nif;
DROP INDEX IF EXISTS idx_clients_nif_unique;

-- Drop modified DNI unique constraint
DROP INDEX IF EXISTS idx_clients_dni_unique;

-- Restore DNI unique constraint
CREATE UNIQUE INDEX idx_clients_dni_unique ON clients(dni) WHERE deleted_at IS NULL;

-- Make DNI required again
ALTER TABLE clients ALTER COLUMN dni SET NOT NULL;

-- Remove NIF column
ALTER TABLE clients DROP COLUMN IF EXISTS nif;
