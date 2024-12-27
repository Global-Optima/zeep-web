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
	warehouseID, err := strconv.ParseUint(c.Param("warehouseId"), 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid warehouse ID")
		return
	}

	pagination := utils.ParsePagination(c)

	stores, err := h.service.GetAllStoresByWarehouse(uint(warehouseID), pagination)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to list stores: "+err.Error())
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

	warehouseID, err := h.service.CreateWarehouse(dto)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessCreatedResponse(c, gin.H{"warehouseId": warehouseID})
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
	pagination := utils.ParsePagination(c)

	warehouses, err := h.service.GetAllWarehouses(pagination)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponseWithPagination(c, warehouses, pagination)
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

func (h *WarehouseHandler) AddToStock(c *gin.Context) {
	var req types.AdjustWarehouseStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.service.AddToStock(req); err != nil {
		utils.SendInternalServerError(c, "Failed to add stock: "+err.Error())
		return
	}

	utils.SendMessageWithStatus(c, "Stock added successfully", http.StatusOK)
}

func (h *WarehouseHandler) DeductFromStock(c *gin.Context) {
	var req types.AdjustWarehouseStockRequest
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

func (h *WarehouseHandler) GetStock(c *gin.Context) {
	queryParams, err := types.ParseStockFilterParamsWithPagination(c)
	if err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	stocks, err := h.service.GetStock(queryParams)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch stock: "+err.Error())
		return
	}

	utils.SendSuccessResponseWithPagination(c, stocks, queryParams.Pagination)
}

func (h *WarehouseHandler) GetStockMaterialDetails(c *gin.Context) {
	stockMaterialID, err := strconv.ParseUint(c.Param("stockMaterialId"), 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "invalid stock material ID")
		return
	}

	warehouseID, err := strconv.ParseUint(c.Query("warehouseId"), 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "invalid warehouse ID")
		return
	}

	details, err := h.service.GetStockMaterialDetails(uint(stockMaterialID), uint(warehouseID))
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, details)
}

func (h *WarehouseHandler) ResetStock(c *gin.Context) {
	var req types.ResetWarehouseStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.service.ResetStock(req); err != nil {
		utils.SendInternalServerError(c, "Failed to reset stock: "+err.Error())
		return
	}

	utils.SendMessageWithStatus(c, "Stock reset successfully", http.StatusOK)
}
