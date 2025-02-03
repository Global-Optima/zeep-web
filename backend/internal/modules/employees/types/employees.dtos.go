package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type CreateEmployeeDTO struct {
	FirstName string                      `json:"firstName" binding:"required"`
	LastName  string                      `json:"lastName" binding:"required"`
	Phone     string                      `json:"phone"`
	Email     string                      `json:"email" binding:"required"`
	Role      data.EmployeeRole           `json:"role" binding:"required"`
	Password  string                      `json:"password" binding:"required"`
	IsActive  bool                        `json:"isActive" binding:"required"`
	Workdays  []CreateOrReplaceWorkdayDTO `json:"workdays" binding:"dive"`
}

type CreateOrReplaceWorkdayDTO struct {
	Day     string `json:"day" binding:"required"`
	StartAt string `json:"startAt" binding:"required"`
	EndAt   string `json:"endAt" binding:"required"`
}

type UpdateEmployeeDTO struct {
	FirstName *string                     `json:"firstName,omitempty"`
	LastName  *string                     `json:"lastName,omitempty"`
	IsActive  *bool                       `json:"isActive"`
	Workdays  []CreateOrReplaceWorkdayDTO `json:"workdays" binding:"dive"`
}

type ReassignEmployeeTypeDTO struct {
	EmployeeType data.EmployeeType `json:"employeeType" binding:"required"`
	Role         data.EmployeeRole `json:"role" binding:"required"`
	WorkplaceID  uint              `json:"workplaceId" binding:"required"`
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

type EmployeeWorkdayDTO struct {
	ID         uint   `json:"id"`
	Day        string `json:"day"`
	StartAt    string `json:"startAt"`
	EndAt      string `json:"endAt"`
	EmployeeID uint   `json:"employeeId"`
}
