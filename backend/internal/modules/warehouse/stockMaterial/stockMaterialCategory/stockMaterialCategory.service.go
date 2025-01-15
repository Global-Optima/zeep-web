package stockMaterialCategory

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialCategory/types"
	"gorm.io/gorm"
)

type StockMaterialCategoryService interface {
	Create(dto types.CreateStockMaterialCategoryDTO) (uint, error)
	GetByID(id uint) (*types.StockMaterialCategoryResponse, error)
	GetAll() ([]types.StockMaterialCategoryResponse, error)
	Update(id uint, dto types.UpdateStockMaterialCategoryDTO) error
	Delete(id uint) error
}

type stockMaterialCategoryService struct {
	repo StockMaterialCategoryRepository
}

func NewStockMaterialCategoryService(repo StockMaterialCategoryRepository) StockMaterialCategoryService {
	return &stockMaterialCategoryService{repo: repo}
}

func (s *stockMaterialCategoryService) Create(dto types.CreateStockMaterialCategoryDTO) (uint, error) {
	category := &data.StockMaterialCategory{
		Name:        dto.Name,
		Description: dto.Description,
	}

	if err := s.repo.Create(category); err != nil {
		return 0, fmt.Errorf("failed to create stock material category: %w", err)
	}
	return category.ID, nil
}

func (s *stockMaterialCategoryService) GetByID(id uint) (*types.StockMaterialCategoryResponse, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("stock material category not found")
		}
		return nil, fmt.Errorf("failed to fetch stock material category: %w", err)
	}

	response := types.ToStockMaterialCategoryResponse(*category)
	return &response, nil
}

func (s *stockMaterialCategoryService) GetAll() ([]types.StockMaterialCategoryResponse, error) {
	categories, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock material categories: %w", err)
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
		return fmt.Errorf("failed to update stock material category: %w", err)
	}
	return nil
}

func (s *stockMaterialCategoryService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete stock material category: %w", err)
	}
	return nil
}
