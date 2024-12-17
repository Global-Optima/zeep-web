package types

type CustomerRegisterDTO struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CustomerLoginDTO struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}
