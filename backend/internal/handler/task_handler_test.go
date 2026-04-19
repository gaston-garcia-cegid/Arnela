package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository/mocks"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTaskHandlerRouter(h *TaskHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/tasks/:id", h.GetTask)
	return r
}

func TestTaskHandler_GetTask_InvalidID(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockEmpRepo := new(mocks.MockEmployeeRepository)
	svc := service.NewTaskService(mockTaskRepo, mockEmpRepo)
	h := NewTaskHandler(svc)
	r := setupTaskHandlerRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/tasks/not-a-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockTaskRepo.AssertNotCalled(t, "GetByID")
}

func TestTaskHandler_GetTask_NotFound(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockEmpRepo := new(mocks.MockEmployeeRepository)
	svc := service.NewTaskService(mockTaskRepo, mockEmpRepo)
	h := NewTaskHandler(svc)
	r := setupTaskHandlerRouter(h)

	id := uuid.New()
	mockTaskRepo.On("GetByID", mock.Anything, id).Return(nil, sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/tasks/"+id.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskHandler_GetTask_Success(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockEmpRepo := new(mocks.MockEmployeeRepository)
	svc := service.NewTaskService(mockTaskRepo, mockEmpRepo)
	h := NewTaskHandler(svc)
	r := setupTaskHandlerRouter(h)

	id := uuid.New()
	task := &domain.Task{
		ID:          id,
		Title:       "Review notes",
		Description: "Desc",
		Status:      domain.TaskStatusPending,
		Priority:    domain.TaskPriorityMedium,
	}
	mockTaskRepo.On("GetByID", mock.Anything, id).Return(task, nil)

	req := httptest.NewRequest(http.MethodGet, "/tasks/"+id.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body domain.Task
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "Review notes", body.Title)
	assert.Equal(t, domain.TaskStatusPending, body.Status)
	mockTaskRepo.AssertExpectations(t)
}
