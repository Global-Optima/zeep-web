package warehouseStock

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type WarehouseStockHandler struct {
	service WarehouseStockService
}

func NewWarehouseStockHandler(service WarehouseStockService) *WarehouseStockHandler {
	return &WarehouseStockHandler{service: service}
}

func (h *WarehouseStockHandler) ReceiveInventory(c *gin.Context) {
	var req types.ReceiveInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	if len(req.NewItems) == 0 && len(req.ExistingItems) == 0 {
		utils.SendBadRequestError(c, "No items provided in the request")
		return
	}

	if err := h.service.ReceiveInventory(req); err != nil {
		utils.SendInternalServerError(c, "failed to receive inventory")
		return
	}

	utils.SendMessageWithStatus(c, "inventory received successfully", http.StatusOK)
}

func (h *WarehouseStockHandler) TransferInventory(c *gin.Context) {
	var req types.TransferInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	if len(req.Items) == 0 {
		utils.SendBadRequestError(c, "No items provided in the request")
		return
	}

	if err := h.service.TransferInventory(req); err != nil {
		utils.SendInternalServerError(c, "failed to transfer inventory")
		return
	}

	utils.SendMessageWithStatus(c, "inventory transferred successfully", http.StatusOK)
}

func (h *WarehouseStockHandler) GetDeliveries(c *gin.Context) {
	var filter types.DeliveryFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters: "+err.Error())
		return
	}

	deliveries, err := h.service.GetDeliveries(filter.WarehouseID, filter.StartDate, filter.EndDate)
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch deliveries")
		return
	}

	utils.SendSuccessResponse(c, deliveries)
}

func (h *WarehouseStockHandler) AddToStock(c *gin.Context) {
	var req types.AdjustWarehouseStock
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.service.AddWarehouseStockMaterial(req); err != nil {
		utils.SendInternalServerError(c, "Failed to add stock: "+err.Error())
		return
	}

	utils.SendMessageWithStatus(c, "Stock added successfully", http.StatusOK)
}

func (h *WarehouseStockHandler) DeductFromStock(c *gin.Context) {
	var req types.AdjustWarehouseStock
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.service.DeductFromStock(req); err != nil {
		utils.SendInternalServerError(c, "Failed to deduct stock: "+err.Error())
		return
	}

	utils.SendMessageWithStatus(c, "Stock deducted successfully", http.StatusOK)
}

func (h *WarehouseStockHandler) GetStocks(c *gin.Context) {
	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	var filter types.GetWarehouseStockFilterQuery
	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.WarehouseStock{})
	if err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}
	filter.WarehouseID = &warehouseID

	stocks, err := h.service.GetStock(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch stock: "+err.Error())
		return
	}

	utils.SendSuccessResponseWithPagination(c, stocks, filter.Pagination)
}

func (h *WarehouseStockHandler) GetStockMaterialDetails(c *gin.Context) {
	stockMaterialID, err := strconv.ParseUint(c.Param("stockMaterialId"), 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "invalid stock material ID")
		return
	}

	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	details, err := h.service.GetStockMaterialDetails(uint(stockMaterialID), warehouseID)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, details)
}

func (h *WarehouseStockHandler) UpdateStock(c *gin.Context) {
	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
	}

	stockMaterialID, err := utils.ParseParam(c, "stockMaterialId")
	if err != nil {
		utils.SendBadRequestError(c, err.Error())
	}

	var dto types.UpdateWarehouseStockDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.service.UpdateStock(warehouseID, stockMaterialID, dto); err != nil {
		utils.SendInternalServerError(c, "Failed to update stock: "+err.Error())
		return
	}

	utils.SendMessageWithStatus(c, "Stock updated successfully", http.StatusOK)
}

func (h *WarehouseStockHandler) AddWarehouseStocks(c *gin.Context) {
	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var req []types.AddWarehouseStockMaterial
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.service.AddWarehouseStocks(warehouseID, req); err != nil {
		utils.SendInternalServerError(c, "Failed to add warehouse stocks: "+err.Error())
		return
	}

	utils.SendMessageWithStatus(c, "Warehouse stocks added successfully", http.StatusCreated)
}
