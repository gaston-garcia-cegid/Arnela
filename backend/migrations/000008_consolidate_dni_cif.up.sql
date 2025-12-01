-- Consolidate NIF and DNI into single DNI/CIF field
-- This migration renames nif to dni_cif and removes the old dni field

-- Step 1: Rename nif column to dni_cif
ALTER TABLE clients RENAME COLUMN nif TO dni_cif;

-- Step 2: Drop the old dni column (it was optional, so data loss is acceptable)
ALTER TABLE clients DROP COLUMN IF EXISTS dni;

-- Step 3: Update the unique index name
DROP INDEX IF EXISTS idx_clients_nif_unique;
CREATE UNIQUE INDEX idx_clients_dni_cif_unique ON clients(dni_cif) WHERE deleted_at IS NULL;

-- Step 4: Add comment to column
COMMENT ON COLUMN clients.dni_cif IS 'Spanish DNI/CIF (Tax ID) - Required and unique';
