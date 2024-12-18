package types

type CreateIngredientDTO struct {
	Name      string  `json:"name" binding:"required"`
	Calories  float64 `json:"calories" binding:"required,gte=0"`
	Fat       float64 `json:"fat" binding:"required,gte=0"`
	Carbs     float64 `json:"carbs" binding:"required,gte=0"`
	Proteins  float64 `json:"proteins" binding:"required,gte=0"`
	ExpiresAt *string `json:"expiresAt,omitempty"` // ISO-8601 formatted string (e.g., "2024-12-18T00:00:00Z")
}

type UpdateIngredientDTO struct {
	ID        uint     `json:"id" binding:"required"`
	Name      string   `json:"name,omitempty"`
	Calories  *float64 `json:"calories,omitempty"` // Nullable fields for partial updates
	Fat       *float64 `json:"fat,omitempty"`
	Carbs     *float64 `json:"carbs,omitempty"`
	Proteins  *float64 `json:"proteins,omitempty"`
	ExpiresAt *string  `json:"expiresAt,omitempty"` // ISO-8601 formatted string
}

type IngredientResponseDTO struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Calories  float64 `json:"calories"`
	Fat       float64 `json:"fat"`
	Carbs     float64 `json:"carbs"`
	Proteins  float64 `json:"proteins"`
	ExpiresAt *string `json:"expiresAt,omitempty"`
}

type IngredientFilter struct {
	Name        *string  `form:"name"`
	MinCalories *float64 `form:"minCalories"`
	MaxCalories *float64 `form:"maxCalories"`
}
