package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

type UpdateStoreEmployeeModels struct {
	*employeesTypes.UpdateEmployeeModels
	StoreEmployee *data.StoreEmployee
}

func StoreEmployeeUpdateFields(input *UpdateStoreEmployeeDTO, role data.EmployeeRole) (*UpdateStoreEmployeeModels, error) {
	storeEmployee := &data.StoreEmployee{}
	if input.StoreID != nil {
		storeEmployee.StoreID = *input.StoreID
	}

	if input.Role != nil {
		if !data.IsAllowableRole(data.StoreEmployeeType, *input.Role) {
			return nil, employeesTypes.ErrUnsupportedEmployeeType
		}
		if !data.CanManageRole(role, *input.Role) {
			return nil, fmt.Errorf("%s %w %s", role, employeesTypes.ErrNotAllowedToManageTheRole, *input.Role)
		}
		storeEmployee.Role = *input.Role
	}

	employeeUpdateModels, err := employeesTypes.PrepareUpdateFields(&input.UpdateEmployeeDTO)
	if err != nil {
		return nil, err
	}

	return &UpdateStoreEmployeeModels{
		StoreEmployee:        storeEmployee,
		UpdateEmployeeModels: employeeUpdateModels,
	}, nil
}
