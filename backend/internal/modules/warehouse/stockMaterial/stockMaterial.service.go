package stockMaterial

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
)

type StockMaterialService interface {
	GetAllStockMaterials(filter *types.StockMaterialFilter) ([]types.StockMaterialsDTO, error)
	GetStockMaterialByID(stockMaterialID uint) (*types.StockMaterialsDTO, error)
	CreateStockMaterial(req *types.CreateStockMaterialDTO) (*types.StockMaterialsDTO, error)
	UpdateStockMaterial(stockMaterialID uint, req *types.UpdateStockMaterialDTO) (*types.StockMaterialsDTO, error)
	DeleteStockMaterial(stockMaterialID uint) error
	DeactivateStockMaterial(stockMaterialID uint) error
}

type stockMaterialService struct {
	repo StockMaterialRepository
}

func NewStockMaterialService(repo StockMaterialRepository) StockMaterialService {
	return &stockMaterialService{
		repo: repo,
	}
}

func (s *stockMaterialService) GetAllStockMaterials(filter *types.StockMaterialFilter) ([]types.StockMaterialsDTO, error) {
	stockMaterials, err := s.repo.GetAllStockMaterials(filter)
	if err != nil {
		return nil, err
	}

	var stockMaterialResponses []types.StockMaterialsDTO
	for _, stockMaterial := range stockMaterials {
		stockMaterialResponses = append(stockMaterialResponses, *types.ConvertStockMaterialToStockMaterialResponse(&stockMaterial))
	}

	return stockMaterialResponses, nil
}

func (s *stockMaterialService) GetStockMaterialByID(stockMaterialID uint) (*types.StockMaterialsDTO, error) {
	stockMaterial, err := s.repo.GetStockMaterialByID(stockMaterialID)
	if err != nil {
		return nil, err
	}

	if stockMaterial == nil {
		return nil, errors.New("StockMaterial not found")
	}

	stockMaterialResponse := types.ConvertStockMaterialToStockMaterialResponse(stockMaterial)

	return stockMaterialResponse, nil
}

func (s *stockMaterialService) CreateStockMaterial(req *types.CreateStockMaterialDTO) (*types.StockMaterialsDTO, error) {
	stockMaterial := types.ConvertCreateStockMaterialRequestToStockMaterial(req)

	err := s.repo.CreateStockMaterial(stockMaterial)
	if err != nil {
		return nil, err
	}

	supplierMaterial := &data.SupplierMaterial{
		StockMaterialID: stockMaterial.ID,
		SupplierID:      req.SupplierID,
	}
	err = s.repo.CreateSupplierMaterial(supplierMaterial)
	if err != nil {
		return nil, fmt.Errorf("failed to create supplier-material association: %w", err)
	}

	stockMaterialResponse := types.ConvertStockMaterialToStockMaterialResponse(stockMaterial)
	return stockMaterialResponse, nil
}

func (s *stockMaterialService) UpdateStockMaterial(stockMaterialID uint, req *types.UpdateStockMaterialDTO) (*types.StockMaterialsDTO, error) {
	updated, err := s.repo.UpdateStockMaterialFields(stockMaterialID, *req)
	if err != nil {
		return nil, err
	}
	updatedStockMaterial := types.ConvertStockMaterialToStockMaterialResponse(updated)
	return updatedStockMaterial, nil
}

func (s *stockMaterialService) DeleteStockMaterial(stockMaterialID uint) error {
	return s.repo.DeleteStockMaterial(stockMaterialID)
}

func (s *stockMaterialService) DeactivateStockMaterial(stockMaterialID uint) error {
	return s.repo.DeactivateStockMaterial(stockMaterialID)
}
