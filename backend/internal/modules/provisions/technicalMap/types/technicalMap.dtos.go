package types

import (
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
)

type ProvisionTechnicalMap struct {
	Ingredients []ProvisionIngredientDTO `json:"ingredients"`
}

type ProvisionIngredientDTO struct {
	Ingredients ingredientTypes.IngredientDTO `json:"ingredient"`
	Quantity    float64                       `json:"quantity"`
}
