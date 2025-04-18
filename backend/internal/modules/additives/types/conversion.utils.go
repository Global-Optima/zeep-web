package types

import (
	provisionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"
	"github.com/pkg/errors"

	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

type AdditiveModels struct {
	Additive    *data.Additive
	Ingredients []data.AdditiveIngredient
	Provisions  []data.AdditiveProvision
}

func ConvertToAdditiveModel(dto *CreateAdditiveDTO) *data.Additive {
	additive := &data.Additive{
		Name:               dto.Name,
		Description:        *dto.Description,
		BasePrice:          dto.BasePrice,
		UnitID:             dto.UnitID,
		Size:               dto.Size,
		AdditiveCategoryID: dto.AdditiveCategoryID,
		MachineId:          dto.MachineId,
	}

	additive.Ingredients = mapIngredientsToAdditiveIngredients(dto.Ingredients)
	additive.AdditiveProvisions = mapProvisionsToAdditiveProvisions(dto.Provisions)

	return additive
}

func ConvertToUpdatedAdditiveModels(dto *UpdateAdditiveDTO, additive *data.Additive) (*AdditiveModels, error) {
	if dto == nil {
		return nil, errors.New("dto cannot be nil")
	}

	if dto.Name != nil {
		additive.Name = *dto.Name
	}
	if dto.Description != nil {
		additive.Description = *dto.Description
	}
	if dto.BasePrice != nil {
		additive.BasePrice = *dto.BasePrice
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
	if dto.MachineId != nil {
		additive.MachineId = *dto.MachineId
	}

	ingredients := mapIngredientsToAdditiveIngredients(dto.Ingredients)
	provisions := mapProvisionsToAdditiveProvisions(dto.Provisions)

	return &AdditiveModels{
		Additive:    additive,
		Ingredients: ingredients,
		Provisions:  provisions,
	}, nil
}

func mapIngredientsToAdditiveIngredients(ingredientsDTO []ingredientTypes.SelectedIngredientDTO) []data.AdditiveIngredient {
	if ingredientsDTO == nil {
		return nil
	}

	ingredients := make([]data.AdditiveIngredient, 0, len(ingredientsDTO))

	for _, dto := range ingredientsDTO {
		ingredient := data.AdditiveIngredient{
			IngredientID: dto.IngredientID,
			Quantity:     dto.Quantity,
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients
}

func mapProvisionsToAdditiveProvisions(provisionsDTO []provisionsTypes.SelectedProvisionDTO) []data.AdditiveProvision {
	if provisionsDTO == nil {
		return nil
	}

	provisions := make([]data.AdditiveProvision, 0, len(provisionsDTO))

	for _, dto := range provisionsDTO {
		provision := data.AdditiveProvision{
			ProvisionID: dto.ProvisionID,
			Volume:      dto.Volume,
		}
		provisions = append(provisions, provision)
	}

	return provisions
}

func ConvertToAdditiveCategoryModel(dto *CreateAdditiveCategoryDTO) *data.AdditiveCategory {
	return &data.AdditiveCategory{
		Name:             dto.Name,
		Description:      *dto.Description,
		IsMultipleSelect: dto.IsMultipleSelect,
		IsRequired:       dto.IsRequired,
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
	if dto.IsRequired != nil {
		existing.IsRequired = *dto.IsRequired
	}
	return existing
}

func ConvertToAdditiveCategoryDTO(model *data.AdditiveCategory) *AdditiveCategoryDTO {
	return &AdditiveCategoryDTO{
		ID:                      model.ID,
		BaseAdditiveCategoryDTO: *ConvertToBaseAdditiveCategoryDTO(model),
	}
}

func ConvertToAdditiveDTO(additive *data.Additive) *AdditiveDTO {
	return &AdditiveDTO{
		ID:              additive.ID,
		BaseAdditiveDTO: *ConvertToBaseAdditiveDTO(additive),
	}
}

func ConvertToAdditiveIngredientDTO(additiveIngredient *data.AdditiveIngredient) *AdditiveIngredientDTO {
	return &AdditiveIngredientDTO{
		Ingredient: *ingredientTypes.ConvertToIngredientResponseDTO(&additiveIngredient.Ingredient),
		Quantity:   additiveIngredient.Quantity,
	}
}

func ConvertToAdditiveProvisionDTO(additiveProvision *data.AdditiveProvision) *AdditiveProvisionDTO {
	return &AdditiveProvisionDTO{
		Provision: *provisionsTypes.MapToProvisionDTO(&additiveProvision.Provision),
		Volume:    additiveProvision.Volume,
	}
}

func ConvertToAdditiveDetailsDTO(additive *data.Additive) *AdditiveDetailsDTO {
	ingredients := make([]AdditiveIngredientDTO, len(additive.Ingredients))
	provisions := make([]AdditiveProvisionDTO, len(additive.AdditiveProvisions))

	for i, additiveIngredient := range additive.Ingredients {
		ingredients[i] = *ConvertToAdditiveIngredientDTO(&additiveIngredient)
	}

	for i, additiveProvision := range additive.AdditiveProvisions {
		provisions[i] = *ConvertToAdditiveProvisionDTO(&additiveProvision)
	}

	return &AdditiveDetailsDTO{
		AdditiveDTO: *ConvertToAdditiveDTO(additive),
		Ingredients: ingredients,
		Provisions:  provisions,
	}
}

func ConvertToBaseAdditiveDTO(additive *data.Additive) *BaseAdditiveDTO {
	return &BaseAdditiveDTO{
		Name:        additive.Name,
		Description: additive.Description,
		BasePrice:   additive.BasePrice,
		ImageURL:    additive.ImageKey.GetURL(),
		Size:        additive.Size,
		Unit:        unitTypes.ToUnitResponse(additive.Unit),
		Category:    *ConvertToAdditiveCategoryDTO(&additive.Category),
		MachineId:   additive.MachineId,
	}
}

func ConvertToBaseAdditiveCategoryDTO(category *data.AdditiveCategory) *BaseAdditiveCategoryDTO {
	return &BaseAdditiveCategoryDTO{
		Name:             category.Name,
		Description:      category.Description,
		IsMultipleSelect: category.IsMultipleSelect,
		IsRequired:       category.IsRequired,
	}
}

func ConvertToAdditiveCategoryDetailsDTO(category *data.AdditiveCategory) *AdditiveCategoryDetailsDTO {
	additivesCount := 0
	if category.Additives != nil {
		additivesCount = len(category.Additives)
	}

	return &AdditiveCategoryDetailsDTO{
		AdditiveCategoryDTO: *ConvertToAdditiveCategoryDTO(category),
		AdditivesCount:      additivesCount,
	}
}
