package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientService_CreateClient_ReactivateDeleted(t *testing.T) {
	deletedClientID := uuid.New()
	deletedUserID := uuid.New()
	dob := "1990-01-15"

	tests := []struct {
		name          string
		request       CreateClientRequest
		mockSetup     func(*mocks.MockClientRepository, *mocks.MockUserRepository)
		expectedError string
	}{
		{
			name: "reactivate deleted client with inactive user",
			request: CreateClientRequest{
				FirstName:   "Juan",
				LastName:    "Pérez",
				Email:       "juan.deleted@example.com",
				Phone:       "612345678",
				DNICIF:      "12345678Z",
				DateOfBirth: &dob,
				Address:     "Calle Nueva 123",
				City:        "Madrid",
				Province:    "Madrid",
				PostalCode:  "28001",
				Notes:       "Cliente reactivado",
			},
			mockSetup: func(mr *mocks.MockClientRepository, mu *mocks.MockUserRepository) {
				deletedClient := &domain.Client{
					ID:        deletedClientID,
					UserID:    deletedUserID,
					Email:     "juan.deleted@example.com",
					FirstName: "Juan",
					LastName:  "Pérez",
					Phone:     "600000000",
					DNICIF:    "12345678Z",
					IsActive:  false,
					CreatedAt: time.Now().Add(-48 * time.Hour),
					UpdatedAt: time.Now().Add(-24 * time.Hour),
					DeletedAt: sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				}

				inactiveUser := &domain.User{
					ID:        deletedUserID,
					Email:     "juan.deleted@example.com",
					FirstName: "Juan",
					LastName:  "Pérez",
					Role:      domain.RoleClient,
					IsActive:  false,
				}

				// FindDeletedByEmailOrDNI returns the deleted client
				mr.On("FindDeletedByEmailOrDNI", mock.Anything, "juan.deleted@example.com", "12345678Z").Return(deletedClient, nil)

				// Reactivate client
				mr.On("Reactivate", mock.Anything, deletedClientID).Return(nil)

				// Update client with new data
				mr.On("Update", mock.Anything, mock.AnythingOfType("*domain.Client")).Return(nil)

				// GetByIDAll to check user status (including inactive users)
				mu.On("GetByIDAll", mock.Anything, deletedUserID).Return(inactiveUser, nil)

				// Reactivate user
				mu.On("Reactivate", mock.Anything, deletedUserID).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "reactivate deleted client - no user found",
			request: CreateClientRequest{
				FirstName: "María",
				LastName:  "García",
				Email:     "maria.deleted@example.com",
				Phone:     "623456789",
				DNICIF:    "87654321X",
			},
			mockSetup: func(mr *mocks.MockClientRepository, mu *mocks.MockUserRepository) {
				deletedClient := &domain.Client{
					ID:        deletedClientID,
					UserID:    uuid.Nil, // No associated user
					Email:     "maria.deleted@example.com",
					FirstName: "María",
					LastName:  "García",
					Phone:     "600000000",
					DNICIF:    "87654321X",
					IsActive:  false,
					DeletedAt: sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				}

				mr.On("FindDeletedByEmailOrDNI", mock.Anything, "maria.deleted@example.com", "87654321X").Return(deletedClient, nil)
				mr.On("Reactivate", mock.Anything, deletedClientID).Return(nil)
				mr.On("Update", mock.Anything, mock.AnythingOfType("*domain.Client")).Return(nil)
				// No user operations expected
			},
			expectedError: "",
		},
		{
			name: "reactivate deleted client - reactivation fails",
			request: CreateClientRequest{
				FirstName: "Carlos",
				LastName:  "López",
				Email:     "carlos.deleted@example.com",
				Phone:     "634567890",
				DNICIF:    "11111111A",
			},
			mockSetup: func(mr *mocks.MockClientRepository, mu *mocks.MockUserRepository) {
				deletedClient := &domain.Client{
					ID:        deletedClientID,
					Email:     "carlos.deleted@example.com",
					DNICIF:    "11111111A",
					DeletedAt: sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				}

				mr.On("FindDeletedByEmailOrDNI", mock.Anything, "carlos.deleted@example.com", "11111111A").Return(deletedClient, nil)
				mr.On("Reactivate", mock.Anything, deletedClientID).Return(errors.New("failed to reactivate"))
			},
			expectedError: "failed to reactivate",
		},
		{
			name: "no deleted client found - proceed with normal creation",
			request: CreateClientRequest{
				FirstName: "Ana",
				LastName:  "Martínez",
				Email:     "ana.new@example.com",
				Phone:     "645678901",
				DNICIF:    "22222222B",
			},
			mockSetup: func(mr *mocks.MockClientRepository, mu *mocks.MockUserRepository) {
				// No deleted client found
				mr.On("FindDeletedByEmailOrDNI", mock.Anything, "ana.new@example.com", "22222222B").Return(nil, nil)

				// Proceed with normal validations
				mr.On("EmailExists", mock.Anything, "ana.new@example.com", mock.Anything).Return(false, nil)
				mu.On("EmailExists", mock.Anything, "ana.new@example.com").Return(false, nil)
				mr.On("DNICIFExists", mock.Anything, "22222222B", mock.Anything).Return(false, nil)

				// Create new user and client
				mu.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
				mr.On("Create", mock.Anything, mock.AnythingOfType("*domain.Client")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "FindDeletedByEmailOrDNI fails",
			request: CreateClientRequest{
				FirstName: "Pedro",
				LastName:  "Sánchez",
				Email:     "pedro@example.com",
				Phone:     "656789012",
				DNICIF:    "33333333C",
			},
			mockSetup: func(mr *mocks.MockClientRepository, mu *mocks.MockUserRepository) {
				mr.On("FindDeletedByEmailOrDNI", mock.Anything, "pedro@example.com", "33333333C").Return(nil, errors.New("database error"))
			},
			expectedError: "failed to check for deleted client",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClientRepo := new(mocks.MockClientRepository)
			mockUserRepo := new(mocks.MockUserRepository)
			tt.mockSetup(mockClientRepo, mockUserRepo)

			clientService := NewClientService(mockClientRepo, mockUserRepo)

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
				assert.Equal(t, tt.request.FirstName, client.FirstName)
				assert.Equal(t, tt.request.LastName, client.LastName)
			}

			mockClientRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}
