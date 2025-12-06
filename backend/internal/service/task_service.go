package service

import (
	"context"
	"errors"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrAssigneeNotFound  = errors.New("assignee not found")
	ErrAssigneeNotActive = errors.New("assignee is not active")
	ErrInvalidTaskStatus = errors.New("invalid task status")
)

type TaskService struct {
	taskRepo     repository.TaskRepository
	employeeRepo repository.EmployeeRepository
}

func NewTaskService(taskRepo repository.TaskRepository, employeeRepo repository.EmployeeRepository) *TaskService {
	return &TaskService{
		taskRepo:     taskRepo,
		employeeRepo: employeeRepo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task *domain.Task) error {
	// Validate Assignee
	employee, err := s.employeeRepo.GetByID(ctx, task.AssigneeID)
	if err != nil {
		return ErrAssigneeNotFound
	}
	if !employee.IsActive {
		return ErrAssigneeNotActive
	}

	// Set defaults
	if task.Status == "" {
		task.Status = domain.TaskStatusPending
	}
	if task.Priority == "" {
		task.Priority = domain.TaskPriorityMedium
	}

	now := time.Now()
	task.ID = uuid.New()
	task.CreatedAt = now
	task.UpdatedAt = now

	return s.taskRepo.Create(ctx, task)
}

func (s *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}

func (s *TaskService) UpdateTask(ctx context.Context, task *domain.Task) error {
	existing, err := s.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		return err
	}

	// Update allowed fields
	existing.Title = task.Title
	existing.Description = task.Description
	existing.Priority = task.Priority
	existing.DueDate = task.DueDate
	existing.UpdatedAt = time.Now()

	// Handle status change
	if task.Status != "" {
		existing.Status = task.Status
		if existing.Status == domain.TaskStatusCompleted && existing.CompletedAt == nil {
			now := time.Now()
			existing.CompletedAt = &now
		} else if existing.Status != domain.TaskStatusCompleted {
			existing.CompletedAt = nil
		}
	}

	// Handle reassignment
	if task.AssigneeID != uuid.Nil && task.AssigneeID != existing.AssigneeID {
		employee, err := s.employeeRepo.GetByID(ctx, task.AssigneeID)
		if err != nil {
			return ErrAssigneeNotFound
		}
		if !employee.IsActive {
			return ErrAssigneeNotActive
		}
		existing.AssigneeID = task.AssigneeID
	}

	return s.taskRepo.Update(ctx, existing)
}

func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	_, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.taskRepo.Delete(ctx, id)
}

func (s *TaskService) ListTasks(ctx context.Context, filter repository.TaskFilter) ([]domain.Task, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	return s.taskRepo.List(ctx, filter)
}
