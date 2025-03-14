package stockMaterial

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StockMaterialHandler struct {
	service      StockMaterialService
	auditService audit.AuditService
}

func NewStockMaterialHandler(service StockMaterialService, auditService audit.AuditService) *StockMaterialHandler {
	return &StockMaterialHandler{
		service:      service,
		auditService: auditService,
	}
}

func (h *StockMaterialHandler) GetAllStockMaterials(c *gin.Context) {
	var filter types.StockMaterialFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StockMaterial{}); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	stockMaterialResponses, err := h.service.GetAllStockMaterials(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve stock materials")
		return
	}

	utils.SendSuccessResponseWithPagination(c, stockMaterialResponses, filter.Pagination)
}

func (h *StockMaterialHandler) GetStockMaterialByID(c *gin.Context) {
	idParam := c.Param("id")
	stockMaterialID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid StockMaterial ID")
		return
	}

	stockMaterialResponse, err := h.service.GetStockMaterialByID(uint(stockMaterialID))
	if err != nil {
		if errors.Is(err, types.ErrStockMaterialNotFound) {
			utils.SendNotFoundError(c, "StockMaterial not found")
		} else {
			utils.SendInternalServerError(c, "failed to retrieve stockMaterial")
		}
		return
	}

	utils.SendSuccessResponse(c, stockMaterialResponse)
}

func (h *StockMaterialHandler) CreateStockMaterial(c *gin.Context) {
	var req types.CreateStockMaterialDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	stockMaterialResponse, err := h.service.CreateStockMaterial(&req)
	if err != nil {
		utils.SendInternalServerError(c, "failed to create stockMaterial")
		return
	}

	action := types.CreateStockMaterialAuditFactory(
		&data.BaseDetails{
			ID:   stockMaterialResponse.ID,
			Name: stockMaterialResponse.Name,
		})

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, stockMaterialResponse)
}

func (h *StockMaterialHandler) UpdateStockMaterial(c *gin.Context) {
	idParam := c.Param("id")
	stockMaterialID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid StockMaterial ID")
		return
	}

	var req types.UpdateStockMaterialDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	stockMaterialResponse, err := h.service.GetStockMaterialByID(uint(stockMaterialID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve stockMaterial")
		return
	}

	err = h.service.UpdateStockMaterial(uint(stockMaterialID), &req)
	if err != nil {
		if errors.Is(err, types.ErrStockMaterialNotFound) {
			utils.SendNotFoundError(c, "StockMaterial not found")
		} else {
			utils.SendInternalServerError(c, "failed to update stockMaterial")
		}
		return
	}

	action := types.UpdateStockMaterialAuditFactory(
		&data.BaseDetails{
			ID:   stockMaterialResponse.ID,
			Name: stockMaterialResponse.Name,
		}, &req)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, "Updated stock material successfully")
}

func (h *StockMaterialHandler) DeleteStockMaterial(c *gin.Context) {
	idParam := c.Param("id")
	stockMaterialID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid StockMaterial ID")
		return
	}

	stockMaterialResponse, err := h.service.GetStockMaterialByID(uint(stockMaterialID))
	if err != nil {
		if errors.Is(err, types.ErrStockMaterialNotFound) {
			utils.SendNotFoundError(c, "stock material not found")
		}
		utils.SendInternalServerError(c, "failed to retrieve stockMaterial")
		return
	}

	err = h.service.DeleteStockMaterial(uint(stockMaterialID))
	if err != nil {
		if errors.Is(err, types.ErrStockMaterialNotFound) {
			utils.SendNotFoundError(c, "StockMaterial not found")
		} else {
			utils.SendInternalServerError(c, "failed to delete stockMaterial")
		}
		return
	}

	action := types.DeleteStockMaterialAuditFactory(
		&data.BaseDetails{
			ID:   stockMaterialResponse.ID,
			Name: stockMaterialResponse.Name,
		})

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	c.Status(http.StatusNoContent)
}

func (h *StockMaterialHandler) DeactivateStockMaterial(c *gin.Context) {
	idParam := c.Param("id")
	stockMaterialID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid StockMaterial ID")
		return
	}

	err = h.service.DeactivateStockMaterial(uint(stockMaterialID))
	if err != nil {
		if errors.Is(err, types.ErrStockMaterialNotFound) {
			utils.SendNotFoundError(c, "StockMaterial not found")
		} else {
			utils.SendInternalServerError(c, "failed to deactivate stockMaterial")
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *StockMaterialHandler) GetStockMaterialBarcode(c *gin.Context) {
	stockMaterialID, errH := utils.ParseParam(c, "id")
	if errH != nil {
		utils.SendBadRequestError(c, errH.Error())
		return
	}

	barcodeImage, err := h.service.GenerateStockMaterialBarcodePDF(stockMaterialID)
	if err != nil {
		utils.SendInternalServerError(c, "failed to get barcode image")
		return
	}

	filename := fmt.Sprintf("barcode-%d.pdf", stockMaterialID)

	// Add headers for downloading the barcode image
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", len(barcodeImage)))
	c.Data(http.StatusOK, "application/pdf", barcodeImage)
}

func (h *StockMaterialHandler) GenerateBarcode(c *gin.Context) {
	response, err := h.service.GenerateBarcode()
	if err != nil {
		utils.SendInternalServerError(c, "failed to generate barcode")
		return
	}

	utils.SendSuccessResponse(c, response)
}

func (h *StockMaterialHandler) RetrieveStockMaterialByBarcode(c *gin.Context) {
	barcode := c.Param("barcode")
	if barcode == "" {
		utils.SendBadRequestError(c, "Barcode is required")
		return
	}

	response, err := h.service.RetrieveStockMaterialByBarcode(barcode)
	if err != nil {
		if errors.Is(err, types.ErrStockMaterialNotFound) {
			utils.SendNotFoundError(c, "StockMaterial not found with the provided barcode")
		} else {
			utils.SendInternalServerError(c, "failed to retrieve stockMaterial barcode")
		}
		return
	}

	utils.SendSuccessResponse(c, response)
}
