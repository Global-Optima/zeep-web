package ingredients

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"go.uber.org/zap"
)

type IngredientService interface {
	CreateIngredient(dto *types.CreateIngredientDTO) error
	UpdateIngredient(ingredientID uint, dto *types.UpdateIngredientDTO) error
	DeleteIngredient(ingredientID uint) error
	GetIngredientByID(ingredientID uint) (*data.Ingredient, error)
	GetIngredients(filter *types.IngredientFilter) ([]data.Ingredient, error)
}

type ingredientService struct {
	repo   IngredientRepository
	logger *zap.SugaredLogger
}

func NewIngredientService(repo IngredientRepository, logger *zap.SugaredLogger) IngredientService {
	return &ingredientService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ingredientService) CreateIngredient(dto *types.CreateIngredientDTO) error {
	ingredient, err := types.ConvertToIngredientModel(dto)
	if err != nil {
		return err
	}

	if err := s.repo.CreateIngredient(ingredient); err != nil {
		s.logger.Error("Failed to create ingredient:", err)
		return err
	}
	return nil
}

func (s *ingredientService) UpdateIngredient(ingredientID uint, dto *types.UpdateIngredientDTO) error {
	existing, err := s.repo.GetIngredientByID(ingredientID)
	if err != nil {
		return err
	}

	ingredient, err := types.ConvertToUpdateIngredientModel(dto, existing)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateIngredient(ingredient); err != nil {
		s.logger.Error("Failed to update ingredient:", err)
		return err
	}
	return nil
}

func (s *ingredientService) DeleteIngredient(ingredientID uint) error {
	if err := s.repo.DeleteIngredient(ingredientID); err != nil {
		s.logger.Error("Failed to delete ingredient:", err)
		return err
	}
	return nil
}

func (s *ingredientService) GetIngredientByID(ingredientID uint) (*data.Ingredient, error) {
	ingredient, err := s.repo.GetIngredientByID(ingredientID)
	if err != nil {
		s.logger.Error("Failed to fetch ingredient by ID:", err)
		return nil, err
	}
	return ingredient, nil
}

func (s *ingredientService) GetIngredients(filter *types.IngredientFilter) ([]data.Ingredient, error) {
	ingredients, err := s.repo.GetIngredients(filter)
	if err != nil {
		s.logger.Error("Failed to fetch ingredients:", err)
		return nil, err
	}
	return ingredients, nil
}
