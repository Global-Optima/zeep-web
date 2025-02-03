package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	regionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/regions/types"
)

func MapToRegionEmployeeDTO(regionEmployee *data.RegionEmployee) *RegionEmployeeDTO {
	regionEmployee.Employee.RegionEmployee = regionEmployee
	return &RegionEmployeeDTO{
		ID:              regionEmployee.ID,
		BaseEmployeeDTO: *employeesTypes.MapToBaseEmployeeDTO(&regionEmployee.Employee),
		EmployeeID:      regionEmployee.Employee.ID,
	}
}

func MapToRegionEmployeeDetailsDTO(regionEmployee *data.RegionEmployee) *RegionEmployeeDetailsDTO {
	regionEmployee.Employee.RegionEmployee = regionEmployee
	dto := &RegionEmployeeDetailsDTO{
		ID:                     regionEmployee.ID,
		BaseEmployeeDetailsDTO: *employeesTypes.MapToBaseEmployeeDetailsDTO(&regionEmployee.Employee),
		EmployeeID:             regionEmployee.EmployeeID,
		Region:                 *regionsTypes.MapRegionToDTO(&regionEmployee.Region),
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
