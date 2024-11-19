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
