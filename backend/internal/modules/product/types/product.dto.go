package types

import (
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

type ProductSizeIngredientDTO struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Calories float64 `json:"calories"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
	Proteins float64 `json:"proteins"`
}

type ProductDetailsDTO struct {
	ProductDTO
	Sizes []ProductSizeDTO `json:"sizes"`
}

type BaseProductSizeDTO struct {
	Name      string                 `json:"name"`
	BasePrice float64                `json:"basePrice"`
	ProductID uint                   `json:"productId"`
	Unit      unitTypes.UnitResponse `json:"unit"`
	Size      int                    `json:"size"`
	IsDefault bool                   `json:"isDefault"`
}

type ProductSizeDTO struct {
	ID uint `json:"id"`
	BaseProductSizeDTO
}

type ProductSizeDetailsDTO struct {
	ProductSizeDTO
	Additives   []ProductSizeAdditiveDTO        `json:"additives"`
	Ingredients []ingredientTypes.IngredientDTO `json:"ingredients"`
}

type ProductSizeAdditiveDTO struct {
	additiveTypes.AdditiveDTO
	IsDefault bool `json:"isDefault"`
}

type CreateProductDTO struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
	ImageURL    string `json:"imageUrl" binding:"omitempty,url"`
	CategoryID  uint   `json:"categoryId" binding:"omitempty"`
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
	IsDefault   bool                    `json:"isDefault"`
	Additives   []SelectedAdditiveDTO   `json:"additives" binding:"omitempty,dive"`
	Ingredients []SelectedIngredientDTO `json:"ingredients" binding:"required,dive"`
}

type UpdateProductDTO struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	ImageURL    string `json:"imageUrl" binding:"omitempty,url"`
	CategoryID  uint   `json:"categoryId" binding:"omitempty,gt=0"`
}

type UpdateProductSizeDTO struct {
	Name        *string                 `json:"name" binding:"omitempty,max=100"`
	BasePrice   *float64                `json:"basePrice" binding:"omitempty,gt=0"`
	Size        *int                    `json:"size" binding:"omitempty,gt=0"`
	UnitID      *uint                   `json:"unitId" binding:"omitempty,gt=0"`
	IsDefault   *bool                   `json:"isDefault"`
	Additives   []SelectedAdditiveDTO   `json:"additives" binding:"omitempty,dive"`
	Ingredients []SelectedIngredientDTO `json:"ingredients" binding:"omitempty,dive"`
}

type ProductsFilterDto struct {
	utils.BaseFilter
	CategoryID *uint   `form:"categoryId" binding:"omitempty,gt=0"`
	Search     *string `form:"search"`
}
