package types

type CreateEmployeeDTO struct {
	Name     string       `json:"name" binding:"required"`
	Phone    string       `json:"phone"`
	Email    string       `json:"email" binding:"required"`
	Role     EmployeeRole `json:"role" binding:"required"`
	StoreID  uint         `json:"storeId"`
	Password string       `json:"password" binding:"required"`
}

type UpdateEmployeeDTO struct {
	Name    *string       `json:"name"`
	Phone   *string       `json:"phone"`
	Email   *string       `json:"email"`
	Role    *EmployeeRole `json:"role"`
	StoreID *uint         `json:"store_id"`
}

type EmployeeDTO struct {
	ID      uint         `json:"id"`
	Name    string       `json:"name"`
	Phone   string       `json:"phone"`
	Email   string       `json:"email"`
	Role    EmployeeRole `json:"role"`
	StoreID uint         `json:"storeId"`
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
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
