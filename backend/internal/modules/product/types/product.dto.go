package types

import (
	"mime/multipart"

	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	categoriesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type BaseProductDTO struct {
	Name        string                             `json:"name"`
	Description string                             `json:"description"`
	ImageURL    string                             `json:"imageUrl"`
	VideoURL    string                             `json:"videoUrl"`
	Category    categoriesTypes.ProductCategoryDTO `json:"category"`
}

type ProductDTO struct {
	ID uint `json:"id"`
	BaseProductDTO
	ProductSizeCount int     `json:"productSizeCount"`
	BasePrice        float64 `json:"basePrice"`
}

type ProductDetailsDTO struct {
	ProductDTO
	Sizes []ProductSizeDTO `json:"sizes"`
}

type BaseProductSizeDTO struct {
	Name      string             `json:"name"`
	BasePrice float64            `json:"basePrice"`
	ProductID uint               `json:"productId"`
	Unit      unitTypes.UnitsDTO `json:"unit"`
	Size      int                `json:"size"`
}

type ProductSizeDTO struct {
	ID uint `json:"id"`
	BaseProductSizeDTO
}

type ProductSizeIngredientDTO struct {
	Quantity   float64                       `json:"quantity"`
	Ingredient ingredientTypes.IngredientDTO `json:"ingredient"`
}

type ProductSizeDetailsDTO struct {
	ProductSizeDTO
	Additives   []ProductSizeAdditiveDTO   `json:"additives"`
	Ingredients []ProductSizeIngredientDTO `json:"ingredients"`
}

type ProductSizeAdditiveDTO struct {
	additiveTypes.AdditiveDTO
	IsDefault bool `json:"isDefault"`
}

type CreateProductDTO struct {
	Name        string `form:"name" binding:"required,min=2,max=100"`
	Description string `form:"description" binding:"max=500"`
	CategoryID  uint   `form:"categoryId" binding:"required"`
	Image       *multipart.FileHeader
	Video       *multipart.FileHeader
}

type SelectedAdditiveDTO struct {
	AdditiveID uint `json:"additiveId" binding:"required"`
	IsDefault  bool `json:"isDefault"`
}

type SelectedIngredientDTO struct {
	IngredientID uint    `json:"ingredientId" binding:"required,gt=0"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
}

type CreateProductSizeDTO struct {
	ProductID   uint                    `json:"productId" binding:"required,gt=0"`
	Name        string                  `json:"name" binding:"required,oneof=S M L"`
	Size        int                     `json:"size" binding:"required,gte=0"`
	UnitID      uint                    `json:"unitId" binding:"required,gt=0"`
	BasePrice   float64                 `json:"basePrice" binding:"required,gt=0"`
	Additives   []SelectedAdditiveDTO   `json:"additives" binding:"omitempty,dive"`
	Ingredients []SelectedIngredientDTO `json:"ingredients" binding:"required,dive"`
}

type UpdateProductDTO struct {
	Name        string `form:"name" binding:"omitempty,min=2,max=100"`
	Description string `form:"description" binding:"omitempty,max=500"`
	CategoryID  uint   `form:"categoryId" binding:"omitempty,gt=0"`
	Image       *multipart.FileHeader
	Video       *multipart.FileHeader
}

type UpdateProductSizeDTO struct {
	Name        *string                 `json:"name" binding:"omitempty,max=100"`
	BasePrice   *float64                `json:"basePrice" binding:"omitempty,gt=0"`
	Size        *int                    `json:"size" binding:"omitempty,gt=0"`
	UnitID      *uint                   `json:"unitId" binding:"omitempty,gt=0"`
	Additives   []SelectedAdditiveDTO   `json:"additives" binding:"omitempty,dive"`
	Ingredients []SelectedIngredientDTO `json:"ingredients" binding:"omitempty,dive"`
}

type ProductsFilterDto struct {
	utils.BaseFilter
	CategoryID *uint   `form:"categoryId" binding:"omitempty,gt=0"`
	Search     *string `form:"search"`
}
