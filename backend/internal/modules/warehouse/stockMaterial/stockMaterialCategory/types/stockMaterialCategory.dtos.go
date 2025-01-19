package types

type CreateStockMaterialCategoryDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateStockMaterialCategoryDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type StockMaterialCategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}
