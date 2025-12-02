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

type invoiceRepository struct {
	db *sqlx.DB
}

// NewInvoiceRepository creates a new invoice repository
func NewInvoiceRepository(db *sqlx.DB) repository.InvoiceRepository {
	return &invoiceRepository{db: db}
}

// Create creates a new invoice
func (r *invoiceRepository) Create(ctx context.Context, invoice *domain.Invoice) error {
	query := `
		INSERT INTO invoices (
			id, invoice_number, client_id, appointment_id, issue_date, due_date, description,
			base_amount, vat_rate, vat_amount, total_amount, status, notes, created_at, updated_at
		) VALUES (
			:id, :invoice_number, :client_id, :appointment_id, :issue_date, :due_date, :description,
			:base_amount, :vat_rate, :vat_amount, :total_amount, :status, :notes, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, invoice)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "invoice_number") {
				return errors.NewConflictError("invoice number already exists", errors.CodeConflict)
			}
		}
		return fmt.Errorf("failed to create invoice: %w", err)
	}

	return nil
}

// GetByID retrieves an invoice by ID
func (r *invoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	var invoice domain.Invoice
	query := `SELECT * FROM invoices WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, &invoice, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("invoice not found")
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	return &invoice, nil
}

// GetByInvoiceNumber retrieves an invoice by its invoice number
func (r *invoiceRepository) GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*domain.Invoice, error) {
	var invoice domain.Invoice
	query := `SELECT * FROM invoices WHERE invoice_number = $1 AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, &invoice, query, invoiceNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("invoice not found")
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	return &invoice, nil
}

// List retrieves a paginated list of invoices with filters
func (r *invoiceRepository) List(ctx context.Context, filters repository.InvoiceFilters) ([]*domain.Invoice, int, error) {
	conditions := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	argCount := 0

	// Apply filters
	if filters.Status != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *filters.Status)
	}

	if filters.ClientID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("client_id = $%d", argCount))
		args = append(args, *filters.ClientID)
	}

	if filters.FromDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("issue_date >= $%d", argCount))
		args = append(args, *filters.FromDate)
	}

	if filters.ToDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("issue_date <= $%d", argCount))
		args = append(args, *filters.ToDate)
	}

	if filters.Search != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("invoice_number ILIKE $%d", argCount))
		args = append(args, "%"+filters.Search+"%")
	}

	whereClause := strings.Join(conditions, " AND ")

	// Get total count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM invoices WHERE %s", whereClause)
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count invoices: %w", err)
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

	// Get invoices
	query := fmt.Sprintf(`
		SELECT * FROM invoices 
		WHERE %s 
		ORDER BY issue_date DESC, created_at DESC 
		LIMIT $%d OFFSET $%d`,
		whereClause, limitArg, offsetArg)

	var invoices []*domain.Invoice
	err = r.db.SelectContext(ctx, &invoices, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list invoices: %w", err)
	}

	return invoices, total, nil
}

// Update updates an existing invoice
func (r *invoiceRepository) Update(ctx context.Context, invoice *domain.Invoice) error {
	query := `
		UPDATE invoices SET
			invoice_number = :invoice_number,
			client_id = :client_id,
			appointment_id = :appointment_id,
			issue_date = :issue_date,
			due_date = :due_date,
			description = :description,
			base_amount = :base_amount,
			vat_rate = :vat_rate,
			vat_amount = :vat_amount,
			total_amount = :total_amount,
			status = :status,
			notes = :notes,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL`

	result, err := r.db.NamedExecContext(ctx, query, invoice)
	if err != nil {
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.NewNotFoundError("invoice not found")
	}

	return nil
}

// Delete soft deletes an invoice
func (r *invoiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE invoices SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete invoice: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.NewNotFoundError("invoice not found")
	}

	return nil
}

// GetNextInvoiceNumber generates the next invoice number for the given year
func (r *invoiceRepository) GetNextInvoiceNumber(ctx context.Context, year int) (string, error) {
	var maxNumber sql.NullInt64
	query := `
		SELECT MAX(CAST(SUBSTRING(invoice_number FROM 8) AS INTEGER))
		FROM invoices
		WHERE invoice_number LIKE $1 AND deleted_at IS NULL`

	pattern := fmt.Sprintf("F_%d_%%", year)
	err := r.db.GetContext(ctx, &maxNumber, query, pattern)
	if err != nil {
		return "", fmt.Errorf("failed to get next invoice number: %w", err)
	}

	nextNumber := 1
	if maxNumber.Valid {
		nextNumber = int(maxNumber.Int64) + 1
	}

	return fmt.Sprintf("F_%d_%04d", year, nextNumber), nil
}

// GetByClientID retrieves all invoices for a specific client
func (r *invoiceRepository) GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*domain.Invoice, error) {
	var invoices []*domain.Invoice
	query := `
		SELECT * FROM invoices 
		WHERE client_id = $1 AND deleted_at IS NULL 
		ORDER BY issue_date DESC`

	err := r.db.SelectContext(ctx, &invoices, query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoices by client: %w", err)
	}

	return invoices, nil
}

// GetByAppointmentID retrieves the invoice associated with an appointment
func (r *invoiceRepository) GetByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*domain.Invoice, error) {
	var invoice domain.Invoice
	query := `SELECT * FROM invoices WHERE appointment_id = $1 AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, &invoice, query, appointmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("invoice not found")
		}
		return nil, fmt.Errorf("failed to get invoice by appointment: %w", err)
	}

	return &invoice, nil
}

// GetTotalRevenueByDateRange calculates total revenue between dates
func (r *invoiceRepository) GetTotalRevenueByDateRange(ctx context.Context, fromDate, toDate time.Time) (float64, error) {
	var total sql.NullFloat64
	query := `
		SELECT COALESCE(SUM(total_amount), 0)
		FROM invoices
		WHERE issue_date >= $1 
		AND issue_date <= $2 
		AND status = 'paid'
		AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, &total, query, fromDate, toDate)
	if err != nil {
		return 0, fmt.Errorf("failed to get total revenue: %w", err)
	}

	if !total.Valid {
		return 0, nil
	}

	return total.Float64, nil
}

// GetUnpaidInvoices retrieves all unpaid invoices
func (r *invoiceRepository) GetUnpaidInvoices(ctx context.Context) ([]*domain.Invoice, error) {
	var invoices []*domain.Invoice
	query := `
		SELECT * FROM invoices 
		WHERE status = 'unpaid' AND deleted_at IS NULL 
		ORDER BY due_date ASC`

	err := r.db.SelectContext(ctx, &invoices, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get unpaid invoices: %w", err)
	}

	return invoices, nil
}
