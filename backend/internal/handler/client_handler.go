package handler

import (
	"net/http"
	"strconv"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	pkgerrors "github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ClientHandler handles client-related endpoints
type ClientHandler struct {
	clientService service.ClientServiceInterface
}

// NewClientHandler creates a new ClientHandler
func NewClientHandler(clientService service.ClientServiceInterface) *ClientHandler {
	return &ClientHandler{
		clientService: clientService,
	}
}

// CreateClient creates a new client
// @Summary      Create a new client
// @Description  Creates a new client (admin/employee only)
// @Tags         clients
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body service.CreateClientRequest true "Client data"
// @Success      201 {object} domain.Client
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      409 {object} map[string]string
// @Router       /api/v1/clients [post]
func (h *ClientHandler) CreateClient(c *gin.Context) {
	var req service.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := pkgerrors.NewValidationError("Datos de entrada inválidos", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	client, err := h.clientService.CreateClient(c.Request.Context(), req)
	if err != nil {
		// Check for specific errors
		errMsg := err.Error()
		if errMsg == "email already registered" {
			appErr := pkgerrors.NewConflictError("El email ya está registrado", pkgerrors.CodeEmailAlreadyExists)
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		if errMsg == "NIF already registered" || errMsg == "DNI already registered" {
			appErr := pkgerrors.NewConflictError("El DNI/NIF ya está registrado", "DNI_ALREADY_EXISTS")
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		appErr := pkgerrors.NewInternalError("Error al crear el cliente")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusCreated, client)
}

// GetClient retrieves a client by ID
// @Summary      Get client by ID
// @Description  Retrieves a client's information (admin/employee only)
// @Tags         clients
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Client ID"
// @Success      200 {object} domain.Client
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/v1/clients/{id} [get]
func (h *ClientHandler) GetClient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		appErr := pkgerrors.NewValidationError("ID de cliente inválido", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	client, err := h.clientService.GetClient(c.Request.Context(), id)
	if err != nil {
		appErr := pkgerrors.NewNotFoundError("Cliente no encontrado")
		appErr.Code = "CLIENT_NOT_FOUND"
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, client)
}

// UpdateClient updates a client's information
// @Summary      Update client
// @Description  Updates a client's information (admin/employee only)
// @Tags         clients
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Client ID"
// @Param        request body service.UpdateClientRequest true "Client update data"
// @Success      200 {object} domain.Client
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      409 {object} map[string]string
// @Router       /api/v1/clients/{id} [put]
func (h *ClientHandler) UpdateClient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		appErr := pkgerrors.NewValidationError("ID de cliente inválido", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	var req service.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := pkgerrors.NewValidationError("Datos de entrada inválidos", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	client, err := h.clientService.UpdateClient(c.Request.Context(), id, req)
	if err != nil {
		// Check for specific errors
		if err.Error() == "email already registered" {
			appErr := pkgerrors.NewConflictError("El email ya está registrado", pkgerrors.CodeEmailAlreadyExists)
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		if err.Error() == "DNI already registered" {
			appErr := pkgerrors.NewConflictError("El DNI ya está registrado", "DNI_ALREADY_EXISTS")
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		if err.Error() == "client not found" {
			appErr := pkgerrors.NewNotFoundError("Cliente no encontrado")
			appErr.Code = "CLIENT_NOT_FOUND"
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		appErr := pkgerrors.NewInternalError("Error al actualizar el cliente")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, client)
}

// DeleteClient soft-deletes a client
// @Summary      Delete client
// @Description  Soft-deletes a client (admin only)
// @Tags         clients
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Client ID"
// @Success      204 "No Content"
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/v1/clients/{id} [delete]
func (h *ClientHandler) DeleteClient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		appErr := pkgerrors.NewValidationError("ID de cliente inválido", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	err = h.clientService.DeleteClient(c.Request.Context(), id)
	if err != nil {
		appErr := pkgerrors.NewNotFoundError("Cliente no encontrado")
		appErr.Code = "CLIENT_NOT_FOUND"
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListClients retrieves a paginated list of clients
// @Summary      List clients
// @Description  Retrieves a paginated list of clients with optional filters (admin/employee only)
// @Tags         clients
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        pageSize query int false "Page size" default(20)
// @Param        search query string false "Search in name, email, phone, DNI"
// @Param        isActive query boolean false "Filter by active status"
// @Param        city query string false "Filter by city"
// @Param        province query string false "Filter by province"
// @Success      200 {object} service.ClientListResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /api/v1/clients [get]
func (h *ClientHandler) ListClients(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// Parse filters
	filters := repository.ClientFilters{
		Search:   c.Query("search"),
		City:     c.Query("city"),
		Province: c.Query("province"),
	}

	// Parse isActive filter
	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filters.IsActive = &isActive
	}

	response, err := h.clientService.ListClients(c.Request.Context(), filters, page, pageSize)
	if err != nil {
		appErr := pkgerrors.NewInternalError("Error al listar clientes")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMyClient retrieves the current client's information (for client role)
// @Summary      Get my client info
// @Description  Retrieves the authenticated client's information (client role only)
// @Tags         clients
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} domain.Client
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/v1/clients/me [get]
func (h *ClientHandler) GetMyClient(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		appErr := pkgerrors.NewUnauthorizedError("Usuario no autenticado")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		appErr := pkgerrors.NewValidationError("ID de usuario inválido", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	// Get client by user ID
	client, err := h.clientService.GetClientByUserID(c.Request.Context(), userID)
	if err != nil {
		appErr := pkgerrors.NewNotFoundError("Perfil de cliente no encontrado")
		appErr.Code = "CLIENT_NOT_FOUND"
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, client)
}
