package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

var (
	CreateAdminEmployeeAuditFactory = shared.NewAuditActionExtendedFactory(
		data.CreateOperation, data.AdminEmployeeComponent, &employeesTypes.CreateEmployeeDTO{})
)
