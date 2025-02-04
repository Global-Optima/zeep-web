package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

func MapToAdminEmployeeDTO(adminEmployee *data.AdminEmployee) *AdminEmployeeDTO {
	adminEmployee.Employee.AdminEmployee = adminEmployee
	dto := &AdminEmployeeDTO{
		ID:              adminEmployee.ID,
		BaseEmployeeDTO: *employeesTypes.MapToBaseEmployeeDTO(&adminEmployee.Employee),
		EmployeeID:      adminEmployee.Employee.ID,
	}
	return dto
}

func MapToAdminEmployeeDetailsDTO(adminEmployee *data.AdminEmployee) *AdminEmployeeDetailsDTO {
	adminEmployee.Employee.AdminEmployee = adminEmployee
	dto := &AdminEmployeeDetailsDTO{
		ID:                     adminEmployee.ID,
		BaseEmployeeDetailsDTO: *employeesTypes.MapToBaseEmployeeDetailsDTO(&adminEmployee.Employee),
		EmployeeID:             adminEmployee.Employee.ID,
	}
	return dto
}

func CreateToAdminEmployee(dto *employeesTypes.CreateEmployeeDTO) (*data.Employee, error) {
	if !data.IsAllowableRole(data.AdminEmployeeType, dto.Role) {
		return nil, fmt.Errorf("%w: %s for type %s", employeesTypes.ErrInvalidEmployeeRole, dto.Role, data.AdminEmployeeType)
	}

	employee, err := employeesTypes.CreateToEmployee(dto)
	if err != nil {
		return nil, err
	}
	employee.AdminEmployee = &data.AdminEmployee{
		Role: dto.Role,
	}

	return employee, nil
}
