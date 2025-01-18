package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

type CreateStockMaterialDTO struct {
	Name                   string  `json:"name" binding:"required"`
	Description            string  `json:"description"`
	SafetyStock            float64 `json:"safetyStock" binding:"required,gt=0"`
	ExpirationFlag         bool    `json:"expirationFlag"`
	UnitID                 uint    `json:"unitId" binding:"required"`
	SupplierID             uint    `json:"supplierId" binding:"required"`
	CategoryID             uint    `json:"categoryId" binding:"required"`
	IngredientID           uint    `json:"ingredientId" binding:"required"`
	Barcode                string  `json:"barcode"`
	ExpirationPeriodInDays int     `json:"expirationPeriodInDays"`
}

type UpdateStockMaterialDTO struct {
	Name                   *string  `json:"name"`
	Description            *string  `json:"description"`
	SafetyStock            *float64 `json:"safetyStock" binding:"omitempty,gt=0"`
	ExpirationFlag         *bool    `json:"expirationFlag"`
	UnitID                 *uint    `json:"unitId"`
	CategoryID             *uint    `json:"categoryId"`
	IngredientID           *uint    `json:"ingredientId"`
	Barcode                *string  `json:"barcode"`
	ExpirationPeriodInDays *int     `json:"expirationPeriodInDays"`
	IsActive               *bool    `json:"isActive"`
}

type StockMaterialsDTO struct {
	ID                     uint    `json:"id"`
	Name                   string  `json:"name"`
	Description            string  `json:"description"`
	SafetyStock            float64 `json:"safetyStock"`
	ExpirationFlag         bool    `json:"expirationFlag"`
	UnitID                 uint    `json:"unitId"`
	UnitName               string  `json:"unitName,omitempty"`
	Category               string  `json:"category"`
	Barcode                string  `json:"barcode"`
	Ingredient             string  `json:"ingredient"`
	ExpirationPeriodInDays int     `json:"expirationPeriodInDays"`
	IsActive               bool    `json:"isActive"`
	CreatedAt              string  `json:"createdAt"`
	UpdatedAt              string  `json:"updatedAt"`
}

type StockMaterialFilter struct {
	Search           *string           `form:"search"`           // Search by name, description, or category
	LowStock         *bool             `form:"lowStock"`         // Filter for materials below safety stock
	ExpirationFlag   *bool             `form:"expirationFlag"`   // Filter by expiration flag
	IsActive         *bool             `form:"isActive"`         // Filter by active/inactive status
	SupplierID       *uint             `form:"supplierId"`       // Filter by supplier
	IngredientID     *uint             `form:"ingredientId"`     // Filter by ingredient
	CategoryID       *uint             `form:"categoryId"`       // Filter by category
	ExpirationInDays *int              `form:"expirationInDays"` // Filter by expiration period in days
	Pagination       *utils.Pagination // Pagination information
}
