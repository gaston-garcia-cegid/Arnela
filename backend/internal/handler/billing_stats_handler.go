package handler

import (
	"net/http"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// BillingStatsHandler handles billing statistics HTTP requests
type BillingStatsHandler struct {
	statsService service.BillingStatsService
}

// NewBillingStatsHandler creates a new billing stats handler
func NewBillingStatsHandler(statsService service.BillingStatsService) *BillingStatsHandler {
	return &BillingStatsHandler{
		statsService: statsService,
	}
}

// GetDashboardStats godoc
// @Summary Get billing dashboard statistics
// @Description Retrieve comprehensive billing statistics for the dashboard
// @Tags billing-stats
// @Security BearerAuth
// @Produce json
// @Param fromDate query string false "From date (YYYY-MM-DD)" default(first day of current month)
// @Param toDate query string false "To date (YYYY-MM-DD)" default(today)
// @Success 200 {object} service.DashboardStats
// @Failure 400 {object} ErrorResponse "Invalid date format"
// @Router /billing/dashboard [get]
func (h *BillingStatsHandler) GetDashboardStats(c *gin.Context) {
	// Parse dates or use defaults (current month)
	now := time.Now()
	fromDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	toDate := now

	if fromDateStr := c.Query("fromDate"); fromDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			fromDate = parsed
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid fromDate format, use YYYY-MM-DD"})
			return
		}
	}

	if toDateStr := c.Query("toDate"); toDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", toDateStr); err == nil {
			toDate = parsed
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid toDate format, use YYYY-MM-DD"})
			return
		}
	}

	stats, err := h.statsService.GetDashboardStats(c.Request.Context(), fromDate, toDate)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetRevenueByMonth godoc
// @Summary Get revenue grouped by month
// @Description Retrieve monthly revenue data for charts
// @Tags billing-stats
// @Security BearerAuth
// @Produce json
// @Param fromDate query string false "From date (YYYY-MM-DD)"
// @Param toDate query string false "To date (YYYY-MM-DD)"
// @Success 200 {array} service.MonthlyRevenue
// @Failure 400 {object} ErrorResponse "Invalid date format"
// @Router /billing/revenue-by-month [get]
func (h *BillingStatsHandler) GetRevenueByMonth(c *gin.Context) {
	// Parse dates or use defaults (last 6 months)
	now := time.Now()
	fromDate := now.AddDate(0, -6, 0)
	toDate := now

	if fromDateStr := c.Query("fromDate"); fromDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			fromDate = parsed
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid fromDate format, use YYYY-MM-DD"})
			return
		}
	}

	if toDateStr := c.Query("toDate"); toDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", toDateStr); err == nil {
			toDate = parsed
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid toDate format, use YYYY-MM-DD"})
			return
		}
	}

	revenue, err := h.statsService.GetRevenueByMonth(c.Request.Context(), fromDate, toDate)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, revenue)
}

// GetExpensesByCategory godoc
// @Summary Get expenses grouped by category
// @Description Retrieve expenses breakdown by category with percentages
// @Tags billing-stats
// @Security BearerAuth
// @Produce json
// @Param fromDate query string false "From date (YYYY-MM-DD)"
// @Param toDate query string false "To date (YYYY-MM-DD)"
// @Success 200 {array} service.CategoryExpense
// @Failure 400 {object} ErrorResponse "Invalid date format"
// @Router /billing/expenses-by-category [get]
func (h *BillingStatsHandler) GetExpensesByCategory(c *gin.Context) {
	// Parse dates or use defaults (current month)
	now := time.Now()
	fromDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	toDate := now

	if fromDateStr := c.Query("fromDate"); fromDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			fromDate = parsed
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid fromDate format, use YYYY-MM-DD"})
			return
		}
	}

	if toDateStr := c.Query("toDate"); toDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", toDateStr); err == nil {
			toDate = parsed
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid toDate format, use YYYY-MM-DD"})
			return
		}
	}

	expenses, err := h.statsService.GetExpensesByCategory(c.Request.Context(), fromDate, toDate)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// GetBalance godoc
// @Summary Get balance (revenue - expenses)
// @Description Calculate the balance for a given period
// @Tags billing-stats
// @Security BearerAuth
// @Produce json
// @Param fromDate query string false "From date (YYYY-MM-DD)"
// @Param toDate query string false "To date (YYYY-MM-DD)"
// @Success 200 {object} map[string]float64 "balance value"
// @Failure 400 {object} ErrorResponse "Invalid date format"
// @Router /billing/balance [get]
func (h *BillingStatsHandler) GetBalance(c *gin.Context) {
	// Parse dates or use defaults (current month)
	now := time.Now()
	fromDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	toDate := now

	if fromDateStr := c.Query("fromDate"); fromDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			fromDate = parsed
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid fromDate format, use YYYY-MM-DD"})
			return
		}
	}

	if toDateStr := c.Query("toDate"); toDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", toDateStr); err == nil {
			toDate = parsed
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid toDate format, use YYYY-MM-DD"})
			return
		}
	}

	balance, err := h.statsService.GetBalance(c.Request.Context(), fromDate, toDate)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance":  balance,
		"fromDate": fromDate.Format("2006-01-02"),
		"toDate":   toDate.Format("2006-01-02"),
	})
}
