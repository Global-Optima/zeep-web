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

type CreateRegionEmployeeDTO struct {
	CreateEmployeeDTO
	RegionID uint `json:"regionId" binding:"required"`
}

type CreateAdminEmployeeDTO struct {
	CreateEmployeeDTO
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
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Email     *string `json:"email,omitempty"`
	IsActive  *bool   `json:"isActive"`
}

type UpdateTypedEmployeeFields struct {
	Role *data.EmployeeRole `json:"role,omitempty"`
}

type UpdateStoreEmployeeDTO struct {
	UpdateTypedEmployeeFields
	StoreID *uint `json:"storeId,omitempty"`
}

type UpdateWarehouseEmployeeDTO struct {
	UpdateTypedEmployeeFields
	WarehouseID *uint `json:"warehouseId,omitempty"`
}

type UpdateFranchiseeEmployeeDTO struct {
	UpdateTypedEmployeeFields
	FranchiseeID *uint `json:"franchiseeId,omitempty"`
}

type UpdateRegionEmployeeDTO struct {
	UpdateTypedEmployeeFields
	RegionID *uint `json:"regionId,omitempty"`
}

type UpdateAdminEmployeeDTO struct {
	UpdateTypedEmployeeFields
}

type BaseEmployeeDTO struct {
	ID        uint              `json:"id"`
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	Phone     string            `json:"phone"`
	Email     string            `json:"email"`
	Type      data.EmployeeType `json:"type"`
	IsActive  bool              `json:"isActive"`
}

type EmployeeDTO struct {
	BaseEmployeeDTO
	Role data.EmployeeRole `json:"role"`
}

type StoreEmployeeDTO struct {
	BaseEmployeeDTO
	StoreID uint              `json:"storeId"`
	Role    data.EmployeeRole `json:"role"`
}

type WarehouseEmployeeDTO struct {
	BaseEmployeeDTO
	WarehouseID uint              `json:"warehouseId"`
	Role        data.EmployeeRole `json:"role"`
}

type FranchiseeEmployeeDTO struct {
	BaseEmployeeDTO
	FranchiseeID uint              `json:"franchiseeId"`
	Role         data.EmployeeRole `json:"role"`
}

type RegionEmployeeDTO struct {
	BaseEmployeeDTO
	RegionID uint              `json:"regionId"`
	Role     data.EmployeeRole `json:"role"`
}

type AdminEmployeeDTO struct {
	BaseEmployeeDTO
	Role data.EmployeeRole `json:"role"`
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

type EmployeeTypeRoles struct {
	EmployeeType data.EmployeeType   `json:"employeeType" binding:"required"`
	Roles        []data.EmployeeRole `json:"roles" binding:"required"`
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
