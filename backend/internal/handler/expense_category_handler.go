package handler

import (
	"net/http"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ExpenseCategoryHandler handles expense category HTTP requests
type ExpenseCategoryHandler struct {
	categoryService service.ExpenseCategoryService
}

// NewExpenseCategoryHandler creates a new expense category handler
func NewExpenseCategoryHandler(categoryService service.ExpenseCategoryService) *ExpenseCategoryHandler {
	return &ExpenseCategoryHandler{
		categoryService: categoryService,
	}
}

// CreateExpenseCategory godoc
// @Summary Create a new expense category
// @Description Create a new expense category or subcategory
// @Tags expense-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body service.CreateExpenseCategoryRequest true "Category creation request"
// @Success 201 {object} domain.ExpenseCategory
// @Failure 400 {object} ErrorResponse "Invalid request or validation error"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /billing/expense-categories [post]
func (h *ExpenseCategoryHandler) CreateExpenseCategory(c *gin.Context) {
	var req service.CreateExpenseCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	category, err := h.categoryService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetExpenseCategory godoc
// @Summary Get an expense category by ID
// @Description Retrieve an expense category by its ID
// @Tags expense-categories
// @Security BearerAuth
// @Produce json
// @Param id path string true "Category ID (UUID)"
// @Success 200 {object} domain.ExpenseCategory
// @Failure 400 {object} ErrorResponse "Invalid ID format"
// @Failure 404 {object} ErrorResponse "Category not found"
// @Router /billing/expense-categories/{id} [get]
func (h *ExpenseCategoryHandler) GetExpenseCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid category ID"})
		return
	}

	category, err := h.categoryService.GetCategory(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, category)
}

// ListExpenseCategories godoc
// @Summary List all expense categories
// @Description Get all active expense categories (flat list)
// @Tags expense-categories
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.ExpenseCategory
// @Router /billing/expense-categories [get]
func (h *ExpenseCategoryHandler) ListExpenseCategories(c *gin.Context) {
	categories, err := h.categoryService.ListCategories(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryTree godoc
// @Summary Get expense category tree
// @Description Get the hierarchical tree of categories with their subcategories
// @Tags expense-categories
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.ExpenseCategoryWithChildren
// @Router /billing/expense-categories/tree [get]
func (h *ExpenseCategoryHandler) GetCategoryTree(c *gin.Context) {
	tree, err := h.categoryService.GetCategoryTree(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tree)
}

// GetParentCategories godoc
// @Summary Get parent categories only
// @Description Get all parent expense categories (categories without parent_id)
// @Tags expense-categories
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.ExpenseCategory
// @Router /billing/expense-categories/parents [get]
func (h *ExpenseCategoryHandler) GetParentCategories(c *gin.Context) {
	categories, err := h.categoryService.GetParentCategories(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetSubcategories godoc
// @Summary Get subcategories for a parent category
// @Description Retrieve all subcategories belonging to a parent category
// @Tags expense-categories
// @Security BearerAuth
// @Produce json
// @Param id path string true "Parent Category ID (UUID)"
// @Success 200 {array} domain.ExpenseCategory
// @Failure 400 {object} ErrorResponse "Invalid parent ID"
// @Failure 404 {object} ErrorResponse "Parent category not found"
// @Router /billing/expense-categories/{id}/subcategories [get]
func (h *ExpenseCategoryHandler) GetSubcategories(c *gin.Context) {
	parentIDParam := c.Param("id")
	parentID, err := uuid.Parse(parentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid parent category ID"})
		return
	}

	subcategories, err := h.categoryService.GetSubcategories(c.Request.Context(), parentID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, subcategories)
}

// UpdateExpenseCategory godoc
// @Summary Update an expense category
// @Description Update an existing expense category
// @Tags expense-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Category ID (UUID)"
// @Param request body service.UpdateExpenseCategoryRequest true "Category update request"
// @Success 200 {object} domain.ExpenseCategory
// @Failure 400 {object} ErrorResponse "Invalid request or validation error"
// @Failure 404 {object} ErrorResponse "Category not found"
// @Router /billing/expense-categories/{id} [put]
func (h *ExpenseCategoryHandler) UpdateExpenseCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid category ID"})
		return
	}

	var req service.UpdateExpenseCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	category, err := h.categoryService.UpdateCategory(c.Request.Context(), id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteExpenseCategory godoc
// @Summary Delete an expense category
// @Description Delete an expense category (cascade deletes subcategories)
// @Tags expense-categories
// @Security BearerAuth
// @Param id path string true "Category ID (UUID)"
// @Success 204 "Category deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid ID or has associated expenses"
// @Failure 404 {object} ErrorResponse "Category not found"
// @Router /billing/expense-categories/{id} [delete]
func (h *ExpenseCategoryHandler) DeleteExpenseCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid category ID"})
		return
	}

	if err := h.categoryService.DeleteCategory(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
