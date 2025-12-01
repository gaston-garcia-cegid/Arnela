-- Drop invoices table
DROP TRIGGER IF EXISTS update_invoices_updated_at ON invoices;
DROP INDEX IF EXISTS idx_invoices_created_at;
DROP INDEX IF EXISTS idx_invoices_invoice_number;
DROP INDEX IF EXISTS idx_invoices_issue_date;
DROP INDEX IF EXISTS idx_invoices_status;
DROP INDEX IF EXISTS idx_invoices_appointment_id;
DROP INDEX IF EXISTS idx_invoices_client_id;
DROP TABLE IF EXISTS invoices;
