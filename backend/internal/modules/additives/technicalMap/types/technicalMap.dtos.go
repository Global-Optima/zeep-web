package types

import (
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
)

type AdditiveTechnicalMap struct {
	Ingredients []AdditiveIngredientDTO `json:"ingredients"`
}

type AdditiveIngredientDTO struct {
	Ingredients ingredientTypes.IngredientDTO `json:"ingredient"`
	Quantity    float64                       `json:"quantity"`
}
