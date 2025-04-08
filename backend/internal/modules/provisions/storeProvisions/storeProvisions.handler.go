package storeProvisions

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions/types"
	provisionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type StoreProvisionHandler struct {
	service           StoreProvisionService
	provisionService  provisions.ProvisionService
	franchiseeService franchisees.FranchiseeService
	auditService      audit.AuditService
	logger            *zap.SugaredLogger
}

func NewStoreProvisionHandler(
	service StoreProvisionService,
	provisionService provisions.ProvisionService,
	franchiseeService franchisees.FranchiseeService,
	auditService audit.AuditService,
	logger *zap.SugaredLogger,
) *StoreProvisionHandler {
	return &StoreProvisionHandler{
		service:           service,
		provisionService:  provisionService,
		franchiseeService: franchiseeService,
		auditService:      auditService,
		logger:            logger,
	}
}

func (h *StoreProvisionHandler) CreateStoreProvision(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var dto types.CreateStoreProvisionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	storeProvision, err := h.service.CreateStoreProvision(storeID, &dto)
	if err != nil {
		switch {
		case errors.Is(err, provisionsTypes.ErrProvisionNotFound):
			localization.SendLocalizedResponseWithKey(c, provisionsTypes.Response404Provision)
			return
		case errors.Is(err, types.ErrStoreProvisionDailyLimitReached):
			localization.SendLocalizedResponseWithKey(c, types.Response409StoreProvisionLimit)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProvisionCreate)
		return
	}

	action := types.CreateStoreProvisionAuditFactory(
		&data.BaseDetails{
			ID:   storeProvision.ID,
			Name: storeProvision.Provision.Name,
		},
		&dto, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201StoreProvision)
}

func (h *StoreProvisionHandler) GetStoreProvisionByID(c *gin.Context) {
	storeProvisionID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreProvision)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	provision, err := h.service.GetStoreProvisionByID(storeID, storeProvisionID)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrStoreProvisionNotFound):
			localization.SendLocalizedResponseWithKey(c, types.Response404StoreProvision)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProvisionGet)
		return
	}

	if provision == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, provision)
}

func (h *StoreProvisionHandler) GetStoreProvisions(c *gin.Context) {
	var filter types.StoreProvisionFilterDTO
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StoreProvision{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	storeProvisions, err := h.service.GetStoreProvisions(storeID, &filter)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrStoreProvisionNotFound):
			localization.SendLocalizedResponseWithKey(c, types.Response404StoreProvision)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProvisionGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, storeProvisions, filter.Pagination)
}

func (h *StoreProvisionHandler) CompleteStoreProvisionByID(c *gin.Context) {
	storeProvisionID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreProvision)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	storeProvision, err := h.service.CompleteStoreProvision(storeID, storeProvisionID)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrStoreProvisionNotFound):
			localization.SendLocalizedResponseWithKey(c, types.Response404StoreProvision)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response409StoreProvisionCompleted)
		return
	}

	action := types.UpdateStoreProvisionAuditFactory(
		&data.BaseDetails{
			ID:   storeProvisionID,
			Name: storeProvision.Provision.Name,
		},
		&types.UpdateStoreProvisionFields{
			Status: &storeProvision.Status,
		},
		storeID,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200StoreProvisionUpdate)
}

func (h *StoreProvisionHandler) UpdateStoreProvisionByID(c *gin.Context) {
	storeProvisionID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreProvision)
		return
	}

	var dto types.UpdateStoreProvisionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	storeProvision, err := h.service.UpdateStoreProvision(storeID, storeProvisionID, &dto)
	if err != nil {
		if errors.Is(err, types.ErrProvisionCompleted) {
			localization.SendLocalizedResponseWithKey(c, types.Response409StoreProvisionCompleted)
			return
		}
		if errors.Is(err, types.ErrStoreProvisionIngredientMismatch) {
			localization.SendLocalizedResponseWithKey(c, types.Response409StoreProvisionIngredientsMismatch)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProvisionUpdate)
		return
	}

	action := types.UpdateStoreProvisionAuditFactory(
		&data.BaseDetails{
			ID:   storeProvisionID,
			Name: storeProvision.Provision.Name,
		},
		&types.UpdateStoreProvisionFields{
			UpdateStoreProvisionDTO: &dto,
		},
		storeID,
	)
	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200StoreProvisionUpdate)
}

func (h *StoreProvisionHandler) DeleteStoreProvisionByID(c *gin.Context) {
	storeProvisionID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreProvision)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	storeProvision, err := h.service.DeleteStoreProvision(storeID, storeProvisionID)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrStoreProvisionNotFound):
			localization.SendLocalizedResponseWithKey(c, types.Response404StoreProvision)
			return
		case errors.Is(err, types.ErrProvisionCompleted):
			localization.SendLocalizedResponseWithKey(c, types.Response409StoreProvisionCompleted)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProvisionDelete)
		return
	}

	action := types.DeleteStoreProvisionAuditFactory(
		&data.BaseDetails{
			ID:   storeProvisionID,
			Name: storeProvision.Provision.Name,
		},
		nil, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200StoreProvisionDelete)
}
