package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

var (
	CreateFranchiseeEmployeeAuditFactory = shared.NewAuditFranchiseeActionExtendedFactory(
		data.CreateOperation, data.FranchiseeEmployeeComponent, &employeesTypes.CreateEmployeeDTO{})

	UpdateFranchiseeEmployeeAuditFactory = shared.NewAuditFranchiseeActionExtendedFactory(
		data.UpdateOperation, data.FranchiseeEmployeeComponent, &UpdateFranchiseeEmployeeDTO{})

	DeleteFranchiseeEmployeeAuditFactory = shared.NewAuditFranchiseeActionExtendedFactory(
		data.DeleteOperation, data.FranchiseeEmployeeComponent, struct{}{})
)
