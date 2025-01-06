package types

import additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"

type CreateStoreAdditiveDTO struct {
	AdditiveID uint    `json:"additiveId" binding:"required,gt=0"`
	StorePrice float64 `json:"storePrice" binding:"required,gte=0"`
}

type UpdateStoreAdditiveDTO struct {
	StorePrice float64 `json:"storePrice" binding:"required,gte=0"`
}

type StoreAdditiveDTO struct {
	additiveTypes.AdditiveDTO
	StoreAdditiveID uint    `json:"storeAdditiveId"`
	StorePrice      float64 `json:"storePrice"`
}

type StoreAdditiveCategoryItemDTO struct {
	additiveTypes.AdditiveCategoryItemDTO
	StoreAdditiveID uint    `json:"storeAdditiveId"`
	StorePrice      float64 `json:"storePrice"`
}

type StoreAdditiveCategoryDTO struct {
	ID               uint                           `json:"id"`
	Name             string                         `json:"name"`
	Description      string                         `json:"description"`
	Additives        []StoreAdditiveCategoryItemDTO `json:"additives"`
	IsMultipleSelect bool                           `json:"isMultipleSelect"`
}
