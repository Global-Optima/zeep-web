package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

type UpdateModels struct {
	Employee       *data.Employee
	RegionEmployee *data.RegionEmployee
}

func RegionEmployeeUpdateFields(input *UpdateRegionEmployeeDTO, role data.EmployeeRole) (*UpdateModels, error) {
	var regionEmployee = &data.RegionEmployee{}
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

	return &UpdateModels{
		Employee:       employeesTypes.PrepareUpdateFields(&input.UpdateEmployeeDTO),
		RegionEmployee: regionEmployee,
	}, nil
}
