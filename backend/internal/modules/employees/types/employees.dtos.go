package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type CreateEmployeeDTO struct {
	FirstName string             `json:"firstName" binding:"required"`
	LastName  string             `json:"lastName" binding:"required"`
	Phone     string             `json:"phone"`
	Email     string             `json:"email" binding:"required"`
	Role      data.EmployeeRole  `json:"role" binding:"required"`
	Password  string             `json:"password" binding:"required"`
	IsActive  bool               `json:"isActive" binding:"required"`
	Workdays  []CreateWorkdayDTO `json:"workdays" binding:"dive"`
}

type CreateStoreEmployeeDTO struct {
	CreateEmployeeDTO
	StoreID uint `json:"storeId" binding:"required"`
}

type CreateWarehouseEmployeeDTO struct {
	CreateEmployeeDTO
	WarehouseID uint `json:"warehouseId" binding:"required"`
}

type CreateFranchiseeEmployeeDTO struct {
	CreateEmployeeDTO
	FranchiseeID uint `json:"franchiseeId" binding:"required"`
}

type CreateRegionManagerDTO struct {
	CreateEmployeeDTO
	RegionID uint `json:"regionId" binding:"required"`
}

type CreateWorkdayDTO struct {
	Day     string `json:"day" binding:"required"`
	StartAt string `json:"startAt" binding:"required"`
	EndAt   string `json:"endAt" binding:"required"`
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
	StoreID *uint `json:"storeId,omitempty"`
}

type UpdateWarehouseEmployeeDTO struct {
	UpdateEmployeeDTO
	WarehouseID *uint `json:"warehouseId,omitempty"`
}

type UpdateFranchiseeEmployeeDTO struct {
	UpdateEmployeeDTO
	FranchiseeID *uint `json:"franchiseeId,omitempty"`
}

type UpdateRegionManagerEmployeeDTO struct {
	UpdateEmployeeDTO
	RegionID *uint `json:"regionId,omitempty"`
}

type EmployeeDTO struct {
	ID        uint              `json:"id"`
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	Phone     string            `json:"phone"`
	Email     string            `json:"email"`
	Type      data.EmployeeType `json:"type"`
	Role      data.EmployeeRole `json:"role"`
	IsActive  bool              `json:"isActive"`
}

type StoreEmployeeDTO struct {
	EmployeeDTO
	StoreID uint `json:"storeId"`
}

type WarehouseEmployeeDTO struct {
	EmployeeDTO
	WarehouseID uint `json:"warehouseId"`
}

type FranchiseeEmployeeDTO struct {
	EmployeeDTO
	FranchiseeID uint `json:"franchiseeId"`
}

type RegionManagerDTO struct {
	EmployeeDTO
	RegionID uint `json:"regionId"`
}

type AdminDTO struct {
	EmployeeDTO
	Roles uint `json:"regionId"`
}

type EmployeeAccountDTO struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type UpdatePasswordDTO struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type RoleDTO struct {
	Name string `json:"name"`
}

type EmployeesFilter struct {
	utils.BaseFilter
	Role     *string `form:"role,omitempty"`
	IsActive *bool   `form:"isActive,omitempty"`
	Search   *string `form:"search,omitempty"`
}

type CreateEmployeeWorkdayDTO struct {
	CreateWorkdayDTO
	EmployeeID uint `json:"employeeId"`
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
