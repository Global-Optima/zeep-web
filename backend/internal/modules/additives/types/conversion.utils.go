package types

import (
	"strings"

	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

type AdditiveModels struct {
	Additive    *data.Additive
	Ingredients []data.ProductSizeIngredient
}

func ConvertToAdditiveModel(dto *CreateAdditiveDTO) *data.Additive {
	additive := &data.Additive{
		Name:               dto.Name,
		Description:        dto.Description,
		BasePrice:          dto.BasePrice,
		ImageURL:           dto.ImageURL,
		UnitID:             dto.UnitID,
		Size:               dto.Size,
		AdditiveCategoryID: dto.AdditiveCategoryID,
	}

	for _, ingredient := range dto.Ingredients {
		additive.Ingredients = append(additive.Ingredients, data.AdditiveIngredient{
			IngredientID: ingredient.IngredientID,
			Quantity:     ingredient.Quantity,
		})
	}

	return additive
}

func ConvertToUpdatedAdditiveModels(dto *UpdateAdditiveDTO) *AdditiveModels {
	additive := &data.Additive{}
	if dto == nil {
		return nil
	}

	if strings.TrimSpace(dto.Name) != "" {
		additive.Name = dto.Name
	}
	if strings.TrimSpace(dto.Description) != "" {
		additive.Description = dto.Description
	}
	if dto.BasePrice != nil {
		additive.BasePrice = *dto.BasePrice
	}
	if dto.ImageURL != nil {
		additive.ImageURL = *dto.ImageURL
	}
	if dto.Size != nil {
		additive.Size = *dto.Size
	}
	if dto.UnitID != nil {
		additive.UnitID = *dto.UnitID
	}
	if dto.AdditiveCategoryID != nil {
		additive.AdditiveCategoryID = *dto.AdditiveCategoryID
	}

	var ingredients []data.ProductSizeIngredient

	for _, ingredient := range dto.Ingredients {
		temp := data.ProductSizeIngredient{
			IngredientID: ingredient.IngredientID,
			Quantity:     ingredient.Quantity,
		}
		ingredients = append(ingredients, temp)
	}

	return &AdditiveModels{
		Additive:    additive,
		Ingredients: ingredients,
	}
}

func ConvertToAdditiveCategoryModel(dto *CreateAdditiveCategoryDTO) *data.AdditiveCategory {
	return &data.AdditiveCategory{
		Name:             dto.Name,
		Description:      dto.Description,
		IsMultipleSelect: dto.IsMultipleSelect,
	}
}

func ConvertToUpdatedAdditiveCategoryModel(dto *UpdateAdditiveCategoryDTO, existing *data.AdditiveCategory) *data.AdditiveCategory {
	if dto.Name != nil {
		existing.Name = *dto.Name
	}
	if dto.Description != nil {
		existing.Description = *dto.Description
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
		ID:              additive.ID,
		BaseAdditiveDTO: *ConvertToBaseAdditiveDTO(additive),
	}
}

func ConvertToAdditiveDetailsDTO(additive *data.Additive) *AdditiveDetailsDTO {
	ingredients := make([]ingredientTypes.IngredientDTO, len(additive.Ingredients))
	for i, additiveIngredient := range additive.Ingredients {
		ingredients[i] = *ingredientTypes.ConvertToIngredientResponseDTO(&additiveIngredient.Ingredient)
	}

	return &AdditiveDetailsDTO{
		AdditiveDTO: *ConvertToAdditiveDTO(additive),
		Ingredients: ingredients,
	}
}

func ConvertToBaseAdditiveDTO(additive *data.Additive) *BaseAdditiveDTO {
	return &BaseAdditiveDTO{
		Name:        additive.Name,
		Description: additive.Description,
		BasePrice:   additive.BasePrice,
		ImageURL:    additive.ImageURL,
		Size:        additive.Size,
		Unit:        unitTypes.ToUnitResponse(additive.Unit),
		Category:    *ConvertToCategoryDTO(&additive.Category),
	}
}

func ConvertToCategoryDTO(category *data.AdditiveCategory) *BaseAdditiveCategoryDTO {
	return &BaseAdditiveCategoryDTO{
		ID:               category.ID,
		Name:             category.Name,
		IsMultipleSelect: category.IsMultipleSelect,
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
		ID:                          additive.ID,
		BaseAdditiveCategoryItemDTO: *ConvertToBaseAdditiveCategoryItem(additive, categoryID),
	}
}

func ConvertToBaseAdditiveCategoryItem(additive *data.Additive, categoryID uint) *BaseAdditiveCategoryItemDTO {
	return &BaseAdditiveCategoryItemDTO{
		Name:        additive.Name,
		Description: additive.Description,
		BasePrice:   additive.BasePrice,
		ImageURL:    additive.ImageURL,
		Size:        additive.Size,
		Unit:        unitTypes.ToUnitResponse(additive.Unit),
		CategoryID:  categoryID,
	}
}
