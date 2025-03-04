package stockMaterial

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
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
		localization.SendLocalizedResponseWithKey(c, types.Response400StockMaterial)
		return
	}

	stockMaterialResponses, err := h.service.GetAllStockMaterials(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, stockMaterialResponses, filter.Pagination)
}

func (h *StockMaterialHandler) GetStockMaterialByID(c *gin.Context) {
	idParam := c.Param("id")
	stockMaterialID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StockMaterial)
		return
	}

	stockMaterialResponse, err := h.service.GetStockMaterialByID(uint(stockMaterialID))
	if err != nil {
		if err.Error() == "StockMaterial not found" {
			localization.SendLocalizedResponseWithKey(c, types.Response404StockMaterial)
		} else {
			localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialGet)
		}
		return
	}

	utils.SendSuccessResponse(c, stockMaterialResponse)
}

func (h *StockMaterialHandler) CreateStockMaterial(c *gin.Context) {
	var req types.CreateStockMaterialDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	stockMaterialResponse, err := h.service.CreateStockMaterial(&req)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialCreate)
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

	localization.SendLocalizedResponseWithKey(c, types.Response201StockMaterial)
}

func (h *StockMaterialHandler) UpdateStockMaterial(c *gin.Context) {
	idParam := c.Param("id")
	stockMaterialID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StockMaterial)
		return
	}

	var req types.UpdateStockMaterialDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StockMaterial)
		return
	}

	stockMaterialResponse, err := h.service.GetStockMaterialByID(uint(stockMaterialID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialGet)
		return
	}

	err = h.service.UpdateStockMaterial(uint(stockMaterialID), &req)
	if err != nil {
		if err.Error() == "StockMaterial not found" {
			localization.SendLocalizedResponseWithKey(c, types.Response404StockMaterial)
		} else {
			localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialUpdate)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200StockMaterialUpdate)
}

func (h *StockMaterialHandler) DeleteStockMaterial(c *gin.Context) {
	idParam := c.Param("id")
	stockMaterialID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StockMaterial)
		return
	}

	stockMaterialResponse, err := h.service.GetStockMaterialByID(uint(stockMaterialID))
	if err != nil {
		if errors.Is(err, types.ErrStockMaterialNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404StockMaterial)
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialGet)
		return
	}

	err = h.service.DeleteStockMaterial(uint(stockMaterialID))
	if err != nil {
		if errors.Is(err, types.ErrStockMaterialNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404StockMaterial)
		} else {
			localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialDelete)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200StockMaterialDelete)
}

func (h *StockMaterialHandler) DeactivateStockMaterial(c *gin.Context) {
	idParam := c.Param("id")
	stockMaterialID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StockMaterial)
		return
	}

	err = h.service.DeactivateStockMaterial(uint(stockMaterialID))
	if err != nil {
		if err.Error() == "StockMaterial not found" {
			localization.SendLocalizedResponseWithKey(c, types.Response404StockMaterial)
		} else {
			localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialDeactivate)
		}
		return
	}

	localization.SendLocalizedResponseWithKey(c, types.Response200StockMaterialDeactivate)
}

func (h *StockMaterialHandler) GetStockMaterialBarcode(c *gin.Context) {
	stockMaterialID, errH := utils.ParseParam(c, "id")
	if errH != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StockMaterial)
		return
	}

	barcodeImage, err := h.service.GenerateStockMaterialBarcodePDF(stockMaterialID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialBarcode)
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
		localization.SendLocalizedResponseWithKey(c, types.Response400StockMaterialBarcodeRequired)
		return
	}

	response, err := h.service.RetrieveStockMaterialByBarcode(barcode)
	if err != nil {
		if errors.Is(err, types.ErrStockMaterialBarcodeNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404StockMaterialBarcode)
		} else {
			localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialBarcodeGet)
		}
		return
	}

	utils.SendSuccessResponse(c, response)
}
