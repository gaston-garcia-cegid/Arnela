package domain

import (
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/google/uuid"
)

// ExpenseCategory represents a category or subcategory for expenses
type ExpenseCategory struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description,omitempty" db:"description"`
	ParentID    *uuid.UUID `json:"parentId,omitempty" db:"parent_id"` // Null for categories, set for subcategories
	IsActive    bool       `json:"isActive" db:"is_active"`
	SortOrder   int        `json:"sortOrder" db:"sort_order"` // For ordering in UI
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`

	// Relationships (populated via joins, not stored in DB)
	Parent        *ExpenseCategory  `json:"parent,omitempty" db:"-"`
	Subcategories []ExpenseCategory `json:"subcategories,omitempty" db:"-"`
}

// IsSubcategory returns true if this is a subcategory
func (c *ExpenseCategory) IsSubcategory() bool {
	return c.ParentID != nil
}

// IsCategory returns true if this is a top-level category
func (c *ExpenseCategory) IsCategory() bool {
	return c.ParentID == nil
}

// Validate performs basic validation on the category
func (c *ExpenseCategory) Validate() error {
	if c.Name == "" {
		return ErrInvalidCategoryName
	}
	return nil
}

// Custom errors
var (
	ErrInvalidCategoryName  = errors.NewValidationError("category name is required", nil)
	ErrCannotDeleteCategory = errors.NewConflictError("cannot delete category with active expenses", errors.CodeConflict)
)

// Predefined expense categories with subcategories
var DefaultExpenseCategories = []struct {
	Category      string
	Subcategories []string
}{
	{
		Category: "Alquiler y Arrendamiento",
		Subcategories: []string{
			"Alquiler de local",
			"Alquiler de equipo",
			"Gastos de comunidad",
		},
	},
	{
		Category: "Suministros",
		Subcategories: []string{
			"Electricidad",
			"Agua",
			"Gas",
			"Internet y teléfono",
		},
	},
	{
		Category: "Material y Consumibles",
		Subcategories: []string{
			"Material clínico",
			"Material de oficina",
			"Productos de limpieza",
			"Material desechable",
		},
	},
	{
		Category: "Personal",
		Subcategories: []string{
			"Salarios y nóminas",
			"Seguridad social",
			"Formación",
			"Dietas y gastos de viaje",
		},
	},
	{
		Category: "Servicios Profesionales",
		Subcategories: []string{
			"Asesoría fiscal",
			"Asesoría laboral",
			"Servicios jurídicos",
			"Consultoría",
		},
	},
	{
		Category: "Seguros",
		Subcategories: []string{
			"Seguro de responsabilidad civil",
			"Seguro de local",
			"Seguro de equipo",
		},
	},
	{
		Category: "Marketing y Publicidad",
		Subcategories: []string{
			"Publicidad online",
			"Publicidad offline",
			"Redes sociales",
			"Diseño gráfico",
			"Web y hosting",
		},
	},
	{
		Category: "Mantenimiento y Reparaciones",
		Subcategories: []string{
			"Mantenimiento de equipos",
			"Reparaciones",
			"Calibración de equipos",
		},
	},
	{
		Category: "Software y Tecnología",
		Subcategories: []string{
			"Licencias de software",
			"Aplicaciones",
			"Servicios cloud",
			"Hardware",
		},
	},
	{
		Category: "Impuestos y Tasas",
		Subcategories: []string{
			"IVA",
			"Tasas municipales",
			"Impuestos locales",
		},
	},
	{
		Category: "Gastos Financieros",
		Subcategories: []string{
			"Comisiones bancarias",
			"Intereses de préstamos",
			"Gastos de transferencias",
		},
	},
	{
		Category: "Otros Gastos",
		Subcategories: []string{
			"Gastos varios",
			"Imprevistos",
		},
	},
}
