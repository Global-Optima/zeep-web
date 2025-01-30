package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	storeTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
)

type UpdateStoreEmployeeDTO struct {
	employeesTypes.UpdateEmployeeDTO
	Role    *data.EmployeeRole `json:"role,omitempty"`
	StoreID *uint              `json:"storeId,omitempty"`
}

type StoreEmployeeDTO struct {
	employeesTypes.EmployeeDTO
	Store storeTypes.StoreDTO
}
