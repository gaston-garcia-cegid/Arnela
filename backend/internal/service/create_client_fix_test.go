package service

import (
	"context"
	"testing"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateClient_AddressHandling(t *testing.T) {
	// Setup Mocks
	mockClientRepo := new(mocks.MockClientRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	// Stub checks to pass
	mockClientRepo.On("EmailExists", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	mockClientRepo.On("DNICIFExists", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	mockUserRepo.On("EmailExists", mock.Anything, mock.Anything).Return(false, nil)
	mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	// Capture the client passed to repository to verify address mapping
	var capturedClient *domain.Client
	mockClientRepo.On("Create", mock.Anything, mock.MatchedBy(func(c *domain.Client) bool {
		capturedClient = c
		return true
	})).Return(nil)

	service := NewClientService(mockClientRepo, mockUserRepo)

	req := CreateClientRequest{
		FirstName: "Juan",
		LastName:  "Test",
		Email:     "juan@test.com",
		Phone:     "600000000",
		DNICIF:    "12345678Z",
		Address:   "Calle Verdadera 123", // Flat string input
		City:      "Madrid",
		Province:  "Madrid",
	}

	// Act
	_, err := service.CreateClient(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, capturedClient)

	// Verify Mapping Logic: Service must convert flat strings to struct fields in domain
	assert.Equal(t, "Calle Verdadera 123", capturedClient.AddressStreet)
	assert.Equal(t, "Madrid", capturedClient.AddressCity)
	assert.Equal(t, "Madrid", capturedClient.AddressProvince)
	assert.Equal(t, "España", capturedClient.AddressCountry) // Default set by service logic
}

func TestDeleteClient_DeactivatesUser(t *testing.T) {
	mockClientRepo := new(mocks.MockClientRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	clientID := uuid.New()
	userID := uuid.New()

	existingClient := &domain.Client{
		ID:     clientID,
		UserID: userID,
	}

	mockClientRepo.On("GetByID", mock.Anything, clientID).Return(existingClient, nil)
	mockClientRepo.On("Delete", mock.Anything, clientID).Return(nil)

	// Verification: User Delete must be called
	mockUserRepo.On("Delete", mock.Anything, userID).Return(nil)

	service := NewClientService(mockClientRepo, mockUserRepo)

	// Act
	err := service.DeleteClient(context.Background(), clientID)

	// Assert
	assert.NoError(t, err)
	mockClientRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
