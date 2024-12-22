package categories

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
)

type CategoryService interface {
	GetCategories(filter *types.CategoriesFilterDTO) ([]types.CategoryDTO, error)
	GetCategoryByID(id uint) (*types.CategoryDTO, error)
	CreateCategory(dto *types.CreateCategoryDTO) (uint, error)
	UpdateCategory(id uint, dto *types.UpdateCategoryDTO) error
	DeleteCategory(id uint) error
}

type categoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

func (s *categoryService) GetCategories(filter *types.CategoriesFilterDTO) ([]types.CategoryDTO, error) {
	categories, err := s.repo.GetCategories(filter)
	if err != nil {
		return nil, err
	}

	dtos := make([]types.CategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = *types.MapCategoryToDTO(category)
	}

	return dtos, nil
}

func (s *categoryService) GetCategoryByID(id uint) (*types.CategoryDTO, error) {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	return types.MapCategoryToDTO(*category), nil
}

func (s *categoryService) CreateCategory(dto *types.CreateCategoryDTO) (uint, error) {
	category := data.ProductCategory{
		Name:        dto.Name,
		Description: dto.Description,
	}

	id, err := s.repo.CreateCategory(&category)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *categoryService) UpdateCategory(id uint, dto *types.UpdateCategoryDTO) error {
	category := types.UpdateToCategory(dto)

	err := s.repo.UpdateCategory(id, category)
	if err != nil {
		return err
	}

	return nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	return s.repo.DeleteCategory(id)
}
