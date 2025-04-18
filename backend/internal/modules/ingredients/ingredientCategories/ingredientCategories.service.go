package ingredientCategories

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/ingredientCategories/types"
)

type IngredientCategoryService interface {
	Create(dto types.CreateIngredientCategoryDTO) (uint, error)
	GetByID(id uint) (*types.IngredientCategoryResponse, error)
	GetTranslatedByID(locale data.LanguageCode, id uint) (*types.IngredientCategoryResponse, error)
	Update(id uint, dto types.UpdateIngredientCategoryDTO) error
	Delete(id uint) error
	GetAll(locale data.LanguageCode, filter *types.IngredientCategoryFilter) ([]types.IngredientCategoryResponse, error)

	UpsertIngredientCategoryTranslations(id uint, dto *types.IngredientCategoryTranslationDTO) error
}

type ingredientCategoryService struct {
	repo               IngredientCategoryRepository
	transactionManager TransactionManager
}

func NewIngredientCategoryService(repo IngredientCategoryRepository, transactionManager TransactionManager) IngredientCategoryService {
	return &ingredientCategoryService{
		repo:               repo,
		transactionManager: transactionManager,
	}
}

func (s *ingredientCategoryService) Create(dto types.CreateIngredientCategoryDTO) (uint, error) {
	category := &data.IngredientCategory{
		Name:        dto.Name,
		Description: *dto.Description,
	}
	if err := s.repo.Create(category); err != nil {
		return 0, fmt.Errorf("failed to create ingredient category: %w", err)
	}
	return category.ID, nil
}

func (s *ingredientCategoryService) GetByID(id uint) (*types.IngredientCategoryResponse, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ingredient category: %w", err)
	}
	return &types.IngredientCategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (s *ingredientCategoryService) GetTranslatedByID(locale data.LanguageCode, id uint) (*types.IngredientCategoryResponse, error) {
	category, err := s.repo.GetTranslatedByID(locale, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ingredient category: %w", err)
	}

	return types.ConvertToIngredientCategoryResponse(category), nil
}

func (s *ingredientCategoryService) Update(id uint, dto types.UpdateIngredientCategoryDTO) error {
	updates, err := s.repo.GetByID(id)
	if err != nil {
		return types.ErrFailedToFetchIngredientCategory
	}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	if dto.Description != nil {
		updates.Description = *dto.Description
	}
	return s.repo.Update(id, updates)
}

func (s *ingredientCategoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ingredientCategoryService) GetAll(locale data.LanguageCode, filter *types.IngredientCategoryFilter) ([]types.IngredientCategoryResponse, error) {
	categories, err := s.repo.GetAll(locale, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ingredient categories: %w", err)
	}

	responses := make([]types.IngredientCategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = *types.ConvertToIngredientCategoryResponse(&category)
	}
	return responses, nil
}

func (s *ingredientCategoryService) UpsertIngredientCategoryTranslations(id uint, dto *types.IngredientCategoryTranslationDTO) error {
	if dto == nil {
		return fmt.Errorf("translations DTO is nil")
	}

	if err := s.transactionManager.UpsertIngredientCategoryTranslations(id, dto); err != nil {
		wrappedErr := fmt.Errorf("failed to upsert ingredient category translations: %w", err)
		return wrappedErr
	}

	return nil
}
