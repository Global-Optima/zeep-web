package ingredients

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type IngredientService interface {
	CreateIngredient(dto *types.CreateIngredientDTO) error
	UpdateIngredient(ingredientID uint, dto *types.UpdateIngredientDTO) error
	DeleteIngredient(ingredientID uint) error
	GetIngredientByID(ingredientID uint) (*types.IngredientDTO, error)
	GetIngredients(filter *types.IngredientFilter) ([]types.IngredientDTO, error)
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
		wrappedErr := utils.WrapError("Failed to create ingredient:", err)
		s.logger.Error(wrappedErr)
		return err
	}
	return nil
}

func (s *ingredientService) UpdateIngredient(ingredientID uint, dto *types.UpdateIngredientDTO) error {
	ingredient, err := types.ConvertToUpdateIngredientModel(dto)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateIngredient(ingredientID, ingredient); err != nil {
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

func (s *ingredientService) GetIngredientByID(ingredientID uint) (*types.IngredientDTO, error) {
	ingredient, err := s.repo.GetIngredientByID(ingredientID)
	if err != nil {
		s.logger.Error("Failed to fetch ingredient by ID:", err)
		return nil, err
	}

	return types.ConvertToIngredientResponseDTO(ingredient), nil
}

func (s *ingredientService) GetIngredients(filter *types.IngredientFilter) ([]types.IngredientDTO, error) {
	ingredients, err := s.repo.GetIngredients(filter)
	if err != nil {
		s.logger.Error("Failed to fetch ingredients:", err)
		return nil, err
	}

	return types.ConvertToIngredientResponseDTOs(ingredients), nil
}
