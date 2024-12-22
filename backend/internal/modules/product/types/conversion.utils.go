package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"strings"
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
					Additives: []additiveTypes.AdditiveCategoryItemDTO{},
				}
			}

			additiveDTO := additiveTypes.AdditiveCategoryItemDTO{
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
		additiveDTO := additiveTypes.AdditiveCategoryItemDTO{
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

	var ingredients []ProductIngredientDTO
	for _, size := range product.ProductSizes {
		for _, productIngredient := range size.ProductIngredients {
			ingredient := productIngredient.Ingredient
			ingredients = append(ingredients, ProductIngredientDTO{
				ID:       ingredient.ID,
				Name:     ingredient.Name,
				Calories: ingredient.Calories,
				Fat:      ingredient.Fat,
				Carbs:    ingredient.Carbs,
				Proteins: ingredient.Proteins,
			})
		}
	}

	// Remove duplicate ingredients (if necessary)
	uniqueIngredients := make(map[uint]ProductIngredientDTO)
	for _, ingredient := range ingredients {
		uniqueIngredients[ingredient.ID] = ingredient
	}
	ingredients = nil
	for _, ingredient := range uniqueIngredients {
		ingredients = append(ingredients, ingredient)
	}

	return StoreProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		ImageURL:    product.ImageURL,
		BasePrice:   basePrice,
		Ingredients: ingredients,
	}
}

func CreateToProductModel(dto *CreateProductDTO) *data.Product {
	product := &data.Product{
		Name:        dto.Name,
		Description: dto.Description,
		ImageURL:    dto.ImageURL,
		CategoryID:  dto.CategoryID,
	}

	for _, additiveID := range dto.DefaultAdditives {
		product.DefaultAdditives = append(product.DefaultAdditives, data.DefaultProductAdditive{
			AdditiveID: additiveID,
		})
	}

	return product
}

func CreateToProductSizeModel(dto *CreateProductSizeDTO) *data.ProductSize {
	productSize := &data.ProductSize{
		ProductID: dto.ProductID,
		Name:      dto.Name,
		Measure:   dto.Measure,
		BasePrice: dto.BasePrice,
		Size:      dto.Size,
		IsDefault: dto.IsDefault,
	}
	for _, additiveID := range dto.Additives {
		productSize.Additives = append(productSize.Additives, data.ProductAdditive{
			AdditiveID: additiveID,
		})
	}

	for _, ingredientID := range dto.Ingredients {
		productSize.ProductIngredients = append(productSize.ProductIngredients, data.ProductIngredient{
			ItemIngredientID: ingredientID,
		})
	}
	return productSize
}

func UpdateProductToModel(dto *UpdateProductDTO) *data.Product {
	product := &data.Product{}
	if strings.TrimSpace(dto.Name) != "" {
		product.Name = dto.Name
	}
	if strings.TrimSpace(dto.Description) != "" {
		product.Description = dto.Description
	}
	if strings.TrimSpace(dto.Description) != "" {
		product.ImageURL = dto.ImageURL
	}
	if dto.CategoryID != nil {
		product.CategoryID = dto.CategoryID
	}
	return product
}

func UpdateProductSizeToModel(dto *UpdateProductSizeDTO) *data.ProductSize {
	updatedProductSize := &data.ProductSize{}

	if dto.Name != nil {
		updatedProductSize.Name = *dto.Name
	}
	if dto.Measure != nil {
		updatedProductSize.Measure = *dto.Measure
	}
	if dto.BasePrice != nil {
		updatedProductSize.BasePrice = *dto.BasePrice
	}
	if dto.Size != nil {
		updatedProductSize.Size = *dto.Size
	}
	if dto.IsDefault != nil {
		updatedProductSize.IsDefault = *dto.IsDefault
	}

	return updatedProductSize
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
