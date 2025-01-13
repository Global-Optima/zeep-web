package types

type CreateIngredientCategoryDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type UpdateIngredientCategoryDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type IngredientCategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
