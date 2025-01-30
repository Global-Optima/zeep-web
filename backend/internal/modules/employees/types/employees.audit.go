package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	UpdateEmployeeAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.EmployeeComponent, &ReassignEmployeeTypeDTO{})
)
