package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockInvoiceRepository struct{ mock.Mock }

func (m *MockInvoiceRepository) Create(ctx context.Context, invoice *domain.Invoice) error {
	return m.Called(ctx, invoice).Error(0)
}
func (m *MockInvoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Invoice), args.Error(1)
}
func (m *MockInvoiceRepository) GetByInvoiceNumber(ctx context.Context, num string) (*domain.Invoice, error) {
	args := m.Called(ctx, num)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Invoice), args.Error(1)
}
func (m *MockInvoiceRepository) List(ctx context.Context, filters repository.InvoiceFilters) ([]*domain.Invoice, int, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*domain.Invoice), args.Int(1), args.Error(2)
}
func (m *MockInvoiceRepository) Update(ctx context.Context, invoice *domain.Invoice) error {
	return m.Called(ctx, invoice).Error(0)
}
func (m *MockInvoiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockInvoiceRepository) GetNextInvoiceNumber(ctx context.Context, year int) (string, error) {
	args := m.Called(ctx, year)
	return args.String(0), args.Error(1)
}
func (m *MockInvoiceRepository) GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*domain.Invoice, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).([]*domain.Invoice), args.Error(1)
}
func (m *MockInvoiceRepository) GetByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*domain.Invoice, error) {
	args := m.Called(ctx, appointmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Invoice), args.Error(1)
}
func (m *MockInvoiceRepository) GetTotalRevenueByDateRange(ctx context.Context, from, to time.Time) (float64, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(float64), args.Error(1)
}
func (m *MockInvoiceRepository) GetUnpaidInvoices(ctx context.Context) ([]*domain.Invoice, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Invoice), args.Error(1)
}
func (m *MockInvoiceRepository) GetRevenueByMonth(ctx context.Context, from, to time.Time) ([]repository.MonthlyRevenueRow, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]repository.MonthlyRevenueRow), args.Error(1)
}

type MockClientRepoForInvoice struct{ mock.Mock }

func (m *MockClientRepoForInvoice) Create(ctx context.Context, c *domain.Client) error {
	return m.Called(ctx, c).Error(0)
}
func (m *MockClientRepoForInvoice) GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}
func (m *MockClientRepoForInvoice) GetByUserID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	return nil, nil
}
func (m *MockClientRepoForInvoice) GetByEmail(ctx context.Context, e string) (*domain.Client, error) {
	return nil, nil
}
func (m *MockClientRepoForInvoice) GetByDNICIF(ctx context.Context, d string) (*domain.Client, error) {
	return nil, nil
}
func (m *MockClientRepoForInvoice) Update(ctx context.Context, c *domain.Client) error { return nil }
func (m *MockClientRepoForInvoice) Delete(ctx context.Context, id uuid.UUID) error     { return nil }
func (m *MockClientRepoForInvoice) List(ctx context.Context, f repository.ClientFilters, o, l int) ([]*domain.Client, error) {
	return nil, nil
}
func (m *MockClientRepoForInvoice) Count(ctx context.Context, f repository.ClientFilters) (int, error) {
	return 0, nil
}
func (m *MockClientRepoForInvoice) EmailExists(ctx context.Context, e string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}
func (m *MockClientRepoForInvoice) DNICIFExists(ctx context.Context, d string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}
func (m *MockClientRepoForInvoice) FindDeletedByEmailOrDNI(ctx context.Context, email, dniCif string) (*domain.Client, error) {
	return nil, nil
}
func (m *MockClientRepoForInvoice) Reactivate(ctx context.Context, id uuid.UUID) error { return nil }

type MockAppointmentRepoForInvoice struct{ mock.Mock }

