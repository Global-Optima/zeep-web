package types

type CreateEmployeeDTO struct {
	Name             string               `json:"name" binding:"required"`
	Phone            string               `json:"phone"`
	Email            string               `json:"email" binding:"required"`
	Role             string               `json:"role" binding:"required"`
	Password         string               `json:"password" binding:"required"`
	Type             EmployeeType         `json:"type" binding:"required"`
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
	Role             *string                    `json:"role,omitempty"`
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
	Role             string                       `json:"role"`
	IsActive         bool                         `json:"isActive"`
	Type             EmployeeType                 `json:"type"`
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

type GetEmployeesQuery struct {
	Type        *string `json:"type,omitempty"`
	Role        *string `json:"role,omitempty"`
	StoreID     *uint   `json:"storeId,omitempty"`
	WarehouseID *uint   `json:"warehouseId,omitempty"`
	Limit       int     `json:"limit,omitempty"`
	Offset      int     `json:"offset,omitempty"`
}
