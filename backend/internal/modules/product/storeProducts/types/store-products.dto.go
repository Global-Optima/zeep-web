package types

import (
	productTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

// StoreProductDTO for sending array of Products
type StoreProductDTO struct {
	productTypes.ProductDTO
	StoreProductID        uint    `json:"storeProductId"`
	StorePrice            float64 `json:"storePrice"`
	StoreProductSizeCount int     `json:"StoreProductSizeCount"`
	IsAvailable           bool    `json:"isAvailable"`
}

// StoreProductDetailsDTO for sending single product with detailed info
type StoreProductDetailsDTO struct {
	StoreProductDTO
	Sizes []StoreProductSizeDTO `json:"sizes"`
}

type StoreProductSizeDTO struct {
	productTypes.ProductSizeDTO
	StoreProductSizeID uint    `json:"storeProductSizeId"`
	StorePrice         float64 `json:"storePrice"`
}

type CreateStoreProductDTO struct {
	ProductID    uint                        `json:"productId" binding:"required,gt=0"`
	IsAvailable  bool                        `json:"isAvailable" binding:"required"`
	ProductSizes []CreateStoreProductSizeDTO `json:"productSizes" binding:"omitempty,gt=0"`
}

type CreateStoreProductSizeDTO struct {
	ProductSizeID uint     `json:"productSizeID" binding:"required,gt=0"`
	Price         *float64 `json:"price" binding:"omitempty,gt=0"`
}

type UpdateStoreProductDTO struct {
	IsAvailable  *bool                       `json:"isAvailable" binding:"omitempty,gt=0"`
	ProductSizes []UpdateStoreProductSizeDTO `json:"productSizes" binding:"omitempty,gt=0"`
}

type UpdateStoreProductSizeDTO struct {
	ProductSizeID uint     `json:"productSizeID" binding:"required,gt=0"`
	Price         *float64 `json:"price" binding:"omitempty,gt=0"`
}

type StoreProductsFilterDTO struct {
	CategoryID  *uint   `form:"categoryId" binding:"omitempty,gt=0"`
	IsAvailable *bool   `form:"isAvailable" binding:"omitempty"`
	Search      *string `form:"search" binding:"omitempty,max=50"`
	utils.BaseFilter
}

type StoreProductSizesFilterDTO struct {
	CategoryID *uint   `form:"categoryId" binding:"omitempty,gt=0"`
	Name       *string `form:"price" binding:"omitempty,max=1"`
	Measure    *string `form:"measure" binding:"omitempty,max=1"`
	Search     *string `form:"search" binding:"omitempty,max=50"`
	IsDefault  *bool   `form:"isDefault" binding:"omitempty"`
	MinSize    *int    `form:"minSize" binding:"omitempty,gt=0"`
	MaxSize    *int    `form:"maxSize" binding:"omitempty,gt=0"`
	utils.BaseFilter
}
