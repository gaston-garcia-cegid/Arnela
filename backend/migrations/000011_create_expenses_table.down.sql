-- Drop expenses table
DROP TRIGGER IF EXISTS validate_expense_subcategory ON expenses;
DROP FUNCTION IF EXISTS check_expense_subcategory();
DROP TRIGGER IF EXISTS update_expenses_updated_at ON expenses;
DROP INDEX IF EXISTS idx_expenses_supplier_invoice;
DROP INDEX IF EXISTS idx_expenses_created_at;
DROP INDEX IF EXISTS idx_expenses_has_invoice;
DROP INDEX IF EXISTS idx_expenses_subcategory_id;
DROP INDEX IF EXISTS idx_expenses_category_id;
DROP INDEX IF EXISTS idx_expenses_supplier;
DROP INDEX IF EXISTS idx_expenses_expense_date;
DROP TABLE IF EXISTS expenses;
