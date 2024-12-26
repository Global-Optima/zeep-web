package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	productTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"sort"
)

type StoreProductModels struct {
	StoreProduct      *data.StoreProduct
	StoreProductSizes []data.StoreProductSize
}

func MapToStoreProductDTO(sp *data.StoreProduct) *StoreProductDTO {
	var productSizes []data.ProductSize
	var spsPrices []float64
	var spsMinPrice float64 = 0
	var spsCount = len(sp.StoreProductSizes)

	if sp.StoreProductSizes != nil && spsCount > 0 {
		for _, storeProductSize := range sp.StoreProductSizes {
			spsPrices = append(spsPrices, storeProductSize.Price)
			productSizes = append(productSizes, storeProductSize.ProductSize)
		}

		sort.Float64s(spsPrices)

		spsMinPrice = spsPrices[0]
	}

	return &StoreProductDTO{
		ProductDTO:            productTypes.MapToProductDTO(sp.Product),
		StoreProductID:        sp.ID,
		StoreProductSizeCount: spsCount,
		IsAvailable:           sp.IsAvailable,
		StorePrice:            spsMinPrice,
	}
}

func MapToStoreProductDetailsDTO(sp *data.StoreProduct) StoreProductDetailsDTO {
	sizes := make([]StoreProductSizeDTO, len(sp.StoreProductSizes))
	for i, size := range sp.StoreProductSizes {
		sizes[i] = MapToStoreProductSizeDTO(size)
	}

	return StoreProductDetailsDTO{
		StoreProductDTO: *MapToStoreProductDTO(sp),
		Sizes:           sizes,
	}
}

func MapToStoreProductSizeDTO(sps data.StoreProductSize) StoreProductSizeDTO {
	return StoreProductSizeDTO{
		ProductSizeDTO: productTypes.MapToProductSizeDTO(sps.ProductSize),
		StorePrice:     sps.Price,
	}
}

func CreateToStoreProduct(dto *CreateStoreProductDTO) *data.StoreProduct {
	storeProductSizes := make([]data.StoreProductSize, len(dto.ProductSizes))

	for i, size := range dto.ProductSizes {
		storeProductSizes[i] = data.StoreProductSize{
			ProductSizeID: size.ProductSizeID,
		}
		if size.StorePrice != nil {
			storeProductSizes[i].Price = *size.StorePrice
		}
	}

	return &data.StoreProduct{
		ProductID:         dto.ProductID,
		IsAvailable:       dto.IsAvailable,
		StoreProductSizes: storeProductSizes,
	}
}

func CreateToStoreProductSize(dto *CreateStoreProductSizeDTO) *data.StoreProductSize {
	model := &data.StoreProductSize{
		ProductSizeID: dto.ProductSizeID,
	}

	if dto.StorePrice != nil {
		model.Price = *dto.StorePrice
	} else {
		model.Price = 0
	}

	return model
}

func UpdateToStoreProduct(dto *UpdateStoreProductDTO) *StoreProductModels {
	storeProduct := &data.StoreProduct{}

	if dto.IsAvailable != nil {
		storeProduct.IsAvailable = *dto.IsAvailable
	}

	storeProductSizes := make([]data.StoreProductSize, len(dto.ProductSizes))
	for i, size := range dto.ProductSizes {
		storeProductSizes[i] = data.StoreProductSize{
			ProductSizeID: size.ProductSizeID,
		}
		if size.StorePrice != nil {
			storeProductSizes[i].Price = *size.StorePrice
		}
	}

	return &StoreProductModels{
		StoreProduct:      storeProduct,
		StoreProductSizes: storeProductSizes,
	}
}

func UpdateToStoreProductSize(dto *UpdateStoreProductSizeDTO) *data.StoreProductSize {
	model := &data.StoreProductSize{}

	if dto.StorePrice != nil {
		model.Price = *dto.StorePrice
	}

	return model
}
