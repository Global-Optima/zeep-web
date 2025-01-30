package types

type GenerateBarcodeResponse struct {
	Barcode string `json:"barcode"`
}

type RetrieveStockMaterialByBarcodeRequest struct {
	Barcode string `json:"barcode" binding:"required"`
}

type RetrieveStockMaterialByBarcodeResponse struct {
	StockMaterialID uint   `json:"stockMaterialId"`
	Name            string `json:"name"`
	Unit            string `json:"unit"`
	Category        string `json:"category"`
	Barcode         string `json:"barcode"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

type PrintAdditionalBarcodesRequest struct {
	StockMaterialID uint `json:"stockMaterialId" binding:"required"`
	Quantity        int  `json:"quantity" binding:"required,gt=0"`
}

type PrintAdditionalBarcodesResponse struct {
	StockMaterialID uint     `json:"stockMaterialId"`
	Barcodes        []string `json:"barcodes"`
	Message         string   `json:"message"`
	PrintedAt       string   `json:"printedAt"`
}

type BarcodeScanRequest struct {
	Barcode  string  `json:"barcode" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
}

type BarcodeScanResponse struct {
	StockMaterialID uint    `json:"stockMaterialId"`
	Name            string  `json:"name"`
	Quantity        float64 `json:"deductedQuantity"`
	Remaining       float64 `json:"remainingQuantity"`
	Message         string  `json:"message"`
	ScannedAt       string  `json:"scannedAt"`
}

type StockMaterialBarcodeResponse struct {
	StockMaterialID uint   `json:"stockMaterialId"`
	Barcode         string `json:"barcode"`
}

type GetBarcodesRequest struct {
	IDs []uint `json:"stockMaterialIDs"`
}
