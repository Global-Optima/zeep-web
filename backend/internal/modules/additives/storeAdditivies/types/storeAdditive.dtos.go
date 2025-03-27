package types

import (
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
)

type CreateStoreAdditiveDTO struct {
	AdditiveID uint     `json:"additiveId" binding:"required,gt=0"`
	StorePrice *float64 `json:"storePrice" binding:"omitempty"`
}

type UpdateStoreAdditiveDTO struct {
	StorePrice *float64 `json:"storePrice"`
}

type StoreAdditiveDTO struct {
	ID uint `json:"id"`
	additiveTypes.BaseAdditiveDTO
	AdditiveID   uint    `json:"additiveId"`
	StorePrice   float64 `json:"storePrice"`
	IsOutOfStock bool    `json:"isOutOfStock"`
}

type StoreAdditiveDetailsDTO struct {
	StoreAdditiveDTO
	Ingredients []additiveTypes.AdditiveIngredientDTO `json:"ingredients"`
}

type StoreAdditiveCategoriesFilter struct {
	Search           *string `form:"search"`
	IsMultipleSelect *bool   `form:"isMultipleSelect"`
}

type StoreAdditiveCategoryItemDTO struct {
	ID uint `json:"id"`
	additiveTypes.BaseAdditiveCategoryItemDTO
	AdditiveID uint    `json:"additiveId"`
	StorePrice float64 `json:"storePrice"`
	IsDefault  bool    `json:"isDefault"`
}

type StoreAdditiveCategoryDTO struct {
	ID               uint                           `json:"id"`
	Name             string                         `json:"name"`
	Description      string                         `json:"description"`
	Additives        []StoreAdditiveCategoryItemDTO `json:"additives"`
	IsMultipleSelect bool                           `json:"isMultipleSelect"`
}
