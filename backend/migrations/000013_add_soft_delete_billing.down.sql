-- Remove soft delete columns
ALTER TABLE expenses DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE invoices DROP COLUMN IF EXISTS deleted_at;
