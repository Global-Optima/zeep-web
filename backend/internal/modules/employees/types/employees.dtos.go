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
	StoreID  uint    `form:"storeId" binding:"required"`
	Role     *string `form:"role,omitempty"`
	IsActive *bool   `form:"isActive,omitempty"`
	Search   *string `form:"search,omitempty"`
	utils.BaseFilter
}

type GetWarehouseEmployeesFilter struct {
	WarehouseID uint    `form:"warehouseId" binding:"required"`
	Role        *string `form:"role,omitempty"`
	IsActive    *bool   `form:"isActive,omitempty"`
	Search      *string `form:"search,omitempty"`
	utils.BaseFilter
}
