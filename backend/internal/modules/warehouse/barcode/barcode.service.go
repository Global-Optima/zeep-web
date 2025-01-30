package barcode

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type BarcodeService interface {
	GenerateBarcode() (*types.GenerateBarcodeResponse, error)
	RetrieveStockMaterialByBarcode(barcode string) (*stockMaterialTypes.StockMaterialsDTO, error)
	PrintAdditionalBarcodes(req *types.PrintAdditionalBarcodesRequest) (*types.PrintAdditionalBarcodesResponse, error)

	GetBarcodesForStockMaterials(stockMaterialIDs []uint) ([]types.StockMaterialBarcodeResponse, error)
	GetBarcodeForStockMaterial(stockMaterialID uint) (*types.StockMaterialBarcodeResponse, error)
}

type barcodeService struct {
	repo              BarcodeRepository
	stockMaterialRepo stockMaterial.StockMaterialRepository

	printerService PrinterService
}

func NewBarcodeService(repo BarcodeRepository, stockMaterialRepo stockMaterial.StockMaterialRepository, printerService PrinterService) BarcodeService {
	return &barcodeService{
		repo:              repo,
		stockMaterialRepo: stockMaterialRepo,
		printerService:    printerService,
	}
}

func (s *barcodeService) GenerateBarcode() (*types.GenerateBarcodeResponse, error) {
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

func (s *barcodeService) RetrieveStockMaterialByBarcode(barcode string) (*stockMaterialTypes.StockMaterialsDTO, error) {
	stockMaterial, err := s.repo.GetStockMaterialByBarcode(barcode)
	if err != nil {
		return nil, err
	}
	if stockMaterial == nil {
		return nil, errors.New("StockMaterial not found with the provided barcode")
	}

	return stockMaterialTypes.ConvertStockMaterialToStockMaterialResponse(stockMaterial), nil
}

func (s *barcodeService) PrintAdditionalBarcodes(req *types.PrintAdditionalBarcodesRequest) (*types.PrintAdditionalBarcodesResponse, error) {
	stockMaterial, err := s.stockMaterialRepo.GetStockMaterialByID(req.StockMaterialID)
	if err != nil {
		return nil, err
	}
	if stockMaterial == nil {
		return nil, errors.New("StockMaterial not found")
	}

	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	var barcodes []string
	for i := 0; i < req.Quantity; i++ {
		err := s.printerService.PrintBarcode(stockMaterial.Barcode)
		if err != nil {
			return nil, fmt.Errorf("failed to print barcode: %w", err)
		}
		barcodes = append(barcodes, stockMaterial.Barcode)
	}

	response := types.ToPrintAdditionalBarcodesResponse(req.StockMaterialID, barcodes)
	return &response, nil
}

func (s *barcodeService) GetBarcodesForStockMaterials(stockMaterialIDs []uint) ([]types.StockMaterialBarcodeResponse, error) {
	stockMaterials, err := s.repo.GetBarcodesForStockMaterials(stockMaterialIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch barcodes for stock materials: %w", err)
	}

	return types.ToStockMaterialBarcodeResponses(stockMaterials), nil
}

func (s *barcodeService) GetBarcodeForStockMaterial(stockMaterialID uint) (*types.StockMaterialBarcodeResponse, error) {
	stockMaterial, err := s.repo.GetBarcodeForStockMaterial(stockMaterialID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch barcode for stock material: %w", err)
	}
	if stockMaterial == nil {
		return nil, errors.New("stock material not found")
	}

	return types.ToStockMaterialBarcodeResponse(stockMaterial), nil
}
