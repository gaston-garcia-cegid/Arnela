package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientService_CreateClient(t *testing.T) {
	dob := "1990-01-15"

	tests := []struct {
		name          string
		request       CreateClientRequest
		mockSetup     func(*mocks.MockClientRepository)
		expectedError string
	}{
		{
			name: "successful creation",
			request: CreateClientRequest{
				FirstName:   "Juan",
				LastName:    "Pérez",
				Email:       "juan.perez@example.com",
				Phone:       "612345678",
				DNI:         "12345678Z",
				DateOfBirth: &dob,
			},
			mockSetup: func(m *mocks.MockClientRepository) {
				m.On("EmailExists", mock.Anything, "juan.perez@example.com", mock.Anything).Return(false, nil)
				m.On("DNIExists", mock.Anything, "12345678Z", mock.Anything).Return(false, nil)
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.Client")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "email already exists",
			request: CreateClientRequest{
				FirstName: "María",
				LastName:  "García",
				Email:     "existing@example.com",
				Phone:     "612345678",
				DNI:       "87654321X",
			},
			mockSetup: func(m *mocks.MockClientRepository) {
				m.On("EmailExists", mock.Anything, "existing@example.com", mock.Anything).Return(true, nil)
			},
			expectedError: "email already registered",
		},
		{
			name: "DNI already exists",
			request: CreateClientRequest{
				FirstName: "Carlos",
				LastName:  "López",
				Email:     "carlos@example.com",
				Phone:     "612345678",
				DNI:       "12345678Z",
			},
			mockSetup: func(m *mocks.MockClientRepository) {
				m.On("EmailExists", mock.Anything, "carlos@example.com", mock.Anything).Return(false, nil)
				m.On("DNIExists", mock.Anything, "12345678Z", mock.Anything).Return(true, nil)
			},
			expectedError: "DNI already registered",
		},
		{
			name: "invalid email format",
			request: CreateClientRequest{
				FirstName: "Ana",
				LastName:  "Martínez",
				Email:     "invalid-email",
				Phone:     "612345678",
				DNI:       "12345678Z",
			},
			mockSetup:     func(m *mocks.MockClientRepository) {},
			expectedError: "invalid email format",
		},
		{
			name: "invalid phone format",
			request: CreateClientRequest{
				FirstName: "Pedro",
				LastName:  "Sánchez",
				Email:     "pedro@example.com",
				Phone:     "123",
				DNI:       "12345678Z",
			},
			mockSetup:     func(m *mocks.MockClientRepository) {},
			expectedError: "invalid phone format",
		},
		{
			name: "invalid DNI format",
			request: CreateClientRequest{
				FirstName: "Laura",
				LastName:  "Fernández",
				Email:     "laura@example.com",
				Phone:     "612345678",
				DNI:       "invalid",
			},
			mockSetup:     func(m *mocks.MockClientRepository) {},
			expectedError: "invalid DNI/NIE format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockClientRepository)
			mockUserRepo := new(mocks.MockUserRepository)
			tt.mockSetup(mockRepo)

			clientService := NewClientService(mockRepo, mockUserRepo)

			ctx := context.Background()
			client, err := clientService.CreateClient(ctx, tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, tt.request.Email, client.Email)
				assert.True(t, client.IsActive)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestClientService_GetClient(t *testing.T) {
	validClient := &domain.Client{
		ID:        uuid.New(),
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "juan@example.com",
		Phone:     "+34612345678",
		DNI:       "12345678Z",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name          string
		clientID      uuid.UUID
		mockSetup     func(*mocks.MockClientRepository)
		expectedError string
	}{
		{
			name:     "successful retrieval",
			clientID: validClient.ID,
			mockSetup: func(m *mocks.MockClientRepository) {
				m.On("GetByID", mock.Anything, validClient.ID).Return(validClient, nil)
			},
			expectedError: "",
		},
		{
			name:     "client not found",
			clientID: uuid.New(),
			mockSetup: func(m *mocks.MockClientRepository) {
				m.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, errors.New("client not found"))
			},
			expectedError: "client not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockClientRepository)
			mockUserRepo := new(mocks.MockUserRepository)
			tt.mockSetup(mockRepo)

			clientService := NewClientService(mockRepo, mockUserRepo)

			ctx := context.Background()
			client, err := clientService.GetClient(ctx, tt.clientID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, tt.clientID, client.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
func TestClientService_DeleteClient(t *testing.T) {
	validClient := &domain.Client{
		ID:        uuid.New(),
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "juan@example.com",
		IsActive:  true,
	}

	tests := []struct {
		name          string
		clientID      uuid.UUID
		mockSetup     func(*mocks.MockClientRepository)
		expectedError string
	}{
		{
			name:     "successful deletion",
			clientID: validClient.ID,
			mockSetup: func(m *mocks.MockClientRepository) {
				m.On("GetByID", mock.Anything, validClient.ID).Return(validClient, nil)
				m.On("Delete", mock.Anything, validClient.ID).Return(nil)
			},
			expectedError: "",
		},
		{
			name:     "client not found",
			clientID: uuid.New(),
			mockSetup: func(m *mocks.MockClientRepository) {
				m.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, errors.New("client not found"))
			},
			expectedError: "client not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockClientRepository)
			mockUserRepo := new(mocks.MockUserRepository)
			tt.mockSetup(mockRepo)

			clientService := NewClientService(mockRepo, mockUserRepo)

			ctx := context.Background()
			err := clientService.DeleteClient(ctx, tt.clientID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
