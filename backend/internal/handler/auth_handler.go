package handler

import (
	"net/http"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	pkgerrors "github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService service.AuthServiceInterface
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// @Summary      Register a new user
// @Description  Creates a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body service.RegisterRequest true "Registration data"
// @Success      201 {object} service.AuthResponse
// @Failure      400 {object} map[string]string
// @Failure      409 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := pkgerrors.NewValidationError("Datos de entrada inválidos", map[string][]string{
			"general": {"Verifica que todos los campos requeridos estén completos"},
		})
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		var appErr *pkgerrors.AppError
		switch err.Error() {
		case "email already registered":
			appErr = pkgerrors.NewConflictError("El email ya está registrado", pkgerrors.CodeEmailAlreadyExists)
		default:
			appErr = pkgerrors.NewInternalError("Error al crear el usuario")
		}
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login handles user login
// @Summary      Login user
// @Description  Authenticates a user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body service.LoginRequest true "Login credentials"
// @Success      200 {object} service.AuthResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := pkgerrors.NewValidationError("Datos de entrada inválidos", map[string][]string{
			"general": {"Verifica el formato del email y contraseña"},
		})
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		var appErr *pkgerrors.AppError
		switch err.Error() {
		case "invalid credentials":
			appErr = pkgerrors.NewUnauthorizedError("Email o contraseña incorrectos", pkgerrors.CodeInvalidCredentials)
		case "user account is inactive":
			appErr = pkgerrors.NewForbiddenError("Usuario inactivo. Contacta al administrador")
			appErr.Code = pkgerrors.CodeUserInactive
		default:
			appErr = pkgerrors.NewInternalError("Error al procesar el login")
		}
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Me retrieves the current user's information
// @Summary      Get current user
// @Description  Returns the authenticated user's information
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} domain.User
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		appErr := pkgerrors.NewUnauthorizedError("Usuario no autenticado")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	// Convert userID from string to UUID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		appErr := pkgerrors.NewValidationError("ID de usuario inválido", nil)
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		appErr := pkgerrors.NewNotFoundError("Usuario no encontrado")
		appErr.Code = pkgerrors.CodeUserNotFound
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, user)
}
