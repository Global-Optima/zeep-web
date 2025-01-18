package barcode

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type BarcodeService interface {
	GenerateBarcode(req *types.GenerateBarcodeRequest) (*types.GenerateBarcodeResponse, error)
	RetrieveStockMaterialByBarcode(req *types.RetrieveStockMaterialByBarcodeRequest) (*types.RetrieveStockMaterialByBarcodeResponse, error)
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

func (s *barcodeService) GenerateBarcode(req *types.GenerateBarcodeRequest) (*types.GenerateBarcodeResponse, error) {
	stockMaterial, err := s.stockMaterialRepo.GetStockMaterialByID(req.StockMaterialID)
	if err != nil {
		return nil, err
	}
	if stockMaterial == nil {
		return nil, errors.New("StockMaterial not found")
	}

	supplierMaterial, err := s.repo.GetSupplierMaterialByStockMaterialID(req.StockMaterialID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch supplier for StockMaterial: %w", err)
	}
	if supplierMaterial == nil {
		return nil, errors.New("supplier not found for the given StockMaterial")
	}

	barcode, err := utils.GenerateUPCBarcode(*stockMaterial, supplierMaterial.SupplierID)
	if err != nil {
		return nil, err
	}

	err = s.repo.AssignBarcode(barcode, req.StockMaterialID)
	if err != nil {
		return nil, err
	}

	err = s.printerService.PrintBarcode(barcode)
	if err != nil {
		return nil, err
	}

	response := types.ToGenerateBarcodeResponse(*stockMaterial, barcode)
	return &response, nil
}

func (s *barcodeService) RetrieveStockMaterialByBarcode(req *types.RetrieveStockMaterialByBarcodeRequest) (*types.RetrieveStockMaterialByBarcodeResponse, error) {
	stockMaterial, err := s.repo.GetStockMaterialByBarcode(req.Barcode)
	if err != nil {
		return nil, err
	}
	if stockMaterial == nil {
		return nil, errors.New("StockMaterial not found with the provided barcode")
	}

	response := types.ToRetrieveStockMaterialByBarcodeResponse(*stockMaterial)
	return &response, nil
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
