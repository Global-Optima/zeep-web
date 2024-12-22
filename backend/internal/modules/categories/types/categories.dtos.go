package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

type CategoryDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoriesFilterDTO struct {
	Search string `form:"search" binding:"omitempty,max=50"`
	utils.BaseFilter
}

type CreateCategoryDTO struct {
	Name        string `json:"name" binding:"required,min=2"`
	Description string `json:"description" binding:"required,min=2"`
}

type UpdateCategoryDTO struct {
	Name        string `json:"name,omitempty" binding:"omitempty,min=2"`
	Description string `json:"description,omitempty" binding:"omitempty,min=2"`
}
