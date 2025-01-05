package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func MapCategoryToDTO(category data.ProductCategory) *CategoryDTO {
	return &CategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
}

func UpdateToCategory(dto *UpdateCategoryDTO) *data.ProductCategory {
	return &data.ProductCategory{
		Name:        dto.Name,
		Description: dto.Description,
	}
}
