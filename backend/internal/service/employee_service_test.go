package service

import (
	"context"
	"testing"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper variables for pointer fields in tests
var (
	position1 = "Fisioterapeuta"
	position2 = "Osteópata"
)

func TestEmployeeService_CreateEmployee_Success(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	ctx := context.Background()
	req := CreateEmployeeRequest{
		FirstName:   "Juan",
		LastName:    "Pérez",
		Email:       "juan.perez@example.com",
		Phone:       "612345678",
		DNI:         "12345678Z",
		Specialty:   "Fisioterapeuta",
		HireDate:    "2024-01-15",
		Notes:       "Especialista en deportiva y traumatología",
		AvatarColor: "#FF5733",
	}

	mockRepo.On("EmailExists", ctx, "juan.perez@example.com").Return(false, nil)
	mockRepo.On("DNIExists", ctx, "12345678Z").Return(false, nil)
	mockUserRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Employee")).Return(nil)

	employee, err := service.CreateEmployee(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, employee)
	assert.Equal(t, "Juan", employee.FirstName)
	assert.Equal(t, "Pérez", employee.LastName)
	assert.Equal(t, "juan.perez@example.com", employee.Email)
	assert.Equal(t, "612345678", employee.Phone)
	assert.Equal(t, "12345678Z", employee.DNI)
	assert.NotNil(t, employee.Position)
	assert.Equal(t, "Fisioterapeuta", *employee.Position)
	assert.Equal(t, 1, len(employee.Specialties))
	assert.Equal(t, "Fisioterapeuta", employee.Specialties[0])
	assert.True(t, employee.IsActive)
	assert.Equal(t, "#FF5733", employee.AvatarColor)
	assert.NotNil(t, employee.UserID) // Verify user was created and linked
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_CreateEmployee_InvalidEmail(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	tests := []struct {
		name  string
		email string
	}{
		{"missing @", "invalidemail.com"},
		{"missing domain", "invalid@"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateEmployeeRequest{
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     tt.email,
				Phone:     "612345678",
				DNI:       "12345678Z",
				Specialty: "Fisioterapeuta",
				HireDate:  "2024-01-15",
			}
			employee, err := service.CreateEmployee(context.Background(), req)

			assert.Error(t, err)
			assert.Nil(t, employee)
			assert.Equal(t, ErrInvalidEmail, err)
		})
	}
}

func TestEmployeeService_CreateEmployee_InvalidPhone(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	tests := []struct {
		name  string
		phone string
	}{
		{"too short", "61234567"},
		{"starts with invalid digit", "512345678"},
		{"non-numeric", "abcdefghi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateEmployeeRequest{
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     "juan@example.com",
				Phone:     tt.phone,
				DNI:       "12345678Z",
				Specialty: "Fisioterapeuta",
				HireDate:  "2024-01-15",
			}

			employee, err := service.CreateEmployee(context.Background(), req)

			assert.Error(t, err)
			assert.Nil(t, employee)
			assert.Equal(t, ErrInvalidPhone, err)
		})
	}
}

func TestEmployeeService_CreateEmployee_InvalidDNI(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	tests := []struct {
		name string
		dni  string
	}{
		{"too short", "1234567Z"},
		{"no letter", "12345678"},
		{"wrong letter", "12345678A"},
		{"with spaces", "12345 678Z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateEmployeeRequest{
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     "juan@example.com",
				Phone:     "612345678",
				DNI:       tt.dni,
				Specialty: "Fisioterapeuta",
				HireDate:  "2024-01-15",
			}
			employee, err := service.CreateEmployee(context.Background(), req)

			assert.Error(t, err)
			assert.Nil(t, employee)
			assert.Equal(t, ErrInvalidDNI, err)
		})
	}
}

