package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type expenseCategoryRepository struct {
	db *sqlx.DB
}

// NewExpenseCategoryRepository creates a new expense category repository
func NewExpenseCategoryRepository(db *sqlx.DB) repository.ExpenseCategoryRepository {
	return &expenseCategoryRepository{db: db}
}

// Create creates a new expense category
func (r *expenseCategoryRepository) Create(ctx context.Context, category *domain.ExpenseCategory) error {
	query := `
		INSERT INTO expense_categories (
			id, name, description, parent_id, is_active, sort_order, created_at, updated_at
		) VALUES (
			:id, :name, :description, :parent_id, :is_active, :sort_order, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, category)
	if err != nil {
		return fmt.Errorf("failed to create expense category: %w", err)
	}

	return nil
}

// GetByID retrieves an expense category by ID
func (r *expenseCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategory, error) {
	var category domain.ExpenseCategory
	query := `SELECT * FROM expense_categories WHERE id = $1`

	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("expense category not found")
		}
		return nil, fmt.Errorf("failed to get expense category: %w", err)
	}

	return &category, nil
}

// GetByName retrieves an expense category by name
func (r *expenseCategoryRepository) GetByName(ctx context.Context, name string) (*domain.ExpenseCategory, error) {
	var category domain.ExpenseCategory
	query := `SELECT * FROM expense_categories WHERE name = $1`

	err := r.db.GetContext(ctx, &category, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("expense category not found")
		}
		return nil, fmt.Errorf("failed to get expense category: %w", err)
	}

	return &category, nil
}

// List retrieves all expense categories
func (r *expenseCategoryRepository) List(ctx context.Context) ([]*domain.ExpenseCategory, error) {
	var categories []*domain.ExpenseCategory
	query := `
		SELECT * FROM expense_categories 
		WHERE is_active = true 
		ORDER BY parent_id NULLS FIRST, sort_order ASC`

	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list expense categories: %w", err)
	}

	return categories, nil
}

// GetCategories retrieves only parent categories (no parent_id)
func (r *expenseCategoryRepository) GetCategories(ctx context.Context) ([]*domain.ExpenseCategory, error) {
	var categories []*domain.ExpenseCategory
	query := `
		SELECT * FROM expense_categories 
		WHERE parent_id IS NULL AND is_active = true 
		ORDER BY sort_order ASC`

	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

// GetSubcategories retrieves subcategories for a parent category
func (r *expenseCategoryRepository) GetSubcategories(ctx context.Context, parentID uuid.UUID) ([]*domain.ExpenseCategory, error) {
	var subcategories []*domain.ExpenseCategory
	query := `
		SELECT * FROM expense_categories 
		WHERE parent_id = $1 AND is_active = true 
		ORDER BY sort_order ASC`

	err := r.db.SelectContext(ctx, &subcategories, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subcategories: %w", err)
	}

	return subcategories, nil
}

// GetCategoryTree retrieves the full hierarchical tree
func (r *expenseCategoryRepository) GetCategoryTree(ctx context.Context) ([]*domain.ExpenseCategoryWithChildren, error) {
	// Get all parent categories
	parents, err := r.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	// Build tree structure
	tree := make([]*domain.ExpenseCategoryWithChildren, len(parents))
	for i, parent := range parents {
		// Get subcategories for each parent
		children, err := r.GetSubcategories(ctx, parent.ID)
		if err != nil {
			return nil, err
		}

		tree[i] = &domain.ExpenseCategoryWithChildren{
			ExpenseCategory: *parent,
			Children:        children,
		}
	}

	return tree, nil
}

// Update updates an existing expense category
func (r *expenseCategoryRepository) Update(ctx context.Context, category *domain.ExpenseCategory) error {
	query := `
		UPDATE expense_categories SET
			name = :name,
			description = :description,
			parent_id = :parent_id,
			is_active = :is_active,
			sort_order = :sort_order,
			updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, category)
	if err != nil {
		return fmt.Errorf("failed to update expense category: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.NewNotFoundError("expense category not found")
	}

	return nil
}

// Delete deletes an expense category (cascade deletes subcategories)
func (r *expenseCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM expense_categories WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete expense category: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.NewNotFoundError("expense category not found")
	}

	return nil
}

// NameExists checks if a category name already exists (excluding given ID)
func (r *expenseCategoryRepository) NameExists(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM expense_categories 
		WHERE name = $1 AND id != $2`

	err := r.db.GetContext(ctx, &count, query, name, excludeID)
	if err != nil {
		return false, fmt.Errorf("failed to check name exists: %w", err)
	}

	return count > 0, nil
}
