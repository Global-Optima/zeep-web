package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

type CreateStockMaterialCategoryDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" sanitize:"soft"`
}

type UpdateStockMaterialCategoryDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description" sanitize:"soft"`
}

type StockMaterialCategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

type StockMaterialCategoryFilter struct {
	utils.BaseFilter
	Search *string `form:"search"`
}
