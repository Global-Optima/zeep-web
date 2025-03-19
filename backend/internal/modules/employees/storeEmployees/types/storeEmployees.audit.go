package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

var (
	CreateStoreEmployeeAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.CreateOperation, data.StoreEmployeeComponent, &employeesTypes.CreateEmployeeDTO{})

	UpdateStoreEmployeeAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.UpdateOperation, data.StoreEmployeeComponent, &UpdateStoreEmployeeDTO{})

	DeleteStoreEmployeeAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.DeleteOperation, data.StoreEmployeeComponent, struct{}{})
)
