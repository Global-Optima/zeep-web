package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

type UpdateRegionEmployeeModels struct {
	*employeesTypes.UpdateEmployeeModels
	RegionEmployee *data.RegionEmployee
}

func RegionEmployeeUpdateFields(regionEmployee *data.RegionEmployee, input *UpdateRegionEmployeeDTO, role data.EmployeeRole) (*UpdateRegionEmployeeModels, error) {
	if regionEmployee == nil {
		return nil, fmt.Errorf("FranchiseeEmployee is nil")
	}

	if regionEmployee.Employee.ID == 0 {
		return nil, fmt.Errorf("FranchiseeEmployee.Employee is not preloaded")
	}

	if input.RegionID != nil {
		regionEmployee.RegionID = *input.RegionID
	}

	if input.Role != nil {
		if !data.IsAllowableRole(data.RegionEmployeeType, *input.Role) {
			return nil, employeesTypes.ErrUnsupportedEmployeeType
		}
		if !data.CanManageRole(role, *input.Role) {
			return nil, fmt.Errorf("%s %w %s", role, employeesTypes.ErrNotAllowedToManageTheRole, *input.Role)
		}
		regionEmployee.Role = *input.Role
	}

	updateEmployeeModels, err := employeesTypes.PrepareUpdateFields(&regionEmployee.Employee, &input.UpdateEmployeeDTO)
	if err != nil {
		return nil, err
	}

	return &UpdateRegionEmployeeModels{
		UpdateEmployeeModels: updateEmployeeModels,
		RegionEmployee:       regionEmployee,
	}, nil
}
