package types

import (
	"encoding/json"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type EmployeeAuditDTO struct {
	ID                         uint                          `json:"id"`
	Timestamp                  time.Time                     `json:"timestamp"`
	OperationType              data.OperationType            `json:"operationType"`
	ComponentName              data.ComponentName            `json:"componentName"`
	LocalizedMessages          localization.LocalizedMessage `json:"localizedMessages"`
	IPAddress                  string                        `json:"ipAddress"`
	ResourceURL                string                        `json:"resourceUrl"`
	Method                     string                        `json:"method"`
	Details                    json.RawMessage               `json:"details"`
	employeesTypes.EmployeeDTO `json:"employee"`
}

type EmployeeAuditFilter struct {
	utils.BaseFilter
	MinTimestamp  *time.Time `form:"minTimestamp" time_format:"2006-01-02T15:04:05Z07:00"`
	MaxTimestamp  *time.Time `form:"maxTimestamp" time_format:"2006-01-02T15:04:05Z07:00"`
	OperationType *string    `form:"operationType"`
	ComponentName *string    `form:"componentName"`
	EmployeeID    *uint      `form:"employeeId"`
	Method        *string    `form:"method"`
	Search        *string    `form:"search"`
}
