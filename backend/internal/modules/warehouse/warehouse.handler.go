package warehouse

import (
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
	"github.com/pkg/errors"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	service       WarehouseService
	regionService regions.RegionService
	auditService  audit.AuditService
}

func NewWarehouseHandler(service WarehouseService, regionService regions.RegionService, auditService audit.AuditService) *WarehouseHandler {
	return &WarehouseHandler{
		service:       service,
		regionService: regionService,
		auditService:  auditService,
	}
}

func (h *WarehouseHandler) AssignStoreToWarehouse(c *gin.Context) {
	var req types.AssignStoreToWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	if err := h.service.AssignStoreToWarehouse(req); err != nil {
		utils.SendInternalServerError(c, "failed to assign store to warehouse")
		return
	}

	utils.SendMessageWithStatus(c, "Store assigned to warehouse successfully", http.StatusOK)
}

func (h *WarehouseHandler) GetAllStoresByWarehouse(c *gin.Context) {
	warehouseID, err := strconv.ParseUint(c.Param("warehouseId"), 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid warehouse ID")
		return
	}

	pagination := utils.ParsePagination(c)

	stores, err := h.service.GetAllStoresByWarehouse(uint(warehouseID), pagination)
	if err != nil {
		utils.SendInternalServerError(c, "failed to list stores")
		return
	}

	utils.SendSuccessResponseWithPagination(c, stores, pagination)
}

func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var dto types.CreateWarehouseDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	_, err := h.service.CreateWarehouse(dto)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessCreatedResponse(c, "warehouse created successfully")
}

func (h *WarehouseHandler) GetWarehouseByID(c *gin.Context) {
	warehouseIDStr := c.Param("warehouseId")
	warehouseID, err := strconv.ParseUint(warehouseIDStr, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	warehouse, err := h.service.GetWarehouseByID(uint(warehouseID))
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, warehouse)
}

func (h *WarehouseHandler) GetAllWarehouses(c *gin.Context) {
	var filter types.WarehouseFilter

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Warehouse{}); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	warehouses, err := h.service.GetAllWarehouses(&filter)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, warehouses)
}

func (h *WarehouseHandler) GetWarehouses(c *gin.Context) {
	var filter types.WarehouseFilter

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Warehouse{}); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	regionID, errH := contexts.GetRegionId(c)
	if errH != nil {
		if errors.Is(errH, contexts.ErrEmptyRegionID) {
			regionID = 0
		} else {
			utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
			return
		}
	}

	warehouses, err := h.service.GetWarehousesByRegion(regionID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve warehouses")
		return
	}

	utils.SendSuccessResponseWithPagination(c, warehouses, filter.Pagination)
}

func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	warehouseIDStr := c.Param("warehouseId")
	warehouseID, err := strconv.ParseUint(warehouseIDStr, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	var dto types.UpdateWarehouseDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	updated, err := h.service.UpdateWarehouse(uint(warehouseID), dto)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
	utils.SendSuccessResponse(c, updated)
}

func (h *WarehouseHandler) DeleteWarehouse(c *gin.Context) {
	warehouseIDStr := c.Param("warehouseId")
	warehouseID, err := strconv.ParseUint(warehouseIDStr, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	if err := h.service.DeleteWarehouse(uint(warehouseID)); err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
