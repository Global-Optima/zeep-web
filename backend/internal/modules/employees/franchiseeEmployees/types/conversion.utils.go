package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	franchiseesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees/types"
)

func MapToFranchiseeEmployeeDTO(franchiseeEmployee *data.FranchiseeEmployee) *FranchiseeEmployeeDTO {
	dto := &FranchiseeEmployeeDTO{
		EmployeeDTO: *employeesTypes.MapToEmployeeDTO(&franchiseeEmployee.Employee),
		Franchisee:  *franchiseesTypes.ConvertFranchiseeToDTO(&franchiseeEmployee.Franchisee),
	}
	return dto
}

func CreateToFranchiseeEmployee(franchiseeID uint, dto *employeesTypes.CreateEmployeeDTO) (*data.Employee, error) {
	if !data.IsAllowableRole(data.FranchiseeEmployeeType, dto.Role) {
		return nil, fmt.Errorf("%w: %s for type %s", employeesTypes.ErrInvalidEmployeeRole, dto.Role, data.FranchiseeEmployeeType)
	}

	employee, err := employeesTypes.CreateToEmployee(dto)
	if err != nil {
		return nil, err
	}
	employee.FranchiseeEmployee = &data.FranchiseeEmployee{
		FranchiseeID: franchiseeID,
		Role:         dto.Role,
	}

	return employee, nil
}
