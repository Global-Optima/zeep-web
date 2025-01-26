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
)

const (
	AUDIT_TRANSLATION_KEY = "audit"
	NAME_KEY              = "Name"
	STORE_NAME_KEY        = "StoreName"
	WAREHOUSE_NAME_KEY    = "WarehouseName"
)

func ConvertToEmployeeAuditDTO(audit *data.EmployeeAudit) (*EmployeeAuditDTO, error) {
	employee := data.Employee{}
	if audit.Employee.ID != 0 {
		employee = audit.Employee
	}

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

	messages, err := MapLocalizedMessages(audit, details)
	if err != nil {
		return nil, err
	}

	return &EmployeeAuditDTO{
		ID:                audit.ID,
		Timestamp:         audit.BaseEntity.CreatedAt,
		OperationType:     audit.OperationType.ToString(),
		ComponentName:     audit.ComponentName.ToString(),
		LocalizedMessages: *messages,
		IPAddress:         audit.IPAddress,
		ResourceURL:       audit.ResourceUrl,
		Details:           json.RawMessage(audit.Details),
		Method:            audit.Method.ToString(),
		BaseEmployeeDTO:   *employeesTypes.MapToBaseEmployeeDTO(&employee),
	}, nil
}

func MapLocalizedMessages(audit *data.EmployeeAudit, details data.AuditDetails) (*localization.LocalizedMessages, error) {
	var err error
	var messages *localization.LocalizedMessages

	key := localization.FormTranslationKey(AUDIT_TRANSLATION_KEY, audit.OperationType.ToString(), audit.ComponentName.ToString())

	switch details := details.(type) {
	case *data.BaseDetails:
		messages, err = localization.Translate(key, map[string]interface{}{
			NAME_KEY: details.GetBaseDetails().Name,
		})
	case *data.ExtendedDetails:
		messages, err = localization.Translate(key, map[string]interface{}{
			NAME_KEY: details.GetBaseDetails().Name,
		})
	case *data.ExtendedDetailsStore:
		messages, err = localization.Translate(key, map[string]interface{}{
			NAME_KEY:       details.GetBaseDetails().Name,
			STORE_NAME_KEY: details.StoreInfo.StoreName,
		})
	case *data.ExtendedDetailsWarehouse:
		messages, err = localization.Translate(key, map[string]interface{}{
			NAME_KEY:           details.GetBaseDetails().Name,
			WAREHOUSE_NAME_KEY: details.WarehouseInfo.WarehouseName,
		})
	default:
		return nil, fmt.Errorf("unsupported type: %T", details)
	}

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
