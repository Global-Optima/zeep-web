package inventory

import (
	"net/http"
	"strconv"
	"time"

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
		utils.SendBadRequestError(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.service.ReceiveInventory(req); err != nil {
		utils.SendInternalServerError(c, "Failed to receive inventory: "+err.Error())
		return
	}

	utils.SendMessageWithStatus(c, "Inventory received successfully", http.StatusOK)
}

func (h *InventoryHandler) TransferInventory(c *gin.Context) {
	var req types.TransferInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.service.TransferInventory(req); err != nil {
		utils.SendInternalServerError(c, "Failed to transfer inventory: "+err.Error())
		return
	}

	utils.SendMessageWithStatus(c, "Inventory transferred successfully", http.StatusOK)
}

func (h *InventoryHandler) GetInventoryLevels(c *gin.Context) {
	warehouseIDStr := c.Param("warehouseID")
	warehouseID, err := strconv.ParseUint(warehouseIDStr, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid warehouse ID")
		return
	}

	levels, err := h.service.GetInventoryLevels(uint(warehouseID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch inventory levels: "+err.Error())
		return
	}

	utils.SendSuccessResponse(c, levels)
}

func (h *InventoryHandler) PickupStock(c *gin.Context) {
	var req types.PickupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.service.PickupStock(req); err != nil {
		utils.SendInternalServerError(c, "Failed to handle store pickup: "+err.Error())
		return
	}

	utils.SendMessageWithStatus(c, "Store pickup completed successfully", http.StatusOK)
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
		utils.SendInternalServerError(c, "Failed to fetch expiring items: "+err.Error())
		return
	}

	utils.SendSuccessResponse(c, items)
}

func (h *InventoryHandler) ExtendExpiration(c *gin.Context) {
	var req types.ExtendExpirationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.service.ExtendExpiration(req); err != nil {
		utils.SendInternalServerError(c, "Failed to extend expiration date: "+err.Error())
		return
	}

	utils.SendMessageWithStatus(c, "Expiration date extended successfully", http.StatusOK)
}

func (h *InventoryHandler) GetDeliveries(c *gin.Context) {
	warehouseIDStr := c.Query("warehouseID")
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	var warehouseID *uint
	if warehouseIDStr != "" {
		parsedID, err := strconv.ParseUint(warehouseIDStr, 10, 32)
		if err != nil {
			utils.SendBadRequestError(c, "Invalid warehouse ID")
			return
		}
		id := uint(parsedID)
		warehouseID = &id
	}

	var startDate, endDate *time.Time
	if startDateStr != "" {
		parsedStartDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			utils.SendBadRequestError(c, "Invalid start date format, use RFC3339")
			return
		}
		startDate = &parsedStartDate
	}
	if endDateStr != "" {
		parsedEndDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			utils.SendBadRequestError(c, "Invalid end date format, use RFC3339")
			return
		}
		endDate = &parsedEndDate
	}

	deliveries, err := h.service.GetDeliveries(warehouseID, startDate, endDate)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch deliveries: "+err.Error())
		return
	}

	utils.SendSuccessResponse(c, deliveries)
}
