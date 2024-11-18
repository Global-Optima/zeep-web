package types

type CreateEmployeeDTO struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Email    string `json:"email" binding:"required"`
	RoleID   uint   `json:"role_id" binding:"required"`
	StoreID  uint   `json:"store_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateEmployeeDTO struct {
	Name    *string `json:"name"`
	Phone   *string `json:"phone"`
	Email   *string `json:"email"`
	RoleID  *uint   `json:"role_id"`
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

type RoleDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
