package barcode

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/sku"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type BarcodeService interface {
	GenerateBarcode(req *types.GenerateBarcodeRequest) (*types.GenerateBarcodeResponse, error)
	RetrieveSKUByBarcode(req *types.RetrieveSKUByBarcodeRequest) (*types.RetrieveSKUByBarcodeResponse, error)
	PrintAdditionalBarcodes(req *types.PrintAdditionalBarcodesRequest) (*types.PrintAdditionalBarcodesResponse, error)
}

type barcodeService struct {
	repo    BarcodeRepository
	skuRepo sku.SKURepository

	printerService PrinterService
}

func NewBarcodeService(repo BarcodeRepository, skuRepo sku.SKURepository, printerService PrinterService) BarcodeService {
	return &barcodeService{
		repo:           repo,
		skuRepo:        skuRepo,
		printerService: printerService,
	}
}

func (s *barcodeService) GenerateBarcode(req *types.GenerateBarcodeRequest) (*types.GenerateBarcodeResponse, error) {
	sku, err := s.skuRepo.GetSKUByID(req.SKU_ID)
	if err != nil {
		return nil, err
	}
	if sku == nil {
		return nil, errors.New("SKU not found")
	}

	supplierMaterial, err := s.repo.GetSupplierMaterialByStockMaterialID(req.SKU_ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch supplier for SKU: %w", err)
	}
	if supplierMaterial == nil {
		return nil, errors.New("supplier not found for the given SKU")
	}

	barcode, err := utils.GenerateUPCBarcode(*sku, supplierMaterial.SupplierID)
	if err != nil {
		return nil, err
	}

	err = s.repo.AssignBarcode(barcode, req.SKU_ID)
	if err != nil {
		return nil, err
	}

	err = s.printerService.PrintBarcode(barcode)
	if err != nil {
		return nil, err
	}

	response := types.ToGenerateBarcodeResponse(*sku, barcode)
	return &response, nil
}

func (s *barcodeService) RetrieveSKUByBarcode(req *types.RetrieveSKUByBarcodeRequest) (*types.RetrieveSKUByBarcodeResponse, error) {
	sku, err := s.repo.GetSKUByBarcode(req.Barcode)
	if err != nil {
		return nil, err
	}
	if sku == nil {
		return nil, errors.New("SKU not found with the provided barcode")
	}

	response := types.ToRetrieveSKUByBarcodeResponse(*sku)
	return &response, nil
}

func (s *barcodeService) PrintAdditionalBarcodes(req *types.PrintAdditionalBarcodesRequest) (*types.PrintAdditionalBarcodesResponse, error) {
	sku, err := s.skuRepo.GetSKUByID(req.SKU_ID)
	if err != nil {
		return nil, err
	}
	if sku == nil {
		return nil, errors.New("SKU not found")
	}

	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	var barcodes []string
	for i := 0; i < req.Quantity; i++ {
		err := s.printerService.PrintBarcode(sku.Barcode)
		if err != nil {
			return nil, fmt.Errorf("failed to print barcode: %w", err)
		}
		barcodes = append(barcodes, sku.Barcode)
	}

	response := types.ToPrintAdditionalBarcodesResponse(req.SKU_ID, barcodes)
	return &response, nil
}
