package domain

import (
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/google/uuid"
)

// InvoiceStatus represents the payment status of an invoice
type InvoiceStatus string

const (
	InvoiceStatusPaid   InvoiceStatus = "paid"   // Cobrado
	InvoiceStatusUnpaid InvoiceStatus = "unpaid" // No Cobrado

	// DefaultVATPercent is the standard Spanish VAT rate as percentage points (21 = 21%).
	// Stored in DB and exposed in JSON as vatRate; CalculateAmounts uses percent (base * rate / 100).
	DefaultVATPercent float64 = 21.0
)

// VatRateAsPercent returns the VAT rate as percentage points (e.g. 21).
// It maps legacy values stored as a fraction (e.g. 0.21) to 21 for display/PDF.
func VatRateAsPercent(rate float64) float64 {
	if rate > 0 && rate < 1 {
		return rate * 100
	}
	return rate
}

// Invoice represents a billing invoice for services rendered
type Invoice struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	InvoiceNumber string        `json:"invoiceNumber" db:"invoice_number"` // F_2025_0001
	ClientID      uuid.UUID     `json:"clientId" db:"client_id"`
	AppointmentID *uuid.UUID    `json:"appointmentId,omitempty" db:"appointment_id"` // Nullable for manual invoices
	IssueDate     time.Time     `json:"issueDate" db:"issue_date"`
	DueDate       time.Time     `json:"dueDate" db:"due_date"` // Payment due date
	Description   string        `json:"description" db:"description"`
	BaseAmount    float64       `json:"baseAmount" db:"base_amount"`                 // Base imponible (sin IVA)
	VATRate       float64       `json:"vatRate" db:"vat_rate"`                       // Percentage points (21 = 21%)
	VATAmount     float64       `json:"vatAmount" db:"vat_amount"`                   // BaseAmount * VATRate / 100
	TotalAmount   float64       `json:"totalAmount" db:"total_amount"`               // BaseAmount + VATAmount
	Status        InvoiceStatus `json:"status" db:"status"`                          // paid/unpaid
	PaymentMethod *string       `json:"paymentMethod,omitempty" db:"payment_method"` // Nullable payment method
	Notes         string        `json:"notes,omitempty" db:"notes"`
	PDFPath       *string       `json:"pdfPath,omitempty" db:"pdf_path"` // Path to generated PDF (nullable)
	CreatedAt     time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time     `json:"updatedAt" db:"updated_at"`
	DeletedAt     *time.Time    `json:"-" db:"deleted_at"` // Soft delete timestamp

	// Relationships (populated via joins, not stored in DB)
	Client      *Client      `json:"client,omitempty" db:"-"`
	Appointment *Appointment `json:"appointment,omitempty" db:"-"`
}

// CalculateAmounts calculates VAT and total amounts based on base amount
func (i *Invoice) CalculateAmounts() {
	i.VATAmount = i.BaseAmount * (i.VATRate / 100)
	i.TotalAmount = i.BaseAmount + i.VATAmount
}

// IsPaid returns true if the invoice has been paid
func (i *Invoice) IsPaid() bool {
	return i.Status == InvoiceStatusPaid
}

// MarkAsPaid marks the invoice as paid
func (i *Invoice) MarkAsPaid() {
	i.Status = InvoiceStatusPaid
	i.UpdatedAt = time.Now()
}

// MarkAsUnpaid marks the invoice as unpaid
func (i *Invoice) MarkAsUnpaid() {
	i.Status = InvoiceStatusUnpaid
	i.UpdatedAt = time.Now()
}

// IsManual returns true if the invoice was created manually (not from appointment)
func (i *Invoice) IsManual() bool {
	return i.AppointmentID == nil
}

// Validate performs basic validation on the invoice
func (i *Invoice) Validate() error {
	if i.InvoiceNumber == "" {
		return ErrInvalidInvoiceNumber
	}
	if i.ClientID == uuid.Nil {
		return ErrInvalidClientID
	}
	if i.BaseAmount <= 0 {
		return ErrInvalidAmount
	}
	if i.Description == "" {
		return ErrInvalidDescription
	}
	return nil
}

// Custom errors
var (
	ErrInvalidInvoiceNumber = errors.NewValidationError("invoice number is required", nil)
	ErrInvalidClientID      = errors.NewValidationError("client ID is required", nil)
	ErrInvalidAmount        = errors.NewValidationError("base amount must be greater than 0", nil)
	ErrInvalidDescription   = errors.NewValidationError("description is required", nil)
)
