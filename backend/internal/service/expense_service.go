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

// CreateExpenseRequest represents the request to create an expense
type CreateExpenseRequest struct {
	ExpenseDate     time.Time  `json:"expenseDate" binding:"required"`
	SupplierInvoice string     `json:"supplierInvoice,omitempty"`
	Supplier        string     `json:"supplier" binding:"required"`
	Amount          float64    `json:"amount" binding:"required,gt=0"`
	CategoryID      uuid.UUID  `json:"categoryId" binding:"required"`
	SubcategoryID   *uuid.UUID `json:"subcategoryId,omitempty"`
	HasInvoice      bool       `json:"hasInvoice"`
	AttachmentPath  string     `json:"attachmentPath,omitempty"`
	Notes           string     `json:"notes,omitempty"`
}

// UpdateExpenseRequest represents the request to update an expense
type UpdateExpenseRequest struct {
	ExpenseDate     time.Time  `json:"expenseDate" binding:"required"`
	SupplierInvoice string     `json:"supplierInvoice,omitempty"`
	Supplier        string     `json:"supplier" binding:"required"`
	Amount          float64    `json:"amount" binding:"required,gt=0"`
	CategoryID      uuid.UUID  `json:"categoryId" binding:"required"`
	SubcategoryID   *uuid.UUID `json:"subcategoryId,omitempty"`
	HasInvoice      bool       `json:"hasInvoice"`
	AttachmentPath  string     `json:"attachmentPath,omitempty"`
	Notes           string     `json:"notes,omitempty"`
}

// ExpenseService handles expense business logic
type ExpenseService interface {
	// CreateExpense creates a new expense
	CreateExpense(ctx context.Context, req *CreateExpenseRequest) (*domain.Expense, error)

	// GetExpense retrieves an expense by ID
	GetExpense(ctx context.Context, id uuid.UUID) (*domain.Expense, error)

	// ListExpenses retrieves a paginated list of expenses with filters
	ListExpenses(ctx context.Context, filters repository.ExpenseFilters) ([]*domain.Expense, int, error)

	// UpdateExpense updates an existing expense
	UpdateExpense(ctx context.Context, id uuid.UUID, req *UpdateExpenseRequest) (*domain.Expense, error)

	// DeleteExpense soft deletes an expense
	DeleteExpense(ctx context.Context, id uuid.UUID) error

	// GetExpensesByCategory retrieves expenses for a category
	GetExpensesByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.Expense, error)

	// GetExpensesBySupplier retrieves expenses by supplier
	GetExpensesBySupplier(ctx context.Context, supplier string) ([]*domain.Expense, error)

	// GetTotalExpenses calculates total expenses in a date range
	GetTotalExpenses(ctx context.Context, fromDate, toDate time.Time) (float64, error)
}

type expenseService struct {
	expenseRepo  repository.ExpenseRepository
	categoryRepo repository.ExpenseCategoryRepository
}

