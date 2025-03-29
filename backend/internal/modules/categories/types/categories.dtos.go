package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type ProductCategoryDTO struct {
	ID              uint                 `json:"id"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	MachineCategory data.MachineCategory `json:"machine_category"`
}

type ProductCategoriesFilterDTO struct {
	utils.BaseFilter
	Search string `form:"search" binding:"omitempty,max=50"`
}

type CreateProductCategoryDTO struct {
	Name            string               `json:"name" binding:"required,min=2"`
	Description     string               `json:"description" binding:"required,min=2"`
	MachineCategory data.MachineCategory `json:"machine_category" binding:"required,oneof=TEA COFFEE ICE_CREAM"`
}

type UpdateProductCategoryDTO struct {
	Name            string               `json:"name,omitempty" binding:"omitempty,min=2"`
	Description     string               `json:"description,omitempty" binding:"omitempty,min=2"`
	MachineCategory data.MachineCategory `json:"machine_category,omitempty" binding:"omitempty,oneof=TEA COFFEE ICE_CREAM"`
}
