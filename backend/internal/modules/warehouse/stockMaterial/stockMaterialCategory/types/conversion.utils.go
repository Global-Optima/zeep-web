package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ToStockMaterialCategoryResponse(category data.StockMaterialCategory) StockMaterialCategoryResponse {
	return StockMaterialCategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   category.UpdatedAt.Format(time.RFC3339),
	}
}
