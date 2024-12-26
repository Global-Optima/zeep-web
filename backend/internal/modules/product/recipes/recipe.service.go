package recipes

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/recipes/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type RecipeService interface {
	CreateOrReplaceRecipeSteps(productID uint, dtos []types.CreateOrReplaceRecipeStepDTO) ([]uint, error)
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

func (s *recipeService) CreateOrReplaceRecipeSteps(productID uint, dtos []types.CreateOrReplaceRecipeStepDTO) ([]uint, error) {

	recipeSteps := types.CreateToRecipeStepModel(productID, dtos)

	ids, err := s.repo.CreateOrReplaceRecipeStepsByProductID(productID, recipeSteps)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create or replace steps for productID %d: %w", productID, err)
		s.logger.Error(wrappedErr)
		return nil, err
	}

	return ids, nil
}

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

func (s *recipeService) GetRecipeStepsByProductID(productID uint) ([]types.RecipeStepDTO, error) {
	recipeSteps, err := s.repo.GetRecipeStepsByProductID(productID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get recipe steps", err)
		s.logger.Error(wrappedErr)
		return nil, err
	}

	dtos := make([]types.RecipeStepDTO, len(recipeSteps))
	for i, step := range recipeSteps {
		dtos[i] = types.MapRecipeStepToDTO(&step)
	}

	return dtos, nil
}

func (s *recipeService) DeleteRecipeStep(productID uint) error {
	err := s.repo.DeleteRecipeStepsByProductID(productID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to delete recipe steps", err)
		s.logger.Error(wrappedErr)
		return err
	}

	return nil
}
