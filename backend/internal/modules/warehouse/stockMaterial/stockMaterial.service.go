package stockMaterial

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type StockMaterialService interface {
	GetAllStockMaterials(filter *types.StockMaterialFilter) ([]types.StockMaterialsDTO, error)
	GetStockMaterialByID(stockMaterialID uint) (*types.StockMaterialsDTO, error)
	CreateStockMaterial(req *types.CreateStockMaterialDTO) (*types.StockMaterialsDTO, error)
	UpdateStockMaterial(stockMaterialID uint, req *types.UpdateStockMaterialDTO) error
	DeleteStockMaterial(stockMaterialID uint) error
	DeactivateStockMaterial(stockMaterialID uint) error
	GetStockMaterialBarcode(stockMaterialID uint) ([]byte, error)
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

	stockMaterialResponses := make([]types.StockMaterialsDTO, 0)
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

	stockMaterialResponse := types.ConvertStockMaterialToStockMaterialResponse(stockMaterial)
	return stockMaterialResponse, nil
}

func (s *stockMaterialService) UpdateStockMaterial(stockMaterialID uint, req *types.UpdateStockMaterialDTO) error {
	stockMaterial, err := s.repo.GetStockMaterialByID(stockMaterialID)
	if err != nil {
		return fmt.Errorf("failed to fetch stock material: %w", err)
	}

	if stockMaterial == nil {
		return fmt.Errorf("stock material with ID %d not found", stockMaterialID)
	}

	updatedStockMaterial, err := types.ValidateAndApplyUpdate(stockMaterial, req)
	if err != nil {
		return err
	}

	err = s.repo.UpdateStockMaterial(stockMaterialID, updatedStockMaterial)
	if err != nil {
		return fmt.Errorf("failed to update stock material: %w", err)
	}

	return nil
}

func (s *stockMaterialService) DeleteStockMaterial(stockMaterialID uint) error {
	return s.repo.DeleteStockMaterial(stockMaterialID)
}

func (s *stockMaterialService) DeactivateStockMaterial(stockMaterialID uint) error {
	return s.repo.DeactivateStockMaterial(stockMaterialID)
}

func (s *stockMaterialService) GetStockMaterialBarcode(stockMaterialID uint) ([]byte, error) {
	stockMaterial, err := s.repo.GetStockMaterialByID(stockMaterialID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock material: %w", err)
	}

	barcodeImage, err := utils.GenerateBarcodeImage(stockMaterial.Barcode)
	if err != nil {
		return nil, fmt.Errorf("failed to generate barcode image: %w", err)
	}

	return barcodeImage, nil
}
