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

func WarehouseEmployeeUpdateFields(warehouseEmployee *data.WarehouseEmployee, input *UpdateWarehouseEmployeeDTO, role data.EmployeeRole) (*UpdateWarehouseEmployeeModels, error) {
	if warehouseEmployee == nil {
		return nil, fmt.Errorf("FranchiseeEmployee is nil")
	}

	if warehouseEmployee.Employee.ID == 0 {
		return nil, fmt.Errorf("FranchiseeEmployee.Employee is not preloaded")
	}

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

	updateEmployeeModels, err := employeesTypes.PrepareUpdateFields(&warehouseEmployee.Employee, &input.UpdateEmployeeDTO)
	if err != nil {
		return nil, err
	}

	return &UpdateWarehouseEmployeeModels{
		UpdateEmployeeModels: updateEmployeeModels,
		WarehouseEmployee:    warehouseEmployee,
	}, nil
}
