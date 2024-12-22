package types

type RecipeStepDTO struct {
	ID          uint   `json:"id"`
	ProductID   uint   `json:"product_id"`
	Step        int    `json:"stepNumber"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
}

type CreateRecipeStepDTO struct {
	ProductID   uint   `json:"product_id" binding:"required"`
	Step        int    `json:"step" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	ImageURL    string `json:"image_url,omitempty"`
}

type UpdateRecipeStepDTO struct {
	Step        int    `json:"step" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	ImageURL    string `json:"image_url,omitempty"`
}
