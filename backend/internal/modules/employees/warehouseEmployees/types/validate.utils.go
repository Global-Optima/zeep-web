package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

type UpdateWarehouseEmployeeModels struct {
	*employeesTypes.UpdateEmployeeModels
	WarehouseEmployee *data.WarehouseEmployee
}

func WarehouseEmployeeUpdateFields(input *UpdateWarehouseEmployeeDTO, role data.EmployeeRole) (*UpdateWarehouseEmployeeModels, error) {
	var warehouseEmployee = &data.WarehouseEmployee{}
	if input.WarehouseID != nil {
		warehouseEmployee.WarehouseID = *input.WarehouseID
	}

	if input.Role != nil {
		if !data.IsAllowableRole(data.WarehouseEmployeeType, *input.Role) {
			return nil, employeesTypes.ErrUnsupportedEmployeeType
		}
		if !data.CanManageRole(role, *input.Role) {
			return nil, fmt.Errorf("%s %w %s", role, employeesTypes.ErrNotAllowedToManageTheRole, *input.Role)
		}
		warehouseEmployee.Role = *input.Role
	}

	updateEmployeeModels, err := employeesTypes.PrepareUpdateFields(&input.UpdateEmployeeDTO)
	if err != nil {
		return nil, err
	}

	return &UpdateWarehouseEmployeeModels{
		UpdateEmployeeModels: updateEmployeeModels,
		WarehouseEmployee:    warehouseEmployee,
	}, nil
}
