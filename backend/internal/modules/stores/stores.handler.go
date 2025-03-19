package stores

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"

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
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	stores, err := h.service.GetAllStores(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreGet)
		return
	}

	utils.SendSuccessResponse(c, stores)
}

func (h *StoreHandler) CreateStore(c *gin.Context) {
	var storeDTO types.CreateStoreDTO

	if err := c.ShouldBindJSON(&storeDTO); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	id, err := h.service.CreateStore(&storeDTO)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreCreate)
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

	localization.SendLocalizedResponseWithKey(c, types.Response201Store)
}

func (h *StoreHandler) GetStoreByID(c *gin.Context) {
	storeID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	store, err := h.service.GetStoreByID(uint(storeID))
	if err != nil {
		if err == types.ErrStoreNotFound {
			localization.SendLocalizedResponseWithKey(c, types.Response404Store)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreGet)
		return
	}

	utils.SendSuccessResponse(c, store)
}

func (h *StoreHandler) GetStoresByFranchisee(c *gin.Context) {
	var filter types.StoreFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Store{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	franchiseeID, errH := contexts.GetFranchiseeId(c)
	if errH != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusUnauthorized)
		return
	}
	if franchiseeID != nil {
		filter.FranchiseeID = franchiseeID
	}

	stores, err := h.service.GetStores(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, stores, filter.Pagination)
}

func (h *StoreHandler) UpdateStore(c *gin.Context) {
	var dto types.UpdateStoreDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	storeID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	existingStore, err := h.service.GetStoreByID(uint(storeID))
	if err != nil {
		if err == types.ErrStoreNotFound {
			localization.SendLocalizedResponseWithKey(c, types.Response404Store)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreGet)
		return
	}

	err = h.service.UpdateStore(uint(storeID), &dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreUpdate)
		return
	}

	action := types.UpdateStoreAuditFactory(
		&data.BaseDetails{
			ID:   uint(storeID),
			Name: existingStore.Name,
		},
		&dto,
	)
	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200StoreUpdate)
}

func (h *StoreHandler) DeleteStore(c *gin.Context) {
	storeID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	hardDelete := c.Query("hardDelete") == "true"

	if err := h.service.DeleteStore(uint(storeID), hardDelete); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreDelete)
		return
	}

	existingStore, err := h.service.GetStoreByID(uint(storeID))
	if err != nil {
		if err == types.ErrStoreNotFound {
			localization.SendLocalizedResponseWithKey(c, types.Response404Store)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreGet)
		return
	}

	action := types.DeleteStoreAuditFactory(
		&data.BaseDetails{
			ID:   uint(storeID),
			Name: existingStore.Name,
		})
	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200StoreDelete)
}
