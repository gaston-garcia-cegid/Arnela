package service

import (
	"context"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
)

// DashboardStats represents billing dashboard statistics
type DashboardStats struct {
	TotalRevenue       float64           `json:"totalRevenue"`
	TotalExpenses      float64           `json:"totalExpenses"`
	Balance            float64           `json:"balance"`
	UnpaidInvoices     int               `json:"unpaidInvoices"`
	UnpaidAmount       float64           `json:"unpaidAmount"`
	RevenueByMonth     []MonthlyRevenue  `json:"revenueByMonth"`
	ExpensesByCategory []CategoryExpense `json:"expensesByCategory"`
	RecentInvoices     []InvoiceSummary  `json:"recentInvoices"`
}

// MonthlyRevenue represents revenue for a specific month
type MonthlyRevenue struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
}

// CategoryExpense represents expenses grouped by category
type CategoryExpense struct {
	CategoryID   uuid.UUID `json:"categoryId"`
	CategoryName string    `json:"categoryName"`
	Total        float64   `json:"total"`
	Percentage   float64   `json:"percentage"`
}

// InvoiceSummary represents a summarized invoice
type InvoiceSummary struct {
	ID            uuid.UUID `json:"id"`
	InvoiceNumber string    `json:"invoiceNumber"`
	ClientID      uuid.UUID `json:"clientId"`
	TotalAmount   float64   `json:"totalAmount"`
	Status        string    `json:"status"`
	IssueDate     time.Time `json:"issueDate"`
}

// BillingStatsService handles billing statistics and dashboard data
type BillingStatsService interface {
	// GetDashboardStats retrieves comprehensive dashboard statistics
	GetDashboardStats(ctx context.Context, fromDate, toDate time.Time) (*DashboardStats, error)

	// GetRevenueByMonth calculates revenue grouped by month
	GetRevenueByMonth(ctx context.Context, fromDate, toDate time.Time) ([]MonthlyRevenue, error)

	// GetExpensesByCategory calculates expenses grouped by category
	GetExpensesByCategory(ctx context.Context, fromDate, toDate time.Time) ([]CategoryExpense, error)

	// GetBalance calculates the balance (revenue - expenses) for a period
	GetBalance(ctx context.Context, fromDate, toDate time.Time) (float64, error)
}

type billingStatsService struct {
	invoiceRepo  repository.InvoiceRepository
	expenseRepo  repository.ExpenseRepository
	categoryRepo repository.ExpenseCategoryRepository
}

// NewBillingStatsService creates a new billing stats service
func NewBillingStatsService(
	invoiceRepo repository.InvoiceRepository,
	expenseRepo repository.ExpenseRepository,
	categoryRepo repository.ExpenseCategoryRepository,
) BillingStatsService {
	return &billingStatsService{
		invoiceRepo:  invoiceRepo,
		expenseRepo:  expenseRepo,
		categoryRepo: categoryRepo,
	}
}

// GetDashboardStats retrieves comprehensive dashboard statistics
func (s *billingStatsService) GetDashboardStats(ctx context.Context, fromDate, toDate time.Time) (*DashboardStats, error) {
	// Get total revenue
	totalRevenue, err := s.invoiceRepo.GetTotalRevenueByDateRange(ctx, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	// Get total expenses
	totalExpenses, err := s.expenseRepo.GetTotalByDateRange(ctx, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	// Calculate balance
	balance := totalRevenue - totalExpenses

	// Get unpaid invoices
	unpaidInvoices, err := s.invoiceRepo.GetUnpaidInvoices(ctx)
	if err != nil {
		return nil, err
	}

	unpaidCount := len(unpaidInvoices)
	unpaidAmount := 0.0
	for _, invoice := range unpaidInvoices {
		unpaidAmount += invoice.TotalAmount
	}

	// Get revenue by month (last 6 months)
	revenueByMonth, err := s.GetRevenueByMonth(ctx, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	// Get expenses by category
	expensesByCategory, err := s.GetExpensesByCategory(ctx, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	// Get recent invoices (last 5)
	recentInvoicesFilters := repository.InvoiceFilters{
		Page:     1,
		PageSize: 5,
	}
	recentInvoicesList, _, err := s.invoiceRepo.List(ctx, recentInvoicesFilters)
	if err != nil {
		return nil, err
	}

	recentInvoices := make([]InvoiceSummary, len(recentInvoicesList))
	for i, invoice := range recentInvoicesList {
		recentInvoices[i] = InvoiceSummary{
			ID:            invoice.ID,
			InvoiceNumber: invoice.InvoiceNumber,
			ClientID:      invoice.ClientID,
			TotalAmount:   invoice.TotalAmount,
			Status:        string(invoice.Status),
			IssueDate:     invoice.IssueDate,
		}
	}

	return &DashboardStats{
		TotalRevenue:       totalRevenue,
		TotalExpenses:      totalExpenses,
		Balance:            balance,
		UnpaidInvoices:     unpaidCount,
		UnpaidAmount:       unpaidAmount,
		RevenueByMonth:     revenueByMonth,
		ExpensesByCategory: expensesByCategory,
		RecentInvoices:     recentInvoices,
	}, nil
}

// GetRevenueByMonth calculates revenue grouped by month
func (s *billingStatsService) GetRevenueByMonth(ctx context.Context, fromDate, toDate time.Time) ([]MonthlyRevenue, error) {
	// This is a simplified implementation
	// For production, you would want to query the database with GROUP BY month

	// For now, return a basic structure
	// TODO: Implement proper monthly aggregation in repository
	return []MonthlyRevenue{}, nil
}

// GetExpensesByCategory calculates expenses grouped by category
func (s *billingStatsService) GetExpensesByCategory(ctx context.Context, fromDate, toDate time.Time) ([]CategoryExpense, error) {
	// Get total expenses by category
	totals, err := s.expenseRepo.GetTotalByCategory(ctx, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	// Calculate grand total
	grandTotal := 0.0
	for _, total := range totals {
		grandTotal += total
	}

	// Build result with category names and percentages
	result := make([]CategoryExpense, 0, len(totals))
	for categoryID, total := range totals {
		category, err := s.categoryRepo.GetByID(ctx, categoryID)
		if err != nil {
			continue // Skip if category not found
		}

		percentage := 0.0
		if grandTotal > 0 {
			percentage = (total / grandTotal) * 100
		}

		result = append(result, CategoryExpense{
			CategoryID:   categoryID,
			CategoryName: category.Name,
			Total:        total,
			Percentage:   percentage,
		})
	}

	return result, nil
}

// GetBalance calculates the balance (revenue - expenses) for a period
func (s *billingStatsService) GetBalance(ctx context.Context, fromDate, toDate time.Time) (float64, error) {
	revenue, err := s.invoiceRepo.GetTotalRevenueByDateRange(ctx, fromDate, toDate)
	if err != nil {
		return 0, err
	}

	expenses, err := s.expenseRepo.GetTotalByDateRange(ctx, fromDate, toDate)
	if err != nil {
		return 0, err
	}

	return revenue - expenses, nil
}
