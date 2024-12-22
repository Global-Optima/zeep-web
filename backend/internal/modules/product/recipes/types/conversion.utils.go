package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func CreateToRecipeStepModel(dto *CreateRecipeStepDTO) *data.RecipeStep {
	return &data.RecipeStep{
		ProductID:   dto.ProductID,
		Step:        dto.Step,
		Name:        dto.Name,
		Description: dto.Description,
		ImageURL:    dto.ImageURL,
	}
}

func UpdateToRecipeStepModel(dto *UpdateRecipeStepDTO) *data.RecipeStep {
	recipeStep := &data.RecipeStep{}

	if dto.Step != 0 {
		recipeStep.Step = dto.Step
	}
	if dto.Name != "" {
		recipeStep.Name = dto.Name
	}
	if dto.Description != "" {
		recipeStep.Description = dto.Description
	}
	if dto.ImageURL != "" {
		recipeStep.ImageURL = dto.ImageURL
	}
	return recipeStep
}

func MapRecipeStepToDTO(recipeStep *data.RecipeStep) RecipeStepDTO {
	return RecipeStepDTO{
		ID:          recipeStep.ID,
		ProductID:   recipeStep.ProductID,
		Step:        recipeStep.Step,
		Description: recipeStep.Description,
		ImageURL:    recipeStep.ImageURL,
	}
}
