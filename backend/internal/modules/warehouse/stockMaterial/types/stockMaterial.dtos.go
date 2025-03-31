package types

import (
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	stockMaterialCategoryTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialCategory/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type GenerateBarcodeResponse struct {
	Barcode string `json:"barcode"`
}

type CreateStockMaterialDTO struct {
	Name                   string  `json:"name" binding:"required"`
	Description            *string `json:"description" binding:"omitempty"`
	SafetyStock            float64 `json:"safetyStock" binding:"required,gt=0"`
	UnitID                 uint    `json:"unitId" binding:"required"`
	CategoryID             uint    `json:"categoryId" binding:"required"`
	IngredientID           uint    `json:"ingredientId" binding:"required"`
	Barcode                string  `json:"barcode"`
	ExpirationPeriodInDays int     `json:"expirationPeriodInDays"`
	IsActive               bool    `json:"isActive" binding:"required"`
	Size                   float64 `json:"size"`
}

type UpdateStockMaterialDTO struct {
	Name                   *string  `json:"name"`
	Description            *string  `json:"description" binding:"omitempty"`
	SafetyStock            *float64 `json:"safetyStock" binding:"omitempty,gt=0"`
	UnitID                 *uint    `json:"unitId"`
	CategoryID             *uint    `json:"categoryId"`
	IngredientID           *uint    `json:"ingredientId"`
	Barcode                *string  `json:"barcode"`
	ExpirationPeriodInDays *int     `json:"expirationPeriodInDays"`
	IsActive               *bool    `json:"isActive"`
	Size                   *float64 `json:"size"`
}

type StockMaterialsDTO struct {
	ID                     uint                                                     `json:"id"`
	Name                   string                                                   `json:"name"`
	Description            string                                                   `json:"description"`
	SafetyStock            float64                                                  `json:"safetyStock"`
	Barcode                string                                                   `json:"barcode"`
	IsActive               bool                                                     `json:"isActive"`
	Unit                   unitTypes.UnitsDTO                                       `json:"unit"`
	Category               stockMaterialCategoryTypes.StockMaterialCategoryResponse `json:"category"`
	Ingredient             ingredientTypes.IngredientDTO                            `json:"ingredient"`
	ExpirationPeriodInDays int                                                      `json:"expirationPeriodInDays"`
	Size                   float64                                                  `json:"size"`
	CreatedAt              string                                                   `json:"createdAt"`
	UpdatedAt              string                                                   `json:"updatedAt"`
}

type StockMaterialFilter struct {
	Search           *string `form:"search"`           // Search by name, description, or category
	LowStock         *bool   `form:"lowStock"`         // Filter for materials below safety stock
	IsActive         *bool   `form:"isActive"`         // Filter by active/inactive status
	IngredientID     *uint   `form:"ingredientId"`     // Filter by ingredient
	CategoryID       *uint   `form:"categoryId"`       // Filter by category
	ExpirationInDays *int    `form:"expirationInDays"` // Filter by expiration period in days
	SupplierID       *uint   `form:"supplierId"`
	utils.BaseFilter
}
