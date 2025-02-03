package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	warehouseTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
)

type UpdateWarehouseEmployeeDTO struct {
	employeesTypes.UpdateEmployeeDTO
	Role        *data.EmployeeRole `json:"role"`
	WarehouseID *uint              `json:"warehouseId"`
}

type WarehouseEmployeeDTO struct {
	ID uint `json:"id"`
	employeesTypes.BaseEmployeeDTO
	EmployeeID uint `json:"employeeId"`
}

type WarehouseEmployeeDetailsDTO struct {
	ID uint `json:"id"`
	employeesTypes.BaseEmployeeDetailsDTO
	EmployeeID uint                             `json:"employeeId"`
	Warehouse  warehouseTypes.WarehouseResponse `json:"warehouse"`
}
