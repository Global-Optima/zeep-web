package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

var (
	CreateWarehouseEmployeeAuditFactory = shared.NewAuditWarehouseActionExtendedFactory(
		data.CreateOperation, data.WarehouseEmployeeComponent, &employeesTypes.CreateEmployeeDTO{})

	UpdateWarehouseEmployeeAuditFactory = shared.NewAuditWarehouseActionExtendedFactory(
		data.UpdateOperation, data.WarehouseEmployeeComponent, &UpdateWarehouseEmployeeDTO{})

	DeleteWarehouseEmployeeAuditFactory = shared.NewAuditWarehouseActionExtendedFactory(
		data.DeleteOperation, data.WarehouseEmployeeComponent, struct{}{})
)
