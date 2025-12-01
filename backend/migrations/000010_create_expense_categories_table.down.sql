-- Drop expense_categories table
DROP TRIGGER IF EXISTS update_expense_categories_updated_at ON expense_categories;
DROP INDEX IF EXISTS idx_expense_categories_unique_name;
DROP INDEX IF EXISTS idx_expense_categories_name;
DROP INDEX IF EXISTS idx_expense_categories_sort_order;
DROP INDEX IF EXISTS idx_expense_categories_is_active;
DROP INDEX IF EXISTS idx_expense_categories_parent_id;
DROP TABLE IF EXISTS expense_categories;
