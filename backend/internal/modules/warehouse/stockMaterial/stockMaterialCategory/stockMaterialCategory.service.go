package stockMaterialCategory

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialCategory/types"
	"go.uber.org/zap"
)

type StockMaterialCategoryService interface {
	Create(dto types.CreateStockMaterialCategoryDTO) (uint, error)
	GetByID(id uint) (*types.StockMaterialCategoryResponse, error)
	GetAll(filter types.StockMaterialCategoryFilter) ([]types.StockMaterialCategoryResponse, error)
	Update(id uint, dto types.UpdateStockMaterialCategoryDTO) error
	Delete(id uint) error
}

type stockMaterialCategoryService struct {
	repo   StockMaterialCategoryRepository
	logger *zap.SugaredLogger
}

func NewStockMaterialCategoryService(repo StockMaterialCategoryRepository, logger *zap.SugaredLogger) StockMaterialCategoryService {
	return &stockMaterialCategoryService{
		repo:   repo,
		logger: logger,
	}
}

func (s *stockMaterialCategoryService) Create(dto types.CreateStockMaterialCategoryDTO) (uint, error) {
	category := &data.StockMaterialCategory{
		Name:        dto.Name,
		Description: dto.Description,
	}

	if err := s.repo.Create(category); err != nil {
		s.logger.Error("failed to create stock material category", zap.Error(err))
		return 0, types.ErrFailedCreateStockMaterialCategory
	}
	return category.ID, nil
}

func (s *stockMaterialCategoryService) GetByID(id uint) (*types.StockMaterialCategoryResponse, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to fetch stock material category", zap.Error(err))
		return nil, types.ErrFailedRetrieveStockMaterialCategory
	}

	response := types.ToStockMaterialCategoryResponse(*category)
	return &response, nil
}

func (s *stockMaterialCategoryService) GetAll(filter types.StockMaterialCategoryFilter) ([]types.StockMaterialCategoryResponse, error) {
	categories, err := s.repo.GetAll(filter)
	if err != nil {
		s.logger.Error("failed to fetch stock material categories", zap.Error(err))
		return nil, types.ErrFailedRetrieveStockMaterialCategories
	}

	responses := []types.StockMaterialCategoryResponse{}
	for _, category := range categories {
		responses = append(responses, types.ToStockMaterialCategoryResponse(category))
	}
	return responses, nil
}

func (s *stockMaterialCategoryService) Update(id uint, dto types.UpdateStockMaterialCategoryDTO) error {
	category := data.StockMaterialCategory{}

	if dto.Name != nil {
		category.Name = *dto.Name
	}
	if dto.Description != nil {
		category.Description = *dto.Description
	}

	if err := s.repo.Update(id, category); err != nil {
		s.logger.Error("failed to update stock material category", zap.Error(err))
		return types.ErrFailedUpdateStockMaterialCategory
	}
	return nil
}

func (s *stockMaterialCategoryService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("failed to delete stock material category", zap.Error(err))
		return types.ErrFailedDeleteStockMaterialCategory
	}
	return nil
}
