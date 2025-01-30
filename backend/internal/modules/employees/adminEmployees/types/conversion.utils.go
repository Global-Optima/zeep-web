package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

func MapToAdminEmployeeDTO(adminEmployee *data.AdminEmployee) *AdminEmployeeDTO {
	dto := &AdminEmployeeDTO{
		EmployeeDTO: *employeesTypes.MapToEmployeeDTO(&adminEmployee.Employee),
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
