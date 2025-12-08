# 📖 Guía de Desarrollo - Ejemplos de Uso

> Ejemplos prácticos y patrones de código para el proyecto Arnela

---

## 📑 Índice

1. [Backend Examples (Go)](#backend-examples-go)
   - [Crear un Nuevo Endpoint](#1-crear-un-nuevo-endpoint)
   - [Implementar Soft Delete](#2-implementar-soft-delete)
   - [Validaciones Personalizadas](#3-validaciones-personalizadas)
   - [Testing con Mocks](#4-testing-con-mocks)
2. [Frontend Examples (TypeScript/React)](#frontend-examples-typescriptreact)
   - [Crear un Componente Reutilizable](#1-crear-un-componente-reutilizable)
   - [Zustand Store](#2-zustand-store)
   - [Custom Hook con API](#3-custom-hook-con-api)
   - [Formulario con Validación](#4-formulario-con-validación)
3. [API Usage Examples](#api-usage-examples)
4. [Common Patterns](#common-patterns)
5. [Troubleshooting](#troubleshooting)

---

## Backend Examples (Go)

### 1. Crear un Nuevo Endpoint

#### Paso 1: Definir el Dominio

```go
// internal/domain/task.go
package domain

import (
    "time"
    "github.com/google/uuid"
)

// Task representa una tarea del sistema
type Task struct {
    ID          uuid.UUID  `json:"id" db:"id"`
    Title       string     `json:"title" db:"title"`
    Description string     `json:"description" db:"description"`
    Status      string     `json:"status" db:"status"` // pending, in_progress, completed
    AssignedTo  *uuid.UUID `json:"assignedTo" db:"assigned_to"`
    CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
    UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
    DeletedAt   *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}

// IsCompleted verifica si la tarea está completada
func (t *Task) IsCompleted() bool {
    return t.Status == "completed"
}
```

#### Paso 2: Definir Repositorio

```go
// internal/repository/task_repository.go
package repository

import (
    "context"
    "arnela/internal/domain"
    "github.com/google/uuid"
)

// TaskRepository define las operaciones de acceso a datos para tareas
type TaskRepository interface {
    // Create crea una nueva tarea
    Create(ctx context.Context, task *domain.Task) error
    
    // GetByID obtiene una tarea por ID (solo activas)
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
    
    // List lista todas las tareas activas
    List(ctx context.Context, filters TaskFilters) ([]*domain.Task, error)
    
    // Update actualiza una tarea existente
    Update(ctx context.Context, task *domain.Task) error
    
    // Delete elimina lógicamente una tarea (soft delete)
    Delete(ctx context.Context, id uuid.UUID) error
}

// TaskFilters define filtros para listar tareas
type TaskFilters struct {
    Status     string
    AssignedTo *uuid.UUID
    Page       int
    PageSize   int
}
```

#### Paso 3: Implementar Repositorio PostgreSQL

```go
// internal/repository/postgres/task_repository.go
package postgres

import (
    "context"
    "database/sql"
    "arnela/internal/domain"
    "arnela/internal/repository"
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
)

type taskRepository struct {
    db *sqlx.DB
}

// NewTaskRepository crea una nueva instancia del repositorio
func NewTaskRepository(db *sqlx.DB) repository.TaskRepository {
    return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *domain.Task) error {
    query := `
        INSERT INTO tasks (id, title, description, status, assigned_to, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    
    _, err := r.db.ExecContext(ctx, query,
        task.ID,
        task.Title,
        task.Description,
        task.Status,
        task.AssignedTo,
        task.CreatedAt,
        task.UpdatedAt,
    )
    
    return err
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
    query := `
        SELECT id, title, description, status, assigned_to, created_at, updated_at
        FROM tasks
        WHERE id = $1 AND deleted_at IS NULL
    `
    
    var task domain.Task
    err := r.db.GetContext(ctx, &task, query, id)
    if err == sql.ErrNoRows {
        return nil, repository.ErrTaskNotFound
    }
    
    return &task, err
}

func (r *taskRepository) List(ctx context.Context, filters repository.TaskFilters) ([]*domain.Task, error) {
    query := `
        SELECT id, title, description, status, assigned_to, created_at, updated_at
        FROM tasks
        WHERE deleted_at IS NULL
    `
    args := []interface{}{}
    argCount := 1
    
    // Aplicar filtros
    if filters.Status != "" {
        query += ` AND status = $` + fmt.Sprintf("%d", argCount)
        args = append(args, filters.Status)
        argCount++
    }
    
    if filters.AssignedTo != nil {
        query += ` AND assigned_to = $` + fmt.Sprintf("%d", argCount)
        args = append(args, *filters.AssignedTo)
        argCount++
    }
    
    query += ` ORDER BY created_at DESC`
    
    // Paginación
    if filters.PageSize > 0 {
        offset := (filters.Page - 1) * filters.PageSize
        query += fmt.Sprintf(" LIMIT %d OFFSET %d", filters.PageSize, offset)
    }
    
    tasks := []*domain.Task{}
    err := r.db.SelectContext(ctx, &tasks, query, args...)
    
    return tasks, err
}

func (r *taskRepository) Update(ctx context.Context, task *domain.Task) error {
    query := `
        UPDATE tasks
        SET title = $1, description = $2, status = $3, assigned_to = $4, updated_at = $5
        WHERE id = $6 AND deleted_at IS NULL
    `
    
    result, err := r.db.ExecContext(ctx, query,
        task.Title,
        task.Description,
        task.Status,
        task.AssignedTo,
        time.Now(),
        task.ID,
    )
    
    if err != nil {
        return err
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rows == 0 {
        return repository.ErrTaskNotFound
    }
    
    return nil
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
    query := `
        UPDATE tasks
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at IS NULL
    `
    
    result, err := r.db.ExecContext(ctx, query, time.Now(), id)
    if err != nil {
        return err
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rows == 0 {
        return repository.ErrTaskNotFound
    }
    
    return nil
}
```

#### Paso 4: Implementar Servicio

```go
// internal/service/task_service.go
package service

import (
    "context"
    "time"
    "arnela/internal/domain"
    "arnela/internal/repository"
    "github.com/google/uuid"
)

// TaskService define la lógica de negocio para tareas
type TaskService interface {
    CreateTask(ctx context.Context, req CreateTaskRequest) (*domain.Task, error)
    GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error)
    ListTasks(ctx context.Context, filters ListTasksRequest) ([]*domain.Task, error)
    UpdateTask(ctx context.Context, id uuid.UUID, req UpdateTaskRequest) (*domain.Task, error)
    DeleteTask(ctx context.Context, id uuid.UUID) error
    CompleteTask(ctx context.Context, id uuid.UUID) (*domain.Task, error)
}

type taskService struct {
    taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
    return &taskService{taskRepo: taskRepo}
}

// CreateTaskRequest representa la petición de creación
type CreateTaskRequest struct {
    Title       string     `json:"title" binding:"required,min=3,max=100"`
    Description string     `json:"description" binding:"max=500"`
    AssignedTo  *uuid.UUID `json:"assignedTo"`
}

func (s *taskService) CreateTask(ctx context.Context, req CreateTaskRequest) (*domain.Task, error) {
    task := &domain.Task{
        ID:          uuid.New(),
        Title:       req.Title,
        Description: req.Description,
        Status:      "pending",
        AssignedTo:  req.AssignedTo,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    if err := s.taskRepo.Create(ctx, task); err != nil {
        return nil, err
    }
    
    return task, nil
}

// UpdateTaskRequest representa la petición de actualización
type UpdateTaskRequest struct {
    Title       *string    `json:"title" binding:"omitempty,min=3,max=100"`
    Description *string    `json:"description" binding:"omitempty,max=500"`
    Status      *string    `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
    AssignedTo  *uuid.UUID `json:"assignedTo"`
}

func (s *taskService) UpdateTask(ctx context.Context, id uuid.UUID, req UpdateTaskRequest) (*domain.Task, error) {
    // Obtener tarea existente
    task, err := s.taskRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Aplicar cambios (solo campos presentes)
    if req.Title != nil {
        task.Title = *req.Title
    }
    if req.Description != nil {
        task.Description = *req.Description
    }
    if req.Status != nil {
        task.Status = *req.Status
    }
    if req.AssignedTo != nil {
        task.AssignedTo = req.AssignedTo
    }
    
    task.UpdatedAt = time.Now()
    
    if err := s.taskRepo.Update(ctx, task); err != nil {
        return nil, err
    }
    
    return task, nil
}

func (s *taskService) CompleteTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
    return s.UpdateTask(ctx, id, UpdateTaskRequest{
        Status: stringPtr("completed"),
    })
}

// Helper
func stringPtr(s string) *string {
    return &s
}
```

#### Paso 5: Implementar Handler

```go
// internal/handler/task_handler.go
package handler

import (
    "net/http"
    "arnela/internal/service"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type TaskHandler struct {
    taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
    return &TaskHandler{taskService: taskService}
}

// Create godoc
// @Summary      Crear nueva tarea
// @Description  Crea una nueva tarea en el sistema
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body service.CreateTaskRequest true "Datos de la tarea"
// @Success      201 {object} domain.Task
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Router       /tasks [post]
// @Security     BearerAuth
func (h *TaskHandler) Create(c *gin.Context) {
    var req service.CreateTaskRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    task, err := h.taskService.CreateTask(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, task)
}

// GetByID godoc
// @Summary      Obtener tarea por ID
// @Description  Obtiene los detalles de una tarea específica
// @Tags         tasks
// @Produce      json
// @Param        id path string true "Task ID" format(uuid)
// @Success      200 {object} domain.Task
// @Failure      404 {object} ErrorResponse
// @Router       /tasks/{id} [get]
// @Security     BearerAuth
func (h *TaskHandler) GetByID(c *gin.Context) {
    idParam := c.Param("id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
        return
    }
    
    task, err := h.taskService.GetTask(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
        return
    }
    
    c.JSON(http.StatusOK, task)
}

// Update godoc
// @Summary      Actualizar tarea
// @Description  Actualiza los datos de una tarea existente
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Task ID" format(uuid)
// @Param        request body service.UpdateTaskRequest true "Datos a actualizar"
// @Success      200 {object} domain.Task
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Router       /tasks/{id} [put]
// @Security     BearerAuth
func (h *TaskHandler) Update(c *gin.Context) {
    idParam := c.Param("id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
        return
    }
    
    var req service.UpdateTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    task, err := h.taskService.UpdateTask(c.Request.Context(), id, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, task)
}

// Complete godoc
// @Summary      Completar tarea
// @Description  Marca una tarea como completada
// @Tags         tasks
// @Param        id path string true "Task ID" format(uuid)
// @Success      200 {object} domain.Task
// @Failure      404 {object} ErrorResponse
// @Router       /tasks/{id}/complete [post]
// @Security     BearerAuth
func (h *TaskHandler) Complete(c *gin.Context) {
    idParam := c.Param("id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
        return
    }
    
    task, err := h.taskService.CompleteTask(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, task)
}
```

#### Paso 6: Registrar Rutas

```go
// cmd/api/main.go
func setupRoutes(router *gin.Engine, handlers *Handlers, authMiddleware *middleware.AuthMiddleware) {
    api := router.Group("/api/v1")
    {
        // Auth endpoints (públicos)
        auth := api.Group("/auth")
        {
            auth.POST("/register", handlers.Auth.Register)
            auth.POST("/login", handlers.Auth.Login)
        }
        
        // Protected endpoints
        protected := api.Group("")
        protected.Use(authMiddleware.RequireAuth())
        {
            // Tasks endpoints
            tasks := protected.Group("/tasks")
            {
                tasks.POST("", handlers.Task.Create)
                tasks.GET("/:id", handlers.Task.GetByID)
                tasks.PUT("/:id", handlers.Task.Update)
                tasks.DELETE("/:id", handlers.Task.Delete)
                tasks.POST("/:id/complete", handlers.Task.Complete)
            }
        }
    }
}
```

---

### 2. Implementar Soft Delete

```go
// Patrón completo de soft delete con reactivación

// 1. Migración SQL
CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    deleted_at TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT TRUE
);

// 2. Repository con soporte soft delete
func (r *taskRepository) FindDeletedByTitle(ctx context.Context, title string) (*domain.Task, error) {
    query := `
        SELECT * FROM tasks
        WHERE title = $1 AND deleted_at IS NOT NULL
    `
    
    var task domain.Task
    err := r.db.GetContext(ctx, &task, query, title)
    if err == sql.ErrNoRows {
        return nil, nil // No encontrado
    }
    
    return &task, err
}

func (r *taskRepository) Reactivate(ctx context.Context, id uuid.UUID) error {
    query := `
        UPDATE tasks
        SET deleted_at = NULL, is_active = TRUE
        WHERE id = $1
    `
    
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}

// 3. Service con lógica de reactivación
func (s *taskService) CreateTask(ctx context.Context, req CreateTaskRequest) (*domain.Task, error) {
    // PASO 1: Buscar task eliminada con mismo título
    deletedTask, err := s.taskRepo.FindDeletedByTitle(ctx, req.Title)
    if err != nil {
        return nil, err
    }
    
    if deletedTask != nil {
        // PASO 2: Actualizar campos en memoria
        deletedTask.Title = req.Title
        deletedTask.Description = req.Description
        deletedTask.IsActive = true // ⚠️ CRÍTICO: Sincronizar memoria
        
        // PASO 3: Reactivar en BD
        if err := s.taskRepo.Reactivate(ctx, deletedTask.ID); err != nil {
            return nil, err
        }
        
        // PASO 4: Actualizar campos
        if err := s.taskRepo.Update(ctx, deletedTask); err != nil {
            return nil, err
        }
        
        return deletedTask, nil
    }
    
    // PASO 5: Si no existe, crear nuevo
    newTask := &domain.Task{
        ID:        uuid.New(),
        Title:     req.Title,
        IsActive:  true,
        CreatedAt: time.Now(),
    }
    
    if err := s.taskRepo.Create(ctx, newTask); err != nil {
        return nil, err
    }
    
    return newTask, nil
}
```

---

### 3. Validaciones Personalizadas

```go
// pkg/validator/spanish_validator.go
package validator

import (
    "regexp"
    "strconv"
    "strings"
)

// ValidateDNI valida un DNI español
func ValidateDNI(dni string) error {
    dni = strings.ToUpper(strings.TrimSpace(dni))
    
    // Formato: 8 dígitos + 1 letra
    if !regexp.MustCompile(`^\d{8}[A-Z]$`).MatchString(dni) {
        return errors.New("invalid DNI format")
    }
    
    // Validar letra de control
    number, _ := strconv.Atoi(dni[:8])
    letters := "TRWAGMYFPDXBNJZSQVHLCKE"
    expectedLetter := letters[number%23]
    
    if dni[8] != byte(expectedLetter) {
        return errors.New("invalid DNI checksum letter")
    }
    
    return nil
}

// ValidateNIE valida un NIE español
func ValidateNIE(nie string) error {
    nie = strings.ToUpper(strings.TrimSpace(nie))
    
    // Formato: X/Y/Z + 7 dígitos + 1 letra
    if !regexp.MustCompile(`^[XYZ]\d{7}[A-Z]$`).MatchString(nie) {
        return errors.New("invalid NIE format")
    }
    
    // Convertir primera letra a número para validación
    first := nie[0]
    var prefix int
    switch first {
    case 'X':
        prefix = 0
    case 'Y':
        prefix = 1
    case 'Z':
        prefix = 2
    }
    
    // Validar letra de control
    numberStr := strconv.Itoa(prefix) + nie[1:8]
    number, _ := strconv.Atoi(numberStr)
    letters := "TRWAGMYFPDXBNJZSQVHLCKE"
    expectedLetter := letters[number%23]
    
    if nie[8] != byte(expectedLetter) {
        return errors.New("invalid NIE checksum letter")
    }
    
    return nil
}

// NormalizePhone normaliza un teléfono español
func NormalizePhone(phone string) string {
    // Eliminar espacios y guiones
    phone = strings.ReplaceAll(phone, " ", "")
    phone = strings.ReplaceAll(phone, "-", "")
    
    // Eliminar prefijos internacionales
    phone = strings.TrimPrefix(phone, "+34")
    phone = strings.TrimPrefix(phone, "0034")
    
    return phone
}

// ValidatePhone valida un teléfono español
func ValidatePhone(phone string) error {
    normalized := NormalizePhone(phone)
    
    // Formato: 9 dígitos empezando con 6, 7, 8 o 9
    if !regexp.MustCompile(`^[6-9]\d{8}$`).MatchString(normalized) {
        return errors.New("invalid Spanish phone format")
    }
    
    return nil
}

// Uso en servicio
func (s *clientService) CreateClient(ctx context.Context, req CreateClientRequest) (*domain.Client, error) {
    // Validar DNI
    if err := validator.ValidateDNI(req.DNICIF); err != nil {
        return nil, err
    }
    
    // Normalizar teléfono
    req.Phone = validator.NormalizePhone(req.Phone)
    
    // Validar teléfono
    if err := validator.ValidatePhone(req.Phone); err != nil {
        return nil, err
    }
    
    // Continuar con creación...
}
```

---

### 4. Testing con Mocks

```go
// internal/service/task_service_test.go
package service_test

import (
    "context"
    "testing"
    "time"
    "arnela/internal/domain"
    "arnela/internal/repository/mocks"
    "arnela/internal/service"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestTaskService_CreateTask_Success(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.MockTaskRepository)
    taskService := service.NewTaskService(mockRepo)
    
    mockRepo.On("FindDeletedByTitle", mock.Anything, "Test Task").
        Return(nil, nil) // No existe tarea eliminada
    
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).
        Return(nil)
    
    req := service.CreateTaskRequest{
        Title:       "Test Task",
        Description: "Test Description",
    }
    
    // Act
    task, err := taskService.CreateTask(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, task)
    assert.Equal(t, "Test Task", task.Title)
    assert.Equal(t, "pending", task.Status)
    mockRepo.AssertExpectations(t)
}

func TestTaskService_CreateTask_ReactivatesDeleted(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.MockTaskRepository)
    taskService := service.NewTaskService(mockRepo)
    
    deletedTime := time.Now().Add(-24 * time.Hour)
    deletedTask := &domain.Task{
        ID:        uuid.New(),
        Title:     "Test Task",
        Status:    "pending",
        IsActive:  false,
        DeletedAt: &deletedTime,
    }
    
    mockRepo.On("FindDeletedByTitle", mock.Anything, "Test Task").
        Return(deletedTask, nil)
    
    mockRepo.On("Reactivate", mock.Anything, deletedTask.ID).
        Return(nil)
    
    mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(task *domain.Task) bool {
        return task.ID == deletedTask.ID && task.IsActive == true
    })).Return(nil)
    
    req := service.CreateTaskRequest{
        Title:       "Test Task",
        Description: "Updated Description",
    }
    
    // Act
    task, err := taskService.CreateTask(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, task)
    assert.Equal(t, deletedTask.ID, task.ID) // Mismo ID (reactivada)
    assert.True(t, task.IsActive)
    assert.Nil(t, task.DeletedAt)
    assert.Equal(t, "Updated Description", task.Description)
    mockRepo.AssertExpectations(t)
}

func TestTaskService_CompleteTask_NotFound(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.MockTaskRepository)
    taskService := service.NewTaskService(mockRepo)
    
    taskID := uuid.New()
    mockRepo.On("GetByID", mock.Anything, taskID).
        Return(nil, repository.ErrTaskNotFound)
    
    // Act
    task, err := taskService.CompleteTask(context.Background(), taskID)
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, task)
    assert.Equal(t, repository.ErrTaskNotFound, err)
    mockRepo.AssertExpectations(t)
}
```

---

## Frontend Examples (TypeScript/React)

### 1. Crear un Componente Reutilizable

```tsx
// components/common/StatusBadge.tsx

/**
 * StatusBadge - Componente para mostrar badges de estado
 * 
 * @component
 * @example
 * <StatusBadge status="completed" />
 * <StatusBadge status="pending" label="Pendiente de revisión" />
 */

import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';

interface StatusBadgeProps {
  /** Estado de la entidad */
  status: 'pending' | 'in_progress' | 'completed' | 'cancelled';
  
  /** Label personalizado (opcional, usa traducción por defecto) */
  label?: string;
  
  /** Clases CSS adicionales */
  className?: string;
}

const STATUS_CONFIG = {
  pending: {
    className: 'bg-yellow-100 text-yellow-800 border-yellow-300',
    defaultLabel: 'Pendiente',
    icon: '⏳',
  },
  in_progress: {
    className: 'bg-blue-100 text-blue-800 border-blue-300',
    defaultLabel: 'En Progreso',
    icon: '🔄',
  },
  completed: {
    className: 'bg-green-100 text-green-800 border-green-300',
    defaultLabel: 'Completado',
    icon: '✅',
  },
  cancelled: {
    className: 'bg-red-100 text-red-800 border-red-300',
    defaultLabel: 'Cancelado',
    icon: '❌',
  },
} as const;

export function StatusBadge({ status, label, className }: StatusBadgeProps) {
  const config = STATUS_CONFIG[status];
  const displayLabel = label || config.defaultLabel;
  
  return (
    <Badge
      variant="outline"
      className={cn(config.className, className)}
    >
      <span className="mr-1">{config.icon}</span>
      {displayLabel}
    </Badge>
  );
}

// Uso en otro componente
<StatusBadge status="completed" />
<StatusBadge status="pending" label="Esperando aprobación" />
```

---

### 2. Zustand Store

```typescript
// stores/useTaskStore.ts

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { api } from '@/lib/api';

/**
 * Task type definition
 */
export interface Task {
  id: string;
  title: string;
  description: string;
  status: 'pending' | 'in_progress' | 'completed';
  assignedTo?: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * Task store state
 */
interface TaskStore {
  // State
  tasks: Task[];
  loading: boolean;
  error: string | null;
  
  // Actions
  fetchTasks: (token: string) => Promise<void>;
  createTask: (token: string, task: Partial<Task>) => Promise<Task>;
  updateTask: (token: string, id: string, updates: Partial<Task>) => Promise<Task>;
  deleteTask: (token: string, id: string) => Promise<void>;
  completeTask: (token: string, id: string) => Promise<Task>;
  
  // Helpers
  getTaskById: (id: string) => Task | undefined;
  clearTasks: () => void;
}

/**
 * useTaskStore - Global state management for tasks
 * 
 * @example
 * const { tasks, fetchTasks, createTask } = useTaskStore();
 * 
 * useEffect(() => {
 *   fetchTasks(token);
 * }, []);
 */
export const useTaskStore = create<TaskStore>()(
  persist(
    (set, get) => ({
      // Initial state
      tasks: [],
      loading: false,
      error: null,
      
      // Fetch all tasks
      fetchTasks: async (token: string) => {
        set({ loading: true, error: null });
        try {
          const response = await api.tasks.list(token);
          set({ tasks: response.tasks, loading: false });
        } catch (error) {
          set({ error: error.message, loading: false });
          throw error;
        }
      },
      
      // Create new task
      createTask: async (token: string, taskData: Partial<Task>) => {
        set({ loading: true, error: null });
        try {
          const newTask = await api.tasks.create(token, taskData);
          set((state) => ({
            tasks: [...state.tasks, newTask],
            loading: false,
          }));
          return newTask;
        } catch (error) {
          set({ error: error.message, loading: false });
          throw error;
        }
      },
      
      // Update existing task
      updateTask: async (token: string, id: string, updates: Partial<Task>) => {
        set({ loading: true, error: null });
        try {
          const updatedTask = await api.tasks.update(token, id, updates);
          set((state) => ({
            tasks: state.tasks.map((t) => (t.id === id ? updatedTask : t)),
            loading: false,
          }));
          return updatedTask;
        } catch (error) {
          set({ error: error.message, loading: false });
          throw error;
        }
      },
      
      // Delete task (soft delete)
      deleteTask: async (token: string, id: string) => {
        set({ loading: true, error: null });
        try {
          await api.tasks.delete(token, id);
          set((state) => ({
            tasks: state.tasks.filter((t) => t.id !== id),
            loading: false,
          }));
        } catch (error) {
          set({ error: error.message, loading: false });
          throw error;
        }
      },
      
      // Complete task
      completeTask: async (token: string, id: string) => {
        return get().updateTask(token, id, { status: 'completed' });
      },
      
      // Get task by ID
      getTaskById: (id: string) => {
        return get().tasks.find((t) => t.id === id);
      },
      
      // Clear all tasks
      clearTasks: () => {
        set({ tasks: [], error: null });
      },
    }),
    {
      name: 'task-storage', // localStorage key
      partialize: (state) => ({ tasks: state.tasks }), // Solo persiste tasks
    }
  )
);
```

---

### 3. Custom Hook con API

```typescript
// hooks/useTasks.ts

import { useState, useEffect, useCallback } from 'react';
import { useAuthStore } from '@/stores/useAuthStore';
import { api } from '@/lib/api';
import type { Task } from '@/types/task';

/**
 * useTasks - Custom hook for task management
 * 
 * @param {Object} options - Hook options
 * @param {boolean} options.autoFetch - Auto-fetch on mount (default: true)
 * @param {string} options.status - Filter by status
 * @returns {Object} Tasks state and actions
 * 
 * @example
 * const { tasks, loading, error, refetch, createTask } = useTasks();
 * 
 * @example
 * // With filters
 * const { tasks } = useTasks({ status: 'pending' });
 */
export function useTasks(options: {
  autoFetch?: boolean;
  status?: string;
} = {}) {
  const { autoFetch = true, status } = options;
  
  const token = useAuthStore((state) => state.token);
  
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  /**
   * Fetch tasks from API
   */
  const fetchTasks = useCallback(async () => {
    if (!token) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await api.tasks.list(token, { status });
      setTasks(response.tasks || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error loading tasks');
    } finally {
      setLoading(false);
    }
  }, [token, status]);
  
  /**
   * Create new task
   */
  const createTask = useCallback(async (taskData: Partial<Task>) => {
    if (!token) throw new Error('No token available');
    
    setLoading(true);
    setError(null);
    
    try {
      const newTask = await api.tasks.create(token, taskData);
      setTasks((prev) => [...prev, newTask]);
      return newTask;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error creating task');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [token]);
  
  /**
   * Update existing task
   */
  const updateTask = useCallback(async (id: string, updates: Partial<Task>) => {
    if (!token) throw new Error('No token available');
    
    setLoading(true);
    setError(null);
    
    try {
      const updatedTask = await api.tasks.update(token, id, updates);
      setTasks((prev) => prev.map((t) => (t.id === id ? updatedTask : t)));
      return updatedTask;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error updating task');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [token]);
  
  /**
   * Delete task
   */
  const deleteTask = useCallback(async (id: string) => {
    if (!token) throw new Error('No token available');
    
    setLoading(true);
    setError(null);
    
    try {
      await api.tasks.delete(token, id);
      setTasks((prev) => prev.filter((t) => t.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error deleting task');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [token]);
  
  /**
   * Complete task (shorthand)
   */
  const completeTask = useCallback(async (id: string) => {
    return updateTask(id, { status: 'completed' });
  }, [updateTask]);
  
  // Auto-fetch on mount
  useEffect(() => {
    if (autoFetch) {
      fetchTasks();
    }
  }, [autoFetch, fetchTasks]);
  
  return {
    tasks,
    loading,
    error,
    refetch: fetchTasks,
    createTask,
    updateTask,
    deleteTask,
    completeTask,
  };
}
```

---

### 4. Formulario con Validación

```tsx
// components/tasks/CreateTaskForm.tsx

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { useTasks } from '@/hooks/useTasks';
import { toast } from 'sonner';

/**
 * Form validation schema with Zod
 */
const taskSchema = z.object({
  title: z
    .string()
    .min(3, 'El título debe tener al menos 3 caracteres')
    .max(100, 'El título no puede exceder 100 caracteres'),
  description: z
    .string()
    .max(500, 'La descripción no puede exceder 500 caracteres')
    .optional(),
  assignedTo: z.string().uuid('ID de usuario inválido').optional(),
});

type TaskFormData = z.infer<typeof taskSchema>;

interface CreateTaskFormProps {
  onSuccess?: () => void;
  onCancel?: () => void;
}

/**
 * CreateTaskForm - Formulario para crear nuevas tareas
 * 
 * @component
 * @example
 * <CreateTaskForm
 *   onSuccess={() => {
 *     toast.success('Tarea creada');
 *     closeModal();
 *   }}
 * />
 */
export function CreateTaskForm({ onSuccess, onCancel }: CreateTaskFormProps) {
  const { createTask } = useTasks({ autoFetch: false });
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<TaskFormData>({
    resolver: zodResolver(taskSchema),
  });
  
  const onSubmit = async (data: TaskFormData) => {
    setIsSubmitting(true);
    
    try {
      await createTask(data);
      toast.success('Tarea creada correctamente');
      reset();
      onSuccess?.();
    } catch (error) {
      toast.error('Error al crear la tarea');
      console.error(error);
    } finally {
      setIsSubmitting(false);
    }
  };
  
  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      {/* Title */}
      <div>
        <label htmlFor="title" className="block text-sm font-medium mb-1">
          Título *
        </label>
        <Input
          id="title"
          {...register('title')}
          placeholder="Ej: Revisar documentación"
          className={errors.title ? 'border-red-500' : ''}
        />
        {errors.title && (
          <p className="text-sm text-red-500 mt-1">{errors.title.message}</p>
        )}
      </div>
      
      {/* Description */}
      <div>
        <label htmlFor="description" className="block text-sm font-medium mb-1">
          Descripción
        </label>
        <Textarea
          id="description"
          {...register('description')}
          placeholder="Describe la tarea..."
          rows={4}
          className={errors.description ? 'border-red-500' : ''}
        />
        {errors.description && (
          <p className="text-sm text-red-500 mt-1">{errors.description.message}</p>
        )}
      </div>
      
      {/* Actions */}
      <div className="flex gap-2 justify-end">
        {onCancel && (
          <Button
            type="button"
            variant="outline"
            onClick={onCancel}
            disabled={isSubmitting}
          >
            Cancelar
          </Button>
        )}
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? 'Creando...' : 'Crear Tarea'}
        </Button>
      </div>
    </form>
  );
}
```

---

## API Usage Examples

### Curl Examples

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePass123"
  }'

# Response:
# {
#   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#   "user": {
#     "id": "uuid",
#     "email": "admin@example.com",
#     "role": "admin"
#   }
# }

# Create Client
curl -X POST http://localhost:8080/api/v1/clients \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria@test.com",
    "firstName": "Maria",
    "lastName": "Lopez",
    "dniCif": "12345678Z",
    "phone": "612345678"
  }'

# Create Appointment
curl -X POST http://localhost:8080/api/v1/appointments \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "clientId": "client-uuid",
    "employeeId": "employee-uuid",
    "startTime": "2025-01-20T14:00:00Z",
    "endTime": "2025-01-20T15:00:00Z",
    "room": "sala_1",
    "notes": "Primera sesión"
  }'

# List Clients with Pagination
curl -X GET "http://localhost:8080/api/v1/clients?page=1&pageSize=10" \
  -H "Authorization: Bearer <token>"

# Update Client
curl -X PUT http://localhost:8080/api/v1/clients/uuid \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "698765432",
    "address": "Calle Nueva 123"
  }'

# Delete Client (Soft Delete)
curl -X DELETE http://localhost:8080/api/v1/clients/uuid \
  -H "Authorization: Bearer <token>"
```

---

## Common Patterns

### 1. Error Handling Pattern

```typescript
// lib/errorHandler.ts

export class APIError extends Error {
  constructor(
    message: string,
    public statusCode: number,
    public code?: string
  ) {
    super(message);
    this.name = 'APIError';
  }
}

export function handleAPIError(error: unknown): string {
  if (error instanceof APIError) {
    switch (error.statusCode) {
      case 400:
        return `Datos inválidos: ${error.message}`;
      case 401:
        return 'Sesión expirada. Por favor, inicia sesión nuevamente.';
      case 403:
        return 'No tienes permisos para realizar esta acción.';
      case 404:
        return 'Recurso no encontrado.';
      case 409:
        return `Conflicto: ${error.message}`;
      case 500:
        return 'Error del servidor. Intenta nuevamente.';
      default:
        return 'Error inesperado. Contacta soporte.';
    }
  }
  
  return 'Error desconocido.';
}

// Uso
try {
  await api.clients.create(token, data);
} catch (error) {
  const message = handleAPIError(error);
  toast.error(message);
}
```

---

### 2. Optimistic Updates Pattern

```typescript
// hooks/useOptimisticUpdate.ts

export function useOptimisticUpdate<T extends { id: string }>() {
  const updateItem = async (
    items: T[],
    setItems: (items: T[]) => void,
    id: string,
    updates: Partial<T>,
    apiCall: () => Promise<T>
  ) => {
    // Guardar estado anterior para rollback
    const previousItems = items;
    
    // Actualización optimista
    setItems(items.map((item) =>
      item.id === id ? { ...item, ...updates } : item
    ));
    
    try {
      // Llamada API real
      const updatedItem = await apiCall();
      
      // Actualizar con datos reales
      setItems(items.map((item) =>
        item.id === id ? updatedItem : item
      ));
      
      return updatedItem;
    } catch (error) {
      // Rollback en caso de error
      setItems(previousItems);
      throw error;
    }
  };
  
  return { updateItem };
}

// Uso
const { updateItem } = useOptimisticUpdate<Task>();

const handleComplete = async (taskId: string) => {
  await updateItem(
    tasks,
    setTasks,
    taskId,
    { status: 'completed' }, // Actualización optimista
    () => api.tasks.complete(token, taskId) // API call
  );
};
```

---

## Troubleshooting

### Backend Issues

#### "Token expired" constantly

**Problema:** Token expira muy rápido

**Solución:**
```go
// config/config.go
type Config struct {
    JWTExpiryHours int `env:"JWT_EXPIRY_HOURS" envDefault:"24"`
}

// Aumentar a 72 horas para desarrollo
JWT_EXPIRY_HOURS=72
```

#### "Database connection refused"

**Problema:** No se puede conectar a PostgreSQL

**Solución:**
```bash
# Verificar que Docker está corriendo
docker ps

# Reiniciar contenedor
docker-compose restart postgres

# Verificar logs
docker-compose logs postgres
```

#### "Port 8080 already in use"

**Problema:** Puerto ocupado

**Solución:**
```powershell
# Windows PowerShell
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Cambiar puerto en .env
API_PORT=8081
```

---

### Frontend Issues

#### "Module not found: '@/components/ui/button'"

**Problema:** Alias de TypeScript no configurado

**Solución:**
```json
// tsconfig.json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  }
}
```

#### "Hydration error" en Next.js

**Problema:** Desincronización server/client

**Solución:**
```tsx
// Usar useEffect para código client-only
const [mounted, setMounted] = useState(false);

useEffect(() => {
  setMounted(true);
}, []);

if (!mounted) return null;
```

---

**Última actualización:** Diciembre 2025  
**Versión:** 1.0.0  
**Autor:** gaston-garcia-cegid
