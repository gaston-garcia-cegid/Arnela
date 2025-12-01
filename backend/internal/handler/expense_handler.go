package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ExpenseHandler handles expense HTTP requests
type ExpenseHandler struct {
	expenseService service.ExpenseService
}

// NewExpenseHandler creates a new expense handler
func NewExpenseHandler(expenseService service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
	}
}

// CreateExpense godoc
// @Summary Create a new expense
// @Description Create a new expense with category and subcategory validation
// @Tags expenses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body service.CreateExpenseRequest true "Expense creation request"
// @Success 201 {object} domain.Expense
// @Failure 400 {object} ErrorResponse "Invalid request or validation error"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /billing/expenses [post]
func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var req service.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	expense, err := h.expenseService.CreateExpense(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, expense)
}

// GetExpense godoc
// @Summary Get an expense by ID
// @Description Retrieve an expense by its ID
// @Tags expenses
// @Security BearerAuth
// @Produce json
// @Param id path string true "Expense ID (UUID)"
// @Success 200 {object} domain.Expense
// @Failure 400 {object} ErrorResponse "Invalid ID format"
// @Failure 404 {object} ErrorResponse "Expense not found"
// @Router /billing/expenses/{id} [get]
func (h *ExpenseHandler) GetExpense(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid expense ID"})
		return
	}

	expense, err := h.expenseService.GetExpense(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, expense)
}

// ListExpenses godoc
// @Summary List expenses
// @Description Get a paginated list of expenses with optional filters
// @Tags expenses
// @Security BearerAuth
// @Produce json
// @Param categoryId query string false "Category ID (UUID)"
// @Param subcategoryId query string false "Subcategory ID (UUID)"
// @Param fromDate query string false "From date (YYYY-MM-DD)"
// @Param toDate query string false "To date (YYYY-MM-DD)"
// @Param hasInvoice query bool false "Has invoice (true/false)"
// @Param supplier query string false "Supplier name (partial match)"
// @Param search query string false "Search in supplier, invoice, notes"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} ErrorResponse "Invalid query parameters"
// @Router /billing/expenses [get]
func (h *ExpenseHandler) ListExpenses(c *gin.Context) {
	filters := repository.ExpenseFilters{
		Page:     1,
		PageSize: 20,
	}

	// Parse page and pageSize
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		}
	}

	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			filters.PageSize = pageSize
		}
	}

	// Parse categoryId
	if categoryIDStr := c.Query("categoryId"); categoryIDStr != "" {
		if categoryID, err := uuid.Parse(categoryIDStr); err == nil {
			filters.CategoryID = &categoryID
		}
	}

	// Parse subcategoryId
	if subcategoryIDStr := c.Query("subcategoryId"); subcategoryIDStr != "" {
		if subcategoryID, err := uuid.Parse(subcategoryIDStr); err == nil {
			filters.SubcategoryID = &subcategoryID
		}
	}

	// Parse dates
	if fromDateStr := c.Query("fromDate"); fromDateStr != "" {
		if fromDate, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			filters.FromDate = &fromDate
		}
	}

	if toDateStr := c.Query("toDate"); toDateStr != "" {
		if toDate, err := time.Parse("2006-01-02", toDateStr); err == nil {
			filters.ToDate = &toDate
		}
	}

	// Parse hasInvoice
	if hasInvoiceStr := c.Query("hasInvoice"); hasInvoiceStr != "" {
		if hasInvoice, err := strconv.ParseBool(hasInvoiceStr); err == nil {
			filters.HasInvoice = &hasInvoice
		}
	}

	// Parse supplier and search
	filters.Supplier = c.Query("supplier")
	filters.Search = c.Query("search")

	expenses, total, err := h.expenseService.ListExpenses(c.Request.Context(), filters)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       expenses,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		Total:      int64(total),
		TotalPages: (total + filters.PageSize - 1) / filters.PageSize,
	})
}

// UpdateExpense godoc
// @Summary Update an expense
// @Description Update an existing expense with validation
// @Tags expenses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Expense ID (UUID)"
// @Param request body service.UpdateExpenseRequest true "Expense update request"
// @Success 200 {object} domain.Expense
// @Failure 400 {object} ErrorResponse "Invalid request or validation error"
// @Failure 404 {object} ErrorResponse "Expense not found"
// @Router /billing/expenses/{id} [put]
func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid expense ID"})
		return
	}

	var req service.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	expense, err := h.expenseService.UpdateExpense(c.Request.Context(), id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, expense)
}

// DeleteExpense godoc
// @Summary Delete an expense
// @Description Soft delete an expense
// @Tags expenses
// @Security BearerAuth
// @Param id path string true "Expense ID (UUID)"
// @Success 204 "Expense deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid ID"
// @Failure 404 {object} ErrorResponse "Expense not found"
// @Router /billing/expenses/{id} [delete]
func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid expense ID"})
		return
	}

	if err := h.expenseService.DeleteExpense(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetExpensesByCategory godoc
// @Summary Get expenses by category
// @Description Retrieve all expenses for a specific category
// @Tags expenses
// @Security BearerAuth
// @Produce json
// @Param categoryId path string true "Category ID (UUID)"
// @Success 200 {array} domain.Expense
// @Failure 400 {object} ErrorResponse "Invalid category ID"
// @Failure 404 {object} ErrorResponse "Category not found"
// @Router /billing/expenses/category/{categoryId} [get]
func (h *ExpenseHandler) GetExpensesByCategory(c *gin.Context) {
	categoryIDParam := c.Param("categoryId")
	categoryID, err := uuid.Parse(categoryIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid category ID"})
		return
	}

	expenses, err := h.expenseService.GetExpensesByCategory(c.Request.Context(), categoryID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// GetExpensesBySupplier godoc
// @Summary Get expenses by supplier
// @Description Retrieve all expenses for a specific supplier
// @Tags expenses
// @Security BearerAuth
// @Produce json
// @Param supplier path string true "Supplier name"
// @Success 200 {array} domain.Expense
// @Router /billing/expenses/supplier/{supplier} [get]
func (h *ExpenseHandler) GetExpensesBySupplier(c *gin.Context) {
	supplier := c.Param("supplier")

	expenses, err := h.expenseService.GetExpensesBySupplier(c.Request.Context(), supplier)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, expenses)
}
