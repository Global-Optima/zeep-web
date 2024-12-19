package types

import (
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type StoreProductDTO struct {
	ID          uint                   `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	ImageURL    string                 `json:"imageUrl"`
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

type StoreProductDetailsDTO struct {
	ID               uint                                    `json:"id"`
	Name             string                                  `json:"name"`
	Description      string                                  `json:"description"`
	ImageURL         string                                  `json:"imageUrl"`
	Sizes            []ProductSizeDTO                        `json:"sizes"`
	DefaultAdditives []additiveTypes.AdditiveCategoryItemDTO `json:"defaultAdditives"`
}

type ProductSizeDTO struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	BasePrice float64 `json:"basePrice"`
	Measure   string  `json:"measure"`
}

type CreateStoreProduct struct {
	Name         string                 `json:"name" binding:"required"`
	Description  string                 `json:"description"`
	ImageURL     string                 `json:"imageUrl"`
	CategoryID   *uint                  `json:"categoryId"`
	ProductSizes []CreateProductSizeDTO `json:"productSizes"`
	Additives    []SelectedAdditiveDTO  `json:"additives"`
}

type UpdateStoreProduct struct {
	ID           uint                   `json:"id" binding:"required"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	ImageURL     string                 `json:"imageUrl"`
	CategoryID   *uint                  `json:"categoryId"`
	ProductSizes []UpdateProductSizeDTO `json:"productSizes"`
	Additives    []SelectedAdditiveDTO  `json:"additives"`
}

type SelectedAdditiveDTO struct {
	AdditiveID uint `json:"additiveId" binding:"required"`
	IsDefault  bool `json:"isDefault"`
}

type CreateProductSizeDTO struct {
	Name      string  `json:"name" binding:"required"`
	Measure   string  `json:"measure"`
	BasePrice float64 `json:"basePrice" binding:"required"`
	Size      int     `json:"size"`
	IsDefault bool    `json:"isDefault"`
}

type UpdateProductSizeDTO struct {
	ID        uint    `json:"id" binding:"required"`
	Name      string  `json:"name"`
	Measure   string  `json:"measure"`
	BasePrice float64 `json:"basePrice"`
	Size      int     `json:"size"`
	IsDefault bool    `json:"isDefault"`
}

type ProductsFilterDto struct {
	StoreID    *uint   `form:"storeId"`
	CategoryID *uint   `form:"categoryId"`
	Search     *string `form:"search"`
	utils.BaseFilter
}
