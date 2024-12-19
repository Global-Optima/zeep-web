package storeWarehouses

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/pkg/errors"
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

	storeID, err := fetchStoreId(c)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse store id")
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	id, err := h.service.AddStock(storeID, &dto)
	if err != nil {
		utils.SendInternalServerError(c, "failed to add new stock")
		return
	}

	utils.SendSuccessResponse(c, gin.H{
		"message": fmt.Sprintf("store warehouse stock with id %d successfully created", id),
	})
}

func (h *StoreWarehouseHandler) AddMultipleStoreWarehouseStock(c *gin.Context) {
	var dto types.AddMultipleStockDTO

	storeID, err := fetchStoreId(c)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse store id")
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.AddMultipleStock(storeID, &dto)
	if err != nil {
		utils.SendInternalServerError(c, "failed to add new multiple stocks")
		return
	}

	utils.SendSuccessResponse(c, gin.H{
		"message": "success",
	})
}

func (h *StoreWarehouseHandler) GetStoreWarehouseStockList(c *gin.Context) {
	storeID, err := fetchStoreId(c)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse store id")
		return
	}

	stockFilter := &types.GetStockFilterQuery{}
	if err := utils.ParseQueryWithBaseFilter(c, stockFilter, &data.StoreWarehouseStock{}); err != nil {
		utils.SendBadRequestError(c, "failed to parse pagination parameters")
		return
	}

	stockList, err := h.service.GetStockList(storeID, stockFilter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to to retrieve stock list")
		return
	}

	utils.SendSuccessResponseWithPagination(c, stockList, stockFilter.GetPagination())
}

func (h *StoreWarehouseHandler) GetStoreWarehouseStockById(c *gin.Context) {
	storeID, err := fetchStoreId(c)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse store id")
		return
	}

	stockId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse stock id")
		return
	}

	ingredients, err := h.service.GetStockById(storeID, uint(stockId))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve stock")
		return
	}

	utils.SendSuccessResponse(c, ingredients)
}

func (h *StoreWarehouseHandler) UpdateStoreWarehouseStockById(c *gin.Context) {
	var input types.UpdateStockDTO

	storeID, err := fetchStoreId(c)
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
		utils.SendBadRequestError(c, "failed to bind json body")
		return
	}

	err = h.service.UpdateStockById(storeID, uint(stockId), &input)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update stock")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "stock updated successfully"})
}

func (h *StoreWarehouseHandler) DeleteStoreWarehouseStockById(c *gin.Context) {
	storeID, err := fetchStoreId(c)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse store id")
		return
	}

	stockId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse stock id")
		return
	}

	err = h.service.DeleteStockById(storeID, uint(stockId))
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete stock")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "stock deleted successfully"})
}

func fetchStoreId(c *gin.Context) (uint, error) {
	employeeClaims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return 0, err
	}

	if employeeClaims.Role == data.RoleAdmin || employeeClaims.Role == data.RoleDirector {
		storeId, err := strconv.ParseUint(c.Query("store_id"), 10, 64)
		if err != nil {
			return 0, err
		}
		if storeId == 0 {
			return 0, errors.New("invalid store id: id cannot be 0")
		}
		return uint(storeId), nil
	}

	if employeeClaims.EmployeeType == data.StoreEmployeeType {
		if employeeClaims.WorkplaceID == 0 {
			return 0, errors.New("invalid store id: id cannot be 0")
		}
		return employeeClaims.WorkplaceID, nil
	}
	return 0, fmt.Errorf("invalid employee type: forbidden for %s", data.StoreEmployeeType)
}
