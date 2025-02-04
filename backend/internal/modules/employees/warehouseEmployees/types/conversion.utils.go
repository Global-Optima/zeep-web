package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	warehouseTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
)

func MapToWarehouseEmployeeDTO(warehouseEmployee *data.WarehouseEmployee) *WarehouseEmployeeDTO {
	warehouseEmployee.Employee.WarehouseEmployee = warehouseEmployee
	return &WarehouseEmployeeDTO{
		ID:              warehouseEmployee.ID,
		BaseEmployeeDTO: *employeesTypes.MapToBaseEmployeeDTO(&warehouseEmployee.Employee),
		EmployeeID:      warehouseEmployee.Employee.ID,
	}
}

func MapToWarehouseEmployeeDetailsDTO(warehouseEmployee *data.WarehouseEmployee) *WarehouseEmployeeDetailsDTO {
	warehouseEmployee.Employee.WarehouseEmployee = warehouseEmployee
	dto := &WarehouseEmployeeDetailsDTO{
		ID:                     warehouseEmployee.ID,
		BaseEmployeeDetailsDTO: *employeesTypes.MapToBaseEmployeeDetailsDTO(&warehouseEmployee.Employee),
		EmployeeID:             warehouseEmployee.EmployeeID,
		Warehouse:              *warehouseTypes.ToWarehouseDTO(warehouseEmployee.Warehouse),
	}
	return dto
}

func CreateToWarehouseEmployee(warehouseID uint, dto *employeesTypes.CreateEmployeeDTO) (*data.Employee, error) {
	if !data.IsAllowableRole(data.WarehouseEmployeeType, dto.Role) {
		return nil, fmt.Errorf("%w: %s for type %s", employeesTypes.ErrInvalidEmployeeRole, dto.Role, data.WarehouseEmployeeType)
	}

	employee, err := employeesTypes.CreateToEmployee(dto)
	if err != nil {
		return nil, err
	}
	employee.WarehouseEmployee = &data.WarehouseEmployee{
		WarehouseID: warehouseID,
		Role:        dto.Role,
	}

	return employee, nil
}
