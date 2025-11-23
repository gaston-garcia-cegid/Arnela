package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository/mocks"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name          string
		request       RegisterRequest
		mockSetup     func(*mocks.MockUserRepository)
		expectedError string
	}{
		{
			name: "successful registration",
			request: RegisterRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "Password123!",
				Role:      string(domain.RoleClient),
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.On("EmailExists", mock.Anything, "john@example.com").Return(false, nil)
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "email already exists",
			request: RegisterRequest{
				FirstName: "Jane",
				LastName:  "Smith",
				Email:     "existing@example.com",
				Password:  "Password123!",
				Role:      string(domain.RoleClient),
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.On("EmailExists", mock.Anything, "existing@example.com").Return(true, nil)
			},
			expectedError: "email already registered",
		},
		{
			name: "database error on email check",
			request: RegisterRequest{
				FirstName: "Test",
				LastName:  "User",
				Email:     "test@example.com",
				Password:  "Password123!",
				Role:      string(domain.RoleClient),
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.On("EmailExists", mock.Anything, "test@example.com").Return(false, errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockUserRepository)
			mockClientRepo := new(mocks.MockClientRepository)
			tt.mockSetup(mockRepo)

			tokenManager := jwt.NewTokenManager("test-secret", "test-issuer")
			authService := NewAuthService(mockRepo, mockClientRepo, tokenManager, 24*time.Hour)

			ctx := context.Background()
			response, err := authService.Register(ctx, tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Token)
				assert.NotNil(t, response.User)
				assert.Equal(t, tt.request.Email, response.User.Email)
				assert.Equal(t, tt.request.FirstName, response.User.FirstName)
				assert.Equal(t, tt.request.LastName, response.User.LastName)
				assert.Equal(t, domain.UserRole(tt.request.Role), response.User.Role)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	// Create a valid password hash for testing
	validPassword := "Password123!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)

	validUser := &domain.User{
		ID:           uuid.New(),
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		PasswordHash: string(hashedPassword),
		Role:         domain.RoleClient,
		IsActive:     true, // Active user
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tests := []struct {
		name          string
		request       LoginRequest
		mockSetup     func(*mocks.MockUserRepository)
		expectedError string
	}{
		{
			name: "successful login",
			request: LoginRequest{
				Email:    "john@example.com",
				Password: validPassword,
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "john@example.com").Return(validUser, nil)
			},
			expectedError: "",
		},
		{
			name: "user not found",
			request: LoginRequest{
				Email:    "nonexistent@example.com",
				Password: validPassword,
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return(nil, errors.New("user not found"))
			},
			expectedError: "invalid credentials",
		},
		{
			name: "invalid password",
			request: LoginRequest{
				Email:    "john@example.com",
				Password: "WrongPassword123!",
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "john@example.com").Return(validUser, nil)
			},
			expectedError: "invalid credentials",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockUserRepository)
			mockClientRepo := new(mocks.MockClientRepository)
			tt.mockSetup(mockRepo)

			tokenManager := jwt.NewTokenManager("test-secret", "test-issuer")
			authService := NewAuthService(mockRepo, mockClientRepo, tokenManager, 24*time.Hour)

			ctx := context.Background()
			response, err := authService.Login(ctx, tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Token)
				assert.NotNil(t, response.User)
				assert.Equal(t, tt.request.Email, response.User.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_GetUserByID(t *testing.T) {
	validUser := &domain.User{
		ID:        uuid.New(),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Role:      domain.RoleClient,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name          string
		userID        uuid.UUID
		mockSetup     func(*mocks.MockUserRepository)
		expectedError string
	}{
		{
			name:   "successful retrieval",
			userID: validUser.ID,
			mockSetup: func(m *mocks.MockUserRepository) {
				m.On("GetByID", mock.Anything, validUser.ID).Return(validUser, nil)
			},
			expectedError: "",
		},
		{
			name:   "user not found",
			userID: uuid.New(),
			mockSetup: func(m *mocks.MockUserRepository) {
				m.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, errors.New("user not found"))
			},
			expectedError: "user not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockUserRepository)
			mockClientRepo := new(mocks.MockClientRepository)
			tt.mockSetup(mockRepo)

			tokenManager := jwt.NewTokenManager("test-secret", "test-issuer")
			authService := NewAuthService(mockRepo, mockClientRepo, tokenManager, 24*time.Hour)

			ctx := context.Background()
			user, err := authService.GetUserByID(ctx, tt.userID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.userID, user.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
