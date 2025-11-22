package errors

import "github.com/gin-gonic/gin"

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string              `json:"error"`
	Code    string              `json:"code,omitempty"`
	Details map[string][]string `json:"details,omitempty"`
}

// Error codes
const (
	CodeValidationFailed   = "VALIDATION_FAILED"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeForbidden          = "FORBIDDEN"
	CodeNotFound           = "NOT_FOUND"
	CodeConflict           = "CONFLICT"
	CodeInternalError      = "INTERNAL_ERROR"
	CodeEmailAlreadyExists = "EMAIL_ALREADY_EXISTS"
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
)

// NewErrorResponse creates a new error response
func NewErrorResponse(message, code string) ErrorResponse {
	return ErrorResponse{
		Error: message,
		Code:  code,
	}
}

// NewValidationError creates a validation error response
func NewValidationError(details map[string][]string) ErrorResponse {
	return ErrorResponse{
		Error:   "Validation failed",
		Code:    CodeValidationFailed,
		Details: details,
	}
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, statusCode int, message, code string) {
	c.JSON(statusCode, NewErrorResponse(message, code))
}

// RespondWithValidationError sends a validation error response
func RespondWithValidationError(c *gin.Context, details map[string][]string) {
	c.JSON(400, NewValidationError(details))
}
