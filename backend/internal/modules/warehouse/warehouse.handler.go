package warehouse

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	service WarehouseService
}

func NewWarehouseHandler(service WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

func (h *WarehouseHandler) AssignStoreToWarehouse(c *gin.Context) {
	var req types.AssignStoreToWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.service.AssignStoreToWarehouse(req); err != nil {
		utils.SendInternalServerError(c, "Failed to assign store to warehouse: "+err.Error())
		return
	}

	utils.SendMessageWithStatus(c, "Store assigned to warehouse successfully", http.StatusOK)
}

func (h *WarehouseHandler) ReassignStore(c *gin.Context) {
	storeIDStr := c.Param("storeId")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store ID")
		return
	}

	var req types.ReassignStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.service.ReassignStore(uint(storeID), req); err != nil {
		utils.SendInternalServerError(c, "Failed to reassign store: "+err.Error())
		return
	}

	utils.SendMessageWithStatus(c, "Store reassigned successfully", http.StatusOK)
}

func (h *WarehouseHandler) GetAllStoresByWarehouse(c *gin.Context) {
	warehouseIDStr := c.Param("warehouseId")
	warehouseID, err := strconv.ParseUint(warehouseIDStr, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid warehouse ID")
		return
	}

	stores, err := h.service.GetAllStoresByWarehouse(uint(warehouseID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to list stores: "+err.Error())
		return
	}

	utils.SendSuccessResponse(c, stores)
}
