package stores

import (
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	service           StoreService
	franchiseeService franchisees.FranchiseeService
	auditService      audit.AuditService
}

func NewStoreHandler(service StoreService, franchiseeService franchisees.FranchiseeService, auditService audit.AuditService) *StoreHandler {
	return &StoreHandler{
		service:           service,
		franchiseeService: franchiseeService,
		auditService:      auditService,
	}
}

func (h *StoreHandler) GetAllStores(c *gin.Context) {
	var filter types.StoreFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Store{}); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	stores, err := h.service.GetAllStores(&filter)
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

	id, err := h.service.CreateStore(&storeDTO)
	if err != nil {
		utils.SendInternalServerError(c, "failed to create store")
		return
	}

	action := types.CreateStoreAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: storeDTO.Name,
		})

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessCreatedResponse(c, "store created successfully")
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

func (h *StoreHandler) GetStoresByFranchisee(c *gin.Context) {
	var filter types.StoreFilter

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Store{}); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	franchiseeID, errH := contexts.GetFranchiseeId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}
	if franchiseeID != nil {
		filter.FranchiseeID = franchiseeID
	}

	store, err := h.service.GetStores(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve stores")
		return
	}

	utils.SendSuccessResponseWithPagination(c, store, filter.Pagination)
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

	err = h.service.UpdateStore(uint(storeID), &dto)
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

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendMessageWithStatus(c, "store updated successfully", http.StatusOK)
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

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{"message": "store deleted successfully"})
}
