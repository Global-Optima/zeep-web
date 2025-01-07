package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

// Converts CreateIngredientDTO to Ingredient model
func ConvertToIngredientModel(dto *CreateIngredientDTO) (*data.Ingredient, error) {
	var expiresAt *time.Time
	if dto.ExpiresAt != nil {
		parsedTime, err := time.Parse(time.RFC3339, *dto.ExpiresAt)
		if err != nil {
			return nil, err
		}
		expiresAt = &parsedTime
	}

	return &data.Ingredient{
		Name:      dto.Name,
		Calories:  dto.Calories,
		Fat:       dto.Fat,
		Carbs:     dto.Carbs,
		Proteins:  dto.Proteins,
		ExpiresAt: expiresAt,
	}, nil
}

// Converts UpdateIngredientDTO to Ingredient model
func ConvertToUpdateIngredientModel(dto *UpdateIngredientDTO, existing *data.Ingredient) (*data.Ingredient, error) {
	if dto.Name != nil && *dto.Name != "" {
		existing.Name = *dto.Name
	}
	if dto.Calories != nil {
		existing.Calories = *dto.Calories
	}
	if dto.Fat != nil {
		existing.Fat = *dto.Fat
	}
	if dto.Carbs != nil {
		existing.Carbs = *dto.Carbs
	}
	if dto.Proteins != nil {
		existing.Proteins = *dto.Proteins
	}
	if dto.ExpiresAt != nil {
		parsedTime, err := time.Parse(time.RFC3339, *dto.ExpiresAt)
		if err != nil {
			return nil, err
		}
		existing.ExpiresAt = &parsedTime
	}
	return existing, nil
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
