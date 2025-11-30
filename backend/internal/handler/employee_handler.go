package handler

import (
	"net/http"
	"strconv"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	pkgerrors "github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// EmployeeHandler handles employee-related endpoints
type EmployeeHandler struct {
	employeeService service.EmployeeService
}

// NewEmployeeHandler creates a new EmployeeHandler
func NewEmployeeHandler(employeeService service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
	}
}

// CreateEmployee creates a new employee
// @Summary      Create a new employee
// @Description  Creates a new employee/therapist (admin only)
// @Tags         employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body service.CreateEmployeeRequest true "Employee data"
// @Success      201 {object} domain.Employee
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      409 {object} map[string]string
// @Router       /api/v1/employees [post]
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req service.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := pkgerrors.NewValidationError("Datos de entrada inválidos", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	employee, err := h.employeeService.CreateEmployee(c.Request.Context(), req)
	if err != nil {
		// Check for specific errors
		switch err {
		case service.ErrInvalidEmail:
			appErr := pkgerrors.NewValidationError("Formato de email inválido", nil)
			pkgerrors.RespondWithAppError(c, appErr)
			return
		case service.ErrInvalidPhone:
			appErr := pkgerrors.NewValidationError("Formato de teléfono inválido (debe ser formato español)", nil)
			pkgerrors.RespondWithAppError(c, appErr)
			return
		case service.ErrInvalidDNI:
			appErr := pkgerrors.NewValidationError("Formato de DNI inválido", nil)
			pkgerrors.RespondWithAppError(c, appErr)
			return
		case service.ErrEmailInUse:
			appErr := pkgerrors.NewConflictError("El email ya está registrado", pkgerrors.CodeEmailAlreadyExists)
			pkgerrors.RespondWithAppError(c, appErr)
			return
		case service.ErrDNIInUse:
			appErr := pkgerrors.NewConflictError("El DNI ya está registrado", pkgerrors.CodeDNIAlreadyExists)
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		appErr := pkgerrors.NewInternalError("Error al crear el empleado")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusCreated, employee)
}

// GetEmployee retrieves an employee by ID
// @Summary      Get employee by ID
// @Description  Retrieves an employee's information (admin/employee)
// @Tags         employees
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Employee ID"
// @Success      200 {object} domain.Employee
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/v1/employees/{id} [get]
func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		appErr := pkgerrors.NewValidationError("ID inválido", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	employee, err := h.employeeService.GetEmployee(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrEmployeeNotFound {
			appErr := pkgerrors.NewNotFoundError("Empleado no encontrado")
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		appErr := pkgerrors.NewInternalError("Error al obtener el empleado")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, employee)
}

// GetMyEmployee retrieves the employee profile for the current logged-in user
// @Summary      Get my employee profile
// @Description  Retrieves the employee profile associated with the current user (employee only)
// @Tags         employees
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} domain.Employee
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/v1/employees/me [get]
func (h *EmployeeHandler) GetMyEmployee(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		appErr := pkgerrors.NewUnauthorizedError("Usuario no autenticado")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		appErr := pkgerrors.NewInternalError("Error al obtener ID de usuario")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	employee, err := h.employeeService.GetEmployeeByUserID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrEmployeeNotFound {
			appErr := pkgerrors.NewNotFoundError("Empleado no encontrado")
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		appErr := pkgerrors.NewInternalError("Error al obtener el empleado")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, employee)
}

// ListEmployees retrieves a paginated list of employees
// @Summary      List employees
// @Description  Retrieves a list of employees with pagination (admin/employee)
// @Tags         employees
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Limit" default(10)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/v1/employees [get]
func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	employees, total, err := h.employeeService.ListEmployees(c.Request.Context(), limit, offset)
	if err != nil {
		appErr := pkgerrors.NewInternalError("Error al listar empleados")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"employees": employees,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

// UpdateEmployee updates an employee's information
// @Summary      Update employee
// @Description  Updates an employee's information (admin only)
// @Tags         employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Employee ID"
// @Param        request body service.UpdateEmployeeRequest true "Employee update data"
// @Success      200 {object} domain.Employee
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      409 {object} map[string]string
// @Router       /api/v1/employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		appErr := pkgerrors.NewValidationError("ID inválido", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	var req service.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := pkgerrors.NewValidationError("Datos de entrada inválidos", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	employee, err := h.employeeService.UpdateEmployee(c.Request.Context(), id, req)
	if err != nil {
		switch err {
		case service.ErrEmployeeNotFound:
			appErr := pkgerrors.NewNotFoundError("Empleado no encontrado")
			pkgerrors.RespondWithAppError(c, appErr)
			return
		case service.ErrInvalidEmail:
			appErr := pkgerrors.NewValidationError("Formato de email inválido", nil)
			pkgerrors.RespondWithAppError(c, appErr)
			return
		case service.ErrInvalidPhone:
			appErr := pkgerrors.NewValidationError("Formato de teléfono inválido", nil)
			pkgerrors.RespondWithAppError(c, appErr)
			return
		case service.ErrInvalidDNI:
			appErr := pkgerrors.NewValidationError("Formato de DNI inválido", nil)
			pkgerrors.RespondWithAppError(c, appErr)
			return
		case service.ErrEmailInUse:
			appErr := pkgerrors.NewConflictError("El email ya está registrado", pkgerrors.CodeEmailAlreadyExists)
			pkgerrors.RespondWithAppError(c, appErr)
			return
		case service.ErrDNIInUse:
			appErr := pkgerrors.NewConflictError("El DNI ya está registrado", pkgerrors.CodeDNIAlreadyExists)
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		appErr := pkgerrors.NewInternalError("Error al actualizar el empleado")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, employee)
}

// DeleteEmployee soft-deletes an employee
// @Summary      Delete employee
// @Description  Soft-deletes an employee (admin only)
// @Tags         employees
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Employee ID"
// @Success      204
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/v1/employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		appErr := pkgerrors.NewValidationError("ID inválido", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	err = h.employeeService.DeleteEmployee(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrEmployeeNotFound {
			appErr := pkgerrors.NewNotFoundError("Empleado no encontrado")
			pkgerrors.RespondWithAppError(c, appErr)
			return
		}
		appErr := pkgerrors.NewInternalError("Error al eliminar el empleado")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetEmployeesBySpecialty retrieves employees with a specific specialty
// @Summary      Get employees by specialty
// @Description  Retrieves employees with a specific specialty (admin/employee/client)
// @Tags         employees
// @Produce      json
// @Security     BearerAuth
// @Param        specialty path string true "Specialty name"
// @Success      200 {array} domain.Employee
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/v1/employees/specialty/{specialty} [get]
func (h *EmployeeHandler) GetEmployeesBySpecialty(c *gin.Context) {
	specialty := c.Param("specialty")

	employees, err := h.employeeService.GetEmployeesBySpecialty(c.Request.Context(), specialty)
	if err != nil {
		appErr := pkgerrors.NewInternalError("Error al obtener empleados por especialidad")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, employees)
}
