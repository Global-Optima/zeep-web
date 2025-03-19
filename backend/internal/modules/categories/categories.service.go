package categories

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	"go.uber.org/zap"
)

type CategoryService interface {
	GetCategories(filter *types.ProductCategoriesFilterDTO) ([]types.ProductCategoryDTO, error)
	GetCategoryByID(id uint) (*types.ProductCategoryDTO, error)
	CreateCategory(dto *types.CreateProductCategoryDTO) (uint, error)
	UpdateCategory(id uint, dto *types.UpdateProductCategoryDTO) error
	DeleteCategory(id uint) error
}

type categoryService struct {
	repo   CategoryRepository
	logger *zap.SugaredLogger
}

func NewCategoryService(repo CategoryRepository, logger *zap.SugaredLogger) CategoryService {
	return &categoryService{
		repo:   repo,
		logger: logger,
	}
}

func (s *categoryService) GetCategories(filter *types.ProductCategoriesFilterDTO) ([]types.ProductCategoryDTO, error) {
	categories, err := s.repo.GetCategories(filter)
	if err != nil {
		return nil, err
	}

	dtos := make([]types.ProductCategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = *types.MapCategoryToDTO(category)
	}

	return dtos, nil
}

func (s *categoryService) GetCategoryByID(id uint) (*types.ProductCategoryDTO, error) {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	return types.MapCategoryToDTO(*category), nil
}

func (s *categoryService) CreateCategory(dto *types.CreateProductCategoryDTO) (uint, error) {
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

func (s *categoryService) UpdateCategory(id uint, dto *types.UpdateProductCategoryDTO) error {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		s.logger.Error("Error updating category during finding", zap.Error(err))
		return err
	}

	types.UpdateToCategory(dto, category)

	err = s.repo.UpdateCategory(id, category)
	if err != nil {
		s.logger.Error("Error updating category", zap.Error(err))
		return err
	}

	return nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	return s.repo.DeleteCategory(id)
}
