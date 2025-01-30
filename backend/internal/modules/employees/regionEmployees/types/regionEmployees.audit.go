package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

var (
	CreateRegionEmployeeAuditFactory = shared.NewAuditRegionActionExtendedFactory(
		data.CreateOperation, data.RegionEmployeeComponent, &employeesTypes.CreateEmployeeDTO{})

	UpdateRegionEmployeeAuditFactory = shared.NewAuditRegionActionExtendedFactory(
		data.UpdateOperation, data.RegionEmployeeComponent, &UpdateRegionEmployeeDTO{})

	DeleteRegionEmployeeAuditFactory = shared.NewAuditRegionActionExtendedFactory(
		data.DeleteOperation, data.RegionEmployeeComponent, struct{}{})
)
