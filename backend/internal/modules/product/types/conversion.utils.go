package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"sort"
	"strings"
)

type ProductSizeModels struct {
	ProductSize *data.ProductSize
	Additives   []data.ProductSizeAdditive
	Ingredients []data.ProductIngredient
}

func MapBaseProductDTO(product *data.Product) BaseProductDTO {
	return BaseProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		ImageURL:    product.ImageURL,
	}
}

func MapToProductDetailsDTO(product *data.Product) *ProductDetailsDTO {
	var sizes []ProductSizeDTO

	for _, size := range product.ProductSizes {
		sizes = append(sizes, MapToProductSizeDTO(size))
	}

	return &ProductDetailsDTO{
		BaseProductDTO: MapBaseProductDTO(product),
		Sizes:          sizes,
	}
}

func MapToProductSizeAdditives(productSizeAdditives []data.ProductSizeAdditive) []additiveTypes.AdditiveCategoryItemDTO {
	var result []additiveTypes.AdditiveCategoryItemDTO
	for _, psa := range productSizeAdditives {
		additive := psa.Additive
		additiveDTO := additiveTypes.AdditiveCategoryItemDTO{
			ID:          additive.ID,
			Name:        additive.Name,
			Description: additive.Description,
			Price:       additive.BasePrice,
			ImageURL:    additive.ImageURL,
		}
		result = append(result, additiveDTO)
	}
	return result
}

func MapToProductDTO(product data.Product) ProductDTO {
	var basePrice float64 = 0
	var productSizesPrices []float64
	var productSizeCount = len(product.ProductSizes)

	if productSizeCount > 0 {
		for _, ps := range product.ProductSizes {
			productSizesPrices = append(productSizesPrices, ps.BasePrice)
		}

		sort.Float64s(productSizesPrices)
		basePrice = productSizesPrices[0]
	}

	return ProductDTO{
		BaseProductDTO:   MapBaseProductDTO(&product),
		BasePrice:        basePrice,
		ProductSizeCount: productSizeCount,
	}
}

func MapToProductSizeDTO(productSize data.ProductSize) ProductSizeDTO {
	return ProductSizeDTO{
		ID:        productSize.ID,
		Name:      productSize.Name,
		IsDefault: productSize.IsDefault,
		Measure:   productSize.Measure,
		Size:      productSize.Size,
		BasePrice: productSize.BasePrice,
	}
}

func MapToIngredients(sizes []data.ProductSize) []ProductSizeIngredientDTO {
	var ingredients []ProductSizeIngredientDTO
	for _, size := range sizes {
		for _, productIngredient := range size.ProductIngredients {
			ingredient := productIngredient.Ingredient
			ingredients = append(ingredients, ProductSizeIngredientDTO{
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
	uniqueIngredients := make(map[uint]ProductSizeIngredientDTO)
	for _, ingredient := range ingredients {
		uniqueIngredients[ingredient.ID] = ingredient
	}
	ingredients = nil
	for _, ingredient := range uniqueIngredients {
		ingredients = append(ingredients, ingredient)
	}

	return ingredients
}

func CreateToProductModel(dto *CreateProductDTO) *data.Product {
	product := &data.Product{
		Name:        dto.Name,
		Description: dto.Description,
		ImageURL:    dto.ImageURL,
		CategoryID:  dto.CategoryID,
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

	for _, additive := range dto.Additives {
		productSize.Additives = append(productSize.Additives, data.ProductSizeAdditive{
			AdditiveID: additive.AdditiveID,
			IsDefault:  additive.IsDefault,
		})
	}

	for _, ingredientID := range dto.Ingredients {
		productSize.ProductIngredients = append(productSize.ProductIngredients, data.ProductIngredient{
			IngredientID: ingredientID,
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

func UpdateProductSizeToModels(dto *UpdateProductSizeDTO) *ProductSizeModels {
	productSize := &data.ProductSize{}

	if dto.Name != nil {
		productSize.Name = *dto.Name
	}
	if dto.Measure != nil {
		productSize.Measure = *dto.Measure
	}
	if dto.BasePrice != nil {
		productSize.BasePrice = *dto.BasePrice
	}
	if dto.Size != nil {
		productSize.Size = *dto.Size
	}
	if dto.IsDefault != nil {
		productSize.IsDefault = *dto.IsDefault
	}

	var additives []data.ProductSizeAdditive
	var ingredients []data.ProductIngredient

	for _, additive := range dto.Additives {
		var temp data.ProductSizeAdditive
		temp = data.ProductSizeAdditive{
			AdditiveID: additive.AdditiveID,
			IsDefault:  additive.IsDefault,
		}
		additives = append(additives, temp)
	}

	for _, ingredientID := range dto.Ingredients {
		var temp data.ProductIngredient
		temp = data.ProductIngredient{
			IngredientID: ingredientID,
		}
		ingredients = append(ingredients, temp)
	}

	return &ProductSizeModels{
		ProductSize: productSize,
		Additives:   additives,
		Ingredients: ingredients,
	}
}
