package types

import (
	"encoding/json"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/gin-gonic/gin"
)

func ConvertToEmployeeAuditDTO(audit *data.EmployeeAudit) (*EmployeeAuditDTO, error) {
	details, err := json.Marshal(audit.Details)
	if err != nil {
		return nil, err
	}

	return &EmployeeAuditDTO{
		ID:            audit.ID,
		Timestamp:     audit.BaseEntity.CreatedAt,
		EmployeeID:    audit.EmployeeID,
		OperationType: audit.OperationType.ToString(),
		ComponentName: audit.ComponentName.ToString(),
		Details:       details,
		IPAddress:     audit.IPAddress,
		ResourceURL:   audit.ResourceUrl,
		Method:        audit.Method.ToString(),
	}, nil
}

func MapToEmployeeAudit(c *gin.Context, action AuditAction) (*data.EmployeeAudit, error) {
	claims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return nil, err
	}

	core := action.GetActionCore()

	action.GetActionDetails()

	method, err := data.ToHTTPMethod(c.Request.Method)
	if err != nil {
		return nil, err
	}

	detailsJSONB, err := action.GetActionDetails().ToDetails()
	if err != nil {
		return nil, err
	}

	return &data.EmployeeAudit{
		EmployeeID:    claims.EmployeeClaimsData.ID,
		OperationType: core.OperationType,
		ComponentName: core.ComponentName,
		Details:       detailsJSONB,
		IPAddress:     c.ClientIP(),
		ResourceUrl:   c.Request.URL.String(),
		Method:        method,
	}, nil
}
