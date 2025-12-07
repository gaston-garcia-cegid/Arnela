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

// Test to verify the bug fix: is_active must be true after reactivation
// Bug: deletedClient.IsActive was false (from deleted record), and Update()
// overwrote the is_active = true set by Reactivate()
// Fix: Set deletedClient.IsActive = true in memory before calling Update()
func TestClientService_CreateClient_ClientIsActiveAfterReactivation(t *testing.T) {
	deletedClientID := uuid.New()
	deletedUserID := uuid.New()

	mockClientRepo := new(mocks.MockClientRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	// Setup: deleted client with is_active = false
	deletedClient := &domain.Client{
		ID:        deletedClientID,
		UserID:    deletedUserID,
		Email:     "maria.lopez@test.com",
		FirstName: "Maria",
		LastName:  "Lopez",
		Phone:     "666-555-444",
		DNICIF:    "83249595Z",
		IsActive:  false, // ❌ This is the state from deleted record
		DeletedAt: sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
	}

	inactiveUser := &domain.User{
		ID:       deletedUserID,
		Email:    "maria.lopez@test.com",
		IsActive: false,
	}

	// Mock expectations
	mockClientRepo.On("FindDeletedByEmailOrDNI", mock.Anything, "maria.lopez@test.com", "83249595Z").Return(deletedClient, nil)
	mockClientRepo.On("Reactivate", mock.Anything, deletedClientID).Return(nil)

	// CRITICAL ASSERTION: Verify that Update() is called with is_active = true
	mockClientRepo.On("Update", mock.Anything, mock.MatchedBy(func(c *domain.Client) bool {
		// This is the key assertion: after reactivation, is_active must be true
		if !c.IsActive {
			t.Errorf("BUG: Update() called with is_active = false. Expected is_active = true after reactivation")
			return false
		}
		return c.ID == deletedClientID && c.IsActive == true
	})).Return(nil)

	mockUserRepo.On("GetByIDAll", mock.Anything, deletedUserID).Return(inactiveUser, nil)
	mockUserRepo.On("Reactivate", mock.Anything, deletedUserID).Return(nil)

	clientService := NewClientService(mockClientRepo, mockUserRepo)

	req := CreateClientRequest{
		FirstName: "Maria",
		LastName:  "Lopez",
		Email:     "maria.lopez@test.com",
		Phone:     "666-555-444",
		DNICIF:    "83249595Z",
	}

	// Act
	client, err := clientService.CreateClient(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, deletedClientID, client.ID)

	// CRITICAL: Verify client is active after reactivation
	assert.True(t, client.IsActive, "Client must be active (is_active = true) after reactivation")

	mockClientRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test the exact scenario from the bug report
func TestClientService_CreateClient_MariaLopezScenario(t *testing.T) {
	clientID := uuid.New()
	userID := uuid.New()

	mockClientRepo := new(mocks.MockClientRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	// Simulate Maria Lopez as deleted client (from screenshot)
	mariaDeleted := &domain.Client{
		ID:        clientID,
		UserID:    userID,
		Email:     "mlopez@test.com",
		FirstName: "Maria",
		LastName:  "Lopez",
		Phone:     "666-555-444",
		DNICIF:    "83249595Z",
		IsActive:  false, // Was false when deleted
		DeletedAt: sql.NullTime{Time: time.Now().Add(-1 * time.Hour), Valid: true},
	}

	mariaUser := &domain.User{
		ID:       userID,
		Email:    "mlopez@test.com",
		IsActive: false, // Was deactivated when client was deleted
	}

	mockClientRepo.On("FindDeletedByEmailOrDNI", mock.Anything, "mlopez@test.com", "83249595Z").Return(mariaDeleted, nil)
	mockClientRepo.On("Reactivate", mock.Anything, clientID).Return(nil)

	// Capture the client passed to Update to verify is_active = true
	var capturedClient *domain.Client
	mockClientRepo.On("Update", mock.Anything, mock.MatchedBy(func(c *domain.Client) bool {
		capturedClient = c
		return true
	})).Return(nil)

	mockUserRepo.On("GetByIDAll", mock.Anything, userID).Return(mariaUser, nil)
	mockUserRepo.On("Reactivate", mock.Anything, userID).Return(nil)

	clientService := NewClientService(mockClientRepo, mockUserRepo)

	// Re-create Maria Lopez
	req := CreateClientRequest{
		FirstName: "Maria",
		LastName:  "Lopez",
		Email:     "mlopez@test.com",
		Phone:     "666-555-444",
		DNICIF:    "83249595Z",
	}

	// Act
	client, err := clientService.CreateClient(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Verify captured client has is_active = true
	assert.NotNil(t, capturedClient, "Update should have been called")
	assert.True(t, capturedClient.IsActive, "BUG FIX VERIFICATION: Client must have is_active = true when Update() is called")

	// Verify returned client is active
	assert.True(t, client.IsActive, "Returned client must be active")

	mockClientRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test to verify Update() doesn't overwrite is_active from Reactivate()
func TestClientService_CreateClient_UpdateDoesNotOverwriteIsActive(t *testing.T) {
	clientID := uuid.New()

	mockClientRepo := new(mocks.MockClientRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	deletedClient := &domain.Client{
		ID:        clientID,
		UserID:    uuid.Nil, // No user
		Email:     "test@example.com",
		DNICIF:    "11111111A",
		IsActive:  false, // From deleted record
		DeletedAt: sql.NullTime{Valid: true},
	}

	mockClientRepo.On("FindDeletedByEmailOrDNI", mock.Anything, "test@example.com", "11111111A").Return(deletedClient, nil)
	mockClientRepo.On("Reactivate", mock.Anything, clientID).Return(nil)

	updateCallCount := 0
	mockClientRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Client")).Run(func(args mock.Arguments) {
		updateCallCount++
		client := args.Get(1).(*domain.Client)

		// On the first (and only) call to Update, is_active must be true
		if updateCallCount == 1 && !client.IsActive {
			t.Errorf("REGRESSION: Update() called with is_active = false on call #%d. This would overwrite Reactivate()", updateCallCount)
		}
	}).Return(nil)

	clientService := NewClientService(mockClientRepo, mockUserRepo)

	req := CreateClientRequest{
		Email:  "test@example.com",
		DNICIF: "11111111A",
	}

	// Act
	_, err := clientService.CreateClient(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, updateCallCount, "Update should be called exactly once")

	mockClientRepo.AssertExpectations(t)
}
