package service

import (
	"context"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAppointmentRepository is a mock implementation of AppointmentRepository
type MockAppointmentRepository struct {
	mock.Mock
}

func (m *MockAppointmentRepository) Create(ctx context.Context, appointment *domain.Appointment) error {
	args := m.Called(ctx, appointment)
	return args.Error(0)
}

func (m *MockAppointmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) GetByIDWithRelations(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) Update(ctx context.Context, appointment *domain.Appointment) error {
	args := m.Called(ctx, appointment)
	return args.Error(0)
}

func (m *MockAppointmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAppointmentRepository) List(ctx context.Context, filters domain.AppointmentFilter) ([]*domain.Appointment, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) ListWithRelations(ctx context.Context, filters domain.AppointmentFilter) ([]*domain.Appointment, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) Count(ctx context.Context, filters domain.AppointmentFilter) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *MockAppointmentRepository) GetByClientID(ctx context.Context, clientID uuid.UUID, limit, offset int) ([]*domain.Appointment, error) {
	args := m.Called(ctx, clientID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) GetByTherapistID(ctx context.Context, therapistID string, limit, offset int) ([]*domain.Appointment, error) {
	args := m.Called(ctx, therapistID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, therapistID *string) ([]*domain.Appointment, error) {
	args := m.Called(ctx, startDate, endDate, therapistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) CheckOverlap(ctx context.Context, therapistID string, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error) {
	args := m.Called(ctx, therapistID, startTime, endTime, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockAppointmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AppointmentStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// MockClientRepository for testing
type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

func (m *MockClientRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Client, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

func (m *MockClientRepository) Create(ctx context.Context, client *domain.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockClientRepository) GetByEmail(ctx context.Context, email string) (*domain.Client, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

func (m *MockClientRepository) GetByDNI(ctx context.Context, dni string) (*domain.Client, error) {
	args := m.Called(ctx, dni)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

func (m *MockClientRepository) Update(ctx context.Context, client *domain.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockClientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClientRepository) List(ctx context.Context, filters repository.ClientFilters, offset, limit int) ([]*domain.Client, error) {
	args := m.Called(ctx, filters, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Client), args.Error(1)
}

func (m *MockClientRepository) Count(ctx context.Context, filters repository.ClientFilters) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *MockClientRepository) EmailExists(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	args := m.Called(ctx, email, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockClientRepository) NIFExists(ctx context.Context, nif string, excludeID *uuid.UUID) (bool, error) {
	args := m.Called(ctx, nif, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockClientRepository) DNIExists(ctx context.Context, dni string, excludeID *uuid.UUID) (bool, error) {
	args := m.Called(ctx, dni, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockClientRepository) GetByNIF(ctx context.Context, nif string) (*domain.Client, error) {
	args := m.Called(ctx, nif)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

// Helper function to create a valid appointment time (Monday 10:00 AM, future date)
func getValidAppointmentTime() time.Time {
	now := time.Now()
	// Find next Monday
	daysUntilMonday := (8 - int(now.Weekday())) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}
	nextMonday := now.AddDate(0, 0, daysUntilMonday)
	return time.Date(nextMonday.Year(), nextMonday.Month(), nextMonday.Day(), 10, 0, 0, 0, nextMonday.Location())
}

func TestCreateAppointment_Success(t *testing.T) {
	mockAppointmentRepo := new(MockAppointmentRepository)
	mockClientRepo := new(MockClientRepository)
	service := NewAppointmentService(mockAppointmentRepo, mockClientRepo)

	ctx := context.Background()
	clientID := uuid.New()
	createdBy := uuid.New()
	startTime := getValidAppointmentTime()

	// Mock client exists and is active
	mockClientRepo.On("GetByID", ctx, clientID).Return(&domain.Client{
		ID:       clientID,
		IsActive: true,
	}, nil)

	// Mock no overlap
	mockAppointmentRepo.On("CheckOverlap", ctx, "therapist-1", mock.Anything, mock.Anything, (*uuid.UUID)(nil)).Return(false, nil)

	// Mock create
	mockAppointmentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Appointment")).Return(nil)

	// Mock get with relations
	mockAppointmentRepo.On("GetByIDWithRelations", ctx, mock.AnythingOfType("uuid.UUID")).Return(&domain.Appointment{
		ID:              uuid.New(),
		ClientID:        clientID,
		TherapistID:     "therapist-1",
		Title:           "Consulta",
		StartTime:       startTime,
		EndTime:         startTime.Add(60 * time.Minute),
		DurationMinutes: 60,
		Status:          domain.AppointmentStatusPending,
	}, nil)

	req := domain.CreateAppointmentRequest{
		ClientID:        clientID,
		TherapistID:     "therapist-1",
		Title:           "Consulta",
		Description:     "Primera consulta",
		StartTime:       startTime,
		DurationMinutes: 60,
	}

	appointment, err := service.CreateAppointment(ctx, req, createdBy)

	assert.NoError(t, err)
	assert.NotNil(t, appointment)
	assert.Equal(t, domain.AppointmentStatusPending, appointment.Status)
	mockClientRepo.AssertExpectations(t)
	mockAppointmentRepo.AssertExpectations(t)
}

func TestCreateAppointment_InvalidTherapist(t *testing.T) {
	mockAppointmentRepo := new(MockAppointmentRepository)
	mockClientRepo := new(MockClientRepository)
	service := NewAppointmentService(mockAppointmentRepo, mockClientRepo)

	ctx := context.Background()
	clientID := uuid.New()
	createdBy := uuid.New()
	startTime := getValidAppointmentTime()

	mockClientRepo.On("GetByID", ctx, clientID).Return(&domain.Client{
		ID:       clientID,
		IsActive: true,
	}, nil)

	req := domain.CreateAppointmentRequest{
		ClientID:        clientID,
		TherapistID:     "invalid-therapist",
		Title:           "Consulta",
		StartTime:       startTime,
		DurationMinutes: 60,
	}

	appointment, err := service.CreateAppointment(ctx, req, createdBy)

	assert.Error(t, err)
	assert.Nil(t, appointment)
	assert.Contains(t, err.Error(), "terapeuta no válido")
	mockClientRepo.AssertExpectations(t)
}

func TestCreateAppointment_WeekendRejected(t *testing.T) {
	mockAppointmentRepo := new(MockAppointmentRepository)
	mockClientRepo := new(MockClientRepository)
	service := NewAppointmentService(mockAppointmentRepo, mockClientRepo)

	ctx := context.Background()
	clientID := uuid.New()
	createdBy := uuid.New()

	// Create a Saturday date
	now := time.Now()
	daysUntilSaturday := (6 - int(now.Weekday()) + 7) % 7
	if daysUntilSaturday == 0 {
		daysUntilSaturday = 7
	}
	saturday := now.AddDate(0, 0, daysUntilSaturday)
	saturdayAt10 := time.Date(saturday.Year(), saturday.Month(), saturday.Day(), 10, 0, 0, 0, saturday.Location())

	mockClientRepo.On("GetByID", ctx, clientID).Return(&domain.Client{
		ID:       clientID,
		IsActive: true,
	}, nil)

	req := domain.CreateAppointmentRequest{
		ClientID:        clientID,
		TherapistID:     "therapist-1",
		Title:           "Consulta",
		StartTime:       saturdayAt10,
		DurationMinutes: 60,
	}

	appointment, err := service.CreateAppointment(ctx, req, createdBy)

	assert.Error(t, err)
	assert.Nil(t, appointment)
	assert.Contains(t, err.Error(), "de lunes a viernes")
	mockClientRepo.AssertExpectations(t)
}

func TestGetAvailableSlots_Success(t *testing.T) {
	mockAppointmentRepo := new(MockAppointmentRepository)
	mockClientRepo := new(MockClientRepository)
	service := NewAppointmentService(mockAppointmentRepo, mockClientRepo)

	ctx := context.Background()
	therapistID := "therapist-1"

	// Get next Monday
	date := getValidAppointmentTime()

	// Mock existing appointment from 10:00 to 11:00
	existingAppointment := &domain.Appointment{
		ID:          uuid.New(),
		TherapistID: therapistID,
		StartTime:   time.Date(date.Year(), date.Month(), date.Day(), 10, 0, 0, 0, date.Location()),
		EndTime:     time.Date(date.Year(), date.Month(), date.Day(), 11, 0, 0, 0, date.Location()),
		Status:      domain.AppointmentStatusConfirmed,
	}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location())
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 18, 0, 0, 0, date.Location())

	mockAppointmentRepo.On("GetByDateRange", ctx, startOfDay, endOfDay, &therapistID).
		Return([]*domain.Appointment{existingAppointment}, nil)

	slots, err := service.GetAvailableSlots(ctx, therapistID, date, 60)

	assert.NoError(t, err)
	assert.NotNil(t, slots)
	// Should have slots before 9:45 (10:00 - 15min buffer) and after 11:15 (11:00 + 15min buffer)
	assert.Greater(t, len(slots), 0)
	mockAppointmentRepo.AssertExpectations(t)
}

func TestConfirmAppointment_Success(t *testing.T) {
	mockAppointmentRepo := new(MockAppointmentRepository)
	mockClientRepo := new(MockClientRepository)
	service := NewAppointmentService(mockAppointmentRepo, mockClientRepo)

	ctx := context.Background()
	appointmentID := uuid.New()

	// Mock get appointment
	mockAppointmentRepo.On("GetByID", ctx, appointmentID).Return(&domain.Appointment{
		ID:     appointmentID,
		Status: domain.AppointmentStatusPending,
	}, nil)

	// Mock update
	mockAppointmentRepo.On("Update", ctx, mock.AnythingOfType("*domain.Appointment")).Return(nil)

	// Mock get with relations
	mockAppointmentRepo.On("GetByIDWithRelations", ctx, appointmentID).Return(&domain.Appointment{
		ID:     appointmentID,
		Status: domain.AppointmentStatusConfirmed,
	}, nil)

	req := domain.ConfirmAppointmentRequest{
		Notes: "Confirmado por admin",
	}

	appointment, err := service.ConfirmAppointment(ctx, appointmentID, req)

	assert.NoError(t, err)
	assert.NotNil(t, appointment)
	assert.Equal(t, domain.AppointmentStatusConfirmed, appointment.Status)
	mockAppointmentRepo.AssertExpectations(t)
}
