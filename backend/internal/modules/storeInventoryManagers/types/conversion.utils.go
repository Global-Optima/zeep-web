package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	provisionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"
)

func ConvertToFrozenInventoryDTO(
	frozenInventory *FrozenInventory,
	ingredients []data.Ingredient,
	provisions []data.Provision,
) *FrozenInventoryDTO {
	ingByID := make(map[uint]*data.Ingredient, len(ingredients))
	for i := range ingredients {
		ing := &ingredients[i]
		ingByID[ing.ID] = ing
	}
	provByID := make(map[uint]*data.Provision, len(provisions))
	for i := range provisions {
		p := &provisions[i]
		provByID[p.ID] = p
	}

	frozenIngredientsDTOs := make([]FrozenIngredientDTO, 0, len(frozenInventory.Ingredients))
	frozenProvisionsDTOs := make([]FrozenProvisionDTO, 0, len(frozenInventory.Provisions))

	for id, qty := range frozenInventory.Ingredients {
		if ing, ok := ingByID[id]; ok {
			frozenIngredientsDTOs = append(frozenIngredientsDTOs, FrozenIngredientDTO{
				IngredientDTO: *ingredientTypes.ConvertToIngredientResponseDTO(ing),
				Quantity:      qty,
			})
		}
	}

	for id, vol := range frozenInventory.Provisions {
		if prov, ok := provByID[id]; ok {
			frozenProvisionsDTOs = append(frozenProvisionsDTOs, FrozenProvisionDTO{
				ProvisionDTO: *provisionsTypes.MapToProvisionDTO(prov),
				Volume:       vol,
			})
		}
	}

	return &FrozenInventoryDTO{
		Ingredients: frozenIngredientsDTOs,
		Provisions:  frozenProvisionsDTOs,
	}
}
