-- Create invoices table for billing management
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(50) UNIQUE NOT NULL, -- F_2025_0001
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE RESTRICT,
    appointment_id UUID REFERENCES appointments(id) ON DELETE SET NULL,
    issue_date DATE NOT NULL,
    due_date DATE NOT NULL,
    description TEXT NOT NULL,
    base_amount DECIMAL(10,2) NOT NULL, -- Base imponible (sin IVA)
    vat_rate DECIMAL(5,2) NOT NULL DEFAULT 21.00, -- IVA rate (21% by default)
    vat_amount DECIMAL(10,2) NOT NULL, -- Calculated VAT amount
    total_amount DECIMAL(10,2) NOT NULL, -- Total amount (base + VAT)
    status VARCHAR(20) NOT NULL CHECK (status IN ('paid', 'unpaid')), -- Cobrado/No Cobrado
    payment_method VARCHAR(50),
    notes TEXT,
    pdf_path VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_invoices_client_id ON invoices(client_id);
CREATE INDEX idx_invoices_appointment_id ON invoices(appointment_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_issue_date ON invoices(issue_date);
CREATE INDEX idx_invoices_invoice_number ON invoices(invoice_number);
CREATE INDEX idx_invoices_created_at ON invoices(created_at);

-- Trigger to auto-update updated_at
CREATE TRIGGER update_invoices_updated_at 
BEFORE UPDATE ON invoices
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at_column();

-- Comments for documentation
COMMENT ON TABLE invoices IS 'Billing invoices for services rendered';
COMMENT ON COLUMN invoices.invoice_number IS 'Unique invoice number format: F_YYYY_NNNN (e.g., F_2025_0001)';
COMMENT ON COLUMN invoices.base_amount IS 'Base amount before VAT';
COMMENT ON COLUMN invoices.vat_rate IS 'VAT rate percentage (21% by default in Spain)';
COMMENT ON COLUMN invoices.total_amount IS 'Total amount including VAT';
COMMENT ON COLUMN invoices.status IS 'Payment status: paid (Cobrado) or unpaid (No Cobrado)';
