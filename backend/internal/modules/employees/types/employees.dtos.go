package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type CreateEmployeeDTO struct {
	FirstName string            `json:"firstName" binding:"required"`
	LastName  string            `json:"lastName" binding:"required"`
	Phone     string            `json:"phone"`
	Email     string            `json:"email" binding:"required"`
	Role      data.EmployeeRole `json:"role" binding:"required"`
	Password  string            `json:"password" binding:"required"`
	IsActive  bool              `json:"isActive" binding:"required"`
}

type CreateStoreEmployeeDTO struct {
	CreateEmployeeDTO
	StoreID     uint `json:"storeId" binding:"required"`
	IsFranchise bool `json:"isFranchise,omitempty"`
}

type CreateWarehouseEmployeeDTO struct {
	CreateEmployeeDTO
	WarehouseID uint `json:"warehouseId" binding:"required"`
}

type StoreDetailsDTO struct {
	StoreID     uint `json:"storeId" binding:"required"`
	IsFranchise bool `json:"isFranchise,omitempty"`
}

type WarehouseDetailsDTO struct {
	WarehouseID uint `json:"warehouseId" binding:"required"`
}

type UpdateEmployeeDTO struct {
	FirstName *string            `json:"firstName,omitempty"`
	LastName  *string            `json:"lastName,omitempty"`
	Phone     *string            `json:"phone,omitempty"`
	Email     *string            `json:"email,omitempty"`
	Role      *data.EmployeeRole `json:"role,omitempty"`
	IsActive  *bool              `json:"isActive"`
}

type UpdateStoreEmployeeDTO struct {
	UpdateEmployeeDTO
	StoreID     *uint `json:"storeId,omitempty"`
	IsFranchise *bool `json:"isFranchise,omitempty"`
}

type UpdateWarehouseEmployeeDTO struct {
	UpdateEmployeeDTO
	WarehouseID *uint `json:"warehouseId,omitempty"`
}

type EmployeeDTO struct {
	ID        uint              `json:"id"`
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	Phone     string            `json:"phone"`
	Email     string            `json:"email"`
	Role      data.EmployeeRole `json:"role"`
	IsActive  bool              `json:"isActive"`
	Type      data.EmployeeType `json:"type"`
}

type StoreEmployeeDTO struct {
	EmployeeDTO
	StoreID     uint `json:"storeId"`
	IsFranchise bool `json:"isFranchise"`
}

type WarehouseEmployeeDTO struct {
	EmployeeDTO
	WarehouseID uint `json:"warehouseId"`
}

type StoreEmployeeDetailsDTO struct {
	StoreID     uint `json:"storeId"`
	IsFranchise bool `json:"isFranchise"`
}

type WarehouseEmployeeDetailsDTO struct {
	WarehouseID uint `json:"warehouseId"`
}

type UpdatePasswordDTO struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RoleDTO struct {
	Name string `json:"name"`
}

type GetStoreEmployeesFilter struct {
	utils.BaseFilter
	StoreID  uint    `form:"storeId"`
	Role     *string `form:"role,omitempty"`
	IsActive *bool   `form:"isActive,omitempty"`
	Search   *string `form:"search,omitempty"`
}

type GetWarehouseEmployeesFilter struct {
	utils.BaseFilter
	WarehouseID *uint   `form:"warehouseId"`
	Role        *string `form:"role,omitempty"`
	IsActive    *bool   `form:"isActive,omitempty"`
	Search      *string `form:"search,omitempty"`
}

type CreateEmployeeWorkdayDTO struct {
	Day        string `json:"day" binding:"required"`
	StartAt    string `json:"startAt"`
	EndAt      string `json:"endAt"`
	EmployeeID uint   `json:"employeeId"`
}

type EmployeeWorkdayDTO struct {
	ID         uint   `json:"id"`
	Day        string `json:"day"`
	StartAt    string `json:"startAt"`
	EndAt      string `json:"endAt"`
	EmployeeID uint   `json:"employeeId"`
}

type UpdateEmployeeWorkdayDTO struct {
	StartAt *string `json:"startAt,omitempty"`
	EndAt   *string `json:"endAt,omitempty"`
}
