package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	productTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"sort"
)

func MapToStoreProductDTO(sp *data.StoreProduct) *StoreProductDTO {
	var productSizes []data.ProductSize
	var spsPrices []float64
	var spsMinPrice float64 = 0
	var spsCount = len(sp.Store.ProductSizes)

	if sp.Store.ProductSizes != nil && spsCount > 0 {
		for _, storeProductSize := range sp.Store.ProductSizes {
			spsPrices = append(spsPrices, storeProductSize.Price)
			productSizes = append(productSizes, storeProductSize.ProductSize)
		}

		sort.Float64s(spsPrices)

		spsMinPrice = spsPrices[0]
	}

	return &StoreProductDTO{
		ProductDTO: productTypes.ProductDTO{
			BaseProductDTO:   productTypes.MapBaseProductDTO(&sp.Product),
			ProductSizeCount: len(sp.Store.ProductSizes),
			BasePrice:        sp.Store.ProductSizes[0].Price, //TODO
		},
		StoreProductID:        sp.ID,
		StoreProductSizeCount: spsCount,
		IsAvailable:           sp.IsAvailable,
		StorePrice:            spsMinPrice,
	}
}

func MapToStoreProductDetailsDTO(sp *data.StoreProduct) StoreProductDetailsDTO {
	return StoreProductDetailsDTO{
		StoreProductDTO: *MapToStoreProductDTO(sp),
		Sizes:           MapToStoreProductSizeDTOs(sp.Store.ProductSizes, sp.StoreID),
	}
}

func MapToStoreProductSizeDTOs(sizes []data.StoreProductSize, storeID uint) []StoreProductSizeDTO {
	var result []StoreProductSizeDTO
	for _, size := range sizes {
		if size.StoreID == storeID {
			result = append(result, MapToStoreProductSizeDTO(size))
		}
	}

	return result
}

func MapToStoreProductSizeDTO(sps data.StoreProductSize) StoreProductSizeDTO {
	return StoreProductSizeDTO{
		ProductSizeDTO: productTypes.MapToProductSizeDTO(sps.ProductSize),
		StorePrice:     sps.Price,
	}
}

func CreateToStoreProduct(dto *CreateStoreProductDTO) *data.StoreProduct {
	return &data.StoreProduct{
		ProductID:   dto.ProductID,
		IsAvailable: dto.IsAvailable,
	}
}

func CreateToStoreProductSize(dto *CreateStoreProductSizeDTO) *data.StoreProductSize {
	model := &data.StoreProductSize{
		ProductSizeID: dto.ProductSizeID,
	}

	if dto.Price != nil {
		model.Price = *dto.Price
	} else {
		model.Price = 0
	}

	return model
}

func UpdateToStoreProduct(dto *UpdateStoreProductDTO) *data.StoreProduct {
	model := &data.StoreProduct{}

	if dto.IsAvailable != nil {
		model.IsAvailable = *dto.IsAvailable
	}
	return model
}

func UpdateToStoreProductSize(dto *UpdateStoreProductSizeDTO) *data.StoreProductSize {
	model := &data.StoreProductSize{}

	if dto.Price != nil {
		model.Price = *dto.Price
	}

	return model
}
