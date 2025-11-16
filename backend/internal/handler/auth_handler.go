package handler

import (
	"net/http"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
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
// @Failure      500 {object} map[string]string
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
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
// @Router       /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert userID from string to UUID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
