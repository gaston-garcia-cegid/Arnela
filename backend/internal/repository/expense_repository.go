package repository

import (
	"context"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/google/uuid"
)

// ExpenseFilters contains filters for listing expenses
type ExpenseFilters struct {
	CategoryID    *uuid.UUID
	SubcategoryID *uuid.UUID
	FromDate      *time.Time
	ToDate        *time.Time
	HasInvoice    *bool
	Supplier      string
	Search        string
	Page          int
	PageSize      int
}

// ExpenseRepository defines the interface for expense data access
type ExpenseRepository interface {
	// Create creates a new expense
	Create(ctx context.Context, expense *domain.Expense) error

	// GetByID retrieves an expense by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error)

	// List retrieves a paginated list of expenses with filters
	List(ctx context.Context, filters ExpenseFilters) ([]*domain.Expense, int, error)

	// Update updates an existing expense
	Update(ctx context.Context, expense *domain.Expense) error

	// Delete soft deletes an expense
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByCategory retrieves all expenses for a specific category
	GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.Expense, error)

	// GetTotalByDateRange calculates total expenses between dates
	GetTotalByDateRange(ctx context.Context, fromDate, toDate time.Time) (float64, error)

	// GetTotalByCategory calculates total expenses by category in a date range
	GetTotalByCategory(ctx context.Context, fromDate, toDate time.Time) (map[uuid.UUID]float64, error)

	// GetBySupplier retrieves expenses by supplier name
	GetBySupplier(ctx context.Context, supplier string) ([]*domain.Expense, error)
}
