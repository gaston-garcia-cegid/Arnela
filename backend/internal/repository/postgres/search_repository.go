package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
)

// SearchRepository implements search operations across all entities
type SearchRepository struct {
	db *sql.DB
}

// NewSearchRepository creates a new search repository
func NewSearchRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{
		db: db,
	}
}

// SearchClients searches for clients by name, email, DNI/CIF, or phone
func (r *SearchRepository) SearchClients(ctx context.Context, query string, limit int) ([]domain.SearchClient, error) {
	queryPattern := "%" + query + "%"

	sqlQuery := `
		SELECT 
			id, 
			first_name, 
			last_name, 
			email, 
			COALESCE(phone, '') as phone,
			COALESCE(dni_cif, '') as dni_cif
		FROM clients
		WHERE 
			deleted_at IS NULL
			AND is_active = true
			AND (
				LOWER(first_name) LIKE LOWER($1)
				OR LOWER(last_name) LIKE LOWER($1)
				OR LOWER(email) LIKE LOWER($1)
				OR LOWER(dni_cif) LIKE LOWER($1)
				OR LOWER(phone) LIKE LOWER($1)
			)
		ORDER BY 
			CASE 
				WHEN LOWER(first_name) = LOWER($2) THEN 1
				WHEN LOWER(last_name) = LOWER($2) THEN 2
				WHEN LOWER(email) = LOWER($2) THEN 3
				ELSE 4
			END,
			first_name, last_name
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, sqlQuery, queryPattern, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search clients: %w", err)
	}
	defer rows.Close()

	var clients []domain.SearchClient
	for rows.Next() {
		var client domain.SearchClient
		err := rows.Scan(
			&client.ID,
			&client.FirstName,
			&client.LastName,
			&client.Email,
			&client.Phone,
			&client.DNICIF,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan client: %w", err)
		}
		clients = append(clients, client)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating client rows: %w", err)
	}

	return clients, nil
}

// SearchEmployees searches for employees by name, email, phone, or specialties
func (r *SearchRepository) SearchEmployees(ctx context.Context, query string, limit int) ([]domain.SearchEmployee, error) {
	queryPattern := "%" + query + "%"

	sqlQuery := `
		SELECT 
			e.id,
			CONCAT(u.first_name, ' ', u.last_name) as name,
			u.email,
			COALESCE(e.phone, '') as phone,
			COALESCE(e.specialties, '{}') as specialties,
			COALESCE(e.avatar_color, '#000000') as avatar_color
		FROM employees e
		INNER JOIN users u ON e.user_id = u.id
		WHERE 
			e.deleted_at IS NULL
			AND e.is_active = true
			AND (
				LOWER(u.first_name) LIKE LOWER($1)
				OR LOWER(u.last_name) LIKE LOWER($1)
				OR LOWER(u.email) LIKE LOWER($1)
				OR LOWER(e.phone) LIKE LOWER($1)
				OR EXISTS (
					SELECT 1 FROM unnest(e.specialties) s 
					WHERE LOWER(s) LIKE LOWER($1)
				)
			)
		ORDER BY 
			CASE 
				WHEN LOWER(u.first_name) = LOWER($2) OR LOWER(u.last_name) = LOWER($2) THEN 1
				ELSE 2
			END,
			u.first_name, u.last_name
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, sqlQuery, queryPattern, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search employees: %w", err)
	}
	defer rows.Close()

	var employees []domain.SearchEmployee
	for rows.Next() {
		var employee domain.SearchEmployee
		
		err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Email,
			&employee.Phone,
			pq.Array(&employee.Specialties),
			&employee.AvatarColor,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan employee: %w", err)
		}
		employees = append(employees, employee)
	}	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating employee rows: %w", err)
	}

	return employees, nil
}

// SearchAppointments searches for appointments by title, client name, or employee name
func (r *SearchRepository) SearchAppointments(ctx context.Context, query string, limit int) ([]domain.SearchAppointment, error) {
	queryPattern := "%" + query + "%"

	sqlQuery := `
		SELECT 
			a.id,
			a.title,
			a.start_time,
			a.end_time,
			a.status,
			CONCAT(c.first_name, ' ', c.last_name) as client_name,
			CONCAT(u.first_name, ' ', u.last_name) as employee_name
		FROM appointments a
		INNER JOIN clients c ON a.client_id = c.id
		INNER JOIN employees e ON a.employee_id = e.id
		INNER JOIN users u ON e.user_id = u.id
		WHERE 
			a.deleted_at IS NULL
			AND (
				LOWER(a.title) LIKE LOWER($1)
				OR LOWER(c.first_name) LIKE LOWER($1)
				OR LOWER(c.last_name) LIKE LOWER($1)
				OR LOWER(u.first_name) LIKE LOWER($1)
				OR LOWER(u.last_name) LIKE LOWER($1)
			)
		ORDER BY a.start_time DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, sqlQuery, queryPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search appointments: %w", err)
	}
	defer rows.Close()

	var appointments []domain.SearchAppointment
	for rows.Next() {
		var appointment domain.SearchAppointment
		err := rows.Scan(
			&appointment.ID,
			&appointment.Title,
			&appointment.StartTime,
			&appointment.EndTime,
			&appointment.Status,
			&appointment.ClientName,
			&appointment.EmployeeName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}
		appointments = append(appointments, appointment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating appointment rows: %w", err)
	}

	return appointments, nil
}

// SearchInvoices searches for invoices by invoice number or client name
func (r *SearchRepository) SearchInvoices(ctx context.Context, query string, limit int) ([]domain.SearchInvoice, error) {
	queryPattern := "%" + query + "%"

	sqlQuery := `
		SELECT 
			i.id,
			i.invoice_number,
			CONCAT(c.first_name, ' ', c.last_name) as client_name,
			i.total_amount,
			i.status,
			i.issue_date
		FROM invoices i
		INNER JOIN clients c ON i.client_id = c.id
		WHERE 
			i.deleted_at IS NULL
			AND (
				LOWER(i.invoice_number) LIKE LOWER($1)
				OR LOWER(c.first_name) LIKE LOWER($1)
				OR LOWER(c.last_name) LIKE LOWER($1)
			)
		ORDER BY 
			CASE 
				WHEN LOWER(i.invoice_number) = LOWER($2) THEN 1
				ELSE 2
			END,
			i.issue_date DESC
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, sqlQuery, queryPattern, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search invoices: %w", err)
	}
	defer rows.Close()

	var invoices []domain.SearchInvoice
	for rows.Next() {
		var invoice domain.SearchInvoice
		err := rows.Scan(
			&invoice.ID,
			&invoice.InvoiceNumber,
			&invoice.ClientName,
			&invoice.TotalAmount,
			&invoice.Status,
			&invoice.IssueDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invoice: %w", err)
		}
		invoices = append(invoices, invoice)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating invoice rows: %w", err)
	}

	return invoices, nil
}
