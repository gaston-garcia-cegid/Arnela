package domain

import (
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/google/uuid"
)

// Expense represents a business expense
type Expense struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	ExpenseDate     time.Time  `json:"expenseDate" db:"expense_date"`
	SupplierInvoice string     `json:"supplierInvoice,omitempty" db:"supplier_invoice"` // Nº Factura emisor (opcional)
	Supplier        string     `json:"supplier" db:"supplier"`                          // Nombre del proveedor
	Amount          float64    `json:"amount" db:"amount"`                              // Importe total
	CategoryID      uuid.UUID  `json:"categoryId" db:"category_id"`
	SubcategoryID   *uuid.UUID `json:"subcategoryId,omitempty" db:"subcategory_id"`   // Nullable
	HasInvoice      bool       `json:"hasInvoice" db:"has_invoice"`                   // Si/No
	AttachmentPath  string     `json:"attachmentPath,omitempty" db:"attachment_path"` // PDF path (si existe)
	Description     string     `json:"description,omitempty" db:"description"`
	PaymentMethod   string     `json:"paymentMethod,omitempty" db:"payment_method"`
	Notes           string     `json:"notes,omitempty" db:"notes"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time  `json:"updatedAt" db:"updated_at"`

	// Relationships (populated via joins, not stored in DB)
	Category    *ExpenseCategory `json:"category,omitempty" db:"-"`
	Subcategory *ExpenseCategory `json:"subcategory,omitempty" db:"-"`
}

// Validate performs basic validation on the expense
func (e *Expense) Validate() error {
	if e.Supplier == "" {
		return ErrInvalidSupplier
	}
	if e.Amount <= 0 {
		return ErrInvalidExpenseAmount
	}
	if e.CategoryID == uuid.Nil {
		return ErrInvalidCategory
	}
	if e.ExpenseDate.IsZero() {
		return ErrInvalidExpenseDate
	}
	return nil
}

// HasAttachment returns true if the expense has an attached invoice PDF
func (e *Expense) HasAttachment() bool {
	return e.AttachmentPath != ""
}

// ExpenseCategoryWithChildren represents a category with its subcategories
type ExpenseCategoryWithChildren struct {
	ExpenseCategory
	Children []*ExpenseCategory `json:"children,omitempty"`
}

// Custom errors
var (
	ErrInvalidSupplier      = errors.NewValidationError("supplier name is required", nil)
	ErrInvalidExpenseAmount = errors.NewValidationError("amount must be greater than 0", nil)
	ErrInvalidCategory      = errors.NewValidationError("category is required", nil)
	ErrInvalidExpenseDate   = errors.NewValidationError("expense date is required", nil)
)
