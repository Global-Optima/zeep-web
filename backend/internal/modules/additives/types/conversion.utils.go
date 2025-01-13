package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"strings"
)

func ConvertToAdditiveModel(dto *CreateAdditiveDTO) *data.Additive {
	return &data.Additive{
		Name:               dto.Name,
		Description:        dto.Description,
		BasePrice:          dto.Price,
		ImageURL:           dto.ImageURL,
		Size:               dto.Size,
		AdditiveCategoryID: dto.AdditiveCategoryID,
	}
}

func ConvertToUpdatedAdditiveModel(dto *UpdateAdditiveDTO) *data.Additive {
	additive := &data.Additive{}
	if dto == nil {
		return additive
	}

	if strings.TrimSpace(dto.Name) != "" {
		additive.Name = dto.Name
	}
	if strings.TrimSpace(dto.Description) != "" {
		additive.Description = dto.Description
	}
	if dto.Price != nil {
		additive.BasePrice = *dto.Price
	}
	if dto.ImageURL != nil {
		additive.ImageURL = *dto.ImageURL
	}
	if dto.Size != nil {
		additive.Size = *dto.Size
	}
	if dto.AdditiveCategoryID != nil {
		additive.AdditiveCategoryID = *dto.AdditiveCategoryID
	}

	return additive
}

func ConvertToAdditiveCategoryModel(dto *CreateAdditiveCategoryDTO) *data.AdditiveCategory {
	return &data.AdditiveCategory{
		Name:             dto.Name,
		Description:      dto.Description,
		IsMultipleSelect: dto.IsMultipleSelect,
	}
}

func ConvertToUpdatedAdditiveCategoryModel(dto *UpdateAdditiveCategoryDTO, existing *data.AdditiveCategory) *data.AdditiveCategory {
	if dto.Name != "" {
		existing.Name = dto.Name
	}
	if dto.Description != "" {
		existing.Description = dto.Description
	}
	if dto.IsMultipleSelect != nil {
		existing.IsMultipleSelect = *dto.IsMultipleSelect
	}
	return existing
}

func ConvertToAdditiveCategoryResponseDTO(model *data.AdditiveCategory) *AdditiveCategoryResponseDTO {
	return &AdditiveCategoryResponseDTO{
		ID:               model.ID,
		Name:             model.Name,
		Description:      model.Description,
		IsMultipleSelect: model.IsMultipleSelect,
	}
}

func ConvertToAdditiveDTO(additive *data.Additive) *AdditiveDTO {
	return &AdditiveDTO{
		ID:          additive.ID,
		Name:        additive.Name,
		Description: additive.Description,
		BasePrice:   additive.BasePrice,
		ImageURL:    additive.ImageURL,
		Size:        additive.Size,
		Category: struct {
			ID               uint   `json:"id"`
			Name             string `json:"name"`
			IsMultipleSelect bool   `json:"isMultipleSelect"`
		}{
			ID:               additive.Category.ID,
			Name:             additive.Category.Name,
			IsMultipleSelect: additive.Category.IsMultipleSelect,
		},
	}
}

func ConvertToAdditiveCategoryDTO(category *data.AdditiveCategory) *AdditiveCategoryDTO {
	additives := ConvertToAdditiveCategoryItemDTOs(category)

	return &AdditiveCategoryDTO{
		ID:               category.ID,
		Name:             category.Name,
		Description:      category.Description,
		IsMultipleSelect: category.IsMultipleSelect,
		Additives:        additives, // Always initialized as a slice
	}
}

func ConvertToAdditiveCategoryItemDTOs(category *data.AdditiveCategory) []AdditiveCategoryItemDTO {
	additives := make([]AdditiveCategoryItemDTO, 0)

	// Populate additives if present
	for _, additive := range category.Additives {
		additives = append(additives, *ConvertToAdditiveCategoryItem(&additive, category.ID))
	}

	return additives
}

func ConvertToAdditiveCategoryItem(additive *data.Additive, categoryID uint) *AdditiveCategoryItemDTO {
	return &AdditiveCategoryItemDTO{
		ID:          additive.ID,
		Name:        additive.Name,
		Description: additive.Description,
		Price:       additive.BasePrice,
		ImageURL:    additive.ImageURL,
		Size:        additive.Size,
		CategoryID:  categoryID,
	}
}
