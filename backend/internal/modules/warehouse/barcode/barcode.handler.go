package barcode

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type BarcodeHandler struct {
	service BarcodeService
}

func NewBarcodeHandler(service BarcodeService) *BarcodeHandler {
	return &BarcodeHandler{service: service}
}

func (h *BarcodeHandler) GenerateBarcode(c *gin.Context) {
	var req types.GenerateBarcodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	response, err := h.service.GenerateBarcode(&req)
	if err != nil {
		utils.SendInternalServerError(c, "failed to generate barcode")
		return
	}

	utils.SendSuccessResponse(c, response)
}

func (h *BarcodeHandler) RetrieveStockMaterialByBarcode(c *gin.Context) {
	barcode := c.Param("barcode")
	if barcode == "" {
		utils.SendBadRequestError(c, "Barcode is required")
		return
	}

	req := types.RetrieveStockMaterialByBarcodeRequest{
		Barcode: barcode,
	}

	response, err := h.service.RetrieveStockMaterialByBarcode(&req)
	if err != nil {
		if err.Error() == "StockMaterial not found with the provided barcode" {
			utils.SendNotFoundError(c, "StockMaterial not found with the provided barcode")
		} else {
			utils.SendInternalServerError(c, "failed to retrieve stockMaterial barcode")
		}
		return
	}

	utils.SendSuccessResponse(c, response)
}

func (h *BarcodeHandler) PrintAdditionalBarcodes(c *gin.Context) {
	var req types.PrintAdditionalBarcodesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	response, err := h.service.PrintAdditionalBarcodes(&req)
	if err != nil {
		utils.SendInternalServerError(c, "failed to print additional barcodes")
		return
	}

	utils.SendSuccessResponse(c, response)
}
