package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/sirupsen/logrus"
)

func ConvertToStoreAdditiveDTO(storeAdditive *data.StoreAdditive) *StoreAdditiveDTO {
	return &StoreAdditiveDTO{
		AdditiveDTO:     *additiveTypes.ConvertToAdditiveDTO(&storeAdditive.Additive),
		StoreAdditiveID: storeAdditive.ID,
		StorePrice:      storeAdditive.Price,
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
				AdditiveCategoryItemDTO: *additiveTypes.ConvertToAdditiveCategoryItem(&additive, category.ID),
				StoreAdditiveID:         additive.StoreAdditives[0].ID,
				StorePrice:              additive.StoreAdditives[0].Price,
				IsDefault:               additive.ProductSizeAdditives[0].IsDefault,
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
