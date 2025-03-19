package audit

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
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
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	audits, err := h.service.GetAuditRecords(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve Audits")
		return
	}

	utils.SendSuccessResponseWithPagination(c, audits, filter.Pagination)
}
