package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ToStockMaterialCategoryResponse(category data.StockMaterialCategory) StockMaterialCategoryResponse {
	return StockMaterialCategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
