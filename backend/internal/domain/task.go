package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskPriority string
type TaskStatus string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
	TaskPriorityUrgent TaskPriority = "urgent"

	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

// Task represents a task assigned to an employee
type Task struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	Title       string       `json:"title" db:"title"`
	Description string       `json:"description" db:"description"`
	CreatorID   uuid.UUID    `json:"creatorId" db:"creator_id"`
	AssigneeID  uuid.UUID    `json:"assigneeId" db:"assignee_id"`
	Status      TaskStatus   `json:"status" db:"status"`
	Priority    TaskPriority `json:"priority" db:"priority"`
	DueDate     *time.Time   `json:"dueDate" db:"due_date"`
	CreatedAt   time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time    `json:"updatedAt" db:"updated_at"`
	CompletedAt *time.Time   `json:"completedAt,omitempty" db:"completed_at"`
}

// IsOverdue checks if the task is past its due date and not completed
func (t *Task) IsOverdue() bool {
	if t.DueDate == nil || t.Status == TaskStatusCompleted || t.Status == TaskStatusCancelled {
		return false
	}
	return time.Now().After(*t.DueDate)
}
