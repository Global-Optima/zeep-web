package warehouseStock

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock/types"

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
		utils.SendErrorWithStatus(c, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	var req types.ReceiveWarehouseDelivery
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	if len(req.Materials) == 0 {
		utils.SendBadRequestError(c, "No items provided in the request")
		return
	}

	if err := h.service.ReceiveInventory(warehouseID, req); err != nil {
		utils.SendInternalServerError(c, "failed to receive inventory: "+err.Error())
		return
	}

	h.recordUpdateWarehouseStockAudit(c, warehouseID, &types.WarehouseStockPayloads{ReceiveWarehouseDelivery: &req})

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
	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	var filter types.WarehouseDeliveryFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.SupplierWarehouseDelivery{}); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters: "+err.Error())
		return
	}
	filter.WarehouseID = &warehouseID

	deliveries, err := h.service.GetDeliveries(filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch deliveries")
		return
	}

	utils.SendSuccessResponseWithPagination(c, deliveries, filter.Pagination)
}

func (h *WarehouseStockHandler) GetDeliveryByID(c *gin.Context) {
	deliveryID, err := utils.ParseParam(c, "id")
	if err != nil {
		utils.SendBadRequestError(c, "invalid delivery ID")
		return
	}

	delivery, err := h.service.GetDeliveryByID(deliveryID)
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch deliveries")
		return
	}

	utils.SendSuccessResponse(c, delivery)
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

func (h *WarehouseStockHandler) GetAvailableToAddStockMaterials(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var filter types.AvailableStockMaterialFilter
	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StockMaterial{})
	if err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	stocks, err := h.service.GetAvailableToAddStockMaterials(storeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch stock materials: "+err.Error())
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

	h.recordUpdateWarehouseStockAudit(c, warehouseID, &types.WarehouseStockPayloads{UpdateWarehouseStockDTO: &dto})

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

	h.recordUpdateWarehouseStockAudit(c, warehouseID, &types.WarehouseStockPayloads{AddWarehouseStockMaterial: req})

	utils.SendMessageWithStatus(c, "Warehouse stocks added successfully", http.StatusCreated)
}

func (h *WarehouseStockHandler) recordUpdateWarehouseStockAudit(c *gin.Context, warehouseID uint, payload *types.WarehouseStockPayloads) {
	warehouse, err := h.warehouseService.GetWarehouseByID(warehouseID)
	if err != nil {
		logger.GetZapSugaredLogger().Errorf("failed to fetch warehouse with ID: %d", warehouseID)
	}

	action := types.UpdateAddWarehouseStockAuditFactory(
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
