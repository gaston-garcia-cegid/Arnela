package service

import (
	"context"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock for ExpenseRepository ---

type MockExpenseRepository struct{ mock.Mock }

func (m *MockExpenseRepository) Create(ctx context.Context, e *domain.Expense) error { return nil }
func (m *MockExpenseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	return nil, nil
}
func (m *MockExpenseRepository) List(ctx context.Context, f repository.ExpenseFilters) ([]*domain.Expense, int, error) {
	return nil, 0, nil
}
func (m *MockExpenseRepository) Update(ctx context.Context, e *domain.Expense) error { return nil }
func (m *MockExpenseRepository) Delete(ctx context.Context, id uuid.UUID) error      { return nil }
func (m *MockExpenseRepository) GetByCategory(ctx context.Context, id uuid.UUID) ([]*domain.Expense, error) {
	return nil, nil
}
func (m *MockExpenseRepository) GetBySupplier(ctx context.Context, s string) ([]*domain.Expense, error) {
	return nil, nil
}
func (m *MockExpenseRepository) GetTotalByDateRange(ctx context.Context, from, to time.Time) (float64, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(float64), args.Error(1)
}
func (m *MockExpenseRepository) GetTotalByCategory(ctx context.Context, from, to time.Time) (map[uuid.UUID]float64, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(map[uuid.UUID]float64), args.Error(1)
}

// --- Mock for ExpenseCategoryRepository ---

type MockExpenseCategoryRepository struct{ mock.Mock }

func (m *MockExpenseCategoryRepository) Create(ctx context.Context, c *domain.ExpenseCategory) error {
	return nil
}
func (m *MockExpenseCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ExpenseCategory), args.Error(1)
}
func (m *MockExpenseCategoryRepository) GetByName(ctx context.Context, n string) (*domain.ExpenseCategory, error) {
	return nil, nil
}
func (m *MockExpenseCategoryRepository) List(ctx context.Context) ([]*domain.ExpenseCategory, error) {
	return nil, nil
}
func (m *MockExpenseCategoryRepository) GetCategories(ctx context.Context) ([]*domain.ExpenseCategory, error) {
	return nil, nil
}
func (m *MockExpenseCategoryRepository) GetSubcategories(ctx context.Context, id uuid.UUID) ([]*domain.ExpenseCategory, error) {
	return nil, nil
}
func (m *MockExpenseCategoryRepository) GetCategoryTree(ctx context.Context) ([]*domain.ExpenseCategoryWithChildren, error) {
	return nil, nil
}
func (m *MockExpenseCategoryRepository) Update(ctx context.Context, c *domain.ExpenseCategory) error {
	return nil
}
func (m *MockExpenseCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockExpenseCategoryRepository) NameExists(ctx context.Context, n string, id uuid.UUID) (bool, error) {
	return false, nil
}

// --- Tests ---

func TestBillingStatsService_GetRevenueByMonth(t *testing.T) {
	ctx := context.Background()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)

	invoiceRepo := new(MockInvoiceRepository)
	expenseRepo := new(MockExpenseRepository)
	categoryRepo := new(MockExpenseCategoryRepository)

	svc := NewBillingStatsService(invoiceRepo, expenseRepo, categoryRepo)

	invoiceRepo.On("GetRevenueByMonth", ctx, from, to).Return([]repository.MonthlyRevenueRow{
		{Month: "2026-01", Revenue: 1500.0},
		{Month: "2026-02", Revenue: 2200.0},
		{Month: "2026-03", Revenue: 1800.0},
	}, nil)

	result, err := svc.GetRevenueByMonth(ctx, from, to)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "2026-01", result[0].Month)
	assert.Equal(t, 1500.0, result[0].Revenue)
	assert.Equal(t, "2026-03", result[2].Month)
	assert.Equal(t, 1800.0, result[2].Revenue)
	invoiceRepo.AssertExpectations(t)
}

func TestBillingStatsService_GetRevenueByMonth_Empty(t *testing.T) {
	ctx := context.Background()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)

	invoiceRepo := new(MockInvoiceRepository)
	expenseRepo := new(MockExpenseRepository)
	categoryRepo := new(MockExpenseCategoryRepository)

	svc := NewBillingStatsService(invoiceRepo, expenseRepo, categoryRepo)

	invoiceRepo.On("GetRevenueByMonth", ctx, from, to).Return([]repository.MonthlyRevenueRow{}, nil)

	result, err := svc.GetRevenueByMonth(ctx, from, to)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestBillingStatsService_GetBalance(t *testing.T) {
	ctx := context.Background()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC)

	invoiceRepo := new(MockInvoiceRepository)
	expenseRepo := new(MockExpenseRepository)
	categoryRepo := new(MockExpenseCategoryRepository)

	svc := NewBillingStatsService(invoiceRepo, expenseRepo, categoryRepo)

	invoiceRepo.On("GetTotalRevenueByDateRange", ctx, from, to).Return(5000.0, nil)
	expenseRepo.On("GetTotalByDateRange", ctx, from, to).Return(3200.0, nil)

	balance, err := svc.GetBalance(ctx, from, to)

	assert.NoError(t, err)
	assert.Equal(t, 1800.0, balance)
}

func TestBillingStatsService_GetExpensesByCategory(t *testing.T) {
	ctx := context.Background()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC)

	invoiceRepo := new(MockInvoiceRepository)
	expenseRepo := new(MockExpenseRepository)
	categoryRepo := new(MockExpenseCategoryRepository)

	svc := NewBillingStatsService(invoiceRepo, expenseRepo, categoryRepo)

	catID1 := uuid.New()
	catID2 := uuid.New()

	expenseRepo.On("GetTotalByCategory", ctx, from, to).Return(map[uuid.UUID]float64{
		catID1: 1000.0,
		catID2: 500.0,
	}, nil)

	categoryRepo.On("GetByID", ctx, catID1).Return(&domain.ExpenseCategory{
		ID: catID1, Name: "Material",
	}, nil)
	categoryRepo.On("GetByID", ctx, catID2).Return(&domain.ExpenseCategory{
		ID: catID2, Name: "Servicios",
	}, nil)

	result, err := svc.GetExpensesByCategory(ctx, from, to)

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	totalPercent := 0.0
	for _, ce := range result {
		totalPercent += ce.Percentage
		assert.True(t, ce.Total > 0)
	}
	assert.InDelta(t, 100.0, totalPercent, 0.01)
}
