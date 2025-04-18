package types

import (
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
)

type CreateStoreAdditiveDTO struct {
	AdditiveID uint     `json:"additiveId" binding:"required,gte=0"`
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
	Provisions  []additiveTypes.AdditiveProvisionDTO  `json:"provisions"`
}

type StoreAdditiveCategoriesFilter struct {
	Search           *string `form:"search"`
	IsOutOfStock     *bool   `form:"isOutOfStock"`
	IsMultipleSelect *bool   `form:"isMultipleSelect"`
}

type StoreAdditiveCategoryItemDTO struct {
	ID uint `json:"id"`
	additiveTypes.BaseAdditiveDTO
	AdditiveID   uint    `json:"additiveId"`
	StorePrice   float64 `json:"storePrice"`
	IsOutOfStock bool    `json:"isOutOfStock"`
	IsDefault    bool    `json:"isDefault"`
	IsHidden     bool    `json:"isHidden"`
}

type StoreAdditiveCategoryDTO struct {
	ID                    uint                           `json:"id"`
	Name                  string                         `json:"name"`
	TranslatedName        string                         `json:"translatedName,omitempty"`
	Description           string                         `json:"description"`
	TranslatedDescription string                         `json:"translatedDescription,omitempty"`
	Additives             []StoreAdditiveCategoryItemDTO `json:"additives"`
	IsMultipleSelect      bool                           `json:"isMultipleSelect"`
	IsRequired            bool                           `json:"isRequired"`
}
