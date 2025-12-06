package service

import (
	"context"
	"errors"
	"testing"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskService_CreateTask_Success(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := NewTaskService(mockTaskRepo, mockEmployeeRepo)

	ctx := context.Background()
	employeeID := uuid.New()
	creatorID := uuid.New()

	task := &domain.Task{
		Title:       "Test Task",
		Description: "Description",
		AssigneeID:  employeeID,
		CreatorID:   creatorID,
	}

	employee := &domain.Employee{
		ID:       employeeID,
		IsActive: true,
	}

	mockEmployeeRepo.On("GetByID", ctx, employeeID).Return(employee, nil)
	mockTaskRepo.On("Create", ctx, mock.MatchedBy(func(t *domain.Task) bool {
		return t.Title == "Test Task" && t.Status == domain.TaskStatusPending
	})).Return(nil)

	err := service.CreateTask(ctx, task)

	assert.NoError(t, err)
	mockEmployeeRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_CreateTask_AssigneeNotFound(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := NewTaskService(mockTaskRepo, mockEmployeeRepo)

	ctx := context.Background()
	employeeID := uuid.New()

	task := &domain.Task{AssigneeID: employeeID}

	mockEmployeeRepo.On("GetByID", ctx, employeeID).Return(nil, errors.New("not found"))

	err := service.CreateTask(ctx, task)

	assert.ErrorIs(t, err, ErrAssigneeNotFound)
	mockTaskRepo.AssertNotCalled(t, "Create")
}

func TestTaskService_UpdateTask_Complete(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := NewTaskService(mockTaskRepo, mockEmployeeRepo)

	ctx := context.Background()
	taskID := uuid.New()

	existingTask := &domain.Task{
		ID:     taskID,
		Status: domain.TaskStatusPending,
		Title:  "Old Title",
	}

	updateReq := &domain.Task{
		ID:     taskID,
		Status: domain.TaskStatusCompleted,
		Title:  "New Title",
	}

	mockTaskRepo.On("GetByID", ctx, taskID).Return(existingTask, nil)
	mockTaskRepo.On("Update", ctx, mock.MatchedBy(func(t *domain.Task) bool {
		return t.Status == domain.TaskStatusCompleted && t.CompletedAt != nil && t.Title == "New Title"
	})).Return(nil)

	err := service.UpdateTask(ctx, updateReq)

	assert.NoError(t, err)
	mockTaskRepo.AssertExpectations(t)
}
