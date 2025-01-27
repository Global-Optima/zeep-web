package types

import (
	"fmt"
	"sort"
	"strings"

	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	categoriesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

type ProductSizeModels struct {
	ProductSize *data.ProductSize
	Additives   []data.ProductSizeAdditive
	Ingredients []data.ProductSizeIngredient
}

func MapToBaseProductDTO(product *data.Product) BaseProductDTO {
	return BaseProductDTO{
		Name:        product.Name,
		Description: product.Description,
		ImageURL:    product.ImageURL,
		VideoURL:    product.VideoURL,
		Category:    *categoriesTypes.MapCategoryToDTO(product.Category),
	}
}

func MapToProductDetailsDTO(product *data.Product) *ProductDetailsDTO {
	var sizes []ProductSizeDTO

	for _, size := range product.ProductSizes {
		sizes = append(sizes, MapToProductSizeDTO(size))
	}

	return &ProductDetailsDTO{
		ProductDTO: MapToProductDTO(*product),
		Sizes:      sizes,
	}
}

func MapToProductDTO(product data.Product) ProductDTO {
	basePrice, productSizeCount := ProductAdditionalInfo(product)

	return ProductDTO{
		ID:               product.ID,
		BaseProductDTO:   MapToBaseProductDTO(&product),
		BasePrice:        basePrice,
		ProductSizeCount: productSizeCount,
	}
}

func ProductAdditionalInfo(product data.Product) (float64, int) {
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

	return basePrice, productSizeCount
}

func MapToBaseProductSizeDTO(productSize data.ProductSize) BaseProductSizeDTO {
	return BaseProductSizeDTO{
		Name:      productSize.Name,
		IsDefault: productSize.IsDefault,
		Unit:      unitTypes.ToUnitResponse(productSize.Unit),
		ProductID: productSize.ProductID,
		Size:      productSize.Size,
		BasePrice: productSize.BasePrice,
	}
}

func MapToProductSizeDTO(productSize data.ProductSize) ProductSizeDTO {
	return ProductSizeDTO{
		ID:                 productSize.ID,
		BaseProductSizeDTO: MapToBaseProductSizeDTO(productSize),
	}
}

func ConvertToProductSizeAdditiveDTO(productSizeAdditive *data.ProductSizeAdditive) ProductSizeAdditiveDTO {
	return ProductSizeAdditiveDTO{
		AdditiveDTO: *additiveTypes.ConvertToAdditiveDTO(&productSizeAdditive.Additive),
		IsDefault:   productSizeAdditive.IsDefault,
	}
}

func MapToProductSizeDetails(productSize data.ProductSize) ProductSizeDetailsDTO {
	var additives = make([]ProductSizeAdditiveDTO, len(productSize.Additives))
	var ingredients = make([]ingredientTypes.IngredientDTO, len(productSize.ProductSizeIngredients))

	for i, productSizeAdditive := range productSize.Additives {
		additives[i] = ConvertToProductSizeAdditiveDTO(&productSizeAdditive)
	}

	for i, productSizeIngredient := range productSize.ProductSizeIngredients {
		ingredients[i] = *ingredientTypes.ConvertToIngredientResponseDTO(&productSizeIngredient.Ingredient)
	}

	return ProductSizeDetailsDTO{
		ProductSizeDTO: MapToProductSizeDTO(productSize),
		Additives:      additives,
		Ingredients:    ingredients,
	}
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

	for _, ingredient := range dto.Ingredients {
		productSize.ProductSizeIngredients = append(productSize.ProductSizeIngredients, data.ProductSizeIngredient{
			IngredientID: ingredient.IngredientID,
			Quantity:     ingredient.Quantity,
		})
	}
	return productSize
}

func UpdateProductToModel(dto *UpdateProductDTO) *data.Product {
	product := &data.Product{}
	if dto == nil {
		return product
	}

	if strings.TrimSpace(dto.Name) != "" {
		product.Name = dto.Name
	}
	if strings.TrimSpace(dto.Description) != "" {
		product.Description = dto.Description
	}
	if strings.TrimSpace(dto.Description) != "" {
		product.ImageURL = dto.ImageURL
	}
	if dto.CategoryID != 0 {
		product.CategoryID = dto.CategoryID
	}
	return product
}

func UpdateProductSizeToModels(dto *UpdateProductSizeDTO) *ProductSizeModels {
	productSize := &data.ProductSize{}

	if dto.Name != nil {
		productSize.Name = *dto.Name
	}
	if dto.UnitID != nil {
		productSize.UnitID = *dto.UnitID
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
	var ingredients []data.ProductSizeIngredient

	for _, additive := range dto.Additives {
		temp := data.ProductSizeAdditive{
			AdditiveID: additive.AdditiveID,
			IsDefault:  additive.IsDefault,
		}
		additives = append(additives, temp)
	}

	for _, ingredient := range dto.Ingredients {
		temp := data.ProductSizeIngredient{
			IngredientID: ingredient.IngredientID,
			Quantity:     ingredient.Quantity,
		}
		ingredients = append(ingredients, temp)
	}

	return &ProductSizeModels{
		ProductSize: productSize,
		Additives:   additives,
		Ingredients: ingredients,
	}
}

func GenerateProductChanges(before *data.Product, dto *UpdateProductDTO) string {
	changes := []string{}

	if dto.Name != "" && dto.Name != before.Name {
		changes = append(changes, fmt.Sprintf("Name: '%s' -> '%s'", before.Name, dto.Name))
	}

	if dto.Description != "" && dto.Description != before.Description {
		changes = append(changes, fmt.Sprintf("Description: '%s' -> '%s'", before.Description, dto.Description))
	}

	if dto.ImageURL != "" && dto.ImageURL != before.ImageURL {
		changes = append(changes, fmt.Sprintf("ImageURL: '%s' -> '%s'", before.ImageURL, dto.ImageURL))
	}

	if dto.CategoryID != 0 && dto.CategoryID != before.CategoryID {
		changes = append(changes, fmt.Sprintf("CategoryID: %d -> %d", before.CategoryID, dto.CategoryID))
	}

	return strings.Join(changes, "; ")
}
