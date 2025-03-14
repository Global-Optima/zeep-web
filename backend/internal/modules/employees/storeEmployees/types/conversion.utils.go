package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	storeTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
)

func MapToStoreEmployeeDTO(storeEmployee *data.StoreEmployee) *StoreEmployeeDTO {
	storeEmployee.Employee.StoreEmployee = storeEmployee
	return &StoreEmployeeDTO{
		ID:              storeEmployee.ID,
		BaseEmployeeDTO: *employeesTypes.MapToBaseEmployeeDTO(&storeEmployee.Employee),
		EmployeeID:      storeEmployee.Employee.ID,
	}
}

func MapToStoreEmployeeDetailsDTO(storeEmployee *data.StoreEmployee) *StoreEmployeeDetailsDTO {
	storeEmployee.Employee.StoreEmployee = storeEmployee
	dto := &StoreEmployeeDetailsDTO{
		ID:                     storeEmployee.ID,
		BaseEmployeeDetailsDTO: *employeesTypes.MapToBaseEmployeeDetailsDTO(&storeEmployee.Employee),
		EmployeeID:             storeEmployee.EmployeeID,
		Store:                  *storeTypes.MapToStoreDTO(&storeEmployee.Store),
	}
	return dto
}

func CreateToStoreEmployee(storeID uint, dto *employeesTypes.CreateEmployeeDTO) (*data.Employee, error) {
	if !data.IsAllowableRole(data.StoreEmployeeType, dto.Role) {
		return nil, fmt.Errorf("%w: %s for type %s", employeesTypes.ErrInvalidEmployeeRole, dto.Role, data.StoreEmployeeType)
	}

	employee, err := employeesTypes.CreateToEmployee(dto)
	if err != nil {
		return nil, err
	}
	employee.StoreEmployee = &data.StoreEmployee{
		StoreID: storeID,
		Role:    dto.Role,
	}

	return employee, nil
}
