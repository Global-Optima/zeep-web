package types

import (
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type StoreStockDTO struct {
	ID                uint                          `json:"id"`
	Name              string                        `json:"name"`
	Quantity          float64                       `json:"quantity"`
	LowStockAlert     bool                          `json:"lowStockAlert"`
	LowStockThreshold float64                       `json:"lowStockThreshold"`
	Ingredient        ingredientTypes.IngredientDTO `json:"ingredient"`
}

type GetStockFilterQuery struct {
	utils.BaseFilter
	Search       *string `form:"search"`
	LowStockOnly *bool   `form:"lowStockOnly"`
}

type UpdateStoreStockDTO struct {
	Quantity          *float64 `json:"quantity" binding:"omitempty,gte=0"`
	LowStockThreshold *float64 `json:"lowStockThreshold" binding:"omitempty,gte=0"`
}

type AddMultipleStoreStockDTO struct {
	IngredientStocks []AddStoreStockDTO `json:"ingredientStocks" binding:"required,dive"`
}

type AddStoreStockDTO struct {
	IngredientID      uint    `json:"ingredientId" binding:"required,gt=0"`
	Quantity          float64 `json:"quantity" binding:"required,gte=0"`
	LowStockThreshold float64 `json:"lowStockThreshold" binding:"required,gt=0"`
}
