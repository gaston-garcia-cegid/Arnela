package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InvoiceHandler handles invoice HTTP requests
type InvoiceHandler struct {
	invoiceService service.InvoiceService
}

// NewInvoiceHandler creates a new invoice handler
func NewInvoiceHandler(invoiceService service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceService: invoiceService,
	}
}

// CreateInvoice godoc
// @Summary Create a new invoice
// @Description Create a new invoice with automatic VAT calculation (21%)
// @Tags invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body service.CreateInvoiceRequest true "Invoice creation request"
// @Success 201 {object} domain.Invoice
// @Failure 400 {object} ErrorResponse "Invalid request or validation error"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /billing/invoices [post]
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	var req service.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	invoice, err := h.invoiceService.CreateInvoice(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

// GetInvoice godoc
// @Summary Get an invoice by ID
// @Description Retrieve an invoice by its ID
// @Tags invoices
// @Security BearerAuth
// @Produce json
// @Param id path string true "Invoice ID (UUID)"
// @Success 200 {object} domain.Invoice
// @Failure 400 {object} ErrorResponse "Invalid ID format"
// @Failure 404 {object} ErrorResponse "Invoice not found"
// @Router /billing/invoices/{id} [get]
func (h *InvoiceHandler) GetInvoice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid invoice ID"})
		return
	}

	invoice, err := h.invoiceService.GetInvoice(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// GetInvoiceByNumber godoc
// @Summary Get an invoice by invoice number
// @Description Retrieve an invoice by its invoice number (e.g., F_2025_0001)
// @Tags invoices
// @Security BearerAuth
// @Produce json
// @Param invoiceNumber path string true "Invoice Number (F_YYYY_NNNN)"
// @Success 200 {object} domain.Invoice
// @Failure 404 {object} ErrorResponse "Invoice not found"
// @Router /billing/invoices/number/{invoiceNumber} [get]
func (h *InvoiceHandler) GetInvoiceByNumber(c *gin.Context) {
	invoiceNumber := c.Param("invoiceNumber")

	invoice, err := h.invoiceService.GetInvoiceByNumber(c.Request.Context(), invoiceNumber)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// ListInvoices godoc
// @Summary List invoices
// @Description Get a paginated list of invoices with optional filters
// @Tags invoices
// @Security BearerAuth
// @Produce json
// @Param status query string false "Invoice status (paid/unpaid)"
// @Param clientId query string false "Client ID (UUID)"
// @Param fromDate query string false "From date (YYYY-MM-DD)"
// @Param toDate query string false "To date (YYYY-MM-DD)"
// @Param search query string false "Search in invoice number"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} ErrorResponse "Invalid query parameters"
// @Router /billing/invoices [get]
func (h *InvoiceHandler) ListInvoices(c *gin.Context) {
	filters := repository.InvoiceFilters{
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

	// Parse status
	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.InvoiceStatus(statusStr)
		filters.Status = &status
	}

	// Parse clientId
	if clientIDStr := c.Query("clientId"); clientIDStr != "" {
		if clientID, err := uuid.Parse(clientIDStr); err == nil {
			filters.ClientID = &clientID
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

	// Parse search
	filters.Search = c.Query("search")

	invoices, total, err := h.invoiceService.ListInvoices(c.Request.Context(), filters)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       invoices,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		Total:      int64(total),
		TotalPages: (total + filters.PageSize - 1) / filters.PageSize,
	})
}

// UpdateInvoice godoc
// @Summary Update an invoice
// @Description Update an existing invoice (cannot update paid invoices)
// @Tags invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID (UUID)"
// @Param request body service.UpdateInvoiceRequest true "Invoice update request"
// @Success 200 {object} domain.Invoice
// @Failure 400 {object} ErrorResponse "Invalid request or validation error"
// @Failure 404 {object} ErrorResponse "Invoice not found"
// @Router /billing/invoices/{id} [put]
func (h *InvoiceHandler) UpdateInvoice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid invoice ID"})
		return
	}

	var req service.UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	invoice, err := h.invoiceService.UpdateInvoice(c.Request.Context(), id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// DeleteInvoice godoc
// @Summary Delete an invoice
// @Description Soft delete an invoice (cannot delete paid invoices)
// @Tags invoices
// @Security BearerAuth
// @Param id path string true "Invoice ID (UUID)"
// @Success 204 "Invoice deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid ID or validation error"
// @Failure 404 {object} ErrorResponse "Invoice not found"
// @Router /billing/invoices/{id} [delete]
func (h *InvoiceHandler) DeleteInvoice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid invoice ID"})
		return
	}

	if err := h.invoiceService.DeleteInvoice(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// MarkInvoiceAsPaid godoc
// @Summary Mark an invoice as paid
// @Description Update invoice status to paid
// @Tags invoices
// @Security BearerAuth
// @Produce json
// @Param id path string true "Invoice ID (UUID)"
// @Success 200 {object} domain.Invoice
// @Failure 400 {object} ErrorResponse "Invalid ID or already paid"
// @Failure 404 {object} ErrorResponse "Invoice not found"
// @Router /billing/invoices/{id}/mark-paid [post]
func (h *InvoiceHandler) MarkInvoiceAsPaid(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid invoice ID"})
		return
	}

	invoice, err := h.invoiceService.MarkAsPaid(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// GetClientInvoices godoc
// @Summary Get all invoices for a client
// @Description Retrieve all invoices for a specific client
// @Tags invoices
// @Security BearerAuth
// @Produce json
// @Param clientId path string true "Client ID (UUID)"
// @Success 200 {array} domain.Invoice
// @Failure 400 {object} ErrorResponse "Invalid client ID"
// @Failure 404 {object} ErrorResponse "Client not found"
// @Router /billing/invoices/client/{clientId} [get]
func (h *InvoiceHandler) GetClientInvoices(c *gin.Context) {
	clientIDParam := c.Param("clientId")
	clientID, err := uuid.Parse(clientIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid client ID"})
		return
	}

	invoices, err := h.invoiceService.GetClientInvoices(c.Request.Context(), clientID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, invoices)
}

// GetUnpaidInvoices godoc
// @Summary Get all unpaid invoices
// @Description Retrieve all invoices with unpaid status
// @Tags invoices
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Invoice
// @Router /billing/invoices/unpaid [get]
func (h *InvoiceHandler) GetUnpaidInvoices(c *gin.Context) {
	invoices, err := h.invoiceService.GetUnpaidInvoices(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, invoices)
}
