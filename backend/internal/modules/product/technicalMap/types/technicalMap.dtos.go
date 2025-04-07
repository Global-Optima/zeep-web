package types

import (
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
)

type ProductSizeTechnicalMap struct {
	Ingredients []ProductSizeIngredientDTO `json:"ingredients"`
}

type ProductSizeIngredientDTO struct {
	Ingredients ingredientTypes.IngredientDTO `json:"ingredient"`
	Quantity    float64                       `json:"quantity"`
}
