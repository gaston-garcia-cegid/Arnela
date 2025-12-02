package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/google/uuid"
)

// CreateExpenseCategoryRequest represents the request to create an expense category
type CreateExpenseCategoryRequest struct {
	Name        string     `json:"name" binding:"required"`
	Code        string     `json:"code" binding:"required"`
	Description *string    `json:"description,omitempty"`
	ParentID    *uuid.UUID `json:"parentId,omitempty"` // Null for parent categories
	SortOrder   int        `json:"sortOrder"`
}

// UpdateExpenseCategoryRequest represents the request to update an expense category
type UpdateExpenseCategoryRequest struct {
	Name        string     `json:"name" binding:"required"`
	Code        string     `json:"code" binding:"required"`
	Description *string    `json:"description,omitempty"`
	ParentID    *uuid.UUID `json:"parentId,omitempty"`
	IsActive    bool       `json:"isActive"`
	SortOrder   int        `json:"sortOrder"`
}

// ExpenseCategoryService handles expense category business logic
type ExpenseCategoryService interface {
	// CreateCategory creates a new expense category or subcategory
	CreateCategory(ctx context.Context, req *CreateExpenseCategoryRequest) (*domain.ExpenseCategory, error)

	// GetCategory retrieves a category by ID
	GetCategory(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategory, error)

	// ListCategories retrieves all active categories
	ListCategories(ctx context.Context) ([]*domain.ExpenseCategory, error)

	// GetCategoryTree retrieves the hierarchical tree of categories
	GetCategoryTree(ctx context.Context) ([]*domain.ExpenseCategoryWithChildren, error)

	// GetParentCategories retrieves only parent categories (no parent_id)
	GetParentCategories(ctx context.Context) ([]*domain.ExpenseCategory, error)

	// GetSubcategories retrieves subcategories for a parent
	GetSubcategories(ctx context.Context, parentID uuid.UUID) ([]*domain.ExpenseCategory, error)

	// UpdateCategory updates an existing category
	UpdateCategory(ctx context.Context, id uuid.UUID, req *UpdateExpenseCategoryRequest) (*domain.ExpenseCategory, error)

	// DeleteCategory deletes a category (cascade deletes subcategories)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type expenseCategoryService struct {
	categoryRepo repository.ExpenseCategoryRepository
}

// NewExpenseCategoryService creates a new expense category service
func NewExpenseCategoryService(categoryRepo repository.ExpenseCategoryRepository) ExpenseCategoryService {
	return &expenseCategoryService{
		categoryRepo: categoryRepo,
	}
}

// CreateCategory creates a new expense category or subcategory
func (s *expenseCategoryService) CreateCategory(ctx context.Context, req *CreateExpenseCategoryRequest) (*domain.ExpenseCategory, error) {
	// Check if name already exists
	exists, err := s.categoryRepo.NameExists(ctx, req.Name, uuid.Nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check category name: %w", err)
	}

	if exists {
		return nil, errors.NewConflictError("category name already exists", errors.CodeConflict)
	}

	// If parentID is provided, validate parent exists and is not a subcategory
	if req.ParentID != nil {
		parent, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, errors.NewValidationError("parent category not found", map[string][]string{
				"parentId": {"parent category does not exist"},
			})
		}

		if parent.IsSubcategory() {
			return nil, errors.NewValidationError("cannot create subcategory of a subcategory", map[string][]string{
				"parentId": {"parent must be a top-level category"},
			})
		}
	}

	// Create category
	category := &domain.ExpenseCategory{
		ID:          uuid.New(),
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		ParentID:    req.ParentID,
		IsActive:    true,
		SortOrder:   req.SortOrder,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Validate
	if err := category.Validate(); err != nil {
		return nil, err
	}

	// Save to database
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// GetCategory retrieves a category by ID
func (s *expenseCategoryService) GetCategory(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategory, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

// ListCategories retrieves all active categories
func (s *expenseCategoryService) ListCategories(ctx context.Context) ([]*domain.ExpenseCategory, error) {
	return s.categoryRepo.List(ctx)
}

// GetCategoryTree retrieves the hierarchical tree of categories
func (s *expenseCategoryService) GetCategoryTree(ctx context.Context) ([]*domain.ExpenseCategoryWithChildren, error) {
	return s.categoryRepo.GetCategoryTree(ctx)
}

// GetParentCategories retrieves only parent categories (no parent_id)
func (s *expenseCategoryService) GetParentCategories(ctx context.Context) ([]*domain.ExpenseCategory, error) {
	return s.categoryRepo.GetCategories(ctx)
}

// GetSubcategories retrieves subcategories for a parent
func (s *expenseCategoryService) GetSubcategories(ctx context.Context, parentID uuid.UUID) ([]*domain.ExpenseCategory, error) {
	// Validate parent exists
	parent, err := s.categoryRepo.GetByID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	// Validate parent is not a subcategory
	if parent.IsSubcategory() {
		return nil, errors.NewValidationError("cannot get subcategories of a subcategory", map[string][]string{
			"parentId": {"parent must be a top-level category"},
		})
	}

	return s.categoryRepo.GetSubcategories(ctx, parentID)
}

// UpdateCategory updates an existing category
func (s *expenseCategoryService) UpdateCategory(ctx context.Context, id uuid.UUID, req *UpdateExpenseCategoryRequest) (*domain.ExpenseCategory, error) {
	// Get existing category
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if name already exists (excluding current category)
	exists, err := s.categoryRepo.NameExists(ctx, req.Name, id)
	if err != nil {
		return nil, fmt.Errorf("failed to check category name: %w", err)
	}

	if exists {
		return nil, errors.NewConflictError("category name already exists", errors.CodeConflict)
	}

	// If parentID is provided, validate parent exists and is not a subcategory
	if req.ParentID != nil {
		parent, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, errors.NewValidationError("parent category not found", map[string][]string{
				"parentId": {"parent category does not exist"},
			})
		}

		if parent.IsSubcategory() {
			return nil, errors.NewValidationError("cannot set parent to a subcategory", map[string][]string{
				"parentId": {"parent must be a top-level category"},
			})
		}

		// Cannot set parent to itself
		if parent.ID == id {
			return nil, errors.NewValidationError("cannot set category as its own parent", map[string][]string{
				"parentId": {"invalid parent"},
			})
		}
	}

	// Update fields
	category.Name = req.Name
	category.Code = req.Code
	category.Description = req.Description
	category.ParentID = req.ParentID
	category.IsActive = req.IsActive
	category.SortOrder = req.SortOrder
	category.UpdatedAt = time.Now()

	// Validate
	if err := category.Validate(); err != nil {
		return nil, err
	}

	// Save changes
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// DeleteCategory deletes a category (cascade deletes subcategories)
func (s *expenseCategoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	// Verify category exists
	_, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// TODO: Check if category has associated expenses
	// For now, we allow deletion (database will handle referential integrity)

	return s.categoryRepo.Delete(ctx, id)
}
