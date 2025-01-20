package types

import (
	"encoding/json"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ConvertToEmployeeAuditDTO(audit *data.EmployeeAudit) (*EmployeeAuditDTO, error) {

	logrus.Info(audit.OperationType, audit.ComponentName)

	employee := data.Employee{}
	if &audit.Employee != nil {
		employee = audit.Employee
	}

	messages, err := MapLocalizedMessages(audit)
	if err != nil {
		return nil, err
	}

	return &EmployeeAuditDTO{
		ID:            audit.ID,
		Timestamp:     audit.BaseEntity.CreatedAt,
		OperationType: audit.OperationType.ToString(),
		ComponentName: audit.ComponentName.ToString(),
		Messages:      *messages,
		IPAddress:     audit.IPAddress,
		ResourceURL:   audit.ResourceUrl,
		Method:        audit.Method.ToString(),
		EmployeeDTO:   *employeesTypes.MapToEmployeeDTO(&employee),
	}, nil
}

func MapLocalizedMessages(audit *data.EmployeeAudit) (*Messages, error) {
	messages := &Messages{}

	core := shared.AuditActionCore{
		OperationType: audit.OperationType,
		ComponentName: audit.ComponentName,
	}

	detailsFactory := shared.GetAuditActionDetailsFactory(core)
	if detailsFactory == nil {
		return nil, fmt.Errorf("no details factory found for operationType: %s, componentName: %s", audit.OperationType, audit.ComponentName)
	}

	details := detailsFactory()
	err := json.Unmarshal(audit.Details, details)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal audit details: %w", err)
	}

	messages.En, err = localization.Translate(localization.English, "createAuditMessage", map[string]interface{}{
		"ComponentName": localization.GetLocalizedComponentName(localization.English, audit.ComponentName),
		"Name":          details.GetBaseDetails().Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate localized message: %w", err)
	}

	messages.Ru, err = localization.Translate(localization.Russian, "createAuditMessage", map[string]interface{}{
		"ComponentName": localization.GetLocalizedComponentName(localization.Russian, audit.ComponentName),
		"Name":          details.GetBaseDetails().Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate localized message: %w", err)
	}

	messages.Kk, err = localization.Translate(localization.Kazakh, "createAuditMessage", map[string]interface{}{
		"ComponentName": localization.GetLocalizedComponentName(localization.Kazakh, audit.ComponentName),
		"Name":          details.GetBaseDetails().Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate localized message: %w", err)
	}

	return messages, nil
}

func MapToEmployeeAudit(c *gin.Context, action shared.AuditAction) (*data.EmployeeAudit, error) {
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
