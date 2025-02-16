package types

import (
	"sort"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	productTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
)

type StoreProductModels struct {
	StoreProduct      *data.StoreProduct
	StoreProductSizes []data.StoreProductSize
}

func MapToStoreProductDTO(sp *data.StoreProduct) *StoreProductDTO {
	basePrice, productSizeCount := productTypes.ProductAdditionalInfo(sp.Product)
	storePrice, storeProductSizeCount := StoreProductAdditionalInfo(*sp)

	return &StoreProductDTO{
		ID:                    sp.ID,
		BaseProductDTO:        productTypes.BaseProductDTO{}, //productTypes.MapToBaseProductDTO(&sp.Product, "imageUrl"),
		ProductID:             sp.ProductID,
		ProductSizeCount:      productSizeCount,
		BasePrice:             basePrice,
		StoreProductSizeCount: storeProductSizeCount,
		StorePrice:            storePrice,
		IsAvailable:           sp.IsAvailable,
	}
}

func MapToStoreProductDetailsDTO(sp *data.StoreProduct) StoreProductDetailsDTO {
	sizes := make([]StoreProductSizeDetailsDTO, len(sp.StoreProductSizes))

	for i, size := range sp.StoreProductSizes {
		sizes[i] = MapToStoreProductSizeDetailsDTO(size)
	}

	return StoreProductDetailsDTO{
		StoreProductDTO: *MapToStoreProductDTO(sp),
		Sizes:           sizes,
	}
}

func StoreProductAdditionalInfo(sp data.StoreProduct) (float64, int) {
	var spsPrices []float64
	var spsMinPrice float64 = 0
	var spsCount = len(sp.StoreProductSizes)

	if sp.StoreProductSizes != nil && spsCount > 0 {
		for _, storeProductSize := range sp.StoreProductSizes {
			spsPrices = append(spsPrices, getStorePrice(&storeProductSize))
		}

		sort.Float64s(spsPrices)

		spsMinPrice = spsPrices[0]
	}

	return spsMinPrice, spsCount
}

func MapToStoreProductSizeDetailsDTO(sps data.StoreProductSize) StoreProductSizeDetailsDTO {
	var additives = make([]productTypes.ProductSizeAdditiveDTO, len(sps.ProductSize.Additives))
	var ingredients = make([]productTypes.ProductSizeIngredientDTO, len(sps.ProductSize.ProductSizeIngredients))

	for i, productSizeAdditive := range sps.ProductSize.Additives {
		additives[i] = productTypes.ConvertToProductSizeAdditiveDTO(&productSizeAdditive)
	}

	for i, productSizeIngredient := range sps.ProductSize.ProductSizeIngredients {
		ingredients[i].Ingredient = *ingredientTypes.ConvertToIngredientResponseDTO(&productSizeIngredient.Ingredient)
		ingredients[i].Quantity = productSizeIngredient.Quantity
	}

	return StoreProductSizeDetailsDTO{
		ID:                 sps.ID,
		BaseProductSizeDTO: productTypes.MapToBaseProductSizeDTO(sps.ProductSize),
		ProductSizeID:      sps.ProductSizeID,
		StorePrice:         getStorePrice(&sps),
		Additives:          additives,
		Ingredients:        ingredients,
	}
}

func CreateToStoreProduct(dto *CreateStoreProductDTO) *data.StoreProduct {
	storeProductSizes := make([]data.StoreProductSize, len(dto.ProductSizes))

	for i, size := range dto.ProductSizes {
		storeProductSizes[i] = data.StoreProductSize{
			ProductSizeID: size.ProductSizeID,
		}
		storeProductSizes[i].StorePrice = size.StorePrice

	}

	return &data.StoreProduct{
		ProductID:         dto.ProductID,
		IsAvailable:       dto.IsAvailable,
		StoreProductSizes: storeProductSizes,
	}
}

func CreateToStoreProductSize(dto *CreateStoreProductSizeDTO) *data.StoreProductSize {
	return &data.StoreProductSize{
		ProductSizeID: dto.ProductSizeID,
		StorePrice:    dto.StorePrice,
	}
}

func UpdateToStoreProductModels(dto *UpdateStoreProductDTO) *StoreProductModels {
	storeProduct := &data.StoreProduct{}

	if dto.IsAvailable != nil {
		storeProduct.IsAvailable = *dto.IsAvailable
	}

	storeProductSizes := make([]data.StoreProductSize, len(dto.ProductSizes))
	for i, size := range dto.ProductSizes {
		storeProductSizes[i] = data.StoreProductSize{
			ProductSizeID: size.ProductSizeID,
		}
		storeProductSizes[i].StorePrice = size.StorePrice
	}

	return &StoreProductModels{
		StoreProduct:      storeProduct,
		StoreProductSizes: storeProductSizes,
	}
}

func UpdateToStoreProductSize(dto *UpdateStoreProductSizeDTO) *data.StoreProductSize {
	return &data.StoreProductSize{
		StorePrice: dto.StorePrice,
	}
}

func getStorePrice(storeProductSize *data.StoreProductSize) float64 {
	if storeProductSize.StorePrice != nil {
		return *storeProductSize.StorePrice
	}
	return storeProductSize.ProductSize.BasePrice
}
