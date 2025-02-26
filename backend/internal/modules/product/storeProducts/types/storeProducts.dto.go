package types

import (
	productTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type StoreProductDTO struct {
	ID uint `json:"id"`
	productTypes.BaseProductDTO
	ProductID             uint    `json:"productId"`
	BasePrice             float64 `json:"basePrice"`
	StorePrice            float64 `json:"storePrice"`
	ProductSizeCount      int     `json:"productSizeCount"`
	StoreProductSizeCount int     `json:"storeProductSizeCount"`
	IsAvailable           bool    `json:"isAvailable"`
}

type StoreProductSizeDTO struct {
	ID uint `json:"id"`
	productTypes.BaseProductSizeDTO
	ProductSizeID uint    `json:"productSizeId"`
	StorePrice    float64 `json:"storePrice"`
}

type StoreProductDetailsDTO struct {
	StoreProductDTO
	Sizes []StoreProductSizeDTO `json:"sizes"`
}

type StoreProductSizeDetailsDTO struct {
	StoreProductSizeDTO
	Additives   []productTypes.ProductSizeAdditiveDTO   `json:"additives"`
	Ingredients []productTypes.ProductSizeIngredientDTO `json:"ingredients"`
}

type CreateStoreProductDTO struct {
	ProductID    uint                        `json:"productId" binding:"required,gt=0"`
	IsAvailable  bool                        `json:"isAvailable" binding:"required"`
	ProductSizes []CreateStoreProductSizeDTO `json:"productSizes" binding:"required,min=1,dive"`
}

type CreateStoreProductSizeDTO struct {
	ProductSizeID uint     `json:"productSizeID" binding:"required,gt=0"`
	StorePrice    *float64 `json:"storePrice" binding:"omitempty,gt=0"`
}

type UpdateStoreProductDTO struct {
	IsAvailable  *bool                       `json:"isAvailable"`
	ProductSizes []UpdateStoreProductSizeDTO `json:"productSizes" binding:"required,min=1,dive"`
}

type UpdateStoreProductSizeDTO struct {
	ProductSizeID uint     `json:"productSizeID" binding:"required,gt=0"`
	StorePrice    *float64 `json:"storePrice" binding:"omitempty,gt=0"`
}

type StoreProductsFilterDTO struct {
	utils.BaseFilter
	CategoryID  *uint    `form:"categoryId" binding:"omitempty,gt=0"`
	IsAvailable *bool    `form:"isAvailable" binding:"omitempty"`
	Search      *string  `form:"search" binding:"omitempty,max=50"`
	MaxPrice    *float64 `form:"maxPrice" binding:"omitempty,gt=0"`
	MinPrice    *float64 `form:"minPrice" binding:"omitempty,gt=0"`
}

type ExcludedStoreProductsFilterDTO struct {
	StoreProductIDs []uint `json:"storeProductIds" binding:"required,gt=0"`
}

type StoreProductSizesFilterDTO struct {
	utils.BaseFilter
	CategoryID *uint   `form:"categoryId" binding:"omitempty,gt=0"`
	Name       *string `form:"price" binding:"omitempty,max=1"`
	Measure    *string `form:"measure" binding:"omitempty,max=1"`
	Search     *string `form:"search" binding:"omitempty,max=50"`
	IsDefault  *bool   `form:"isDefault" binding:"omitempty"`
	MinSize    *int    `form:"minSize" binding:"omitempty,gt=0"`
	MaxSize    *int    `form:"maxSize" binding:"omitempty,gt=0"`
}
