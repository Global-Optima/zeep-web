package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

type UpdateFranchiseeEmployeeModels struct {
	*employeesTypes.UpdateEmployeeModels
	FranchiseeEmployee *data.FranchiseeEmployee
}

func FranchiseeEmployeeUpdateFields(franchiseeEmployee *data.FranchiseeEmployee, input *UpdateFranchiseeEmployeeDTO, role data.EmployeeRole) (*UpdateFranchiseeEmployeeModels, error) {
	if franchiseeEmployee == nil {
		return nil, fmt.Errorf("FranchiseeEmployee is nil")
	}

	if franchiseeEmployee.Employee.ID == 0 {
		return nil, fmt.Errorf("FranchiseeEmployee.Employee is not preloaded")
	}

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

	updateEmployeeModels, err := employeesTypes.PrepareUpdateFields(&franchiseeEmployee.Employee, &input.UpdateEmployeeDTO)
	if err != nil {
		return nil, err
	}

	return &UpdateFranchiseeEmployeeModels{
		UpdateEmployeeModels: updateEmployeeModels,
		FranchiseeEmployee:   franchiseeEmployee,
	}, nil
}
