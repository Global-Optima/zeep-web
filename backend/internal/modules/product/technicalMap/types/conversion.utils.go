package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
)

func ConvertToProductSizeTechnicalMapDTO(productSizeIngredients []data.ProductSizeIngredient) ProductSizeTechnicalMap {
	technicalMap := ProductSizeTechnicalMap{
		Ingredients: make([]ProductSizeIngredientDTO, len(productSizeIngredients)),
	}

	for i, productSizeIngredient := range productSizeIngredients {
		technicalMap.Ingredients[i] = ProductSizeIngredientDTO{
			Ingredients: *ingredientTypes.ConvertToIngredientResponseDTO(&productSizeIngredient.Ingredient),
			Quantity:    productSizeIngredient.Quantity,
		}
	}

	return technicalMap
}
