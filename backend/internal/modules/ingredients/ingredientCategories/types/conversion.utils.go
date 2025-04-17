package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func ToIngredientCategoryResponse(category *data.IngredientCategory) *IngredientCategoryResponse {
	return &IngredientCategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
}