// NewExpenseService creates a new expense service
func NewExpenseService(expenseRepo repository.ExpenseRepository, categoryRepo repository.ExpenseCategoryRepository) ExpenseService {
	return &expenseService{
		expenseRepo:  expenseRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateExpense creates a new expense
func (s *expenseService) CreateExpense(ctx context.Context, req *CreateExpenseRequest) (*domain.Expense, error) {
	// Validate category exists
	category, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, errors.NewValidationError("category not found", map[string][]string{
			"categoryId": {"category does not exist"},
		})
	}

	// Validate subcategory if provided
	if req.SubcategoryID != nil {
		subcategory, err := s.categoryRepo.GetByID(ctx, *req.SubcategoryID)
		if err != nil {
			return nil, errors.NewValidationError("subcategory not found", map[string][]string{
				"subcategoryId": {"subcategory does not exist"},
			})
		}

		// Validate subcategory belongs to the category
		if subcategory.ParentID == nil || *subcategory.ParentID != req.CategoryID {
			return nil, errors.NewValidationError("subcategory does not belong to the selected category", map[string][]string{
				"subcategoryId": {"invalid subcategory for this category"},
			})
		}
	}

	// Validate category is a parent (not a subcategory)
	if category.IsSubcategory() {
		return nil, errors.NewValidationError("cannot use a subcategory as the main category", map[string][]string{
			"categoryId": {"must be a parent category, not a subcategory"},
		})
	}

	// Create expense
	expense := &domain.Expense{
		ID:              uuid.New(),
		ExpenseDate:     req.ExpenseDate,
		SupplierInvoice: req.SupplierInvoice,
		Supplier:        req.Supplier,
		Amount:          req.Amount,
		CategoryID:      req.CategoryID,
		SubcategoryID:   req.SubcategoryID,
		HasInvoice:      req.HasInvoice,
		AttachmentPath:  req.AttachmentPath,
		Notes:           req.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Validate expense
	if err := expense.Validate(); err != nil {
		return nil, err
	}

	// Save to database
	if err := s.expenseRepo.Create(ctx, expense); err != nil {
		return nil, fmt.Errorf("failed to create expense: %w", err)
	}

	return expense, nil
}

// GetExpense retrieves an expense by ID
func (s *expenseService) GetExpense(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	return s.expenseRepo.GetByID(ctx, id)
}

// ListExpenses retrieves a paginated list of expenses with filters
func (s *expenseService) ListExpenses(ctx context.Context, filters repository.ExpenseFilters) ([]*domain.Expense, int, error) {
	return s.expenseRepo.List(ctx, filters)
}

// UpdateExpense updates an existing expense
func (s *expenseService) UpdateExpense(ctx context.Context, id uuid.UUID, req *UpdateExpenseRequest) (*domain.Expense, error) {
	// Get existing expense
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate category exists
	category, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, errors.NewValidationError("category not found", map[string][]string{
			"categoryId": {"category does not exist"},
		})
	}

	// Validate subcategory if provided
	if req.SubcategoryID != nil {
		subcategory, err := s.categoryRepo.GetByID(ctx, *req.SubcategoryID)
		if err != nil {
			return nil, errors.NewValidationError("subcategory not found", map[string][]string{
				"subcategoryId": {"subcategory does not exist"},
			})
		}

		// Validate subcategory belongs to the category
		if subcategory.ParentID == nil || *subcategory.ParentID != req.CategoryID {
			return nil, errors.NewValidationError("subcategory does not belong to the selected category", map[string][]string{
				"subcategoryId": {"invalid subcategory for this category"},
			})
		}
	}

	// Validate category is a parent
	if category.IsSubcategory() {
		return nil, errors.NewValidationError("cannot use a subcategory as the main category", map[string][]string{
			"categoryId": {"must be a parent category, not a subcategory"},
		})
	}

	// Update fields
	expense.ExpenseDate = req.ExpenseDate
	expense.SupplierInvoice = req.SupplierInvoice
	expense.Supplier = req.Supplier
	expense.Amount = req.Amount
	expense.CategoryID = req.CategoryID
	expense.SubcategoryID = req.SubcategoryID
	expense.HasInvoice = req.HasInvoice
	expense.AttachmentPath = req.AttachmentPath
	expense.Notes = req.Notes
	expense.UpdatedAt = time.Now()

	// Validate
	if err := expense.Validate(); err != nil {
		return nil, err
	}

	// Save changes
	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return nil, fmt.Errorf("failed to update expense: %w", err)
	}

	return expense, nil
}

// DeleteExpense soft deletes an expense
func (s *expenseService) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	// Verify expense exists
	_, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.expenseRepo.Delete(ctx, id)
}

// GetExpensesByCategory retrieves expenses for a category
func (s *expenseService) GetExpensesByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.Expense, error) {
	// Validate category exists
	_, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, errors.NewNotFoundError("category not found")
	}

	return s.expenseRepo.GetByCategory(ctx, categoryID)
}

// GetExpensesBySupplier retrieves expenses by supplier
func (s *expenseService) GetExpensesBySupplier(ctx context.Context, supplier string) ([]*domain.Expense, error) {
	if supplier == "" {
		return nil, errors.NewValidationError("supplier name is required", nil)
	}

	return s.expenseRepo.GetBySupplier(ctx, supplier)
}

// GetTotalExpenses calculates total expenses in a date range
func (s *expenseService) GetTotalExpenses(ctx context.Context, fromDate, toDate time.Time) (float64, error) {
	if toDate.Before(fromDate) {
		return 0, errors.NewValidationError("end date must be after start date", nil)
	}

	return s.expenseRepo.GetTotalByDateRange(ctx, fromDate, toDate)
}
