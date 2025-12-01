-- Revert consolidation of DNI/CIF back to NIF and DNI

-- Step 1: Rename dni_cif back to nif
ALTER TABLE clients RENAME COLUMN dni_cif TO nif;

-- Step 2: Re-add dni column as optional
ALTER TABLE clients ADD COLUMN dni VARCHAR(20);

-- Step 3: Update the unique index name
DROP INDEX IF EXISTS idx_clients_dni_cif_unique;
CREATE UNIQUE INDEX idx_clients_nif_unique ON clients(nif) WHERE deleted_at IS NULL;

-- Step 4: Remove comment
COMMENT ON COLUMN clients.nif IS NULL;
