package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ToGenerateBarcodeResponse(sku data.StockMaterial, barcode string) GenerateBarcodeResponse {
	return GenerateBarcodeResponse{
		SKU_ID:    sku.ID,
		Barcode:   barcode,
		Message:   "Barcode generated and assigned successfully.",
		CreatedAt: sku.CreatedAt.Format(time.RFC3339),
	}
}

func ToRetrieveSKUByBarcodeResponse(sku data.StockMaterial) RetrieveSKUByBarcodeResponse {
	return RetrieveSKUByBarcodeResponse{
		SKU_ID:    sku.ID,
		Name:      sku.Name,
		Unit:      sku.Unit.Name,
		Category:  sku.Category,
		Barcode:   sku.Barcode,
		CreatedAt: sku.CreatedAt.Format(time.RFC3339),
		UpdatedAt: sku.UpdatedAt.Format(time.RFC3339),
	}
}

func ToPrintAdditionalBarcodesResponse(skuID uint, barcodes []string) PrintAdditionalBarcodesResponse {
	return PrintAdditionalBarcodesResponse{
		SKU_ID:    skuID,
		Barcodes:  barcodes,
		Message:   "Additional barcodes printed successfully.",
		PrintedAt: time.Now().Format(time.RFC3339),
	}
}

func ToBarcodeScanResponse(sku data.StockMaterial, deductedQty, remainingQty float64) BarcodeScanResponse {
	return BarcodeScanResponse{
		SKU_ID:    sku.ID,
		Name:      sku.Name,
		Quantity:  deductedQty,
		Remaining: remainingQty,
		Message:   "Stock deducted successfully.",
		ScannedAt: time.Now().Format(time.RFC3339),
	}
}
