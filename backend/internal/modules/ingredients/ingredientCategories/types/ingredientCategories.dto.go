package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

type CreateIngredientCategoryDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type UpdateIngredientCategoryDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type IngredientCategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type IngredientCategoryFilter struct {
	utils.BaseFilter
	Search *string `form:"search"` // Search by name or description
}
