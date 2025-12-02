-- Add deleted_at column for soft delete support
ALTER TABLE invoices ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;
ALTER TABLE expenses ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;

-- Add indexes for soft delete queries
CREATE INDEX idx_invoices_deleted_at ON invoices(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_expenses_deleted_at ON expenses(deleted_at) WHERE deleted_at IS NULL;

-- Comments
COMMENT ON COLUMN invoices.deleted_at IS 'Soft delete timestamp - NULL means active';
COMMENT ON COLUMN expenses.deleted_at IS 'Soft delete timestamp - NULL means active';
