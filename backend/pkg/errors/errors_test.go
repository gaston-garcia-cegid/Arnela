package errors

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewValidationError(t *testing.T) {
	message := "Datos de entrada inválidos"
	details := map[string][]string{
		"email":    {"Email is required"},
		"password": {"Password is too short"},
	}

	err := NewValidationError(message, details)

	assert.Equal(t, message, err.Message)
	assert.Equal(t, CodeValidationFailed, err.Code)
	assert.Equal(t, 400, err.StatusCode)
	assert.Equal(t, details, err.Details)
}

func TestNewUnauthorizedError(t *testing.T) {
	message := "Email o contraseña incorrectos"
	err := NewUnauthorizedError(message)

	assert.Equal(t, message, err.Message)
	assert.Equal(t, CodeUnauthorized, err.Code)
	assert.Equal(t, 401, err.StatusCode)
}

func TestNewForbiddenError(t *testing.T) {
	message := "Usuario inactivo. Contacta al administrador"
	err := NewForbiddenError(message)

	assert.Equal(t, message, err.Message)
	assert.Equal(t, CodeForbidden, err.Code)
	assert.Equal(t, 403, err.StatusCode)
}

func TestNewNotFoundError(t *testing.T) {
	message := "Usuario no encontrado"
	err := NewNotFoundError(message)

	assert.Equal(t, message, err.Message)
	assert.Equal(t, CodeNotFound, err.Code)
	assert.Equal(t, 404, err.StatusCode)
}

func TestNewConflictError(t *testing.T) {
	message := "El email ya está registrado"
	code := CodeEmailAlreadyExists

	err := NewConflictError(message, code)

	assert.Equal(t, message, err.Message)
	assert.Equal(t, code, err.Code)
	assert.Equal(t, 409, err.StatusCode)
}

func TestNewInternalError(t *testing.T) {
	message := "Error interno del servidor"
	err := NewInternalError(message)

	assert.Equal(t, message, err.Message)
	assert.Equal(t, CodeInternalError, err.Code)
	assert.Equal(t, 500, err.StatusCode)
}

func TestRespondWithAppError(t *testing.T) {
	// Setup Gin test context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test error
	testError := &AppError{
		Message:    "Test error message",
		Code:       "TEST_ERROR",
		StatusCode: 400,
		Details: map[string][]string{
			"field": {"Error detail"},
		},
	}

	// Call the function
	RespondWithAppError(c, testError)

	// Assertions
	assert.Equal(t, 400, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Test error message", response.Error)
	assert.Equal(t, "TEST_ERROR", response.Code)
	assert.Equal(t, testError.Details, response.Details)
}

func TestRespondWithValidationError(t *testing.T) {
	// Setup Gin test context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create validation details
	details := map[string][]string{
		"email":    {"Email is required", "Email must be valid"},
		"password": {"Password is too short"},
	}

	// Call the function
	RespondWithValidationError(c, "Validation failed", details)

	// Assertions
	assert.Equal(t, 400, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Validation failed", response.Error)
	assert.Equal(t, CodeValidationFailed, response.Code)
	assert.Equal(t, details, response.Details)
}

func TestErrorCodes(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{"Validation Failed", CodeValidationFailed},
		{"Unauthorized", CodeUnauthorized},
		{"Forbidden", CodeForbidden},
		{"Not Found", CodeNotFound},
		{"Conflict", CodeConflict},
		{"Internal Error", CodeInternalError},
		{"Email Already Exists", CodeEmailAlreadyExists},
		{"Invalid Credentials", CodeInvalidCredentials},
		{"User Inactive", CodeUserInactive},
		{"User Not Found", CodeUserNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.code)
		})
	}
}
