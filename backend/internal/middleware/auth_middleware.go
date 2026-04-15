package middleware

import (
	"net/http"
	"strings"

	pkgerrors "github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware handles JWT authentication and authorization
type AuthMiddleware struct {
	tokenManager *jwt.TokenManager
	clientRepo   repository.ClientRepository
	employeeRepo repository.EmployeeRepository
}

// NewAuthMiddleware creates a new AuthMiddleware instance
func NewAuthMiddleware(tokenManager *jwt.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
		clientRepo:   nil,
		employeeRepo: nil,
	}
}

// SetClientRepo sets the client repository for the middleware
func (m *AuthMiddleware) SetClientRepo(clientRepo repository.ClientRepository) {
	m.clientRepo = clientRepo
}

// SetEmployeeRepo sets the employee repository for the middleware
func (m *AuthMiddleware) SetEmployeeRepo(employeeRepo repository.EmployeeRepository) {
	m.employeeRepo = employeeRepo
}

// RequireAuth ensures the user is authenticated
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			pkgerrors.RespondWithError(c, http.StatusUnauthorized, "Authorization header required", pkgerrors.CodeUnauthorized)
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			pkgerrors.RespondWithError(c, http.StatusUnauthorized, "Invalid authorization header format", pkgerrors.CodeUnauthorized)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := m.tokenManager.ValidateToken(tokenString)
		if err != nil {
			pkgerrors.RespondWithError(c, http.StatusUnauthorized, "Invalid or expired token", pkgerrors.CodeUnauthorized)
			c.Abort()
			return
		}

		// Parse userID as UUID
		userID, err := uuid.Parse(claims.UserID.String())
		if err != nil {
			pkgerrors.RespondWithError(c, http.StatusInternalServerError, "Invalid user ID format", pkgerrors.CodeInternalError)
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("userID", userID)
		c.Set("userRole", claims.Role)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role) // Keep for backward compatibility

		// If user is a client, fetch and set clientID
		if claims.Role == "client" && m.clientRepo != nil {
			client, err := m.clientRepo.GetByUserID(c.Request.Context(), userID)
			if err == nil && client != nil {
				c.Set("clientID", client.ID)
			}
		}

		// If user is an employee, fetch and set employeeID
		if claims.Role == "employee" && m.employeeRepo != nil {
			employee, err := m.employeeRepo.GetByUserID(c.Request.Context(), userID)
			if err == nil && employee != nil {
				c.Set("employeeID", employee.ID)
			}
		}

		c.Next()
	}
}

// RequireRole creates a middleware that checks if the user has one of the specified roles
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			pkgerrors.RespondWithError(c, http.StatusUnauthorized, "Unauthorized - no role in context", pkgerrors.CodeUnauthorized)
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			pkgerrors.RespondWithError(c, http.StatusInternalServerError, "Invalid role type", pkgerrors.CodeInternalError)
			c.Abort()
			return
		}

		// Check if user has one of the allowed roles
		for _, role := range roles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		pkgerrors.RespondWithError(c, http.StatusForbidden, "Insufficient permissions", pkgerrors.CodeForbidden)
		c.Abort()
	}
}

// OptionalAuth validates the token if present, but doesn't require it
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := m.tokenManager.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		userID, parseErr := uuid.Parse(claims.UserID.String())
		if parseErr != nil {
			c.Next()
			return
		}

		c.Set("userID", userID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}
