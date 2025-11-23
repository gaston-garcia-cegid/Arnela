package middleware

import (
	"net/http"
	"strings"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware handles JWT authentication and authorization
type AuthMiddleware struct {
	tokenManager *jwt.TokenManager
	clientRepo   repository.ClientRepository
}

// NewAuthMiddleware creates a new AuthMiddleware instance
func NewAuthMiddleware(tokenManager *jwt.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
		clientRepo:   nil,
	}
}

// SetClientRepo sets the client repository for the middleware
func (m *AuthMiddleware) SetClientRepo(clientRepo repository.ClientRepository) {
	m.clientRepo = clientRepo
}

// RequireAuth ensures the user is authenticated
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := m.tokenManager.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Parse userID as UUID
		userID, err := uuid.Parse(claims.UserID.String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
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

		c.Next()
	}
}

// RequireRole creates a middleware that checks if the user has one of the specified roles
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - no role in context"})
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role type"})
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

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
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

		// Set user information in context if token is valid
		c.Set("userID", claims.UserID.String())
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}
