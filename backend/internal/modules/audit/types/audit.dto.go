package types

import (
	"encoding/json"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"time"
)

type EmployeeAuditDTO struct {
	ID            uint            `json:"id"`
	Timestamp     time.Time       `json:"timestamp"`
	EmployeeID    uint            `json:"employeeId"`
	OperationType string          `json:"operationType"`
	ComponentName string          `json:"componentName"`
	Details       json.RawMessage `json:"details"`
	IPAddress     string          `json:"ipAddress"`
	ResourceURL   string          `json:"resourceUrl"`
	Method        string          `json:"method"`
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
