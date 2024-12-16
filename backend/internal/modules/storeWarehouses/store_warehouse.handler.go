package storeWarehouses

import (
	"fmt"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StoreWarehouseHandler struct {
	service StoreWarehouseService
}

func NewStoreWarehouseHandler(service StoreWarehouseService) *StoreWarehouseHandler {
	return &StoreWarehouseHandler{service: service}
}

func (h *StoreWarehouseHandler) AddStoreWarehouseStock(c *gin.Context) {
	var dto types.AddStockDTO

	storeID, err := strconv.ParseUint(c.Param("store_id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	id, err := h.service.AddStock(uint(storeID), &dto)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": fmt.Sprintf("store warehouse stock with id %d successfully created", id),
	})
}

func (h *StoreWarehouseHandler) AddMultipleStoreWarehouseStock(c *gin.Context) {
	var dto types.AddMultipleStockDTO

	storeID, err := strconv.ParseUint(c.Param("store_id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	err = h.service.AddMultipleStock(uint(storeID), &dto)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "success",
	})
}

func (h *StoreWarehouseHandler) GetStoreWarehouseStockList(c *gin.Context) {
	storeID, err := strconv.ParseUint(c.Param("store_id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	queryParams, err := types.ParseStockParamsWithPagination(c)
	if err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	stockList, err := h.service.GetStockList(uint(storeID), queryParams)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponseWithPagination(c, stockList, queryParams.Pagination)
}

func (h *StoreWarehouseHandler) GetStoreWarehouseStockById(c *gin.Context) {
	storeId, err := strconv.ParseUint(c.Param("store_id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store id")
		return
	}

	stockId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse stock id")
		return
	}

	ingredients, err := h.service.GetStockById(uint(storeId), uint(stockId))
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, ingredients)
}

func (h *StoreWarehouseHandler) UpdateStoreWarehouseStockById(c *gin.Context) {
	var input types.UpdateStockDTO

	storeId, err := strconv.ParseUint(c.Param("store_id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse store id")
		return
	}

	stockId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse stock id")
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	err = h.service.UpdateStockById(uint(storeId), uint(stockId), &input)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "stock updated successfully"})
}

func (h *StoreWarehouseHandler) DeleteStoreWarehouseStockById(c *gin.Context) {
	storeId, err := strconv.ParseUint(c.Param("store_id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse stock id")
		return
	}

	stockId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse stock id")
		return
	}

	err = h.service.DeleteStockById(uint(storeId), uint(stockId))
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "stock deleted successfully"})
}
