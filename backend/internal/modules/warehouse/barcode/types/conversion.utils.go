package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ToGenerateBarcodeResponse(stockMaterial data.StockMaterial, barcode string) GenerateBarcodeResponse {
	return GenerateBarcodeResponse{
		StockMaterialID: stockMaterial.ID,
		Barcode:         barcode,
		Message:         "Barcode generated and assigned successfully.",
		CreatedAt:       stockMaterial.CreatedAt.Format(time.RFC3339),
	}
}

func ToRetrieveStockMaterialByBarcodeResponse(stockMaterial data.StockMaterial) RetrieveStockMaterialByBarcodeResponse {
	return RetrieveStockMaterialByBarcodeResponse{
		StockMaterialID: stockMaterial.ID,
		Name:            stockMaterial.Name,
		Unit:            stockMaterial.Unit.Name,
		Category:        stockMaterial.Category,
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