func (m *MockAppointmentRepoForInvoice) Create(ctx context.Context, a *domain.Appointment) error {
	return nil
}
func (m *MockAppointmentRepoForInvoice) GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	return nil, nil
}
func (m *MockAppointmentRepoForInvoice) GetByIDWithRelations(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Appointment), args.Error(1)
}
func (m *MockAppointmentRepoForInvoice) Update(ctx context.Context, a *domain.Appointment) error {
	return nil
}
func (m *MockAppointmentRepoForInvoice) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *MockAppointmentRepoForInvoice) List(ctx context.Context, f domain.AppointmentFilter) ([]*domain.Appointment, error) {
	return nil, nil
}
func (m *MockAppointmentRepoForInvoice) ListWithRelations(ctx context.Context, f domain.AppointmentFilter) ([]*domain.Appointment, error) {
	return nil, nil
}
func (m *MockAppointmentRepoForInvoice) Count(ctx context.Context, f domain.AppointmentFilter) (int, error) {
	return 0, nil
}
func (m *MockAppointmentRepoForInvoice) GetByClientID(ctx context.Context, id uuid.UUID, l, o int) ([]*domain.Appointment, error) {
	return nil, nil
}
func (m *MockAppointmentRepoForInvoice) GetByEmployeeID(ctx context.Context, id uuid.UUID, l, o int) ([]*domain.Appointment, error) {
	return nil, nil
}
func (m *MockAppointmentRepoForInvoice) GetByDateRange(ctx context.Context, s, e time.Time, eid *uuid.UUID) ([]*domain.Appointment, error) {
	return nil, nil
}
func (m *MockAppointmentRepoForInvoice) CheckOverlap(ctx context.Context, eid uuid.UUID, s, e time.Time, exID *uuid.UUID) (bool, error) {
	return false, nil
}
func (m *MockAppointmentRepoForInvoice) CheckRoomAvailability(ctx context.Context, r domain.RoomType, s, e time.Time, exID *uuid.UUID) (bool, error) {
	return true, nil
}
func (m *MockAppointmentRepoForInvoice) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AppointmentStatus) error {
	return nil
}

// --- Tests ---

func TestInvoiceService_CreateInvoice(t *testing.T) {
	ctx := context.Background()
	clientID := uuid.New()

	invoiceRepo := new(MockInvoiceRepository)
	clientRepo := new(MockClientRepoForInvoice)

	svc := NewInvoiceService(invoiceRepo, clientRepo)

	clientRepo.On("GetByID", ctx, clientID).Return(&domain.Client{ID: clientID}, nil)
	invoiceRepo.On("GetNextInvoiceNumber", ctx, mock.AnythingOfType("int")).Return("F_2026_0001", nil)
	invoiceRepo.On("Create", ctx, mock.AnythingOfType("*domain.Invoice")).Return(nil)

	req := &CreateInvoiceRequest{
		ClientID:    clientID,
		IssueDate:   time.Now(),
		DueDate:     time.Now().AddDate(0, 1, 0),
		BaseAmount:  100.0,
		Description: "Sesión de terapia",
	}

	invoice, err := svc.CreateInvoice(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, invoice)
	assert.Equal(t, "F_2026_0001", invoice.InvoiceNumber)
	assert.Equal(t, 100.0, invoice.BaseAmount)
	assert.Equal(t, domain.DefaultVATPercent, invoice.VATRate)
	expectedVAT := 100.0 * (domain.DefaultVATPercent / 100)
	assert.Equal(t, expectedVAT, invoice.VATAmount)
	assert.Equal(t, 100.0+expectedVAT, invoice.TotalAmount)
	assert.Equal(t, domain.InvoiceStatusUnpaid, invoice.Status)
	invoiceRepo.AssertExpectations(t)
	clientRepo.AssertExpectations(t)
}

func TestInvoiceService_CreateInvoice_ClientNotFound(t *testing.T) {
	ctx := context.Background()
	clientID := uuid.New()

	invoiceRepo := new(MockInvoiceRepository)
	clientRepo := new(MockClientRepoForInvoice)

	svc := NewInvoiceService(invoiceRepo, clientRepo)

	clientRepo.On("GetByID", ctx, clientID).Return(nil, fmt.Errorf("not found"))

	req := &CreateInvoiceRequest{
		ClientID:    clientID,
		IssueDate:   time.Now(),
		DueDate:     time.Now().AddDate(0, 1, 0),
		BaseAmount:  100.0,
		Description: "Sesión",
	}

	invoice, err := svc.CreateInvoice(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, invoice)
	assert.Contains(t, err.Error(), "client not found")
}

