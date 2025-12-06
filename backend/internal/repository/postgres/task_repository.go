package postgres

import (
	"context"
	"fmt"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) repository.TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (
			id, title, description, creator_id, assignee_id, 
			status, priority, due_date, created_at, updated_at
		) VALUES (
			:id, :title, :description, :creator_id, :assignee_id, 
			:status, :priority, :due_date, :created_at, :updated_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, task)
	return err
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	var task domain.Task
	query := `SELECT * FROM tasks WHERE id = $1`
	err := r.db.GetContext(ctx, &task, query, id)
	return &task, err
}

func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	query := `
		UPDATE tasks SET
			title = :title,
			description = :description,
			assignee_id = :assignee_id,
			status = :status,
			priority = :priority,
			due_date = :due_date,
			updated_at = :updated_at,
			completed_at = :completed_at
		WHERE id = :id
	`
	// Handle completion logic in service layer, here just persist
	_, err := r.db.NamedExecContext(ctx, query, task)
	return err
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *TaskRepository) List(ctx context.Context, filter repository.TaskFilter) ([]domain.Task, int, error) {
	var tasks []domain.Task
	var count int

	baseQuery := `SELECT * FROM tasks WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM tasks WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if filter.AssigneeID != nil {
		clause := fmt.Sprintf(" AND assignee_id = $%d", argIdx)
		baseQuery += clause
		countQuery += clause
		args = append(args, *filter.AssigneeID)
		argIdx++
	}

	if filter.CreatorID != nil {
		clause := fmt.Sprintf(" AND creator_id = $%d", argIdx)
		baseQuery += clause
		countQuery += clause
		args = append(args, *filter.CreatorID)
		argIdx++
	}

	if filter.Status != nil {
		clause := fmt.Sprintf(" AND status = $%d", argIdx)
		baseQuery += clause
		countQuery += clause
		args = append(args, *filter.Status)
		argIdx++
	}

	if filter.Priority != nil {
		clause := fmt.Sprintf(" AND priority = $%d", argIdx)
		baseQuery += clause
		countQuery += clause
		args = append(args, *filter.Priority)
		argIdx++
	}

	// Get total count
	if err := r.db.GetContext(ctx, &count, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// Add ordering and pagination
	baseQuery += " ORDER BY due_date ASC NULLS LAST, created_at DESC"

	limit := filter.PageSize
	offset := (filter.Page - 1) * filter.PageSize

	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	if err := r.db.SelectContext(ctx, &tasks, baseQuery, args...); err != nil {
		return nil, 0, err
	}

	return tasks, count, nil
}
