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

func RegionEmployeeUpdateFields(input *UpdateRegionEmployeeDTO, role data.EmployeeRole) (*UpdateRegionEmployeeModels, error) {
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

	updateEmployeeModels, err := employeesTypes.PrepareUpdateFields(&input.UpdateEmployeeDTO)
	if err != nil {
		return nil, err
	}

	return &UpdateRegionEmployeeModels{
		UpdateEmployeeModels: updateEmployeeModels,
		RegionEmployee:       regionEmployee,
	}, nil
}
