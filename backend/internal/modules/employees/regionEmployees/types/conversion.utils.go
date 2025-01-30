package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	regionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/regions/types"
)

func MapToRegionEmployeeDTO(regionEmployee *data.RegionEmployee) *RegionEmployeeDTO {
	dto := &RegionEmployeeDTO{
		EmployeeDTO: *employeesTypes.MapToEmployeeDTO(&regionEmployee.Employee),
		Region:      *regionsTypes.MapRegionToDTO(&regionEmployee.Region),
	}
	return dto
}

func CreateToRegionEmployee(regionID uint, dto *employeesTypes.CreateEmployeeDTO) (*data.Employee, error) {
	if !data.IsAllowableRole(data.RegionEmployeeType, dto.Role) {
		return nil, fmt.Errorf("%w: %s for type %s", employeesTypes.ErrInvalidEmployeeRole, dto.Role, data.RegionEmployeeType)
	}

	employee, err := employeesTypes.CreateToEmployee(dto)
	if err != nil {
		return nil, err
	}
	employee.RegionEmployee = &data.RegionEmployee{
		RegionID: regionID,
		Role:     dto.Role,
	}

	return employee, nil
}
