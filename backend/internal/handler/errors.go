package handler

import (
	"net/http"

	pkgerrors "github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// PaginatedResponse represents a paginated list response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
}

// handleError converts service errors to appropriate HTTP responses
func handleError(c *gin.Context, err error) {
	// Check if it's an AppError
	appErr, ok := err.(*pkgerrors.AppError)
	if ok {
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	// Default to internal server error
	c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
}
