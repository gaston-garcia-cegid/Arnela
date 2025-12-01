package repository

import (
	"context"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/google/uuid"
)

// ExpenseCategoryRepository defines the interface for expense category data access
type ExpenseCategoryRepository interface {
	// Create creates a new expense category
	Create(ctx context.Context, category *domain.ExpenseCategory) error

	// GetByID retrieves an expense category by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategory, error)

	// GetByName retrieves an expense category by name
	GetByName(ctx context.Context, name string) (*domain.ExpenseCategory, error)

	// List retrieves all expense categories
	List(ctx context.Context) ([]*domain.ExpenseCategory, error)

	// GetCategories retrieves only parent categories (no parent_id)
	GetCategories(ctx context.Context) ([]*domain.ExpenseCategory, error)

	// GetSubcategories retrieves subcategories for a parent category
	GetSubcategories(ctx context.Context, parentID uuid.UUID) ([]*domain.ExpenseCategory, error)

	// GetCategoryTree retrieves the full hierarchical tree
	GetCategoryTree(ctx context.Context) ([]*domain.ExpenseCategoryWithChildren, error)

	// Update updates an existing expense category
	Update(ctx context.Context, category *domain.ExpenseCategory) error

	// Delete deletes an expense category (cascade deletes subcategories)
	Delete(ctx context.Context, id uuid.UUID) error

	// NameExists checks if a category name already exists (excluding given ID)
	NameExists(ctx context.Context, name string, excludeID uuid.UUID) (bool, error)
}
