package types

import (
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"sort"
	"strings"

	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	categoriesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
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
		ImageURL:    product.ImageURL.GetURL(),
		VideoURL:    product.VideoURL.GetURL(),
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
	productSizeCount := len(product.ProductSizes)

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
		Unit:      unitTypes.ToUnitResponse(productSize.Unit),
		ProductID: productSize.ProductID,
		Size:      productSize.Size,
		BasePrice: productSize.BasePrice,
		MachineId: productSize.MachineId,
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
	additives := make([]ProductSizeAdditiveDTO, len(productSize.Additives))
	ingredients := make([]ProductSizeIngredientDTO, len(productSize.ProductSizeIngredients))

	for i, productSizeAdditive := range productSize.Additives {
		additives[i] = ConvertToProductSizeAdditiveDTO(&productSizeAdditive)
	}

	for i, productSizeIngredient := range productSize.ProductSizeIngredients {
		ingredients[i].Ingredient = *ingredientTypes.ConvertToIngredientResponseDTO(&productSizeIngredient.Ingredient)
		ingredients[i].Quantity = productSizeIngredient.Quantity
	}

	return ProductSizeDetailsDTO{
		ProductSizeDTO: MapToProductSizeDTO(productSize),
		TotalNutrition: *CalculateTotalNutrition(&productSize),
		Additives:      additives,
		Ingredients:    ingredients,
	}
}

func CreateToProductModel(dto *CreateProductDTO) *data.Product {
	product := &data.Product{
		Name:        dto.Name,
		Description: dto.Description,
		CategoryID:  dto.CategoryID,
	}

	return product
}

func CreateToProductSizeModel(dto *CreateProductSizeDTO) *data.ProductSize {
	productSize := &data.ProductSize{
		ProductID: dto.ProductID,
		Name:      dto.Name,
		UnitID:    dto.UnitID,
		BasePrice: dto.BasePrice,
		Size:      dto.Size,
		MachineId: dto.MachineId,
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
		return nil
	}

	if strings.TrimSpace(dto.Name) != "" {
		if dto.Name != product.Name {
			product.Name = dto.Name
		}
	}
	if strings.TrimSpace(dto.Description) != "" {
		if dto.Description != product.Description {
			product.Description = dto.Description
		}
	}
	if dto.CategoryID != 0 {
		if dto.CategoryID != product.CategoryID {
			product.CategoryID = dto.CategoryID
		}
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
	if dto.MachineId != nil {
		productSize.MachineId = *dto.MachineId
	}

	var additives []data.ProductSizeAdditive
	var ingredients []data.ProductSizeIngredient

	if dto.Additives != nil {
		if len(dto.Additives) == 0 {
			additives = []data.ProductSizeAdditive{}
		} else {
			for _, additive := range dto.Additives {
				temp := data.ProductSizeAdditive{
					AdditiveID: additive.AdditiveID,
					IsDefault:  additive.IsDefault,
				}
				additives = append(additives, temp)
			}
		}
	}

	if dto.Ingredients != nil {
		if len(dto.Ingredients) == 0 {
			ingredients = []data.ProductSizeIngredient{}
		} else {
			for _, ingredient := range dto.Ingredients {
				temp := data.ProductSizeIngredient{
					IngredientID: ingredient.IngredientID,
					Quantity:     ingredient.Quantity,
				}
				ingredients = append(ingredients, temp)
			}
		}
	}

	return &ProductSizeModels{
		ProductSize: productSize,
		Additives:   additives,
		Ingredients: ingredients,
	}
}

func GenerateProductChanges(before *data.Product, dto *UpdateProductDTO, imageURL data.StorageKey) []details.CentralCatalogChange {
	var changes []details.CentralCatalogChange

	if dto.Name != "" && dto.Name != before.Name {
		key := "notification.centralCatalogUpdateDetails.nameChange"
		changes = append(changes, details.CentralCatalogChange{
			Key: key,
			Params: map[string]interface{}{
				"OldName": before.Name,
				"NewName": dto.Name,
			},
		})
	}

	if dto.Description != "" && dto.Description != before.Description {
		key := "notification.centralCatalogUpdateDetails.descriptionChange"
		changes = append(changes, details.CentralCatalogChange{
			Key: key,
			Params: map[string]interface{}{
				"OldDescription": before.Description,
				"NewDescription": dto.Description,
			},
		})
	}

	if imageURL.ToString() != "" && imageURL != before.ImageURL {
		key := "notification.centralCatalogUpdateDetails.imageUrlChange"
		changes = append(changes, details.CentralCatalogChange{
			Key: key,
			Params: map[string]interface{}{
				"OldImageURL": before.ImageURL,
				"NewImageURL": imageURL.ToString(),
			},
		})
	}

	return changes
}

func CalculateTotalNutrition(productSize *data.ProductSize) *TotalNutrition {
	totalNutrition := &TotalNutrition{}

	for _, psi := range productSize.ProductSizeIngredients {
		totalNutrition.Calories += (psi.Ingredient.Calories * psi.Quantity) / 100
		totalNutrition.Proteins += (psi.Ingredient.Proteins * psi.Quantity) / 100
		totalNutrition.Fats += (psi.Ingredient.Fat * psi.Quantity) / 100
		totalNutrition.Carbs += (psi.Ingredient.Carbs * psi.Quantity) / 100
	}

	for _, psa := range productSize.Additives {
		if psa.IsDefault {
			for _, ai := range psa.Additive.Ingredients {
				totalNutrition.Calories += (ai.Ingredient.Calories * ai.Quantity) / 100
				totalNutrition.Proteins += (ai.Ingredient.Proteins * ai.Quantity) / 100
				totalNutrition.Fats += (ai.Ingredient.Fat * ai.Quantity) / 100
				totalNutrition.Carbs += (ai.Ingredient.Carbs * ai.Quantity) / 100
			}
		}
	}

	totalNutrition.Calories = utils.RoundToOneDecimal(totalNutrition.Calories)
	totalNutrition.Proteins = utils.RoundToOneDecimal(totalNutrition.Proteins)
	totalNutrition.Fats = utils.RoundToOneDecimal(totalNutrition.Fats)
	totalNutrition.Carbs = utils.RoundToOneDecimal(totalNutrition.Carbs)

	return totalNutrition
}
