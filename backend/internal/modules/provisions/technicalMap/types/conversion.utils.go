package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
)

func ConvertToProvisionTechnicalMapDTO(ProvisionIngredients []data.ProvisionIngredient) ProvisionTechnicalMap {
	technicalMap := ProvisionTechnicalMap{
		Ingredients: make([]ProvisionIngredientDTO, len(ProvisionIngredients)),
	}

	for i, ProvisionIngredient := range ProvisionIngredients {
		technicalMap.Ingredients[i] = ProvisionIngredientDTO{
			Ingredients: *ingredientTypes.ConvertToIngredientResponseDTO(&ProvisionIngredient.Ingredient),
			Quantity:    ProvisionIngredient.Quantity,
		}
	}

	return technicalMap
}
