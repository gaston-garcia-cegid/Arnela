package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	pkgerrors "github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AppointmentHandler handles appointment-related endpoints
type AppointmentHandler struct {
	appointmentService service.AppointmentServiceInterface
}

// NewAppointmentHandler creates a new AppointmentHandler
func NewAppointmentHandler(appointmentService service.AppointmentServiceInterface) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}

// CreateAppointment creates a new appointment
// @Summary      Create a new appointment
// @Description  Creates a new appointment (pending status)
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreateAppointmentRequest true "Appointment data"
// @Success      201 {object} domain.Appointment
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /api/v1/appointments [post]
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var req domain.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := pkgerrors.NewValidationError("Datos de entrada inválidos", map[string][]string{
			"general": {"Verifica todos los campos requeridos"},
		})
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		appErr := pkgerrors.NewUnauthorizedError("Usuario no autenticado", pkgerrors.CodeUnauthorized)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	appointment, err := h.appointmentService.CreateAppointment(c.Request.Context(), req, userID.(uuid.UUID))
	if err != nil {
		appErr := pkgerrors.NewValidationError(err.Error(), nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

// GetAppointment retrieves an appointment by ID
// @Summary      Get appointment by ID
// @Description  Retrieves an appointment's information
// @Tags         appointments
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Appointment ID"
// @Success      200 {object} domain.Appointment
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/v1/appointments/{id} [get]
func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		appErr := pkgerrors.NewValidationError("ID de cita inválido", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	appointment, err := h.appointmentService.GetAppointment(c.Request.Context(), id)
	if err != nil {
		appErr := pkgerrors.NewNotFoundError("Cita no encontrada")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, appointment)
}

// UpdateAppointment updates an appointment
// @Summary      Update appointment
// @Description  Updates an appointment (client can only update their own)
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Appointment ID"
// @Param        request body domain.UpdateAppointmentRequest true "Update data"
// @Success      200 {object} domain.Appointment
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/v1/appointments/{id} [put]
func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		appErr := pkgerrors.NewValidationError("ID de cita inválido", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	var req domain.UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := pkgerrors.NewValidationError("Datos de entrada inválidos", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		appErr := pkgerrors.NewUnauthorizedError("Usuario no autenticado", pkgerrors.CodeUnauthorized)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	appointment, err := h.appointmentService.UpdateAppointment(c.Request.Context(), id, req, userID.(uuid.UUID))
	if err != nil {
		if err.Error() == "cita no encontrada" {
			appErr := pkgerrors.NewNotFoundError("Cita no encontrada")
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		if err.Error() == "no tienes permiso para modificar esta cita" {
			appErr := pkgerrors.NewForbiddenError("No tienes permiso para modificar esta cita")
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		appErr := pkgerrors.NewValidationError(err.Error(), nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, appointment)
}

// CancelAppointment cancels an appointment
// @Summary      Cancel appointment
// @Description  Cancels an appointment
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Appointment ID"
// @Param        request body domain.CancelAppointmentRequest true "Cancellation data"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/v1/appointments/{id}/cancel [post]
func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		appErr := pkgerrors.NewValidationError("ID de cita inválido", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	var req domain.CancelAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := pkgerrors.NewValidationError("Datos de entrada inválidos", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		appErr := pkgerrors.NewUnauthorizedError("Usuario no autenticado", pkgerrors.CodeUnauthorized)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	// Check if user is admin/employee
	userRole, _ := c.Get("userRole")
	isAdmin := userRole == string(domain.RoleAdmin) || userRole == string(domain.RoleEmployee)

	err = h.appointmentService.CancelAppointment(c.Request.Context(), id, req, userID.(uuid.UUID), isAdmin)
	if err != nil {
		if err.Error() == "cita no encontrada" {
			appErr := pkgerrors.NewNotFoundError("Cita no encontrada")
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		if err.Error() == "no tienes permiso para cancelar esta cita" {
			appErr := pkgerrors.NewForbiddenError("No tienes permiso para cancelar esta cita")
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		appErr := pkgerrors.NewValidationError(err.Error(), nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cita cancelada exitosamente"})
}

// GetMyAppointments retrieves appointments for the authenticated client
// @Summary      Get my appointments
// @Description  Retrieves all appointments for the authenticated client
// @Tags         appointments
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        pageSize query int false "Page size" default(10)
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]string
// @Router       /api/v1/appointments/me [get]
func (h *AppointmentHandler) GetMyAppointments(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		appErr := pkgerrors.NewUnauthorizedError("Usuario no autenticado", pkgerrors.CodeUnauthorized)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	clientID, exists := c.Get("clientID")
	if !exists {
		appErr := pkgerrors.NewForbiddenError("Solo clientes pueden acceder a sus citas")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	appointments, total, err := h.appointmentService.GetMyAppointments(c.Request.Context(), clientID.(uuid.UUID), page, pageSize)
	if err != nil {
		appErr := pkgerrors.NewInternalError("Error al obtener las citas")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"appointments": appointments,
		"total":        total,
		"page":         page,
		"pageSize":     pageSize,
	})
}

// ConfirmAppointment confirms a pending appointment (admin/employee only)
// @Summary      Confirm appointment
// @Description  Confirms a pending appointment (admin/employee only)
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Appointment ID"
// @Param        request body domain.ConfirmAppointmentRequest true "Confirmation data"
// @Success      200 {object} domain.Appointment
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/v1/appointments/{id}/confirm [post]
func (h *AppointmentHandler) ConfirmAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		appErr := pkgerrors.NewValidationError("ID de cita inválido", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	var req domain.ConfirmAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := pkgerrors.NewValidationError("Datos de entrada inválidos", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	appointment, err := h.appointmentService.ConfirmAppointment(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "cita no encontrada" {
			appErr := pkgerrors.NewNotFoundError("Cita no encontrada")
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		appErr := pkgerrors.NewValidationError(err.Error(), nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, appointment)
}

// ListAppointments lists all appointments with filters (admin/employee only)
// @Summary      List all appointments
// @Description  Lists all appointments with optional filters (admin/employee only)
// @Tags         appointments
// @Produce      json
// @Security     BearerAuth
// @Param        clientId query string false "Filter by client ID"
// @Param        therapistId query string false "Filter by therapist ID"
// @Param        status query string false "Filter by status"
// @Param        startDate query string false "Filter by start date (RFC3339)"
// @Param        endDate query string false "Filter by end date (RFC3339)"
// @Param        page query int false "Page number" default(1)
// @Param        pageSize query int false "Page size" default(20)
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /api/v1/appointments [get]
func (h *AppointmentHandler) ListAppointments(c *gin.Context) {
	filters := domain.AppointmentFilter{
		Page:     1,
		PageSize: 20,
	}

	// Parse query parameters
	if clientIDStr := c.Query("clientId"); clientIDStr != "" {
		if clientID, err := uuid.Parse(clientIDStr); err == nil {
			filters.ClientID = &clientID
		}
	}

	if therapistID := c.Query("therapistId"); therapistID != "" {
		filters.TherapistID = &therapistID
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.AppointmentStatus(statusStr)
		filters.Status = &status
	}

	if startDateStr := c.Query("startDate"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			filters.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("endDate"); endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			filters.EndDate = &endDate
		}
	}

	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filters.Page = page
	}

	if pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20")); err == nil {
		filters.PageSize = pageSize
	}

	appointments, total, err := h.appointmentService.ListAppointments(c.Request.Context(), filters)
	if err != nil {
		appErr := pkgerrors.NewInternalError("Error al listar las citas")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"appointments": appointments,
		"total":        total,
		"page":         filters.Page,
		"pageSize":     filters.PageSize,
	})
}

// GetTherapists returns all available therapists
// @Summary      Get therapists
// @Description  Returns all available therapists
// @Tags         appointments
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Router       /api/v1/appointments/therapists [get]
func (h *AppointmentHandler) GetTherapists(c *gin.Context) {
	therapists := h.appointmentService.GetTherapists(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"therapists": therapists})
}

// GetAvailableSlots returns available time slots for a therapist
// @Summary      Get available slots
// @Description  Returns available time slots for a therapist on a specific date
// @Tags         appointments
// @Produce      json
// @Security     BearerAuth
// @Param        therapistId query string true "Therapist ID"
// @Param        date query string true "Date (YYYY-MM-DD)"
// @Param        duration query int true "Duration in minutes (45 or 60)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Router       /api/v1/appointments/available-slots [get]
func (h *AppointmentHandler) GetAvailableSlots(c *gin.Context) {
	therapistID := c.Query("therapistId")
	dateStr := c.Query("date")
	durationStr := c.Query("duration")

	if therapistID == "" || dateStr == "" || durationStr == "" {
		appErr := pkgerrors.NewValidationError("Parámetros faltantes", map[string][]string{
			"general": {"Se requiere therapistId, date y duration"},
		})
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		appErr := pkgerrors.NewValidationError("Formato de fecha inválido", map[string][]string{
			"date": {"Usa formato YYYY-MM-DD"},
		})
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		appErr := pkgerrors.NewValidationError("Duración inválida", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	slots, err := h.appointmentService.GetAvailableSlots(c.Request.Context(), therapistID, date, duration)
	if err != nil {
		appErr := pkgerrors.NewValidationError(err.Error(), nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"slots": slots})
}
