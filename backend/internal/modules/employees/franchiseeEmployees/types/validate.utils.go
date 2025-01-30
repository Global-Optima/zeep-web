package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

type UpdateModels struct {
	Employee           *data.Employee
	FranchiseeEmployee *data.FranchiseeEmployee
}

func FranchiseeEmployeeUpdateFields(input *UpdateFranchiseeEmployeeDTO, role data.EmployeeRole) (*UpdateModels, error) {
	var franchiseeEmployee = &data.FranchiseeEmployee{}
	if input.FranchiseeID != nil {
		franchiseeEmployee.FranchiseeID = *input.FranchiseeID
	}

	if input.Role != nil {
		if !data.IsAllowableRole(data.FranchiseeEmployeeType, *input.Role) {
			return nil, employeesTypes.ErrUnsupportedEmployeeType
		}
		if !data.CanManageRole(role, *input.Role) {
			return nil, fmt.Errorf("%s %w %s", role, employeesTypes.ErrNotAllowedToManageTheRole, *input.Role)
		}
		franchiseeEmployee.Role = *input.Role
	}

	return &UpdateModels{
		Employee:           employeesTypes.PrepareUpdateFields(&input.UpdateEmployeeDTO),
		FranchiseeEmployee: franchiseeEmployee,
	}, nil
}
