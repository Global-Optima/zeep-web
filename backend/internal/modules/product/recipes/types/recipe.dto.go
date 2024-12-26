package types

type RecipeStepDTO struct {
	ID          uint   `json:"id"`
	ProductID   uint   `json:"productId"`
	Step        int    `json:"stepNumber"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
}

type CreateOrReplaceRecipeStepDTO struct {
	Step        int    `json:"step" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	ImageURL    string `json:"imageUrl,omitempty"`
}
