package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
)

func ConvertToStoreAdditiveDTO(storeAdditive *data.StoreAdditive) *StoreAdditiveDTO {
	return &StoreAdditiveDTO{
		AdditiveDTO:     *additiveTypes.ConvertToAdditiveDTO(&storeAdditive.Additive),
		StoreAdditiveID: storeAdditive.ID,
		StorePrice:      storeAdditive.Price,
	}
}

func ConvertToStoreAdditiveCategoryDTO(category *data.AdditiveCategory) *StoreAdditiveCategoryDTO {
	storeAdditives := ConvertToStoreAdditiveCategoryItemDTOs(category)

	return &StoreAdditiveCategoryDTO{
		ID:               category.ID,
		Name:             category.Name,
		Description:      category.Description,
		IsMultipleSelect: category.IsMultipleSelect,
		Additives:        storeAdditives, // Always initialized as a slice
	}
}

func ConvertToStoreAdditiveCategoryItemDTOs(category *data.AdditiveCategory) []StoreAdditiveCategoryItemDTO {
	var storeAdditives []StoreAdditiveCategoryItemDTO

	// Populate additives if present
	for _, additive := range category.Additives {
		if len(additive.StoreAdditives) > 0 {
			storeAdditives = append(storeAdditives, StoreAdditiveCategoryItemDTO{
				AdditiveCategoryItemDTO: *additiveTypes.ConvertToAdditiveCategoryItem(&additive, category.ID),
				StoreAdditiveID:         additive.StoreAdditives[0].ID,
				StorePrice:              additive.StoreAdditives[0].Price,
			})
		}
	}

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
