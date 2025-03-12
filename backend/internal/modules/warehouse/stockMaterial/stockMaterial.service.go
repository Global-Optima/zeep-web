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
	GenerateStockMaterialBarcodePDF(stockMaterialID uint) ([]byte, error)
	GenerateBarcode() (*types.GenerateBarcodeResponse, error)
	RetrieveStockMaterialByBarcode(barcode string) (*types.StockMaterialsDTO, error)
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

	stockMaterialResponse := types.ConvertStockMaterialToStockMaterialResponse(stockMaterial)

	return stockMaterialResponse, nil
}

func (s *stockMaterialService) CreateStockMaterial(req *types.CreateStockMaterialDTO) (*types.StockMaterialsDTO, error) {
	if req.SafetyStock <= 0 {
		return nil, fmt.Errorf("safety stock must be greater than zero")
	}

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

	return barcodeImage.Bytes(), nil
}

func (s *stockMaterialService) GenerateStockMaterialBarcodePDF(stockMaterialID uint) ([]byte, error) {
	stockMaterial, err := s.repo.GetStockMaterialByID(stockMaterialID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock material: %w", err)
	}

	// 2. Encode the barcode data using Code-128
	barcodeData := stockMaterial.Barcode

	return utils.GenerateBarcodePDF(barcodeData)
}

func (s *stockMaterialService) GenerateBarcode() (*types.GenerateBarcodeResponse, error) {
	maxRetries := 10
	var barcode string
	var exists bool
	var err error

	for i := 0; i < maxRetries; i++ {
		barcode = utils.GenerateRandomEAN13()

		exists, err = s.repo.IsBarcodeExists(barcode)
		if err != nil {
			return nil, fmt.Errorf("failed to check uniqueness of barcode: %w", err)
		}

		if !exists {
			break
		}
	}

	if exists {
		return nil, fmt.Errorf("failed to generate unique barcode after %d attempts", maxRetries)
	}

	response := types.ToGenerateBarcodeResponse(barcode)
	return &response, nil
}

func (s *stockMaterialService) RetrieveStockMaterialByBarcode(barcode string) (*types.StockMaterialsDTO, error) {
	stockMaterial, err := s.repo.GetStockMaterialByBarcode(barcode)
	if err != nil {
		return nil, err
	}
	if stockMaterial == nil {
		return nil, errors.New("StockMaterial not found with the provided barcode")
	}

	return types.ConvertStockMaterialToStockMaterialResponse(stockMaterial), nil
}
