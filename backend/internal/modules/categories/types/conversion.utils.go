package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func MapCategoryToDTO(category data.ProductCategory) *ProductCategoryDTO {
	return &ProductCategoryDTO{
		ID:              category.ID,
		Name:            category.Name,
		Description:     category.Description,
		MachineCategory: category.MachineCategory,
	}
}

func UpdateToCategory(dto *UpdateProductCategoryDTO, category *data.ProductCategory) {
	if dto.Name != nil {
		category.Name = *dto.Name
	}

	if dto.Description != nil {
		category.Description = *dto.Description
	}

	if dto.MachineCategory != nil {
		category.MachineCategory = *dto.MachineCategory
	}
}
