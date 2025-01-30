package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ToGenerateBarcodeResponse(barcode string) GenerateBarcodeResponse {
	return GenerateBarcodeResponse{
		Barcode: barcode,
	}
}

func ToRetrieveStockMaterialByBarcodeResponse(stockMaterial data.StockMaterial) RetrieveStockMaterialByBarcodeResponse {
	return RetrieveStockMaterialByBarcodeResponse{
		StockMaterialID: stockMaterial.ID,
		Name:            stockMaterial.Name,
		Unit:            stockMaterial.Unit.Name,
		Category:        stockMaterial.StockMaterialCategory.Name,
		Barcode:         stockMaterial.Barcode,
		CreatedAt:       stockMaterial.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       stockMaterial.UpdatedAt.Format(time.RFC3339),
	}
}

func ToPrintAdditionalBarcodesResponse(stockMaterialID uint, barcodes []string) PrintAdditionalBarcodesResponse {
	return PrintAdditionalBarcodesResponse{
		StockMaterialID: stockMaterialID,
		Barcodes:        barcodes,
		Message:         "Additional barcodes printed successfully.",
		PrintedAt:       time.Now().Format(time.RFC3339),
	}
}

func ToBarcodeScanResponse(stockMaterial data.StockMaterial, deductedQty, remainingQty float64) BarcodeScanResponse {
	return BarcodeScanResponse{
		StockMaterialID: stockMaterial.ID,
		Name:            stockMaterial.Name,
		Quantity:        deductedQty,
		Remaining:       remainingQty,
		Message:         "Stock deducted successfully.",
		ScannedAt:       time.Now().Format(time.RFC3339),
	}
}

func ToStockMaterialBarcodeResponse(stockMaterial *data.StockMaterial) *StockMaterialBarcodeResponse {
	return &StockMaterialBarcodeResponse{
		StockMaterialID: stockMaterial.ID,
		Barcode:         stockMaterial.Barcode,
	}
}

func ToStockMaterialBarcodeResponses(stockMaterials []data.StockMaterial) []StockMaterialBarcodeResponse {
	responses := make([]StockMaterialBarcodeResponse, len(stockMaterials))
	for i, material := range stockMaterials {
		responses[i] = *ToStockMaterialBarcodeResponse(&material)
	}
	return responses
}
