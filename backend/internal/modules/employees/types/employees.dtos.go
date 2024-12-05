package types

type CreateEmployeeDTO struct {
	Name             string               `json:"name" binding:"required"`
	Phone            string               `json:"phone"`
	Email            string               `json:"email" binding:"required"`
	Role             string               `json:"role" binding:"required"`
	Password         string               `json:"password" binding:"required"`
	Type             EmployeeType         `json:"type" binding:"required"`
	StoreDetails     *StoreDetailsDTO     `json:"store_details,omitempty"`
	WarehouseDetails *WarehouseDetailsDTO `json:"warehouse_details,omitempty"`
}

type StoreDetailsDTO struct {
	StoreID     uint `json:"store_id" binding:"required"`
	IsFranchise bool `json:"is_franchise,omitempty"`
}

type WarehouseDetailsDTO struct {
	WarehouseID uint `json:"warehouse_id" binding:"required"`
}

type UpdateEmployeeDTO struct {
	Name             *string                    `json:"name,omitempty"`
	Phone            *string                    `json:"phone,omitempty"`
	Email            *string                    `json:"email,omitempty"`
	Role             *string                    `json:"role,omitempty"`
	StoreDetails     *UpdateStoreDetailsDTO     `json:"store_details,omitempty"`
	WarehouseDetails *UpdateWarehouseDetailsDTO `json:"warehouse_details,omitempty"`
}

type UpdateStoreDetailsDTO struct {
	StoreID     *uint `json:"store_id,omitempty"`
	IsFranchise *bool `json:"is_franchise,omitempty"`
}

type UpdateWarehouseDetailsDTO struct {
	WarehouseID *uint `json:"warehouse_id,omitempty"`
}

type EmployeeDTO struct {
	ID               uint                         `json:"id"`
	Name             string                       `json:"name"`
	Phone            string                       `json:"phone"`
	Email            string                       `json:"email"`
	Role             string                       `json:"role"`
	IsActive         bool                         `json:"is_active"`
	Type             EmployeeType                 `json:"type"`
	StoreDetails     *StoreEmployeeDetailsDTO     `json:"store_details,omitempty"`
	WarehouseDetails *WarehouseEmployeeDetailsDTO `json:"warehouse_details,omitempty"`
}

type StoreEmployeeDetailsDTO struct {
	StoreID     uint `json:"store_id"`
	IsFranchise bool `json:"is_franchise"`
}

type WarehouseEmployeeDetailsDTO struct {
	WarehouseID uint `json:"warehouse_id"`
}

type UpdatePasswordDTO struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type LoginDTO struct {
	EmployeeId uint   `json:"employeeId" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type RoleDTO struct {
	Name string `json:"name"`
}

type GetEmployeesQuery struct {
	Type        *string `json:"type,omitempty"`
	Role        *string `json:"role,omitempty"`
	StoreID     *uint   `json:"store_id,omitempty"`
	WarehouseID *uint   `json:"warehouse_id,omitempty"`
	Limit       int     `json:"limit,omitempty"`
	Offset      int     `json:"offset,omitempty"`
}
