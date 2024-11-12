package categories

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
)

type CategoryService interface {
	GetCategories() ([]types.CategoryDTO, error)
}

type categoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetCategories() ([]types.CategoryDTO, error) {
	categories, err := s.repo.GetCategories()
	if err != nil {
		return nil, err
	}

	dtos := make([]types.CategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = MapCategoryToDTO(category)
	}

	return dtos, nil
}

func MapCategoryToDTO(category data.ProductCategory) types.CategoryDTO {
	return types.CategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
}
