package stores

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	service      StoreService
	auditService audit.AuditService
}

func NewStoreHandler(service StoreService, auditService audit.AuditService) *StoreHandler {
	return &StoreHandler{
		service:      service,
		auditService: auditService,
	}
}

func (h *StoreHandler) GetAllStores(c *gin.Context) {
	searchTerm := c.Query("searchTerm")

	stores, err := h.service.GetAllStores(searchTerm)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve stores")
		return
	}

	utils.SendSuccessResponse(c, stores)
}

func (h *StoreHandler) CreateStore(c *gin.Context) {
	var storeDTO types.CreateStoreDTO

	if err := c.ShouldBindJSON(&storeDTO); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	createdStore, err := h.service.CreateStore(storeDTO)
	if err != nil {
		utils.SendInternalServerError(c, "failed to create store")
		return
	}

	action := types.CreateStoreAuditFactory(
		&data.BaseDetails{
			ID:   createdStore.ID,
			Name: storeDTO.Name,
		})

	_ = h.auditService.RecordEmployeeAction(c, &action)

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
		utils.SendInternalServerError(c, "failed to retrieve store")
		return
	}

	utils.SendSuccessResponse(c, store)
}

func (h *StoreHandler) UpdateStore(c *gin.Context) {
	var dto types.UpdateStoreDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	storeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	existingProduct, err := h.service.GetStoreByID(uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update product details: product not found")
		return
	}

	updatedStore, err := h.service.UpdateStore(uint(storeID), dto)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update store")
		return
	}

	action := types.UpdateStoreAuditFactory(
		&data.BaseDetails{
			ID:   uint(storeID),
			Name: existingProduct.Name,
		},
		&dto)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, updatedStore)
}

func (h *StoreHandler) DeleteStore(c *gin.Context) {

	storeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	hardDelete := c.Query("hardDelete") == "true"

	if err := h.service.DeleteStore(uint(storeID), hardDelete); err != nil {
		utils.SendInternalServerError(c, "failed to delete store")
		return
	}

	existingProduct, err := h.service.GetStoreByID(uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update product details: product not found")
		return
	}

	action := types.DeleteStoreAuditFactory(
		&data.BaseDetails{
			ID:   uint(storeID),
			Name: existingProduct.Name,
		})

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "store deleted successfully"})
}
