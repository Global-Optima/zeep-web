package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/sirupsen/logrus"
)

func ConvertToStoreAdditiveDTO(storeAdditive *data.StoreAdditive) *StoreAdditiveDTO {
	storePrice := storeAdditive.Price
	if storePrice == 0 {
		storePrice = storeAdditive.Additive.BasePrice
	}

	return &StoreAdditiveDTO{
		ID:              storeAdditive.ID,
		BaseAdditiveDTO: *additiveTypes.ConvertToBaseAdditiveDTO(&storeAdditive.Additive),
		AdditiveID:      storeAdditive.AdditiveID,
		StorePrice:      storePrice,
	}
}

func ConvertToStoreAdditiveDetailsDTO(storeAdditive *data.StoreAdditive) *StoreAdditiveDetailsDTO {
	ingredients := make([]ingredientTypes.IngredientDTO, len(storeAdditive.Additive.Ingredients))
	for i, additiveIngredient := range storeAdditive.Additive.Ingredients {
		ingredients[i] = *ingredientTypes.ConvertToIngredientResponseDTO(&additiveIngredient.Ingredient)
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
				StorePrice:                  additive.StoreAdditives[0].Price,
				IsDefault:                   additive.ProductSizeAdditives[0].IsDefault,
			})
		}
	}

	logrus.Info(storeAdditives)

	return storeAdditives
}

func CreateToStoreAdditive(dto *CreateStoreAdditiveDTO, storeID uint) *data.StoreAdditive {
	return &data.StoreAdditive{
		StoreID:    storeID,
		AdditiveID: dto.AdditiveID,
		Price:      dto.StorePrice,
	}
}

func UpdateToStoreAdditive(dto *UpdateStoreAdditiveDTO) *data.StoreAdditive {
	return &data.StoreAdditive{
		Price: dto.StorePrice,
	}
}
