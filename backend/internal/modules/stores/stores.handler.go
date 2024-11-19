package stores

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
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
	cacheKey := "stores:all"

	cacheUtil := utils.GetCacheInstance()

	var cachedStores []types.StoreDTO
	if err := cacheUtil.Get(cacheKey, &cachedStores); err == nil {
		utils.SuccessResponse(c, cachedStores)
		return
	}

	stores, err := h.service.GetAllStores()
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve stores")
		return
	}

	if err := cacheUtil.Set(cacheKey, stores, 30*time.Minute); err != nil {
		fmt.Printf("Failed to cache stores: %v\n", err)
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

	cacheKey := utils.BuildCacheKey("storeEmployees", map[string]string{
		"storeId": storeIDParam,
	})

	cacheUtil := utils.GetCacheInstance()

	var cachedEmployees []types.EmployeeDTO
	if err := cacheUtil.Get(cacheKey, &cachedEmployees); err == nil {
		utils.SuccessResponse(c, cachedEmployees)
		return
	}

	employees, err := h.service.GetStoreEmployees(uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve employees")
		return
	}

	if err := cacheUtil.Set(cacheKey, employees, 15*time.Minute); err != nil {
		fmt.Printf("Failed to cache employees: %v\n", err)
	}

	utils.SuccessResponse(c, employees)
}
