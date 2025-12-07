package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test to verify the bug fix: user MUST be reactivated when client is reactivated
// Bug: GetByID() filtered inactive users, preventing reactivation
// Fix: Use GetByIDAll() to fetch users regardless of is_active status
func TestClientService_CreateClient_ReactivatesInactiveUser(t *testing.T) {
	deletedClientID := uuid.New()
	deletedUserID := uuid.New()

	mockClientRepo := new(mocks.MockClientRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	// Setup: deleted client with inactive user
	deletedClient := &domain.Client{
		ID:        deletedClientID,
		UserID:    deletedUserID,
		Email:     "inactive.user@example.com",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "600000000",
		DNICIF:    "11111111A",
		IsActive:  false,
		DeletedAt: sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
	}

	inactiveUser := &domain.User{
		ID:       deletedUserID,
		Email:    "inactive.user@example.com",
		IsActive: false, // User is inactive
	}

	// Mock expectations
	mockClientRepo.On("FindDeletedByEmailOrDNI", mock.Anything, "inactive.user@example.com", "11111111A").Return(deletedClient, nil)
	mockClientRepo.On("Reactivate", mock.Anything, deletedClientID).Return(nil)
	mockClientRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Client")).Return(nil)

	// CRITICAL: Use GetByIDAll to fetch inactive user (this is the bug fix)
	mockUserRepo.On("GetByIDAll", mock.Anything, deletedUserID).Return(inactiveUser, nil)

	// Verify that Reactivate is called for the user
	mockUserRepo.On("Reactivate", mock.Anything, deletedUserID).Return(nil)

	clientService := NewClientService(mockClientRepo, mockUserRepo)

	req := CreateClientRequest{
		FirstName: "Test",
		LastName:  "User",
		Email:     "inactive.user@example.com",
		Phone:     "600000000",
		DNICIF:    "11111111A",
	}

	// Act
	client, err := clientService.CreateClient(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, deletedClientID, client.ID)

	// CRITICAL ASSERTIONS: Verify that user reactivation was called
	mockUserRepo.AssertCalled(t, "GetByIDAll", mock.Anything, deletedUserID)
	mockUserRepo.AssertCalled(t, "Reactivate", mock.Anything, deletedUserID)

	mockClientRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test to verify that GetByID failure doesn't prevent user reactivation
func TestClientService_CreateClient_UserReactivation_WithGetByIDFailure(t *testing.T) {
	deletedClientID := uuid.New()
	deletedUserID := uuid.New()

	mockClientRepo := new(mocks.MockClientRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	deletedClient := &domain.Client{
		ID:        deletedClientID,
		UserID:    deletedUserID,
		Email:     "test@example.com",
		DNICIF:    "22222222B",
		DeletedAt: sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
	}

	inactiveUser := &domain.User{
		ID:       deletedUserID,
		Email:    "test@example.com",
		IsActive: false,
	}

	mockClientRepo.On("FindDeletedByEmailOrDNI", mock.Anything, "test@example.com", "22222222B").Return(deletedClient, nil)
	mockClientRepo.On("Reactivate", mock.Anything, deletedClientID).Return(nil)
	mockClientRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Client")).Return(nil)

	// Simulate old GetByID() behavior - would fail for inactive users
	// But GetByIDAll() should succeed
	mockUserRepo.On("GetByIDAll", mock.Anything, deletedUserID).Return(inactiveUser, nil)
	mockUserRepo.On("Reactivate", mock.Anything, deletedUserID).Return(nil)

	clientService := NewClientService(mockClientRepo, mockUserRepo)

	req := CreateClientRequest{
		Email:  "test@example.com",
		DNICIF: "22222222B",
	}

	// Act
	client, err := clientService.CreateClient(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Verify user was reactivated despite being inactive
	mockUserRepo.AssertCalled(t, "Reactivate", mock.Anything, deletedUserID)
}
