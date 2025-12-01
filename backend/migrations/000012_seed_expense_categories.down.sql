-- Remove seeded expense categories
-- This will cascade delete all subcategories due to ON DELETE CASCADE
DELETE FROM expense_categories WHERE parent_id IS NULL;
