package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/google/uuid"
)

// CreateInvoiceRequest represents the request to create an invoice
type CreateInvoiceRequest struct {
	ClientID      uuid.UUID  `json:"clientId" binding:"required"`
	AppointmentID *uuid.UUID `json:"appointmentId,omitempty"`
	IssueDate     time.Time  `json:"issueDate" binding:"required"`
	DueDate       time.Time  `json:"dueDate" binding:"required"`
	BaseAmount    float64    `json:"baseAmount" binding:"required,gt=0"`
	Notes         string     `json:"notes,omitempty"`
}

// UpdateInvoiceRequest represents the request to update an invoice
type UpdateInvoiceRequest struct {
	IssueDate  time.Time `json:"issueDate" binding:"required"`
	DueDate    time.Time `json:"dueDate" binding:"required"`
	BaseAmount float64   `json:"baseAmount" binding:"required,gt=0"`
	Notes      string    `json:"notes,omitempty"`
}

// InvoiceService handles invoice business logic
type InvoiceService interface {
	// CreateInvoice creates a new invoice with automatic VAT calculation
	CreateInvoice(ctx context.Context, req *CreateInvoiceRequest) (*domain.Invoice, error)

	// CreateInvoiceFromAppointment creates an invoice automatically from an appointment
	CreateInvoiceFromAppointment(ctx context.Context, appointmentID uuid.UUID, baseAmount float64) (*domain.Invoice, error)

	// GetInvoice retrieves an invoice by ID
	GetInvoice(ctx context.Context, id uuid.UUID) (*domain.Invoice, error)

	// GetInvoiceByNumber retrieves an invoice by its invoice number
	GetInvoiceByNumber(ctx context.Context, invoiceNumber string) (*domain.Invoice, error)

	// ListInvoices retrieves a paginated list of invoices with filters
	ListInvoices(ctx context.Context, filters repository.InvoiceFilters) ([]*domain.Invoice, int, error)

	// UpdateInvoice updates an existing invoice
	UpdateInvoice(ctx context.Context, id uuid.UUID, req *UpdateInvoiceRequest) (*domain.Invoice, error)

	// DeleteInvoice soft deletes an invoice
	DeleteInvoice(ctx context.Context, id uuid.UUID) error

	// MarkAsPaid marks an invoice as paid
	MarkAsPaid(ctx context.Context, id uuid.UUID) (*domain.Invoice, error)

	// GetClientInvoices retrieves all invoices for a client
	GetClientInvoices(ctx context.Context, clientID uuid.UUID) ([]*domain.Invoice, error)

	// GetUnpaidInvoices retrieves all unpaid invoices
	GetUnpaidInvoices(ctx context.Context) ([]*domain.Invoice, error)
}

type invoiceService struct {
	invoiceRepo repository.InvoiceRepository
	clientRepo  repository.ClientRepository
}

// NewInvoiceService creates a new invoice service
func NewInvoiceService(invoiceRepo repository.InvoiceRepository, clientRepo repository.ClientRepository) InvoiceService {
	return &invoiceService{
		invoiceRepo: invoiceRepo,
		clientRepo:  clientRepo,
	}
}

