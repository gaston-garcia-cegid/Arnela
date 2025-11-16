package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthServiceInterface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req service.RegisterRequest) (*service.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req service.LoginRequest) (*service.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.AuthResponse), args.Error(1)
}

func (m *MockAuthService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful registration",
			requestBody: map[string]interface{}{
				"email":     "john@example.com",
				"password":  "Password123!",
				"firstName": "John",
				"lastName":  "Doe",
				"role":      "client",
			},
			mockSetup: func(m *MockAuthService) {
				validUser := &domain.User{
					ID:        uuid.New(),
					Email:     "john@example.com",
					FirstName: "John",
					LastName:  "Doe",
					Role:      domain.RoleClient,
					IsActive:  true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				response := &service.AuthResponse{
					Token: "test-token",
					User:  validUser,
				}
				m.On("Register", mock.Anything, mock.AnythingOfType("service.RegisterRequest")).Return(response, nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response service.AuthResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Token)
				assert.Equal(t, "john@example.com", response.User.Email)
			},
		},
		{
			name: "missing required fields",
			requestBody: map[string]interface{}{
				"email": "john@example.com",
			},
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "email already exists",
			requestBody: map[string]interface{}{
				"email":     "existing@example.com",
				"password":  "Password123!",
				"firstName": "Jane",
				"lastName":  "Smith",
				"role":      "client",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Register", mock.Anything, mock.AnythingOfType("service.RegisterRequest")).Return(nil, errors.New("email already registered"))
			},
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], "email already registered")
			},
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid-json",
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.mockSetup(mockService)

			handler := NewAuthHandler(mockService)
			router := setupTestRouter()
			router.POST("/register", handler.Register)

			var reqBody []byte
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			tt.checkResponse(t, rec)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful login",
			requestBody: map[string]interface{}{
				"email":    "john@example.com",
				"password": "Password123!",
			},
			mockSetup: func(m *MockAuthService) {
				validUser := &domain.User{
					ID:        uuid.New(),
					Email:     "john@example.com",
					FirstName: "John",
					LastName:  "Doe",
					Role:      domain.RoleClient,
					IsActive:  true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				response := &service.AuthResponse{
					Token: "test-token",
					User:  validUser,
				}
				m.On("Login", mock.Anything, mock.AnythingOfType("service.LoginRequest")).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response service.AuthResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Token)
				assert.Equal(t, "john@example.com", response.User.Email)
			},
		},
		{
			name: "invalid credentials",
			requestBody: map[string]interface{}{
				"email":    "john@example.com",
				"password": "WrongPassword",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Login", mock.Anything, mock.AnythingOfType("service.LoginRequest")).Return(nil, errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], "invalid credentials")
			},
		},
		{
			name: "missing password",
			requestBody: map[string]interface{}{
				"email": "john@example.com",
			},
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.mockSetup(mockService)

			handler := NewAuthHandler(mockService)
			router := setupTestRouter()
			router.POST("/login", handler.Login)

			var reqBody []byte
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			tt.checkResponse(t, rec)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Me(t *testing.T) {
	validUser := &domain.User{
		ID:        uuid.New(),
		Email:     "john@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      domain.RoleClient,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		mockSetup      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful retrieval",
			setupContext: func(c *gin.Context) {
				c.Set("userID", validUser.ID.String())
			},
			mockSetup: func(m *MockAuthService) {
				m.On("GetUserByID", mock.Anything, validUser.ID).Return(validUser, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response domain.User
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "john@example.com", response.Email)
			},
		},
		{
			name: "missing user ID in context",
			setupContext: func(c *gin.Context) {
				// Don't set userID
			},
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "user not found",
			setupContext: func(c *gin.Context) {
				c.Set("userID", uuid.New().String())
			},
			mockSetup: func(m *MockAuthService) {
				m.On("GetUserByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], "user not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.mockSetup(mockService)

			handler := NewAuthHandler(mockService)
			router := setupTestRouter()

			// Middleware to setup context
			router.Use(func(c *gin.Context) {
				tt.setupContext(c)
				c.Next()
			})

			router.GET("/me", handler.Me)

			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			tt.checkResponse(t, rec)
			mockService.AssertExpectations(t)
		})
	}
}
