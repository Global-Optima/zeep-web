package warehouse

import (
	"errors"
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
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
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := h.service.AssignStoreToWarehouse(req); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseAssign)
		return
	}

	localization.SendLocalizedResponseWithKey(c, types.Response200WarehouseAssign)
}

func (h *WarehouseHandler) GetAllStoresByWarehouse(c *gin.Context) {
	warehouseID, err := utils.ParseParam(c, "warehouseId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	pagination := utils.ParsePagination(c)

	stores, err := h.service.GetAllStoresByWarehouse(uint(warehouseID), pagination)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, stores, pagination)
}

func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var dto types.CreateWarehouseDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	warehouse, err := h.service.CreateWarehouse(dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseCreate)
		return
	}

	action := types.CreateWarehouseAuditFactory(
		&data.BaseDetails{
			ID:   warehouse.ID,
			Name: warehouse.Name,
		})
	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201Warehouse)
}

func (h *WarehouseHandler) GetWarehouseByID(c *gin.Context) {
	warehouseID, err := utils.ParseParam(c, "warehouseId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	warehouse, err := h.service.GetWarehouseByID(uint(warehouseID))
	if err != nil {
		if err == types.ErrWarehouseNotFound {
			localization.SendLocalizedResponseWithKey(c, types.Response404Warehouse)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseGet)
		return
	}

	utils.SendSuccessResponse(c, warehouse)
}

func (h *WarehouseHandler) GetAllWarehouses(c *gin.Context) {
	var filter types.WarehouseFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Warehouse{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	warehouses, err := h.service.GetAllWarehouses(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseGet)
		return
	}

	utils.SendSuccessResponse(c, warehouses)
}

func (h *WarehouseHandler) GetWarehouses(c *gin.Context) {
	var filter types.WarehouseFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Warehouse{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	regionID, errH := contexts.GetRegionId(c)
	if errH != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusUnauthorized)
		return
	}
	if regionID != nil {
		filter.RegionID = regionID
	}

	warehouses, err := h.service.GetWarehouses(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, warehouses, filter.Pagination)
}

func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	warehouseID, err := utils.ParseParam(c, "warehouseId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	var dto types.UpdateWarehouseDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	warehouse, err := h.service.GetWarehouseByID(uint(warehouseID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseUpdate)
		return
	}

	updated, err := h.service.UpdateWarehouse(uint(warehouseID), dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseUpdate)
		return
	}

	action := types.UpdateWarehouseAuditFactory(
		&data.BaseDetails{
			ID:   warehouse.ID,
			Name: warehouse.Name,
		}, &dto)
	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	// localization.SendLocalizedResponseWithKey(c, types.Response200WarehouseUpdate)
	utils.SendSuccessResponse(c, updated)
}

func (h *WarehouseHandler) DeleteWarehouse(c *gin.Context) {
	warehouseID, err := utils.ParseParam(c, "warehouseId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	warehouse, err := h.service.GetWarehouseByID(uint(warehouseID))
	if err != nil {
		if errors.Is(err, types.ErrWarehouseNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404Warehouse)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseDelete)
		return
	}

	if err := h.service.DeleteWarehouse(uint(warehouseID)); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseDelete)
		return
	}

	action := types.DeleteWarehouseAuditFactory(
		&data.BaseDetails{
			ID:   warehouse.ID,
			Name: warehouse.Name,
		})
	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200WarehouseDelete)
}
