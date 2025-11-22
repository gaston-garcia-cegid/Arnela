package errors

import (
	"github.com/gin-gonic/gin"
)

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
	CodeDNIAlreadyExists   = "DNI_ALREADY_EXISTS"
	CodeNIFAlreadyExists   = "NIF_ALREADY_EXISTS"
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeUserInactive       = "USER_INACTIVE"
	CodeUserNotFound       = "USER_NOT_FOUND"
)

// AppError represents an application-level error with HTTP status
type AppError struct {
	Message    string
	Code       string
	StatusCode int
	Details    map[string][]string
}

func (e *AppError) Error() string {
	return e.Message
}

// Constructor functions for AppError
func NewValidationError(message string, details map[string][]string) *AppError {
	return &AppError{
		Message:    message,
		Code:       CodeValidationFailed,
		StatusCode: 400,
		Details:    details,
	}
}

func NewUnauthorizedError(message string, code ...string) *AppError {
	errCode := CodeUnauthorized
	if len(code) > 0 {
		errCode = code[0]
	}
	return &AppError{
		Message:    message,
		Code:       errCode,
		StatusCode: 401,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Message:    message,
		Code:       CodeForbidden,
		StatusCode: 403,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Message:    message,
		Code:       CodeNotFound,
		StatusCode: 404,
	}
}

func NewConflictError(message string, code string) *AppError {
	return &AppError{
		Message:    message,
		Code:       code,
		StatusCode: 409,
	}
}

func NewInternalError(message string) *AppError {
	return &AppError{
		Message:    message,
		Code:       CodeInternalError,
		StatusCode: 500,
	}
}

// Helper functions for creating ErrorResponse
func NewErrorResponse(message, code string) ErrorResponse {
	return ErrorResponse{
		Error: message,
		Code:  code,
	}
}

func NewErrorResponseWithDetails(message, code string, details map[string][]string) ErrorResponse {
	return ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
	}
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, statusCode int, message, code string) {
	c.JSON(statusCode, NewErrorResponse(message, code))
}

// RespondWithAppError sends an AppError as response
func RespondWithAppError(c *gin.Context, err *AppError) {
	response := ErrorResponse{
		Error:   err.Message,
		Code:    err.Code,
		Details: err.Details,
	}
	c.JSON(err.StatusCode, response)
}

// RespondWithValidationError sends a validation error response
func RespondWithValidationError(c *gin.Context, message string, details map[string][]string) {
	c.JSON(400, NewErrorResponseWithDetails(message, CodeValidationFailed, details))
}
