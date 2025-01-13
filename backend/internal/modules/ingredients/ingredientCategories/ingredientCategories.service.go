package ingredientCategories

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/ingredientCategories/types"
)

type IngredientCategoryService interface {
	Create(dto types.CreateIngredientCategoryDTO) (uint, error)
	GetByID(id uint) (*types.IngredientCategoryResponse, error)
	Update(id uint, dto types.UpdateIngredientCategoryDTO) error
	Delete(id uint) error
	GetAll() ([]types.IngredientCategoryResponse, error)
}

type ingredientCategoryService struct {
	repo IngredientCategoryRepository
}

func NewIngredientCategoryService(repo IngredientCategoryRepository) IngredientCategoryService {
	return &ingredientCategoryService{repo: repo}
}

func (s *ingredientCategoryService) Create(dto types.CreateIngredientCategoryDTO) (uint, error) {
	category := &data.IngredientCategory{
		Name:        dto.Name,
		Description: dto.Description,
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

func (s *ingredientCategoryService) Update(id uint, dto types.UpdateIngredientCategoryDTO) error {
	updates := data.IngredientCategory{}
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

func (s *ingredientCategoryService) GetAll() ([]types.IngredientCategoryResponse, error) {
	categories, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ingredient categories: %w", err)
	}

	responses := make([]types.IngredientCategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = types.IngredientCategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}
	}
	return responses, nil
}