func TestInvoiceService_CreateInvoice_DueDateBeforeIssueDate(t *testing.T) {
	ctx := context.Background()
	clientID := uuid.New()

	invoiceRepo := new(MockInvoiceRepository)
	clientRepo := new(MockClientRepoForInvoice)

	svc := NewInvoiceService(invoiceRepo, clientRepo)

	clientRepo.On("GetByID", ctx, clientID).Return(&domain.Client{ID: clientID}, nil)

	req := &CreateInvoiceRequest{
		ClientID:    clientID,
		IssueDate:   time.Now(),
		DueDate:     time.Now().AddDate(0, 0, -10),
		BaseAmount:  100.0,
		Description: "Sesión",
	}

	invoice, err := svc.CreateInvoice(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, invoice)
	assert.Contains(t, err.Error(), "due date")
}

func TestInvoiceService_CreateInvoiceFromAppointment_Success(t *testing.T) {
	ctx := context.Background()
	appointmentID := uuid.New()
	clientID := uuid.New()

	invoiceRepo := new(MockInvoiceRepository)
	clientRepo := new(MockClientRepoForInvoice)
	appointmentRepo := new(MockAppointmentRepoForInvoice)

	svc := NewInvoiceService(invoiceRepo, clientRepo, appointmentRepo)

	appointment := &domain.Appointment{
		ID:        appointmentID,
		ClientID:  clientID,
		Title:     "Terapia familiar",
		StartTime: time.Now(),
		Status:    domain.AppointmentStatusCompleted,
	}

	appointmentRepo.On("GetByIDWithRelations", ctx, appointmentID).Return(appointment, nil)
	invoiceRepo.On("GetByAppointmentID", ctx, appointmentID).Return(nil, fmt.Errorf("not found"))
	clientRepo.On("GetByID", ctx, clientID).Return(&domain.Client{ID: clientID}, nil)
	invoiceRepo.On("GetNextInvoiceNumber", ctx, mock.AnythingOfType("int")).Return("F_2026_0002", nil)
	invoiceRepo.On("Create", ctx, mock.AnythingOfType("*domain.Invoice")).Return(nil)

	invoice, err := svc.CreateInvoiceFromAppointment(ctx, appointmentID, 75.0)

	assert.NoError(t, err)
	assert.NotNil(t, invoice)
	assert.Equal(t, &appointmentID, invoice.AppointmentID)
	assert.Equal(t, clientID, invoice.ClientID)
	assert.Equal(t, "Terapia familiar", invoice.Description)
	assert.Equal(t, 75.0, invoice.BaseAmount)
	assert.Equal(t, domain.DefaultVATPercent, invoice.VATRate)
	assert.InDelta(t, 75.0*(domain.DefaultVATPercent/100), invoice.VATAmount, 1e-9)
	assert.InDelta(t, 75.0+invoice.VATAmount, invoice.TotalAmount, 1e-9)
	assert.Equal(t, domain.InvoiceStatusUnpaid, invoice.Status)
	appointmentRepo.AssertExpectations(t)
}

func TestInvoiceService_CreateInvoiceFromAppointment_PendingStatus(t *testing.T) {
	ctx := context.Background()
	appointmentID := uuid.New()

	invoiceRepo := new(MockInvoiceRepository)
	clientRepo := new(MockClientRepoForInvoice)
	appointmentRepo := new(MockAppointmentRepoForInvoice)

	svc := NewInvoiceService(invoiceRepo, clientRepo, appointmentRepo)

	appointment := &domain.Appointment{
		ID:     appointmentID,
		Status: domain.AppointmentStatusPending,
	}

	appointmentRepo.On("GetByIDWithRelations", ctx, appointmentID).Return(appointment, nil)

	invoice, err := svc.CreateInvoiceFromAppointment(ctx, appointmentID, 75.0)

	assert.Error(t, err)
	assert.Nil(t, invoice)
	assert.Contains(t, err.Error(), "only completed or confirmed")
}

