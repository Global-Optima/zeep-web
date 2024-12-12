package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

type StockDTO struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	CurrentStock      float64 `json:"currentStock"`
	Unit              string  `json:"unit"`
	LowStockAlert     bool    `json:"lowStockAlert"`
	LowStockThreshold float64 `json:"minimumStockThreshold"`
}

type GetStockQuery struct {
	SearchTerm   *string `json:"searchTerm,omitempty"`
	LowStockOnly *bool   `json:"lowStockAlert,omitempty"`
	Pagination   *utils.Pagination
}

type UpdateStockDTO struct {
	StockId           uint     `json:"stockId"`
	CurrentStock      *float64 `json:"currentStock"`
	LowStockThreshold *float64 `json:"lowStockThreshold"`
}

type AddStockDTO struct {
	IngredientID      uint    `json:"ingredientId" binding:"required"`
	CurrentStock      float64 `json:"currentStock" binding:"required,gte=0"`
	LowStockThreshold float64 `json:"lowStockAlert" binding:"required,gte=0"`
}
