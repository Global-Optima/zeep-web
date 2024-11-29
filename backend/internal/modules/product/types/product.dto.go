package types

import additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"

type StoreProductDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
	BasePrice   float64 `json:"basePrice"`
}

type StoreProductDetailsDTO struct {
	ID               uint                        `json:"id"`
	Name             string                      `json:"name"`
	Description      string                      `json:"description"`
	ImageURL         string                      `json:"imageUrl"`
	Sizes            []ProductSizeDTO            `json:"sizes"`
	DefaultAdditives []additiveTypes.AdditiveDTO `json:"defaultAdditives"`
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
	ImageURL     string                 `json:"image_url"`
	CategoryID   *uint                  `json:"category_id"`
	ProductSizes []CreateProductSizeDTO `json:"product_sizes"`
	Additives    []SelectedAdditiveDTO  `json:"additives"` // Additive IDs with default flag
}

type UpdateStoreProduct struct {
	ID           uint                   `json:"id" binding:"required"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	ImageURL     string                 `json:"image_url"`
	CategoryID   *uint                  `json:"category_id"`
	ProductSizes []UpdateProductSizeDTO `json:"product_sizes"`
	Additives    []SelectedAdditiveDTO  `json:"additives"` // Additive IDs with default flag
}

type SelectedAdditiveDTO struct {
	AdditiveID uint `json:"additive_id" binding:"required"`
	IsDefault  bool `json:"is_default"`
}

type CreateProductSizeDTO struct {
	Name      string  `json:"name" binding:"required"`
	Measure   string  `json:"measure"`
	BasePrice float64 `json:"base_price" binding:"required"`
	Size      int     `json:"size"`
	IsDefault bool    `json:"is_default"`
}

type UpdateProductSizeDTO struct {
	ID        uint    `json:"id" binding:"required"`
	Name      string  `json:"name"`
	Measure   string  `json:"measure"`
	BasePrice float64 `json:"base_price"`
	Size      int     `json:"size"`
	IsDefault bool    `json:"is_default"`
}
