package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type CreateEmployeeDTO struct {
	Name             string               `json:"name" binding:"required"`
	Phone            string               `json:"phone"`
	Email            string               `json:"email" binding:"required"`
	Role             data.EmployeeRole    `json:"role" binding:"required"`
	Password         string               `json:"password" binding:"required"`
	Type             data.EmployeeType    `json:"type" binding:"required"`
	StoreDetails     *StoreDetailsDTO     `json:"storeDetails,omitempty"`
	WarehouseDetails *WarehouseDetailsDTO `json:"warehouseDetails,omitempty"`
}

type StoreDetailsDTO struct {
	StoreID     uint `json:"storeId" binding:"required"`
	IsFranchise bool `json:"isFranchise,omitempty"`
}

type WarehouseDetailsDTO struct {
	WarehouseID uint `json:"warehouseId" binding:"required"`
}

type UpdateEmployeeDTO struct {
	Name             *string                    `json:"name,omitempty"`
	Phone            *string                    `json:"phone,omitempty"`
	Email            *string                    `json:"email,omitempty"`
	Role             *data.EmployeeRole         `json:"role,omitempty"`
	StoreDetails     *UpdateStoreDetailsDTO     `json:"storeDetails,omitempty"`
	WarehouseDetails *UpdateWarehouseDetailsDTO `json:"warehouseDetails,omitempty"`
}

type UpdateStoreDetailsDTO struct {
	StoreID     *uint `json:"storeId,omitempty"`
	IsFranchise *bool `json:"isFranchise,omitempty"`
}

type UpdateWarehouseDetailsDTO struct {
	WarehouseID *uint `json:"warehouseId,omitempty"`
}

type EmployeeDTO struct {
	ID               uint                         `json:"id"`
	Name             string                       `json:"name"`
	Phone            string                       `json:"phone"`
	Email            string                       `json:"email"`
	Role             data.EmployeeRole            `json:"role"`
	IsActive         bool                         `json:"isActive"`
	Type             data.EmployeeType            `json:"type"`
	StoreDetails     *StoreEmployeeDetailsDTO     `json:"storeDetails,omitempty"`
	WarehouseDetails *WarehouseEmployeeDetailsDTO `json:"warehouseDetails,omitempty"`
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

type GetEmployeesFilter struct {
	Type        *string `form:"type,omitempty"`
	Role        *string `form:"role,omitempty"`
	StoreID     *uint   `form:"storeId,omitempty"`
	WarehouseID *uint   `form:"warehouseId,omitempty"`
	Pagination  *utils.Pagination
}
