package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

func MapCategoryToDTO(category data.ProductCategory) *ProductCategoryDTO {
	name := category.Name
	description := category.Description

	if len(category.NameTranslation) > 0 {
		name = category.NameTranslation[0].TranslatedText
	}
	if len(category.DescriptionTranslation) > 0 {
		description = category.DescriptionTranslation[0].TranslatedText
	}

	return &ProductCategoryDTO{
		ID:              category.ID,
		Name:            name,
		Description:     description,
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