func TestInvoiceService_CreateInvoiceFromAppointment_DuplicateInvoice(t *testing.T) {
	ctx := context.Background()
	appointmentID := uuid.New()
	clientID := uuid.New()

	invoiceRepo := new(MockInvoiceRepository)
	clientRepo := new(MockClientRepoForInvoice)
	appointmentRepo := new(MockAppointmentRepoForInvoice)

	svc := NewInvoiceService(invoiceRepo, clientRepo, appointmentRepo)

	appointment := &domain.Appointment{
		ID:       appointmentID,
		ClientID: clientID,
		Status:   domain.AppointmentStatusCompleted,
	}

	appointmentRepo.On("GetByIDWithRelations", ctx, appointmentID).Return(appointment, nil)
	invoiceRepo.On("GetByAppointmentID", ctx, appointmentID).Return(&domain.Invoice{ID: uuid.New()}, nil)

	invoice, err := svc.CreateInvoiceFromAppointment(ctx, appointmentID, 75.0)

	assert.Error(t, err)
	assert.Nil(t, invoice)
	assert.Contains(t, err.Error(), "already exists")
}

func TestInvoiceService_CreateInvoiceFromAppointment_InvalidAmount(t *testing.T) {
	ctx := context.Background()

	invoiceRepo := new(MockInvoiceRepository)
	clientRepo := new(MockClientRepoForInvoice)
	appointmentRepo := new(MockAppointmentRepoForInvoice)

	svc := NewInvoiceService(invoiceRepo, clientRepo, appointmentRepo)

	invoice, err := svc.CreateInvoiceFromAppointment(ctx, uuid.New(), 0)

	assert.Error(t, err)
	assert.Nil(t, invoice)
	assert.Contains(t, err.Error(), "greater than 0")
}

func TestInvoiceService_MarkAsPaid(t *testing.T) {
	ctx := context.Background()
	invoiceID := uuid.New()

	invoiceRepo := new(MockInvoiceRepository)
	clientRepo := new(MockClientRepoForInvoice)

	svc := NewInvoiceService(invoiceRepo, clientRepo)

	existing := &domain.Invoice{
		ID:     invoiceID,
		Status: domain.InvoiceStatusUnpaid,
	}

	invoiceRepo.On("GetByID", ctx, invoiceID).Return(existing, nil)
	invoiceRepo.On("Update", ctx, mock.AnythingOfType("*domain.Invoice")).Return(nil)

	result, err := svc.MarkAsPaid(ctx, invoiceID)

	assert.NoError(t, err)
	assert.Equal(t, domain.InvoiceStatusPaid, result.Status)
}

func TestInvoiceService_MarkAsPaid_AlreadyPaid(t *testing.T) {
	ctx := context.Background()
	invoiceID := uuid.New()

	invoiceRepo := new(MockInvoiceRepository)
	clientRepo := new(MockClientRepoForInvoice)

	svc := NewInvoiceService(invoiceRepo, clientRepo)

	existing := &domain.Invoice{
		ID:     invoiceID,
		Status: domain.InvoiceStatusPaid,
	}

	invoiceRepo.On("GetByID", ctx, invoiceID).Return(existing, nil)

	result, err := svc.MarkAsPaid(ctx, invoiceID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already paid")
}

func TestInvoiceService_DeleteInvoice_CannotDeletePaid(t *testing.T) {
	ctx := context.Background()
	invoiceID := uuid.New()

	invoiceRepo := new(MockInvoiceRepository)
	clientRepo := new(MockClientRepoForInvoice)

	svc := NewInvoiceService(invoiceRepo, clientRepo)

	existing := &domain.Invoice{
		ID:     invoiceID,
		Status: domain.InvoiceStatusPaid,
	}

	invoiceRepo.On("GetByID", ctx, invoiceID).Return(existing, nil)

	err := svc.DeleteInvoice(ctx, invoiceID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete a paid invoice")
}
