package types

import (
	"fmt"
	"strings"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientCategoryTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/ingredientCategories/types"
	unitType "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
)

// Converts CreateIngredientDTO to Ingredient model
func ConvertToIngredientModel(dto *CreateIngredientDTO) (*data.Ingredient, error) {
	if dto == nil {
		return nil, fmt.Errorf("%w: DTO cannot be nil", ErrValidation)
	}

	if strings.TrimSpace(dto.Name) == "" {
		return nil, fmt.Errorf("%w: Ingredient name is required", ErrValidation)
	}

	ingredient := &data.Ingredient{
		Name:             dto.Name,
		CategoryID:       dto.CategoryID,
		UnitID:           dto.UnitID,
		Calories:         dto.Calories,
		Fat:              dto.Fat,
		Carbs:            dto.Carbs,
		Proteins:         dto.Proteins,
		ExpirationInDays: dto.ExpirationInDays,
	}

	return ingredient, nil
}

// Converts UpdateIngredientDTO to Ingredient model
func ConvertToUpdateIngredientModel(dto *UpdateIngredientDTO) (*data.Ingredient, error) {
	if dto == nil {
		return nil, fmt.Errorf("%w: DTO cannot be nil", ErrValidation)
	}

	ingredient := &data.Ingredient{}

	if strings.TrimSpace(dto.Name) != "" {
		ingredient.Name = dto.Name
	}

	if dto.Calories != nil {
		ingredient.Calories = *dto.Calories
	}

	if dto.Fat != nil {
		ingredient.Fat = *dto.Fat
	}

	if dto.Carbs != nil {
		ingredient.Carbs = *dto.Carbs
	}

	if dto.Proteins != nil {
		ingredient.Proteins = *dto.Proteins
	}

	if dto.UnitID != nil {
		if *dto.UnitID == 0 {
			return nil, fmt.Errorf("%w: UnitID must be greater than 0", ErrValidation)
		}
		ingredient.UnitID = *dto.UnitID
	}

	if dto.CategoryID != nil {
		if *dto.CategoryID == 0 {
			return nil, fmt.Errorf("%w: CategoryID must be greater than 0", ErrValidation)
		}
		ingredient.CategoryID = *dto.CategoryID
	}

	if dto.ExpirationInDays != nil {
		if *dto.ExpirationInDays == 0 {
			return nil, fmt.Errorf("%w: Expiration days can not be 0", ErrValidation)
		}
		ingredient.ExpirationInDays = *dto.ExpirationInDays
	}
	return ingredient, nil
}

// Converts Ingredient model to IngredientResponseDTO
func ConvertToIngredientResponseDTO(ingredient *data.Ingredient) *IngredientDTO {
	return &IngredientDTO{
		ID:               ingredient.ID,
		Name:             ingredient.Name,
		Calories:         ingredient.Calories,
		Fat:              ingredient.Fat,
		Carbs:            ingredient.Carbs,
		Proteins:         ingredient.Proteins,
		ExpirationInDays: ingredient.ExpirationInDays,
		Unit: unitType.UnitsDTO{
			ID:               ingredient.Unit.ID,
			Name:             ingredient.Unit.Name,
			ConversionFactor: ingredient.Unit.ConversionFactor,
		},
		Category: ingredientCategoryTypes.IngredientCategoryResponse{
			ID:          ingredient.IngredientCategory.ID,
			Name:        ingredient.IngredientCategory.Name,
			Description: ingredient.IngredientCategory.Description,
		},
	}
}

// Converts list of Ingredient models to list of IngredientResponseDTOs
func ConvertToIngredientResponseDTOs(ingredients []data.Ingredient) []IngredientDTO {
	dtos := make([]IngredientDTO, len(ingredients))
	for i, ingredient := range ingredients {
		dtos[i] = *ConvertToIngredientResponseDTO(&ingredient)
	}
	return dtos
}
