package handler

import (
	"net/http"
	"strconv"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task assigned to an employee
// @Tags tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.Task true "Task creation request"
// @Success 201 {object} domain.Task
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Set creator from context (middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}
	task.CreatorID = uuid.MustParse(userID.(string))

	if err := h.taskService.CreateTask(c.Request.Context(), &task); err != nil {
		if err == service.ErrAssigneeNotFound {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Assignee not found"})
			return
		}
		if err == service.ErrAssigneeNotActive {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Assignee is not active"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// ListTasks godoc
// @Summary List tasks
// @Description List tasks with filters
// @Tags tasks
// @Security BearerAuth
// @Param assigneeId query string false "Filter by assignee ID"
// @Param status query string false "Filter by status"
// @Param page query int false "Page number"
// @Success 200 {array} domain.Task
// @Router /api/v1/tasks [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	var filter repository.TaskFilter

	if assigneeID := c.Query("assigneeId"); assigneeID != "" {
		id, err := uuid.Parse(assigneeID)
		if err == nil {
			filter.AssigneeID = &id
		}
	}

	if status := c.Query("status"); status != "" {
		s := domain.TaskStatus(status)
		filter.Status = &s
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil {
		filter.Page = page
	}

	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil {
		filter.PageSize = pageSize
	}

	tasks, count, err := h.taskService.ListTasks(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
		"meta": gin.H{
			"total": count,
			"page":  filter.Page,
		},
	})
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update task status or details
// @Tags tasks
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param request body domain.Task true "Task update request"
// @Success 200 {object} domain.Task
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}
	task.ID = id

	if err := h.taskService.UpdateTask(c.Request.Context(), &task); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task (Admin only)
// @Tags tasks
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 204 "No Content"
// @Router /api/v1/tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	if err := h.taskService.DeleteTask(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetMyTasks godoc
// @Summary Get my tasks
// @Description Get tasks assigned to the current employee
// @Tags tasks
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param status query string false "Filter by status"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/tasks/me [get]
func (h *TaskHandler) GetMyTasks(c *gin.Context) {
	employeeID, exists := c.Get("employeeID")
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Employee profile not found for current user"})
		return
	}

	empID, ok := employeeID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Invalid employee ID format"})
		return
	}

	filter := repository.TaskFilter{
		AssigneeID: &empID,
	}

	if status := c.Query("status"); status != "" {
		s := domain.TaskStatus(status)
		filter.Status = &s
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil {
		filter.Page = page
	}

	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil {
		filter.PageSize = pageSize
	}

	tasks, count, err := h.taskService.ListTasks(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
		"meta": gin.H{
			"total": count,
			"page":  filter.Page,
		},
	})
}
