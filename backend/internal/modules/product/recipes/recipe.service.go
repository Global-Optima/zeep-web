package recipes

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/recipes/types"
	"go.uber.org/zap"
)

type RecipeService interface {
	CreateRecipeStep(dto *types.CreateRecipeStepDTO) (uint, error)
	UpdateRecipeStep(id uint, dto *types.UpdateRecipeStepDTO) error
	GetRecipeStepByID(id uint) (*types.RecipeStepDTO, error)
	GetRecipeStepsByProductID(productID uint) ([]types.RecipeStepDTO, error)
	DeleteRecipeStep(id uint) error
}

type recipeService struct {
	repo   RecipeRepository
	logger *zap.SugaredLogger
}

func NewRecipeService(repo RecipeRepository, logger *zap.SugaredLogger) RecipeService {
	return &recipeService{
		repo:   repo,
		logger: logger,
	}
}

func (s *recipeService) CreateRecipeStep(dto *types.CreateRecipeStepDTO) (uint, error) {
	// Create a RecipeStep from the DTO
	recipeStep := types.CreateToRecipeStepModel(dto)

	// Insert the new RecipeStep into the database using the repository
	id, err := s.repo.CreateRecipeStep(recipeStep)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *recipeService) UpdateRecipeStep(id uint, dto *types.UpdateRecipeStepDTO) error {
	recipeStep := types.UpdateToRecipeStepModel(dto)

	err := s.repo.UpdateRecipeStep(id, recipeStep)
	if err != nil {
		return err
	}

	return nil
}

// GetRecipeStepByID retrieves a recipe step by ID
func (s *recipeService) GetRecipeStepByID(id uint) (*types.RecipeStepDTO, error) {
	recipeStep, err := s.repo.GetRecipeStepByID(id)
	if err != nil {
		wrappedErr := fmt.Errorf("cannot get recipe step by id %d: %w", id, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dto := types.MapRecipeStepToDTO(recipeStep)

	return &dto, nil
}

// GetRecipeStepsByProductID retrieves all recipe steps for a given product
func (s *recipeService) GetRecipeStepsByProductID(productID uint) ([]types.RecipeStepDTO, error) {
	recipeSteps, err := s.repo.GetRecipeStepsByProductID(productID)
	if err != nil {
		return nil, err
	}

	dtos := make([]types.RecipeStepDTO, len(recipeSteps))
	for i, step := range recipeSteps {
		dtos[i] = types.MapRecipeStepToDTO(&step)
	}

	return dtos, nil
}

// DeleteRecipeStep deletes a recipe step by ID
func (s *recipeService) DeleteRecipeStep(id uint) error {
	// Check if the recipe step exists
	recipeStep, err := s.repo.GetRecipeStepByID(id)
	if err != nil {
		return err
	}

	// Delete the recipe step
	err = s.repo.DeleteRecipeStep(recipeStep)
	if err != nil {
		return err
	}

	return nil
}
