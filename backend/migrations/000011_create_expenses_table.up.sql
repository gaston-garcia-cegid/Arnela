-- Create expenses table for expense tracking
CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expense_date DATE NOT NULL,
    supplier_invoice VARCHAR(100), -- Nº Factura emisor (opcional)
    supplier VARCHAR(255) NOT NULL, -- Nombre del proveedor
    amount DECIMAL(10,2) NOT NULL, -- Importe total
    category_id UUID NOT NULL REFERENCES expense_categories(id) ON DELETE RESTRICT,
    subcategory_id UUID REFERENCES expense_categories(id) ON DELETE RESTRICT,
    has_invoice BOOLEAN DEFAULT false, -- Si/No
    attachment_path VARCHAR(255), -- PDF path
    description TEXT,
    payment_method VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_expenses_expense_date ON expenses(expense_date);
CREATE INDEX idx_expenses_supplier ON expenses(supplier);
CREATE INDEX idx_expenses_category_id ON expenses(category_id);
CREATE INDEX idx_expenses_subcategory_id ON expenses(subcategory_id);
CREATE INDEX idx_expenses_has_invoice ON expenses(has_invoice);
CREATE INDEX idx_expenses_created_at ON expenses(created_at);
CREATE INDEX idx_expenses_supplier_invoice ON expenses(supplier_invoice) WHERE supplier_invoice IS NOT NULL;

-- Trigger to auto-update updated_at
CREATE TRIGGER update_expenses_updated_at 
BEFORE UPDATE ON expenses
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at_column();

-- Constraint: subcategory must belong to the selected category
CREATE OR REPLACE FUNCTION check_expense_subcategory()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.subcategory_id IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1 FROM expense_categories 
            WHERE id = NEW.subcategory_id 
            AND parent_id = NEW.category_id
        ) THEN
            RAISE EXCEPTION 'Subcategory must belong to the selected category';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_expense_subcategory
BEFORE INSERT OR UPDATE ON expenses
FOR EACH ROW
EXECUTE FUNCTION check_expense_subcategory();

-- Comments for documentation
COMMENT ON TABLE expenses IS 'Business expenses tracking';
COMMENT ON COLUMN expenses.supplier_invoice IS 'Supplier invoice number (optional)';
COMMENT ON COLUMN expenses.has_invoice IS 'Whether the expense has an attached invoice PDF';
COMMENT ON COLUMN expenses.attachment_path IS 'Path to attached invoice PDF file';
