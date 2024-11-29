package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
)

func MapToStoreProductDetailsDTO(product *data.Product, defaultAdditives []data.DefaultProductAdditive) *StoreProductDetailsDTO {
	dto := &StoreProductDetailsDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		ImageURL:    product.ImageURL,
	}

	for _, size := range product.ProductSizes {
		sizeDTO := ProductSizeDTO{
			ID:        size.ID,
			Name:      size.Name,
			BasePrice: size.BasePrice,
			Measure:   size.Measure,
		}

		additiveCategoryMap := make(map[uint]*additiveTypes.AdditiveCategoryDTO)
		for _, pa := range size.Additives {
			additive := pa.Additive
			categoryID := additive.AdditiveCategoryID

			if _, exists := additiveCategoryMap[categoryID]; !exists {
				additiveCategoryMap[categoryID] = &additiveTypes.AdditiveCategoryDTO{
					ID:        categoryID,
					Name:      additive.Category.Name,
					Additives: []additiveTypes.AdditiveDTO{},
				}
			}

			additiveDTO := additiveTypes.AdditiveDTO{
				ID:          additive.ID,
				Name:        additive.Name,
				Description: additive.Description,
				Price:       additive.BasePrice,
				ImageURL:    additive.ImageURL,
			}

			additiveCategoryMap[categoryID].Additives = append(additiveCategoryMap[categoryID].Additives, additiveDTO)
		}

		dto.Sizes = append(dto.Sizes, sizeDTO)
	}

	for _, da := range defaultAdditives {
		additive := da.Additive
		additiveDTO := additiveTypes.AdditiveDTO{
			ID:          additive.ID,
			Name:        additive.Name,
			Description: additive.Description,
			Price:       additive.BasePrice,
			ImageURL:    additive.ImageURL,
		}
		dto.DefaultAdditives = append(dto.DefaultAdditives, additiveDTO)
	}

	return dto
}

func MapToStoreProductDTO(product data.Product) StoreProductDTO {
	var basePrice float64 = 0
	if len(product.ProductSizes) > 0 {
		basePrice = product.ProductSizes[0].BasePrice
	}

	return StoreProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		ImageURL:    product.ImageURL,
		BasePrice:   basePrice,
	}
}

func CreateToProductModel(dto *CreateStoreProduct) *data.Product {
	return &data.Product{
		Name:        dto.Name,
		Description: dto.Description,
		ImageURL:    dto.ImageURL,
		CategoryID:  dto.CategoryID,
	}
}

func UpdateToProductModel(dto *UpdateStoreProduct) *data.Product {
	return &data.Product{
		BaseEntity: data.BaseEntity{
			ID: dto.ID,
		},
		Name:        dto.Name,
		Description: dto.Description,
		ImageURL:    dto.ImageURL,
		CategoryID:  dto.CategoryID,
	}
}

func ToProductSizesModels(dtoSizes []CreateProductSizeDTO, productID uint) []data.ProductSize {
	var sizes []data.ProductSize
	for _, size := range dtoSizes {
		sizes = append(sizes, data.ProductSize{
			ProductID: productID,
			Name:      size.Name,
			Measure:   size.Measure,
			BasePrice: size.BasePrice,
			Size:      size.Size,
			IsDefault: size.IsDefault,
		})
	}
	return sizes
}
