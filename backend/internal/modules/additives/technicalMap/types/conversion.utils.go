package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
)

func ConvertToAdditiveTechnicalMapDTO(AdditiveIngredients []data.AdditiveIngredient) AdditiveTechnicalMap {
	technicalMap := AdditiveTechnicalMap{
		Ingredients: make([]AdditiveIngredientDTO, len(AdditiveIngredients)),
	}

	for i, AdditiveIngredient := range AdditiveIngredients {
		technicalMap.Ingredients[i] = AdditiveIngredientDTO{
			Ingredients: *ingredientTypes.ConvertToIngredientResponseDTO(&AdditiveIngredient.Ingredient),
			Quantity:    AdditiveIngredient.Quantity,
		}
	}

	return technicalMap
}
