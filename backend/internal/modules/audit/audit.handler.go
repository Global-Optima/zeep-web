package audit

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	service AuditService
}

func NewAuditHandler(service AuditService) *AuditHandler {
	return &AuditHandler{
		service: service,
	}
}

func (h *AuditHandler) GetAudits(c *gin.Context) {
	var filter types.EmployeeAuditFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.EmployeeAudit{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	audits, err := h.service.GetAuditRecords(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AuditGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, audits, filter.Pagination)
}
