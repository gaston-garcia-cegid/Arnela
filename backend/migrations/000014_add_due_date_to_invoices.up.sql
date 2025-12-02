-- Add due_date column to invoices table
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS due_date DATE NOT NULL DEFAULT CURRENT_DATE;

-- Update existing invoices to have due_date 15 days after issue_date
UPDATE invoices SET due_date = issue_date + INTERVAL '15 days' WHERE due_date IS NULL OR due_date = CURRENT_DATE;

-- Remove default after backfilling
ALTER TABLE invoices ALTER COLUMN due_date DROP DEFAULT;

COMMENT ON COLUMN invoices.due_date IS 'Payment due date';
