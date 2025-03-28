package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func MapCategoryToDTO(category data.ProductCategory) *ProductCategoryDTO {
	return &ProductCategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
}

func UpdateToCategory(dto *UpdateProductCategoryDTO, category *data.ProductCategory) {
	if dto.Name != "" {
		category.Name = dto.Name
	}

	if dto.Description != "" {
		category.Description = dto.Description
	}
}
