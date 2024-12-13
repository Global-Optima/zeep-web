package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

type StockDTO struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	Quantity          float64 `json:"quantity"`
	Unit              string  `json:"unit"`
	LowStockAlert     bool    `json:"lowStockAlert"`
	LowStockThreshold float64 `json:"minimumStockThreshold"`
}

type GetStockQuery struct {
	Search       *string `json:"search,omitempty"`
	LowStockOnly *bool   `json:"lowStockAlert,omitempty"`
	Pagination   *utils.Pagination
}

type UpdateStockDTO struct {
	Quantity          *float64 `json:"quantity"`
	LowStockThreshold *float64 `json:"lowStockThreshold"`
}

type AddStockDTO struct {
	IngredientID      uint    `json:"ingredientId" binding:"required,gt=0"`
	Quantity          float64 `json:"quantity" binding:"required,gte=0"`
	LowStockThreshold float64 `json:"lowStockAlert" binding:"required,gte=0"`
}
