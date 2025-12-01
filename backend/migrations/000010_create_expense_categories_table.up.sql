-- Create expense_categories table for managing expense categories and subcategories
CREATE TABLE expense_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id UUID REFERENCES expense_categories(id) ON DELETE CASCADE, -- NULL for categories, set for subcategories
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_expense_categories_parent_id ON expense_categories(parent_id);
CREATE INDEX idx_expense_categories_is_active ON expense_categories(is_active);
CREATE INDEX idx_expense_categories_sort_order ON expense_categories(sort_order);
CREATE INDEX idx_expense_categories_name ON expense_categories(name);

-- Unique constraint: name must be unique within the same parent level
CREATE UNIQUE INDEX idx_expense_categories_unique_name 
ON expense_categories(name, COALESCE(parent_id, '00000000-0000-0000-0000-000000000000'::UUID));

-- Trigger to auto-update updated_at
CREATE TRIGGER update_expense_categories_updated_at 
BEFORE UPDATE ON expense_categories
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at_column();

-- Comments for documentation
COMMENT ON TABLE expense_categories IS 'Categories and subcategories for expense management';
COMMENT ON COLUMN expense_categories.parent_id IS 'NULL for top-level categories, references parent category for subcategories';
COMMENT ON COLUMN expense_categories.sort_order IS 'Order for displaying categories in UI';
