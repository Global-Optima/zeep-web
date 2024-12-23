package types

import (
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	productTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type StoreProductDTO struct {
	productTypes.ProductDTO
	IsAvailable bool    `json:"isAvailable"`
	Price       float64 `json:"price"`
}

type StoreProductDetailsDTO struct {
	productTypes.BaseProductDTO
	Sizes            []StoreProductSizeDTO `json:"sizes"`
	DefaultAdditives []additiveTypes.AdditiveCategoryItemDTO
	IsAvailable      bool    `json:"isAvailable"`
	Price            float64 `json:"price"`
}

type StoreProductWithSizesDTO struct {
	productTypes.BaseProductDTO
	IsAvailable       bool                  `json:"isAvailable"`
	Price             float64               `json:"price"`
	StoreProductSizes []StoreProductSizeDTO `json:"storeProductSizes"`
}

type StoreProductSizeDTO struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Measure    string  `json:"measure"`
	StorePrice float64 `json:"storePrice"`
	BasePrice  float64 `json:"basePrice"`
	Size       int     `json:"size"`
	IsDefault  bool    `json:"isDefault"`
}

type CreateStoreProductDTO struct {
	StoreID     uint `json:"storeId" binding:"required,gt=0"`
	ProductID   uint `json:"productId" binding:"required,gt=0"`
	IsAvailable bool `json:"isAvailable" binding:"required"`
}

type CreateStoreProductSizeDTO struct {
	StoreID       uint     `json:"storeId" binding:"required,gt=0"`
	ProductSizeID uint     `json:"productSizeID" binding:"required,gt=0"`
	Price         *float64 `json:"price" binding:"omitempty,gt=0"`
}

type UpdateStoreProductDTO struct {
	IsAvailable *bool `json:"isAvailable" binding:"omitempty,gt=0"`
}

type UpdateStoreProductSizeDTO struct {
	Price *float64 `json:"price" binding:"required,gt=0"`
}

type StoreProductsFilterDTO struct {
	Name        *string `form:"name" binding:"omitempty"`
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
