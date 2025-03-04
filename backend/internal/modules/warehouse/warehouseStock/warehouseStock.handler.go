package warehouseStock

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock/types"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/gin-gonic/gin"
)

type WarehouseStockHandler struct {
	service           WarehouseStockService
	franchiseeService franchisees.FranchiseeService
	warehouseService  warehouse.WarehouseService
	auditService      audit.AuditService
}

func NewWarehouseStockHandler(
	service WarehouseStockService,
	franchiseeService franchisees.FranchiseeService,
	warehouseService warehouse.WarehouseService,
	auditService audit.AuditService,
) *WarehouseStockHandler {
	return &WarehouseStockHandler{
		service:           service,
		franchiseeService: franchiseeService,
		warehouseService:  warehouseService,
		auditService:      auditService,
	}
}

func (h *WarehouseStockHandler) ReceiveInventory(c *gin.Context) {
	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusUnauthorized)
		return
	}

	var req types.ReceiveWarehouseDelivery
	if err := c.ShouldBindJSON(&req); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if len(req.Materials) == 0 {
		// Using localized response key for bad request (400) on receiving inventory
		localization.SendLocalizedResponseWithKey(c, types.Response400WarehouseStockReceive)
		return
	}

	if err := h.service.ReceiveInventory(warehouseID, req); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseStockReceive)
		return
	}

	h.recordUpdateWarehouseStockAudit(c, warehouseID, &types.WarehouseStockPayloads{ReceiveWarehouseDelivery: &req})
	localization.SendLocalizedResponseWithKey(c, types.Response200WarehouseStockReceive)
}

func (h *WarehouseStockHandler) GetDeliveries(c *gin.Context) {
	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusUnauthorized)
		return
	}

	var filter types.WarehouseDeliveryFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.SupplierWarehouseDelivery{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}
	filter.WarehouseID = &warehouseID

	deliveries, err := h.service.GetDeliveries(filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseStockFetchDeliveries)
		return
	}

	utils.SendSuccessResponseWithPagination(c, deliveries, filter.Pagination)
}

func (h *WarehouseStockHandler) GetDeliveryByID(c *gin.Context) {
	deliveryID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	delivery, err := h.service.GetDeliveryByID(deliveryID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseStockFetchDelivery)
		return
	}

	utils.SendSuccessResponse(c, delivery)
}

func (h *WarehouseStockHandler) AddToStock(c *gin.Context) {
	var req types.AdjustWarehouseStock
	if err := c.ShouldBindJSON(&req); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := h.service.AddWarehouseStockMaterial(req); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseStockAddMaterial)
		return
	}

	localization.SendLocalizedResponseWithKey(c, types.Response200WarehouseStockAddMaterial)
}

func (h *WarehouseStockHandler) DeductFromStock(c *gin.Context) {
	var req types.AdjustWarehouseStock
	if err := c.ShouldBindJSON(&req); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := h.service.DeductFromStock(req); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseStockDeductStock)
		return
	}

	localization.SendLocalizedResponseWithKey(c, types.Response200WarehouseStockDeductStock)
}

func (h *WarehouseStockHandler) GetStocks(c *gin.Context) {
	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusUnauthorized)
		return
	}

	var filter types.GetWarehouseStockFilterQuery
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.WarehouseStock{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}
	filter.WarehouseID = &warehouseID

	stocks, err := h.service.GetStock(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseStockFetchStock)
		return
	}

	utils.SendSuccessResponseWithPagination(c, stocks, filter.Pagination)
}

func (h *WarehouseStockHandler) GetAvailableToAddStockMaterials(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		localization.SendLocalizedResponseWithStatus(c, errH.Status())
		return
	}

	var filter types.AvailableStockMaterialFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StockMaterial{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	stocks, err := h.service.GetAvailableToAddStockMaterials(storeID, &filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseStockFetchStock)
		return
	}

	utils.SendSuccessResponseWithPagination(c, stocks, filter.Pagination)
}

func (h *WarehouseStockHandler) GetStockMaterialDetails(c *gin.Context) {
	stockMaterialID, err := utils.ParseParam(c, "stockMaterialId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	filter, errH := contexts.GetWarehouseContextFilter(c)
	if errH != nil {
		localization.SendLocalizedResponseWithStatus(c, errH.Status())
		return
	}

	details, err := h.service.GetStockMaterialDetails(uint(stockMaterialID), filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseStockFetchDetails)
		return
	}

	utils.SendSuccessResponse(c, details)
}

func (h *WarehouseStockHandler) UpdateStock(c *gin.Context) {
	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		localization.SendLocalizedResponseWithStatus(c, errH.Status())
		return
	}

	stockMaterialID, err := utils.ParseParam(c, "stockMaterialId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	var dto types.UpdateWarehouseStockDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := h.service.UpdateStock(warehouseID, stockMaterialID, dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseStockUpdate)
		return
	}

	h.recordUpdateWarehouseStockAudit(c, warehouseID, &types.WarehouseStockPayloads{UpdateWarehouseStockDTO: &dto})
	localization.SendLocalizedResponseWithKey(c, types.Response200WarehouseStockUpdate)
}

func (h *WarehouseStockHandler) AddWarehouseStocks(c *gin.Context) {
	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		localization.SendLocalizedResponseWithStatus(c, errH.Status())
		return
	}

	var req []types.AddWarehouseStockMaterial
	if err := c.ShouldBindJSON(&req); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := h.service.AddWarehouseStocks(warehouseID, req); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseStockAddStocks)
		return
	}

	h.recordUpdateWarehouseStockAudit(c, warehouseID, &types.WarehouseStockPayloads{AddWarehouseStockMaterial: req})
	localization.SendLocalizedResponseWithKey(c, types.Response201WarehouseStock)
}

func (h *WarehouseStockHandler) recordUpdateWarehouseStockAudit(c *gin.Context, warehouseID uint, payload *types.WarehouseStockPayloads) {
	warehouse, err := h.warehouseService.GetWarehouseByID(warehouseID)
	if err != nil {
		logger.GetZapSugaredLogger().Errorf("failed to fetch warehouse with ID: %d", warehouseID)
		return
	}

	action := types.UpdateWarehouseStockAuditFactory(
		&data.BaseDetails{
			ID:   warehouseID,
			Name: warehouse.Name,
		},
		payload,
		warehouse.ID,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()
}
