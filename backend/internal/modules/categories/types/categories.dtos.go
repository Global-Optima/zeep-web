package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

type ProductCategoryDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProductCategoriesFilterDTO struct {
	utils.BaseFilter
	Search string `form:"search" binding:"omitempty,max=50"`
}

type CreateProductCategoryDTO struct {
	Name        string `json:"name" binding:"required,min=2"`
	Description string `json:"description" binding:"required,min=2"`
}

type UpdateProductCategoryDTO struct {
	Name        string `json:"name,omitempty" binding:"omitempty,min=2"`
	Description string `json:"description,omitempty" binding:"omitempty,min=2"`
}
