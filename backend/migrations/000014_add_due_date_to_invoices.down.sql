-- Remove due_date column from invoices
ALTER TABLE invoices DROP COLUMN IF EXISTS due_date;
