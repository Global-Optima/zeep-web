package types

import (
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type BaseProductDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
}

type ProductDTO struct {
	BaseProductDTO
	BasePrice   float64                `json:"basePrice"`
	Ingredients []ProductIngredientDTO `json:"ingredients"`
}

type ProductIngredientDTO struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Calories float64 `json:"calories"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
	Proteins float64 `json:"proteins"`
}

type ProductDetailsDTO struct {
	BaseProductDTO
	Sizes            []ProductSizeDTO                        `json:"sizes"`
	DefaultAdditives []additiveTypes.AdditiveCategoryItemDTO `json:"defaultAdditives"`
}

type ProductSizeDTO struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	Price   float64 `json:"basePrice"`
	Measure string  `json:"measure"`
}

type CreateProductDTO struct {
	Name             string `json:"name" binding:"required,min=2,max=100"`
	Description      string `json:"description" binding:"max=500"`
	ImageURL         string `json:"imageUrl" binding:"omitempty,url"`
	CategoryID       *uint  `json:"categoryId" binding:"omitempty"`
	DefaultAdditives []uint `json:"defaultAdditives" binding:"omitempty,dive,gt=0"`
}

type CreateProductWithAttachesDTO struct {
	Product      CreateProductDTO           `json:"product" binding:"dive"`
	ProductSizes []CreateProductSizeDTO     `json:"productSizes" binding:"dive"`
	Additives    []SelectedAdditiveTypesDTO `json:"additives" binding:"dive"`
}

type SelectedAdditiveTypesDTO struct {
	AdditiveID uint `json:"additiveId" binding:"required"`
	IsDefault  bool `json:"isDefault"`
}

type SelectedAdditiveDTO struct {
	AdditiveID uint `json:"additiveId" binding:"required,gt=0"`
}

type CreateProductSizeDTO struct {
	ProductID   uint    `json:"productId" binding:"required,gt=0"`
	Name        string  `json:"name" binding:"required,min=2,max=50"`
	Measure     string  `json:"measure" binding:"required,max=20"`
	BasePrice   float64 `json:"basePrice" binding:"required,gt=0"`
	Size        int     `json:"size" binding:"required,gte=0"`
	IsDefault   bool    `json:"isDefault"`
	Additives   []uint  `json:"additives" binding:"omitempty,dive,gt=0"`
	Ingredients []uint  `json:"ingredients" binding:"omitempty,dive,gt=0"`
}

type UpdateProductDTO struct {
	Name             string `json:"name" binding:"omitempty,min=2,max=100"`
	Description      string `json:"description" binding:"omitempty,max=500"`
	ImageURL         string `json:"imageUrl" binding:"omitempty,url"`
	CategoryID       *uint  `json:"categoryId" binding:"omitempty,gt=0"`
	DefaultAdditives []uint `json:"defaultAdditives" binding:"omitempty,dive,gt=0"`
}

type UpdateProductSizeDTO struct {
	Name        *string  `json:"name" binding:"omitempty,max=100"`
	Measure     *string  `json:"measure" binding:"omitempty,oneof=мл г"`
	BasePrice   *float64 `json:"basePrice" binding:"omitempty,gt=0"`
	Size        *int     `json:"size" binding:"omitempty,min=1,max=3"`
	IsDefault   *bool    `json:"isDefault"`
	Ingredients []uint   `json:"ingredients" binding:"omitempty,dive,gt=0"`
	Additives   []uint   `json:"additives" binding:"omitempty,dive,gt=0"`
}

type ProductsFilterDto struct {
	CategoryID *uint   `form:"categoryId" binding:"omitempty,gt=0"`
	Search     *string `form:"search"`
	utils.BaseFilter
}
