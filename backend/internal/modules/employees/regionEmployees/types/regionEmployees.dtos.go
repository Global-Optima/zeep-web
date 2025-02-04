package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	regionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/regions/types"
)

type UpdateRegionEmployeeDTO struct {
	employeesTypes.UpdateEmployeeDTO
	Role     *data.EmployeeRole `json:"role,omitempty"`
	RegionID *uint              `json:"regionId,omitempty"`
}

type RegionEmployeeDTO struct {
	ID uint `json:"id"`
	employeesTypes.BaseEmployeeDTO
	EmployeeID uint `json:"employeeId"`
}

type RegionEmployeeDetailsDTO struct {
	ID uint `json:"id"`
	employeesTypes.BaseEmployeeDetailsDTO
	EmployeeID uint                   `json:"employeeId"`
	Region     regionsTypes.RegionDTO `json:"region"`
}
