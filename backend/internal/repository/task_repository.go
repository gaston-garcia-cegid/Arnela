package repository

import (
	"context"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/google/uuid"
)

// TaskFilter represents filters for listing tasks
type TaskFilter struct {
	AssigneeID *uuid.UUID
	CreatorID  *uuid.UUID
	Status     *domain.TaskStatus
	Priority   *domain.TaskPriority
	Page       int
	PageSize   int
}

// TaskRepository defines the interface for task persistence
type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter TaskFilter) ([]domain.Task, int, error)
}
