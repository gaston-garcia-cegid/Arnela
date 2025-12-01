package repository

import (
	"context"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/google/uuid"
)

// InvoiceFilters contains filters for listing invoices
type InvoiceFilters struct {
	Status      *domain.InvoiceStatus
	ClientID    *uuid.UUID
	FromDate    *time.Time
	ToDate      *time.Time
	Search      string
	Page        int
	PageSize    int
	IncludePaid bool
}

// InvoiceRepository defines the interface for invoice data access
type InvoiceRepository interface {
	// Create creates a new invoice
	Create(ctx context.Context, invoice *domain.Invoice) error

	// GetByID retrieves an invoice by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error)

	// GetByInvoiceNumber retrieves an invoice by its invoice number
	GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*domain.Invoice, error)

	// List retrieves a paginated list of invoices with filters
	List(ctx context.Context, filters InvoiceFilters) ([]*domain.Invoice, int, error)

	// Update updates an existing invoice
	Update(ctx context.Context, invoice *domain.Invoice) error

	// Delete soft deletes an invoice
	Delete(ctx context.Context, id uuid.UUID) error

	// GetNextInvoiceNumber generates the next invoice number for the given year
	GetNextInvoiceNumber(ctx context.Context, year int) (string, error)

	// GetByClientID retrieves all invoices for a specific client
	GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*domain.Invoice, error)

	// GetByAppointmentID retrieves the invoice associated with an appointment
	GetByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*domain.Invoice, error)

	// GetTotalRevenueByDateRange calculates total revenue between dates
	GetTotalRevenueByDateRange(ctx context.Context, fromDate, toDate time.Time) (float64, error)

	// GetUnpaidInvoices retrieves all unpaid invoices
	GetUnpaidInvoices(ctx context.Context) ([]*domain.Invoice, error)
}