func TestEmployeeService_CreateEmployee_EmailExists(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	ctx := context.Background()
	req := CreateEmployeeRequest{
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "existing@example.com",
		Phone:     "612345678",
		DNI:       "12345678Z",
		Specialty: "Fisioterapeuta",
		HireDate:  "2024-01-15",
	}

	mockRepo.On("EmailExists", ctx, "existing@example.com").Return(true, nil)

	employee, err := service.CreateEmployee(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, employee)
	assert.Equal(t, ErrEmailInUse, err)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_CreateEmployee_DNIExists(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	ctx := context.Background()
	req := CreateEmployeeRequest{
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "juan@example.com",
		Phone:     "612345678",
		DNI:       "12345678Z",
		Specialty: "Fisioterapeuta",
		HireDate:  "2024-01-15",
	}

	mockRepo.On("EmailExists", ctx, "juan@example.com").Return(false, nil)
	mockRepo.On("DNIExists", ctx, "12345678Z").Return(true, nil)

	employee, err := service.CreateEmployee(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, employee)
	assert.Equal(t, ErrDNIInUse, err)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_GetEmployee_Success(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	ctx := context.Background()
	employeeID := uuid.New()
	expectedEmployee := &domain.Employee{
		ID:        employeeID,
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "juan@example.com",
		Position:  &position1,
	}

	mockRepo.On("GetByID", ctx, employeeID).Return(expectedEmployee, nil)

	employee, err := service.GetEmployee(ctx, employeeID)

	assert.NoError(t, err)
	assert.NotNil(t, employee)
	assert.Equal(t, employeeID, employee.ID)
	assert.Equal(t, "Juan", employee.FirstName)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_GetEmployee_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	ctx := context.Background()
	employeeID := uuid.New()

	mockRepo.On("GetByID", ctx, employeeID).Return(nil, repository.ErrEmployeeNotFound)

	employee, err := service.GetEmployee(ctx, employeeID)

	assert.Error(t, err)
	assert.Nil(t, employee)
	assert.Equal(t, ErrEmployeeNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_UpdateEmployee_Success(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	ctx := context.Background()
	employeeID := uuid.New()
	existingEmployee := &domain.Employee{
		ID:        employeeID,
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "juan@example.com",
		Phone:     "612345678",
		DNI:       "12345678Z",
		Position:  &position1,
		IsActive:  true,
	}

	req := UpdateEmployeeRequest{
		FirstName: "Juan Carlos",
		Specialty: "Fisioterapeuta Senior",
	}

	mockRepo.On("GetByID", ctx, employeeID).Return(existingEmployee, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Employee")).Return(nil)

	employee, err := service.UpdateEmployee(ctx, employeeID, req)

	assert.NoError(t, err)
	assert.NotNil(t, employee)
	assert.Equal(t, "Juan Carlos", employee.FirstName)
	assert.NotNil(t, employee.Position)
	assert.Equal(t, "Fisioterapeuta Senior", *employee.Position)
	assert.Equal(t, "Fisioterapeuta Senior", employee.Specialties[0])
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_UpdateEmployee_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	ctx := context.Background()
	employeeID := uuid.New()
	req := UpdateEmployeeRequest{
		FirstName: "Juan",
	}

	mockRepo.On("GetByID", ctx, employeeID).Return(nil, repository.ErrEmployeeNotFound)

	employee, err := service.UpdateEmployee(ctx, employeeID, req)

	assert.Error(t, err)
	assert.Nil(t, employee)
	assert.Equal(t, ErrEmployeeNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_DeleteEmployee_Success(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	ctx := context.Background()
	employeeID := uuid.New()

	mockRepo.On("Delete", ctx, employeeID).Return(nil)

	err := service.DeleteEmployee(ctx, employeeID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_DeleteEmployee_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	ctx := context.Background()
	employeeID := uuid.New()

	mockRepo.On("Delete", ctx, employeeID).Return(repository.ErrEmployeeNotFound)

	err := service.DeleteEmployee(ctx, employeeID)

	assert.Error(t, err)
	assert.Equal(t, ErrEmployeeNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_ListEmployees_Success(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	ctx := context.Background()
	expectedEmployees := []*domain.Employee{
		{
			ID:        uuid.New(),
			FirstName: "Juan",
			LastName:  "Pérez",
			Email:     "juan@example.com",
			Position:  &position1,
		},
		{
			ID:        uuid.New(),
			FirstName: "María",
			LastName:  "García",
			Email:     "maria@example.com",
			Position:  &position2,
		},
	}

	mockRepo.On("List", ctx, 10, 0).Return(expectedEmployees, nil)
	mockRepo.On("Count", ctx).Return(2, nil)

	employees, count, err := service.ListEmployees(ctx, 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, employees)
	assert.Equal(t, 2, len(employees))
	assert.Equal(t, 2, count)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_GetEmployeesBySpecialty_Success(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	ctx := context.Background()
	specialty := "Fisioterapeuta"
	expectedEmployees := []*domain.Employee{
		{
			ID:          uuid.New(),
			FirstName:   "Juan",
			LastName:    "Pérez",
			Email:       "juan@example.com",
			Position:    &position1,
			Specialties: domain.StringArray{"Fisioterapeuta"},
		},
	}

	mockRepo.On("GetBySpecialty", ctx, specialty).Return(expectedEmployees, nil)

	employees, err := service.GetEmployeesBySpecialty(ctx, specialty)

	assert.NoError(t, err)
	assert.NotNil(t, employees)
	assert.Equal(t, 1, len(employees))
	assert.True(t, employees[0].HasSpecialty("Fisioterapeuta"))
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_CreateEmployee_DefaultAvatarColor(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	ctx := context.Background()
	req := CreateEmployeeRequest{
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "juan@example.com",
		Phone:     "612345678",
		DNI:       "12345678Z",
		Specialty: "Fisioterapeuta",
		HireDate:  "2024-01-15",
		// no AvatarColor provided -> should get default
	}

	mockRepo.On("EmailExists", ctx, "juan@example.com").Return(false, nil)
	mockRepo.On("DNIExists", ctx, "12345678Z").Return(false, nil)
	mockUserRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Employee")).Return(nil)

	employee, err := service.CreateEmployee(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, employee)
	assert.Equal(t, "#6366F1", employee.AvatarColor)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_UpdateEmployee_EmailConflict(t *testing.T) {
	mockRepo := new(mocks.MockEmployeeRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := NewEmployeeService(mockRepo, mockUserRepo)

	ctx := context.Background()
	employeeID := uuid.New()
	otherEmployeeID := uuid.New()

	existingEmployee := &domain.Employee{
		ID:    employeeID,
		Email: "juan@example.com",
	}

	otherEmployee := &domain.Employee{
		ID:    otherEmployeeID,
		Email: "existing@example.com",
	}

	req := UpdateEmployeeRequest{
		Email: "existing@example.com",
	}

	mockRepo.On("GetByID", ctx, employeeID).Return(existingEmployee, nil)
	mockRepo.On("EmailExists", ctx, "existing@example.com").Return(true, nil)
	mockRepo.On("GetByEmail", ctx, "existing@example.com").Return(otherEmployee, nil)

	employee, err := service.UpdateEmployee(ctx, employeeID, req)

	assert.Error(t, err)
	assert.Nil(t, employee)
	assert.Equal(t, ErrEmailInUse, err)
	mockRepo.AssertExpectations(t)
}
