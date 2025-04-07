package provisions

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"net/http"
)

type ProvisionHandler struct {
	service      ProvisionService
	auditService audit.AuditService
	logger       *zap.SugaredLogger
}

func NewProvisionHandler(service ProvisionService, auditService audit.AuditService, logger *zap.SugaredLogger) *ProvisionHandler {
	return &ProvisionHandler{
		service:      service,
		auditService: auditService,
		logger:       logger,
	}
}

func (h *ProvisionHandler) CreateProvision(c *gin.Context) {
	var dto types.CreateProvisionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		logrus.Info(err)
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	id, err := h.service.CreateProvision(&dto)
	if err != nil {
		if errors.Is(err, types.ErrProvisionUniqueName) {
			localization.SendLocalizedResponseWithKey(c, types.Response409ProvisionCreateDuplicate)
			return
		}

		localization.SendLocalizedResponseWithKey(c, types.Response500ProvisionCreate)
		return
	}

	action := types.CreateProvisionAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		})

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201Provision)
}

func (h *ProvisionHandler) GetProvisionByID(c *gin.Context) {
	provisionID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Provision)
		return
	}

	provisionDetails, err := h.service.GetProvisionByID(provisionID)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrProvisionNotFound):
			localization.SendLocalizedResponseWithKey(c, types.Response404Provision)
			return
		default:
			localization.SendLocalizedResponseWithKey(c, types.Response500ProvisionGet)
			return
		}
	}

	if provisionDetails == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, provisionDetails)
}

func (h *ProvisionHandler) GetProvisions(c *gin.Context) {
	var filter types.ProvisionFilterDTO
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Provision{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	provisions, err := h.service.GetProvisions(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500ProvisionGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, provisions, filter.Pagination)
}

func (h *ProvisionHandler) UpdateProvisionByID(c *gin.Context) {
	provisionID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Provision)
		return
	}

	var dto types.UpdateProvisionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	existingProvision, err := h.service.UpdateProvision(provisionID, &dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500ProvisionUpdate)
		return
	}

	action := types.UpdateProvisionAuditFactory(
		&data.BaseDetails{
			ID:   provisionID,
			Name: existingProvision.Name,
		},
		&dto,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200ProvisionUpdate)
}

func (h *ProvisionHandler) DeleteProvisionByID(c *gin.Context) {
	provisionID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Provision)
		return
	}

	existingProvision, err := h.service.DeleteProvision(provisionID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500ProvisionDelete)
		return
	}

	action := types.DeleteProvisionAuditFactory(
		&data.BaseDetails{
			ID:   provisionID,
			Name: existingProvision.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200ProvisionDelete)
}
