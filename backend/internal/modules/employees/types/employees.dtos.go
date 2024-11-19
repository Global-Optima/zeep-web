package types

type CreateEmployeeDTO struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Email    string `json:"email" binding:"required"`
	Role     string `json:"role" binding:"required"`
	StoreID  uint   `json:"store_id"`
	Password string `json:"password" binding:"required"`
}

type UpdateEmployeeDTO struct {
	Name    *string `json:"name"`
	Phone   *string `json:"phone"`
	Email   *string `json:"email"`
	Role    *string `json:"role"`
	StoreID *uint   `json:"store_id"`
}

type EmployeeDTO struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	StoreID uint   `json:"store_id"`
}

type UpdatePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RoleDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
