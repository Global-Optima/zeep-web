package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type ProductCategoryDTO struct {
	ID              uint                 `json:"id"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	MachineCategory data.MachineCategory `json:"machineCategory"`
}

type ProductCategoriesFilterDTO struct {
	utils.BaseFilter
	Search string `form:"search" binding:"omitempty,max=50"`
}

type CreateProductCategoryDTO struct {
	Name            string               `json:"name" binding:"required,min=2"`
	Description     *string              `json:"description" binding:"omitempty"`
	MachineCategory data.MachineCategory `json:"machineCategory" binding:"required,oneof=TEA COFFEE ICE_CREAM"`
}

type UpdateProductCategoryDTO struct {
	Name            *string              `json:"name,omitempty" binding:"min=2,omitempty"`
	Description     *string              `json:"description" binding:"omitempty"`
	MachineCategory data.MachineCategory `json:"machineCategory,omitempty" binding:"omitempty,oneof=TEA COFFEE ICE_CREAM"`
}
