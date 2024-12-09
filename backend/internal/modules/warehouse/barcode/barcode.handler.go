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
		utils.SendBadRequestError(c, err.Error())
		return
	}

	response, err := h.service.GenerateBarcode(&req)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

func (h *BarcodeHandler) RetrieveSKUByBarcode(c *gin.Context) {
	barcode := c.Param("barcode")
	if barcode == "" {
		utils.SendBadRequestError(c, "Barcode is required")
		return
	}

	req := types.RetrieveSKUByBarcodeRequest{
		Barcode: barcode,
	}

	response, err := h.service.RetrieveSKUByBarcode(&req)
	if err != nil {
		if err.Error() == "SKU not found with the provided barcode" {
			utils.SendNotFoundError(c, err.Error())
		} else {
			utils.SendInternalServerError(c, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, response)
}

func (h *BarcodeHandler) PrintAdditionalBarcodes(c *gin.Context) {
	var req types.PrintAdditionalBarcodesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	response, err := h.service.PrintAdditionalBarcodes(&req)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}