// CreateInvoice creates a new invoice with automatic VAT calculation
func (s *invoiceService) CreateInvoice(ctx context.Context, req *CreateInvoiceRequest) (*domain.Invoice, error) {
	// Validate client exists
	_, err := s.clientRepo.GetByID(ctx, req.ClientID)
	if err != nil {
		return nil, errors.NewValidationError("client not found", map[string][]string{
			"clientId": {"client does not exist"},
		})
	}

	// Validate dates
	if req.DueDate.Before(req.IssueDate) {
		return nil, errors.NewValidationError("due date must be after issue date", map[string][]string{
			"dueDate": {"must be after issue date"},
		})
	}

	// Generate invoice number
	year := req.IssueDate.Year()
	invoiceNumber, err := s.invoiceRepo.GetNextInvoiceNumber(ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice number: %w", err)
	}

	// Create invoice
	invoice := &domain.Invoice{
		ID:            uuid.New(),
		InvoiceNumber: invoiceNumber,
		ClientID:      req.ClientID,
		AppointmentID: req.AppointmentID,
		IssueDate:     req.IssueDate,
		DueDate:       req.DueDate,
		BaseAmount:    req.BaseAmount,
		VATRate:       domain.FixedVATRate,
		Status:        domain.InvoiceStatusUnpaid,
		Notes:         req.Notes,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Calculate VAT and total
	invoice.CalculateAmounts()

	// Validate invoice
	if err := invoice.Validate(); err != nil {
		return nil, err
	}

	// Save to database
	if err := s.invoiceRepo.Create(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	return invoice, nil
}

// CreateInvoiceFromAppointment creates an invoice automatically from an appointment
func (s *invoiceService) CreateInvoiceFromAppointment(ctx context.Context, appointmentID uuid.UUID, baseAmount float64) (*domain.Invoice, error) {
	// TODO: Implement once AppointmentRepository is created
	// For now, return not implemented error
	return nil, errors.NewValidationError("appointment invoice creation not yet implemented", nil)
}

// GetInvoice retrieves an invoice by ID
func (s *invoiceService) GetInvoice(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	return s.invoiceRepo.GetByID(ctx, id)
}

// GetInvoiceByNumber retrieves an invoice by its invoice number
func (s *invoiceService) GetInvoiceByNumber(ctx context.Context, invoiceNumber string) (*domain.Invoice, error) {
	return s.invoiceRepo.GetByInvoiceNumber(ctx, invoiceNumber)
}

// ListInvoices retrieves a paginated list of invoices with filters
func (s *invoiceService) ListInvoices(ctx context.Context, filters repository.InvoiceFilters) ([]*domain.Invoice, int, error) {
	return s.invoiceRepo.List(ctx, filters)
}

// UpdateInvoice updates an existing invoice
func (s *invoiceService) UpdateInvoice(ctx context.Context, id uuid.UUID, req *UpdateInvoiceRequest) (*domain.Invoice, error) {
	// Get existing invoice
	invoice, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cannot update paid invoices
	if invoice.IsPaid() {
		return nil, errors.NewValidationError("cannot update a paid invoice", map[string][]string{
			"status": {"invoice is already paid"},
		})
	}

	// Validate dates
	if req.DueDate.Before(req.IssueDate) {
		return nil, errors.NewValidationError("due date must be after issue date", map[string][]string{
			"dueDate": {"must be after issue date"},
		})
	}

	// Update fields
	invoice.IssueDate = req.IssueDate
	invoice.DueDate = req.DueDate
	invoice.BaseAmount = req.BaseAmount
	invoice.Notes = req.Notes
	invoice.UpdatedAt = time.Now()

	// Recalculate VAT and total
	invoice.CalculateAmounts()

	// Validate
	if err := invoice.Validate(); err != nil {
		return nil, err
	}

	// Save changes
	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}

	return invoice, nil
}

// DeleteInvoice soft deletes an invoice
func (s *invoiceService) DeleteInvoice(ctx context.Context, id uuid.UUID) error {
	// Get invoice to check if it exists and is paid
	invoice, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Cannot delete paid invoices
	if invoice.IsPaid() {
		return errors.NewValidationError("cannot delete a paid invoice", map[string][]string{
			"status": {"invoice is already paid"},
		})
	}

	return s.invoiceRepo.Delete(ctx, id)
}

// MarkAsPaid marks an invoice as paid
func (s *invoiceService) MarkAsPaid(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	invoice, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if invoice.IsPaid() {
		return nil, errors.NewValidationError("invoice is already paid", map[string][]string{
			"status": {"invoice is already marked as paid"},
		})
	}

	invoice.MarkAsPaid()
	invoice.UpdatedAt = time.Now()

	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to mark invoice as paid: %w", err)
	}

	return invoice, nil
}

// GetClientInvoices retrieves all invoices for a client
func (s *invoiceService) GetClientInvoices(ctx context.Context, clientID uuid.UUID) ([]*domain.Invoice, error) {
	// Validate client exists
	_, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, errors.NewNotFoundError("client not found")
	}

	return s.invoiceRepo.GetByClientID(ctx, clientID)
}

// GetUnpaidInvoices retrieves all unpaid invoices
func (s *invoiceService) GetUnpaidInvoices(ctx context.Context) ([]*domain.Invoice, error) {
	return s.invoiceRepo.GetUnpaidInvoices(ctx)
}
