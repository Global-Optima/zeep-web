package storeWarehouses

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type StoreWarehouseHandler struct {
	service StoreWarehouseService
}

func NewStoreWarehouseHandler(service StoreWarehouseService) *StoreWarehouseHandler {
	return &StoreWarehouseHandler{service: service}
}

func (h *StoreWarehouseHandler) GetStoreWarehouseStockList(c *gin.Context) {
	queryParams, err := types.ParseStoreWarehouseIngredientParams(c.Request.URL.Query())
	if err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	ingredients, err := h.service.GetStoreWarehouseStockList(*queryParams)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve products")
		return
	}

	utils.SuccessResponse(c, ingredients)
}

func (h *StoreWarehouseHandler) GetStoreWarehouseStockById(c *gin.Context) {

	storeId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse stock id")
		return
	}

	ingredientId, err := strconv.ParseUint(c.Param("ingredientId"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "ingredientId is required")
		return
	}

	ingredients, err := h.service.GetStoreWarehouseStockById(uint(storeId), uint(ingredientId))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve products")
		return
	}

	utils.SuccessResponse(c, ingredients)
}

func (h *StoreWarehouseHandler) UpdateStoreWarehouseIngredient(c *gin.Context) {
	var input types.UpdateStoreWarehouseIngredientDTO

	if c.Param("id") != "" {
		utils.SendBadRequestError(c, "empty store warehouse stock id")
		return
	}

	storeWarehouseStockId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse stock id")
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	err = h.service.UpdateStoreWarehouseStockById(uint(storeWarehouseStockId), input)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve products")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "stock updated successfully"})
}
