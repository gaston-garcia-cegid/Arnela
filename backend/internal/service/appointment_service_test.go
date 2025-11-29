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

func (m *MockAppointmentRepository) GetByEmployeeID(ctx context.Context, employeeID uuid.UUID, limit, offset int) ([]*domain.Appointment, error) {
	args := m.Called(ctx, employeeID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, employeeID *uuid.UUID) ([]*domain.Appointment, error) {
	args := m.Called(ctx, startDate, endDate, employeeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) CheckOverlap(ctx context.Context, employeeID uuid.UUID, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error) {
	args := m.Called(ctx, employeeID, startTime, endTime, excludeID)
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

// MockEmployeeRepository for testing
type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	args := m.Called(ctx, employee)
	return args.Error(0)
}

func (m *MockEmployeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) Update(ctx context.Context, employee *domain.Employee) error {
	args := m.Called(ctx, employee)
	return args.Error(0)
}

func (m *MockEmployeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEmployeeRepository) List(ctx context.Context, limit, offset int) ([]*domain.Employee, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockEmployeeRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmployeeRepository) DNIExists(ctx context.Context, dni string) (bool, error) {
	args := m.Called(ctx, dni)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmployeeRepository) GetByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Employee, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetByDNI(ctx context.Context, dni string) (*domain.Employee, error) {
	args := m.Called(ctx, dni)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetBySpecialty(ctx context.Context, specialty string) ([]*domain.Employee, error) {
	args := m.Called(ctx, specialty)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Employee), args.Error(1)
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
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := NewAppointmentService(mockAppointmentRepo, mockClientRepo, mockEmployeeRepo)

	ctx := context.Background()
	clientID := uuid.New()
	employeeID := uuid.New()
	createdBy := uuid.New()
	startTime := getValidAppointmentTime()

	// Mock client exists and is active (via GetByUserID since ClientID not provided in request)
	mockClientRepo.On("GetByUserID", ctx, createdBy).Return(&domain.Client{
		ID:       clientID,
		IsActive: true,
	}, nil)

	// Mock employee exists and is active
	mockEmployeeRepo.On("GetByID", ctx, employeeID).Return(&domain.Employee{
		ID:        employeeID,
		FirstName: "Juan",
		LastName:  "Pérez",
		IsActive:  true,
	}, nil)

	// Mock no overlap
	mockAppointmentRepo.On("CheckOverlap", ctx, employeeID, mock.Anything, mock.Anything, (*uuid.UUID)(nil)).Return(false, nil)

	// Mock create
	mockAppointmentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Appointment")).Return(nil)

	// Mock get with relations
	mockAppointmentRepo.On("GetByIDWithRelations", ctx, mock.AnythingOfType("uuid.UUID")).Return(&domain.Appointment{
		ID:              uuid.New(),
		ClientID:        clientID,
		EmployeeID:      employeeID,
		Title:           "Consulta",
		StartTime:       startTime,
		EndTime:         startTime.Add(60 * time.Minute),
		DurationMinutes: 60,
		Status:          domain.AppointmentStatusPending,
	}, nil)

	req := domain.CreateAppointmentRequest{
		EmployeeID:      employeeID.String(),
		Title:           "Consulta",
		Description:     "Primera consulta",
		StartTime:       startTime,
		DurationMinutes: 60,
		// ClientID not provided - simulates client self-booking
	}

	appointment, err := service.CreateAppointment(ctx, req, createdBy)

	assert.NoError(t, err)
	assert.NotNil(t, appointment)
	assert.Equal(t, domain.AppointmentStatusPending, appointment.Status)
	mockClientRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
	mockAppointmentRepo.AssertExpectations(t)
}

func TestCreateAppointment_InvalidEmployee(t *testing.T) {
	mockAppointmentRepo := new(MockAppointmentRepository)
	mockClientRepo := new(MockClientRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := NewAppointmentService(mockAppointmentRepo, mockClientRepo, mockEmployeeRepo)

	ctx := context.Background()
	clientID := uuid.New()
	createdBy := uuid.New()
	startTime := getValidAppointmentTime()

	// Mock GetByUserID since ClientID not provided
	mockClientRepo.On("GetByUserID", ctx, createdBy).Return(&domain.Client{
		ID:       clientID,
		IsActive: true,
	}, nil)

	req := domain.CreateAppointmentRequest{
		EmployeeID:      "invalid-uuid",
		Title:           "Consulta",
		StartTime:       startTime,
		DurationMinutes: 60,
	}

	appointment, err := service.CreateAppointment(ctx, req, createdBy)

	assert.Error(t, err)
	assert.Nil(t, appointment)
	assert.Contains(t, err.Error(), "employeeId no válido")
	mockClientRepo.AssertExpectations(t)
}

func TestCreateAppointment_WeekendRejected(t *testing.T) {
	mockAppointmentRepo := new(MockAppointmentRepository)
	mockClientRepo := new(MockClientRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := NewAppointmentService(mockAppointmentRepo, mockClientRepo, mockEmployeeRepo)

	ctx := context.Background()
	clientID := uuid.New()
	employeeID := uuid.New()
	createdBy := uuid.New()

	// Create a Saturday date
	now := time.Now()
	daysUntilSaturday := (6 - int(now.Weekday()) + 7) % 7
	if daysUntilSaturday == 0 {
		daysUntilSaturday = 7
	}
	saturday := now.AddDate(0, 0, daysUntilSaturday)
	saturdayAt10 := time.Date(saturday.Year(), saturday.Month(), saturday.Day(), 10, 0, 0, 0, saturday.Location())

	// Mock GetByUserID since ClientID not provided
	mockClientRepo.On("GetByUserID", ctx, createdBy).Return(&domain.Client{
		ID:       clientID,
		IsActive: true,
	}, nil)

	// Mock employee exists
	mockEmployeeRepo.On("GetByID", ctx, employeeID).Return(&domain.Employee{
		ID:        employeeID,
		FirstName: "Juan",
		LastName:  "Pérez",
		IsActive:  true,
	}, nil)

	req := domain.CreateAppointmentRequest{
		EmployeeID:      employeeID.String(),
		Title:           "Consulta",
		StartTime:       saturdayAt10,
		DurationMinutes: 60,
	}

	appointment, err := service.CreateAppointment(ctx, req, createdBy)

	assert.Error(t, err)
	assert.Nil(t, appointment)
	assert.Contains(t, err.Error(), "de lunes a viernes")
	mockClientRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
}

func TestGetAvailableSlots_Success(t *testing.T) {
	mockAppointmentRepo := new(MockAppointmentRepository)
	mockClientRepo := new(MockClientRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := NewAppointmentService(mockAppointmentRepo, mockClientRepo, mockEmployeeRepo)

	ctx := context.Background()
	employeeID := uuid.New()

	// Get next Monday
	date := getValidAppointmentTime()

	// Mock employee exists and is active
	mockEmployeeRepo.On("GetByID", ctx, employeeID).Return(&domain.Employee{
		ID:        employeeID,
		FirstName: "Juan",
		LastName:  "Pérez",
		IsActive:  true,
	}, nil)

	// Mock existing appointment from 10:00 to 11:00
	existingAppointment := &domain.Appointment{
		ID:         uuid.New(),
		EmployeeID: employeeID,
		StartTime:  time.Date(date.Year(), date.Month(), date.Day(), 10, 0, 0, 0, date.Location()),
		EndTime:    time.Date(date.Year(), date.Month(), date.Day(), 11, 0, 0, 0, date.Location()),
		Status:     domain.AppointmentStatusConfirmed,
	}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location())
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 18, 0, 0, 0, date.Location())

	mockAppointmentRepo.On("GetByDateRange", ctx, startOfDay, endOfDay, &employeeID).
		Return([]*domain.Appointment{existingAppointment}, nil)

	slots, err := service.GetAvailableSlots(ctx, employeeID, date, 60)

	assert.NoError(t, err)
	assert.NotNil(t, slots)
	// Should have slots before 9:45 (10:00 - 15min buffer) and after 11:15 (11:00 + 15min buffer)
	assert.Greater(t, len(slots), 0)
	mockAppointmentRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
}

func TestConfirmAppointment_Success(t *testing.T) {
	mockAppointmentRepo := new(MockAppointmentRepository)
	mockClientRepo := new(MockClientRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := NewAppointmentService(mockAppointmentRepo, mockClientRepo, mockEmployeeRepo)

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
