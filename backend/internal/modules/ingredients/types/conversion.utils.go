package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
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
		Name:                 dto.Name,
		IngredientCategoryID: dto.CategoryID,
		UnitID:               dto.UnitID,
		Calories:             dto.Calories,
		Fat:                  dto.Fat,
		Carbs:                dto.Carbs,
		Proteins:             dto.Proteins,
	}

	if dto.ExpiresAt != nil {
		parsedTime, err := time.Parse(time.RFC3339, *dto.ExpiresAt)
		if err != nil {
			return nil, err
		}
		ingredient.ExpiresAt = &parsedTime
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
		ingredient.IngredientCategoryID = *dto.CategoryID
	}

	if dto.ExpiresAt != nil {
		parsedTime, err := time.Parse(time.RFC3339, *dto.ExpiresAt)
		if err != nil {
			return nil, err
		}
		ingredient.ExpiresAt = &parsedTime
	}
	return ingredient, nil
}

// Converts Ingredient model to IngredientResponseDTO
func ConvertToIngredientResponseDTO(ingredient *data.Ingredient) *IngredientResponseDTO {
	var expiresAt *string
	if ingredient.ExpiresAt != nil {
		formattedTime := ingredient.ExpiresAt.Format(time.RFC3339)
		expiresAt = &formattedTime
	}

	return &IngredientResponseDTO{
		ID:        ingredient.ID,
		Name:      ingredient.Name,
		Calories:  ingredient.Calories,
		Fat:       ingredient.Fat,
		Carbs:     ingredient.Carbs,
		Proteins:  ingredient.Proteins,
		ExpiresAt: expiresAt,
	}
}

// Converts list of Ingredient models to list of IngredientResponseDTOs
func ConvertToIngredientResponseDTOs(ingredients []data.Ingredient) []IngredientResponseDTO {
	dtos := make([]IngredientResponseDTO, len(ingredients))
	for i, ingredient := range ingredients {
		dtos[i] = *ConvertToIngredientResponseDTO(&ingredient)
	}
	return dtos
}

func ConvertToIngredientDetailsDTO(ingredient *data.Ingredient) IngredientDetailsDTO {
	return IngredientDetailsDTO{
		IngredientResponseDTO: *ConvertToIngredientResponseDTO(ingredient),
		Category: struct {
			ID          uint   `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		}{ID: ingredient.IngredientCategory.ID, Name: ingredient.IngredientCategory.Name, Description: ingredient.IngredientCategory.Description},
	}
}
