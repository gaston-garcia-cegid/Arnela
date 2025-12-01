package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type expenseRepository struct {
	db *sqlx.DB
}

// NewExpenseRepository creates a new expense repository
func NewExpenseRepository(db *sqlx.DB) repository.ExpenseRepository {
	return &expenseRepository{db: db}
}

// Create creates a new expense
func (r *expenseRepository) Create(ctx context.Context, expense *domain.Expense) error {
	query := `
		INSERT INTO expenses (
			id, expense_date, supplier_invoice, supplier, amount,
			category_id, subcategory_id, has_invoice, attachment_path, notes,
			created_at, updated_at
		) VALUES (
			:id, :expense_date, :supplier_invoice, :supplier, :amount,
			:category_id, :subcategory_id, :has_invoice, :attachment_path, :notes,
			:created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, expense)
	if err != nil {
		return fmt.Errorf("failed to create expense: %w", err)
	}

	return nil
}

// GetByID retrieves an expense by ID
func (r *expenseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	var expense domain.Expense
	query := `SELECT * FROM expenses WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, &expense, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("expense not found")
		}
		return nil, fmt.Errorf("failed to get expense: %w", err)
	}

	return &expense, nil
}

// List retrieves a paginated list of expenses with filters
func (r *expenseRepository) List(ctx context.Context, filters repository.ExpenseFilters) ([]*domain.Expense, int, error) {
	conditions := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	argCount := 0

	// Apply filters
	if filters.CategoryID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argCount))
		args = append(args, *filters.CategoryID)
	}

	if filters.SubcategoryID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("subcategory_id = $%d", argCount))
		args = append(args, *filters.SubcategoryID)
	}

	if filters.FromDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("expense_date >= $%d", argCount))
		args = append(args, *filters.FromDate)
	}

	if filters.ToDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("expense_date <= $%d", argCount))
		args = append(args, *filters.ToDate)
	}

	if filters.HasInvoice != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("has_invoice = $%d", argCount))
		args = append(args, *filters.HasInvoice)
	}

	if filters.Supplier != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("supplier ILIKE $%d", argCount))
		args = append(args, "%"+filters.Supplier+"%")
	}

	if filters.Search != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("(supplier ILIKE $%d OR supplier_invoice ILIKE $%d OR notes ILIKE $%d)", argCount, argCount, argCount))
		args = append(args, "%"+filters.Search+"%")
	}

	whereClause := strings.Join(conditions, " AND ")

	// Get total count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM expenses WHERE %s", whereClause)
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count expenses: %w", err)
	}

	// Pagination
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}
	if filters.PageSize > 100 {
		filters.PageSize = 100
	}

	offset := (filters.Page - 1) * filters.PageSize
	argCount++
	limitArg := argCount
	argCount++
	offsetArg := argCount
	args = append(args, filters.PageSize, offset)

	// Get expenses
	query := fmt.Sprintf(`
		SELECT * FROM expenses 
		WHERE %s 
		ORDER BY expense_date DESC, created_at DESC 
		LIMIT $%d OFFSET $%d`,
		whereClause, limitArg, offsetArg)

	var expenses []*domain.Expense
	err = r.db.SelectContext(ctx, &expenses, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list expenses: %w", err)
	}

	return expenses, total, nil
}

// Update updates an existing expense
func (r *expenseRepository) Update(ctx context.Context, expense *domain.Expense) error {
	query := `
		UPDATE expenses SET
			expense_date = :expense_date,
			supplier_invoice = :supplier_invoice,
			supplier = :supplier,
			amount = :amount,
			category_id = :category_id,
			subcategory_id = :subcategory_id,
			has_invoice = :has_invoice,
			attachment_path = :attachment_path,
			notes = :notes,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL`

	result, err := r.db.NamedExecContext(ctx, query, expense)
	if err != nil {
		return fmt.Errorf("failed to update expense: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.NewNotFoundError("expense not found")
	}

	return nil
}

// Delete soft deletes an expense
func (r *expenseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE expenses SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.NewNotFoundError("expense not found")
	}

	return nil
}

// GetByCategory retrieves all expenses for a specific category
func (r *expenseRepository) GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.Expense, error) {
	var expenses []*domain.Expense
	query := `
		SELECT * FROM expenses 
		WHERE category_id = $1 AND deleted_at IS NULL 
		ORDER BY expense_date DESC`

	err := r.db.SelectContext(ctx, &expenses, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get expenses by category: %w", err)
	}

	return expenses, nil
}

// GetTotalByDateRange calculates total expenses between dates
func (r *expenseRepository) GetTotalByDateRange(ctx context.Context, fromDate, toDate time.Time) (float64, error) {
	var total sql.NullFloat64
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE expense_date >= $1 
		AND expense_date <= $2 
		AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, &total, query, fromDate, toDate)
	if err != nil {
		return 0, fmt.Errorf("failed to get total expenses: %w", err)
	}

	if !total.Valid {
		return 0, nil
	}

	return total.Float64, nil
}

// GetTotalByCategory calculates total expenses by category in a date range
func (r *expenseRepository) GetTotalByCategory(ctx context.Context, fromDate, toDate time.Time) (map[uuid.UUID]float64, error) {
	type categoryTotal struct {
		CategoryID uuid.UUID `db:"category_id"`
		Total      float64   `db:"total"`
	}

	var results []categoryTotal
	query := `
		SELECT category_id, COALESCE(SUM(amount), 0) as total
		FROM expenses
		WHERE expense_date >= $1 
		AND expense_date <= $2 
		AND deleted_at IS NULL
		GROUP BY category_id`

	err := r.db.SelectContext(ctx, &results, query, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get expenses by category: %w", err)
	}

	totals := make(map[uuid.UUID]float64)
	for _, result := range results {
		totals[result.CategoryID] = result.Total
	}

	return totals, nil
}

// GetBySupplier retrieves expenses by supplier name
func (r *expenseRepository) GetBySupplier(ctx context.Context, supplier string) ([]*domain.Expense, error) {
	var expenses []*domain.Expense
	query := `
		SELECT * FROM expenses 
		WHERE supplier ILIKE $1 AND deleted_at IS NULL 
		ORDER BY expense_date DESC`

	err := r.db.SelectContext(ctx, &expenses, query, "%"+supplier+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get expenses by supplier: %w", err)
	}

	return expenses, nil
}
