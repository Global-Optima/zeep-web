package types

type CustomerDTO struct {
	ID    uint   `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

type CustomerAdminDTO struct {
	CustomerDTO
	IsVerified bool `json:"isVerified" binding:"required"`
	IsBanned   bool `json:"isBanned" binding:"required"`
}
