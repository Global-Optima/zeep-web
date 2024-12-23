package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	productTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
)

func MapToStoreProductDTO(sp *data.StoreProduct) *StoreProductDTO {
	var productSizes []data.ProductSize
	for _, storeProductSize := range sp.Store.ProductSizes {
		productSizes = append(productSizes, storeProductSize.ProductSize)
	}

	return &StoreProductDTO{
		ProductDTO: productTypes.ProductDTO{
			BaseProductDTO: productTypes.MapBaseProductDTO(&sp.Product),
			Ingredients:    productTypes.MapToIngredients(productSizes),
		},
		IsAvailable: sp.IsAvailable,
		Price:       sp.Store.ProductSizes[0].Price,
	}
}

func MapToStoreProductDetailsDTO(sp *data.StoreProduct, da []data.DefaultProductAdditive) StoreProductDetailsDTO {
	return StoreProductDetailsDTO{
		BaseProductDTO:   productTypes.MapBaseProductDTO(&sp.Product),
		Sizes:            MapStoreProductSizes(sp.Store.ProductSizes, sp.StoreID),
		DefaultAdditives: productTypes.MapToDefaultAdditives(da),
		IsAvailable:      sp.IsAvailable,
		Price:            sp.Store.ProductSizes[0].Price,
	}
}

func MapStoreProductSizes(sizes []data.StoreProductSize, storeID uint) []StoreProductSizeDTO {
	var result []StoreProductSizeDTO
	for _, size := range sizes {
		if size.StoreID == storeID {
			result = append(result, StoreProductSizeDTO{
				ID:         size.ID,
				Name:       size.ProductSize.Name,
				Measure:    size.ProductSize.Measure,
				StorePrice: size.Price,
				BasePrice:  size.ProductSize.BasePrice,
				Size:       size.ProductSize.Size,
				IsDefault:  size.ProductSize.IsDefault,
			})
		}
	}

	return result
}

func MapToStoreProductSizeDTO(input data.StoreProductSize) StoreProductSizeDTO {
	return StoreProductSizeDTO{
		ID:         input.ID,
		Name:       input.ProductSize.Name,
		Measure:    input.ProductSize.Measure,
		StorePrice: input.Price,
		BasePrice:  input.ProductSize.BasePrice,
		Size:       input.ProductSize.Size,
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
