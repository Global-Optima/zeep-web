package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

type CreateIngredientDTO struct {
	Name       string  `json:"name" binding:"required"`
	Calories   float64 `json:"calories" binding:"gte=0"`
	Fat        float64 `json:"fat" binding:"gte=0"`
	Carbs      float64 `json:"carbs" binding:"gte=0"`
	Proteins   float64 `json:"proteins" binding:"gte=0"`
	CategoryID uint    `json:"categoryId" binding:"required,gt=0"`
	UnitID     uint    `json:"unitId" binding:"required,gt=0"`
	ExpiresAt  *string `json:"expiresAt,omitempty"` // ISO-8601 formatted string (e.g., "2024-12-18T00:00:00Z")
}

type UpdateIngredientDTO struct {
	Name       string   `json:"name,omitempty"`
	Calories   *float64 `json:"calories" binding:"omitempty,gte=0"` // Nullable fields for partial updates
	Fat        *float64 `json:"fat" binding:"omitempty,gte=0"`
	Carbs      *float64 `json:"carbs" binding:"omitempty,gte=0"`
	Proteins   *float64 `json:"proteins" binding:"omitempty,gte=0"`
	UnitID     *uint    `json:"unitId" binding:"omitempty,gt=0"`
	CategoryID *uint    `json:"categoryId" binding:"omitempty,gt=0"`
	ExpiresAt  *string  `json:"expiresAt,omitempty"` // ISO-8601 formatted string
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

type IngredientDetailsDTO struct {
	IngredientResponseDTO
	Category struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"category"`
}

type IngredientFilter struct {
	utils.BaseFilter
	ProductSizeID *uint    `form:"productSizeId" binding:"omitempty,gt=0"`
	Name          *string  `form:"name"`
	MinCalories   *float64 `form:"minCalories"`
	MaxCalories   *float64 `form:"maxCalories"`
}
