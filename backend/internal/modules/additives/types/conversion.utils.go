package types

import "github.com/Global-Optima/zeep-web/backend/internal/data"

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

func ConvertToUpdatedAdditiveModel(dto *UpdateAdditiveDTO, existing *data.Additive) *data.Additive {
	if dto.Name != "" {
		existing.Name = dto.Name
	}
	if dto.Description != "" {
		existing.Description = dto.Description
	}
	if dto.Price != nil {
		existing.BasePrice = *dto.Price
	}
	if dto.ImageURL != nil {
		existing.ImageURL = *dto.ImageURL
	}
	if dto.Size != nil {
		existing.Size = *dto.Size
	}
	if dto.AdditiveCategoryID != nil {
		existing.AdditiveCategoryID = *dto.AdditiveCategoryID
	}
	return existing
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
		Price:       additive.BasePrice,
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

func ConvertToStoreAdditiveDTO(storeAdditive *data.StoreAdditive) *StoreAdditiveDTO {
	return &StoreAdditiveDTO{
		AdditiveDTO: *ConvertToAdditiveDTO(&storeAdditive.Additive),
		StorePrice:  storeAdditive.Price,
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

func ConvertToStoreAdditiveCategoryDTO(category *data.AdditiveCategory) *StoreAdditiveCategoryDTO {
	storeAdditives := ConvertToStoreAdditiveCategoryItemDTOs(category)

	return &StoreAdditiveCategoryDTO{
		ID:               category.ID,
		Name:             category.Name,
		Description:      category.Description,
		IsMultipleSelect: category.IsMultipleSelect,
		Additives:        storeAdditives, // Always initialized as a slice
	}
}

func ConvertToAdditiveCategoryItemDTOs(category *data.AdditiveCategory) []AdditiveCategoryItemDTO {
	additives := make([]AdditiveCategoryItemDTO, 0)

	// Populate additives if present
	for _, additive := range category.Additives {
		additives = append(additives, *convertToAdditiveCategoryItem(&additive, category.ID))
	}

	return additives
}

func ConvertToStoreAdditiveCategoryItemDTOs(category *data.AdditiveCategory) []StoreAdditiveCategoryItemDTO {
	storeAdditives := make([]StoreAdditiveCategoryItemDTO, len(category.Additives))

	// Populate additives if present
	for i, additive := range category.Additives {
		storeAdditives[i] = StoreAdditiveCategoryItemDTO{
			AdditiveCategoryItemDTO: *convertToAdditiveCategoryItem(&additive, category.ID),
			StorePrice:              additive.StoreAdditives[0].Price,
		}
	}

	return storeAdditives
}

func convertToAdditiveCategoryItem(additive *data.Additive, categoryID uint) *AdditiveCategoryItemDTO {
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
