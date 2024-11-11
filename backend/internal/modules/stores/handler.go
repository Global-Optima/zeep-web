package stores

import (
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	service StoreService
}

func NewStoreHandler(service StoreService) *StoreHandler {
	return &StoreHandler{service: service}
}

func (h *StoreHandler) GetAllStores(c *gin.Context) {
	stores, err := h.service.GetAllStores()
	if err != nil {
		utils.SendInternalError(c, "Failed to retrieve stores")
		return
	}

	utils.SuccessResponse(c, stores)
}

func (h *StoreHandler) GetStoreEmployees(c *gin.Context) {
	storeIDParam := c.Param("storeId")
	storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store ID")
		return
	}

	employees, err := h.service.GetStoreEmployees(uint(storeID))
	if err != nil {
		utils.SendInternalError(c, "Failed to retrieve employees")
		return
	}

	utils.SuccessResponse(c, employees)
}
