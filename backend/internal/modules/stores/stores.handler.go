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
	searchTerm := c.Query("searchTerm")
	cacheKey := "stores:all"

	if searchTerm != "" {
		cacheKey = "stores:" + searchTerm
	}

	cacheUtil := utils.GetCacheInstance()

	var cachedStores []types.StoreDTO
	if err := cacheUtil.Get(cacheKey, &cachedStores); err == nil {
		if !utils.IsEmpty(cachedStores) {
			utils.SendSuccessResponse(c, cachedStores)
			return
		}
	}

	stores, err := h.service.GetAllStores(searchTerm)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve stores")
		return
	}

	// Cache the fresh data
	if err := cacheUtil.Set(cacheKey, stores, 30*time.Minute); err != nil {
		fmt.Printf("Failed to cache stores: %v\n", err)
	}

	utils.SendSuccessResponse(c, stores)
}

func (h *StoreHandler) CreateStore(c *gin.Context) {
	var storeDTO types.CreateStoreDTO

	if err := c.ShouldBindJSON(&storeDTO); err != nil {
		utils.SendBadRequestError(c, "Invalid input: "+err.Error())
		return
	}

	createdStore, err := h.service.CreateStore(storeDTO)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create store: "+err.Error())
		return
	}

	utils.SendSuccessResponse(c, createdStore)
}

func (h *StoreHandler) GetStoreByID(c *gin.Context) {

	storeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store ID")
		return
	}

	store, err := h.service.GetStoreByID(uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve store: "+err.Error())
		return
	}

	utils.SendSuccessResponse(c, store)
}

func (h *StoreHandler) UpdateStore(c *gin.Context) {
	var dto types.UpdateStoreDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid input: "+err.Error())
		return
	}

	storeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store ID")
		return
	}

	updatedStore, err := h.service.UpdateStore(uint(storeID), dto)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update store: "+err.Error())
		return
	}

	utils.SendSuccessResponse(c, updatedStore)
}

func (h *StoreHandler) DeleteStore(c *gin.Context) {

	storeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store ID")
		return
	}

	hardDelete := c.Query("hardDelete") == "true"

	if err := h.service.DeleteStore(uint(storeID), hardDelete); err != nil {
		utils.SendInternalServerError(c, "Failed to delete store: "+err.Error())
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Store deleted successfully"})
}
