package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func CreateToRecipeStepModel(productID uint, dto []CreateOrReplaceRecipeStepDTO) []data.RecipeStep {
	recipeSteps := make([]data.RecipeStep, len(dto))
	for i, step := range dto {
		recipeStep := &data.RecipeStep{
			ProductID: productID,
		}
		recipeStep.Step = step.Step
		recipeStep.Name = step.Name
		recipeStep.Description = step.Description
		recipeStep.ImageKey = step.ImageURL
		recipeSteps[i] = *recipeStep
	}
	return recipeSteps
}

func UpdateToRecipeStepModel(dto *CreateOrReplaceRecipeStepDTO) *data.RecipeStep {
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
		recipeStep.ImageKey = dto.ImageURL
	}
	return recipeStep
}

func MapRecipeStepToDTO(recipeStep *data.RecipeStep) RecipeStepDTO {
	return RecipeStepDTO{
		ID:          recipeStep.ID,
		Name:        recipeStep.Name,
		ProductID:   recipeStep.ProductID,
		Step:        recipeStep.Step,
		Description: recipeStep.Description,
		ImageURL:    recipeStep.ImageKey,
	}
}
