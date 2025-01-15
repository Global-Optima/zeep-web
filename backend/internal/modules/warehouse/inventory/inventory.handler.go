package inventory

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/inventory/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	service InventoryService
}

func NewInventoryHandler(service InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) ReceiveInventory(c *gin.Context) {
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

func (h *InventoryHandler) TransferInventory(c *gin.Context) {
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

func (h *InventoryHandler) GetExpiringItems(c *gin.Context) {
	warehouseIDStr := c.Param("warehouseID")
	thresholdDaysStr := c.Query("thresholdDays")

	warehouseID, err := strconv.ParseUint(warehouseIDStr, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid warehouse ID")
		return
	}

	thresholdDays, err := strconv.Atoi(thresholdDaysStr)
	if err != nil || thresholdDays <= 0 {
		utils.SendBadRequestError(c, "Invalid threshold days")
		return
	}

	items, err := h.service.GetExpiringItems(uint(warehouseID), thresholdDays)
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch expiring items")
		return
	}

	utils.SendSuccessResponse(c, items)
}

func (h *InventoryHandler) ExtendExpiration(c *gin.Context) {
	var req types.ExtendExpirationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	if err := h.service.ExtendExpiration(req); err != nil {
		utils.SendInternalServerError(c, "failed to extend expiration date")
		return
	}

	utils.SendMessageWithStatus(c, "Expiration date extended successfully", http.StatusOK)
}

func (h *InventoryHandler) GetDeliveries(c *gin.Context) {
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
