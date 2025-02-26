package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
)

func ConvertToStoreAdditiveDTO(storeAdditive *data.StoreAdditive) *StoreAdditiveDTO {
	return &StoreAdditiveDTO{
		ID:              storeAdditive.ID,
		BaseAdditiveDTO: *additiveTypes.ConvertToBaseAdditiveDTO(&storeAdditive.Additive),
		AdditiveID:      storeAdditive.AdditiveID,
		StorePrice:      getStorePrice(storeAdditive),
	}
}

func getStorePrice(storeAdditive *data.StoreAdditive) float64 {
	if storeAdditive.StorePrice != nil {
		return *storeAdditive.StorePrice
	}
	return storeAdditive.Additive.BasePrice
}

func ConvertToStoreAdditiveDetailsDTO(storeAdditive *data.StoreAdditive) *StoreAdditiveDetailsDTO {
	ingredients := make([]additiveTypes.AdditiveIngredientDTO, len(storeAdditive.Additive.Ingredients))
	for i, additiveIngredient := range storeAdditive.Additive.Ingredients {
		ingredients[i].Ingredient = *ingredientTypes.ConvertToIngredientResponseDTO(&additiveIngredient.Ingredient)
		ingredients[i].Quantity = additiveIngredient.Quantity
	}

	return &StoreAdditiveDetailsDTO{
		StoreAdditiveDTO: *ConvertToStoreAdditiveDTO(storeAdditive),
		Ingredients:      ingredients,
	}
}

func ConvertToStoreAdditiveCategoryDTO(category *data.AdditiveCategory) *StoreAdditiveCategoryDTO {
	return &StoreAdditiveCategoryDTO{
		ID:               category.ID,
		Name:             category.Name,
		Description:      category.Description,
		IsMultipleSelect: category.IsMultipleSelect,
		Additives:        ConvertToStoreAdditiveCategoryItemDTOs(category), // Always initialized as a slice
	}
}

func ConvertToStoreAdditiveCategoryItemDTOs(category *data.AdditiveCategory) []StoreAdditiveCategoryItemDTO {
	var storeAdditives []StoreAdditiveCategoryItemDTO
	// Populate additives if present
	for _, additive := range category.Additives {
		if len(additive.StoreAdditives) > 0 && len(additive.ProductSizeAdditives) > 0 {
			storeAdditives = append(storeAdditives, StoreAdditiveCategoryItemDTO{
				ID:                          additive.StoreAdditives[0].ID,
				BaseAdditiveCategoryItemDTO: *additiveTypes.ConvertToBaseAdditiveCategoryItem(&additive, category.ID),
				AdditiveID:                  additive.StoreAdditives[0].AdditiveID,
				StorePrice:                  getStorePrice(&additive.StoreAdditives[0]),
			})
		}
	}

	return storeAdditives
}

func CreateToStoreAdditive(dto *CreateStoreAdditiveDTO, storeID uint) *data.StoreAdditive {
	return &data.StoreAdditive{
		StoreID:    storeID,
		AdditiveID: dto.AdditiveID,
		StorePrice: dto.StorePrice,
	}
}

func UpdateToStoreAdditive(dto *UpdateStoreAdditiveDTO) *data.StoreAdditive {
	return &data.StoreAdditive{
		StorePrice: dto.StorePrice,
	}
}
