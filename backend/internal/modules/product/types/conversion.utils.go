package types

import (
	"fmt"
	"sort"

	provisionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"

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
	Provisions  []data.ProductSizeProvision
}

func MapToBaseProductDTO(product *data.Product) BaseProductDTO {
	if product == nil {
		return BaseProductDTO{}
	}

	return BaseProductDTO{
		Name:        product.Name,
		Description: product.Description,
		ImageURL:    product.ImageKey.GetURL(),
		VideoURL:    product.VideoKey.GetURL(),
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
		IsHidden:    productSizeAdditive.IsHidden,
	}
}

func ConvertToProductSizeProvisionDTO(productSizeProvision *data.ProductSizeProvision) ProductSizeProvisionDTO {
	return ProductSizeProvisionDTO{
		Provision: *provisionsTypes.MapToProvisionDTO(&productSizeProvision.Provision),
		Volume:    productSizeProvision.Volume,
	}
}

func ConvertToProductSizeIngredientDTO(productSizeIngredient *data.ProductSizeIngredient) ProductSizeIngredientDTO {
	return ProductSizeIngredientDTO{
		Ingredient: *ingredientTypes.ConvertToIngredientResponseDTO(&productSizeIngredient.Ingredient),
		Quantity:   productSizeIngredient.Quantity,
	}
}

func MapToProductSizeDetails(productSize data.ProductSize) ProductSizeDetailsDTO {
	additives := make([]ProductSizeAdditiveDTO, len(productSize.Additives))
	ingredients := make([]ProductSizeIngredientDTO, len(productSize.ProductSizeIngredients))
	provisions := make([]ProductSizeProvisionDTO, len(productSize.ProductSizeProvisions))

	for i, productSizeIngredient := range productSize.ProductSizeIngredients {
		ingredients[i] = ConvertToProductSizeIngredientDTO(&productSizeIngredient)
	}

	for i, productSizeAdditive := range productSize.Additives {
		additives[i] = ConvertToProductSizeAdditiveDTO(&productSizeAdditive)
	}

	for i, productSizeProvision := range productSize.ProductSizeProvisions {
		provisions[i] = ConvertToProductSizeProvisionDTO(&productSizeProvision)
	}

	return ProductSizeDetailsDTO{
		ProductSizeDTO: MapToProductSizeDTO(productSize),
		TotalNutrition: *CalculateTotalNutrition(&productSize),
		Additives:      additives,
		Ingredients:    ingredients,
		Provisions:     provisions,
	}
}

func CreateToProductModel(dto *CreateProductDTO) *data.Product {
	product := &data.Product{
		Name:        dto.Name,
		Description: utils.DerefString(dto.Description),
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

	productSize.Additives = mapAdditivesToProductSizeAdditives(dto.Additives)
	productSize.ProductSizeIngredients = mapIngredientsToProductSizeIngredients(dto.Ingredients)
	productSize.ProductSizeProvisions = mapProvisionsToProductSizeProvisions(dto.Provisions)

	return productSize
}

func UpdateProductToModel(dto *UpdateProductDTO, product *data.Product) error {
	if dto == nil {
		return fmt.Errorf("cannot update nil product")
	}

	if dto.Name != nil {
		if dto.Name != &product.Name {
			product.Name = *dto.Name
		}
	}

	if dto.Description != nil {
		product.Description = *dto.Description
	}

	if dto.CategoryID != 0 {
		product.CategoryID = dto.CategoryID
	}

	return nil
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

	additives := mapAdditivesToProductSizeAdditives(dto.Additives)
	ingredients := mapIngredientsToProductSizeIngredients(dto.Ingredients)
	provisions := mapProvisionsToProductSizeProvisions(dto.Provisions)

	return &ProductSizeModels{
		ProductSize: productSize,
		Additives:   additives,
		Ingredients: ingredients,
		Provisions:  provisions,
	}
}

func mapAdditivesToProductSizeAdditives(additivesDTO []SelectedAdditiveDTO) []data.ProductSizeAdditive {
	if additivesDTO == nil {
		return nil
	}

	additives := make([]data.ProductSizeAdditive, 0, len(additivesDTO))

	for _, dto := range additivesDTO {
		additive := data.ProductSizeAdditive{
			AdditiveID: dto.AdditiveID,
			IsDefault:  dto.IsDefault,
		}

		if dto.IsDefault {
			additive.IsHidden = dto.IsHidden
		} else {
			additive.IsHidden = false
		}

		additives = append(additives, additive)
	}

	return additives
}

func mapIngredientsToProductSizeIngredients(ingredientsDTO []ingredientTypes.SelectedIngredientDTO) []data.ProductSizeIngredient {
	if ingredientsDTO == nil {
		return nil
	}

	ingredients := make([]data.ProductSizeIngredient, 0, len(ingredientsDTO))

	for _, dto := range ingredientsDTO {
		ingredient := data.ProductSizeIngredient{
			IngredientID: dto.IngredientID,
			Quantity:     dto.Quantity,
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients
}

func mapProvisionsToProductSizeProvisions(provisionsDTO []provisionsTypes.SelectedProvisionDTO) []data.ProductSizeProvision {
	if provisionsDTO == nil {
		return nil
	}

	provisions := make([]data.ProductSizeProvision, 0, len(provisionsDTO))

	for _, dto := range provisionsDTO {
		provision := data.ProductSizeProvision{
			ProvisionID: dto.ProvisionID,
			Volume:      dto.Volume,
		}
		provisions = append(provisions, provision)
	}

	return provisions
}

var emptyStringRu = "пустое значение"

// emptyStringEn = "empty value"
// emptyStringKk = "бос мән"

func GenerateProductChanges(before *data.Product, dto *UpdateProductDTO, imageKey *data.StorageImageKey) []details.CentralCatalogChange {
	var changes []details.CentralCatalogChange

	if dto.Name != nil && dto.Name != &before.Name {
		key := "notification.centralCatalogUpdateDetails.nameChange"

		if *dto.Name == "" {
			changes = append(changes, details.CentralCatalogChange{
				Key: key,
				Params: map[string]interface{}{
					"OldName": before.Name,
					"NewName": emptyStringRu,
				},
			})
		} else {
			changes = append(changes, details.CentralCatalogChange{
				Key: key,
				Params: map[string]interface{}{
					"OldName": before.Name,
					"NewName": dto.Name,
				},
			})
		}
	}

	if dto.Description != nil && dto.Description != &before.Description {
		key := "notification.centralCatalogUpdateDetails.descriptionChange"
		if *dto.Description == "" {
			changes = append(changes, details.CentralCatalogChange{
				Key: key,
				Params: map[string]interface{}{
					"OldDescription": before.Description,
					"NewDescription": emptyStringRu,
				},
			})
		} else {
			changes = append(changes, details.CentralCatalogChange{
				Key: key,
				Params: map[string]interface{}{
					"OldDescription": before.Description,
					"NewDescription": dto.Description,
				},
			})
		}
	}

	if imageKey != before.ImageKey {
		key := "notification.centralCatalogUpdateDetails.imageUrlChange"
		changes = append(changes, details.CentralCatalogChange{
			Key: key,
			Params: map[string]interface{}{
				"OldImageURL": before.ImageKey.GetURL(),
				"NewImageURL": imageKey.GetURL(),
			},
		})
	}

	return changes
}

func CalculateTotalNutrition(productSize *data.ProductSize) *TotalNutrition {
	totalNutrition := &TotalNutrition{}
	ingredientSet := make(map[string]struct{})
	allergenSet := make(map[string]struct{})

	for _, psi := range productSize.ProductSizeIngredients {
		totalNutrition.Calories += (psi.Ingredient.Calories * psi.Quantity) / 100
		totalNutrition.Proteins += (psi.Ingredient.Proteins * psi.Quantity) / 100
		totalNutrition.Fats += (psi.Ingredient.Fat * psi.Quantity) / 100
		totalNutrition.Carbs += (psi.Ingredient.Carbs * psi.Quantity) / 100

		ingredientSet[psi.Ingredient.Name] = struct{}{}

		if psi.Ingredient.IsAllergen {
			allergenSet[psi.Ingredient.Name] = struct{}{}
		}
	}

	for _, psa := range productSize.Additives {
		if psa.IsDefault {
			for _, ai := range psa.Additive.Ingredients {
				totalNutrition.Calories += (ai.Ingredient.Calories * ai.Quantity) / 100
				totalNutrition.Proteins += (ai.Ingredient.Proteins * ai.Quantity) / 100
				totalNutrition.Fats += (ai.Ingredient.Fat * ai.Quantity) / 100
				totalNutrition.Carbs += (ai.Ingredient.Carbs * ai.Quantity) / 100

				ingredientSet[ai.Ingredient.Name] = struct{}{}

				if ai.Ingredient.IsAllergen {
					allergenSet[ai.Ingredient.Name] = struct{}{}
				}
			}
		}
	}

	totalNutrition.Ingredients = make([]string, 0, len(ingredientSet))
	for name := range ingredientSet {
		totalNutrition.Ingredients = append(totalNutrition.Ingredients, name)
	}

	totalNutrition.AllergenIngredients = make([]string, 0, len(allergenSet))
	for name := range allergenSet {
		totalNutrition.AllergenIngredients = append(totalNutrition.AllergenIngredients, name)
	}

	totalNutrition.Calories = utils.RoundToDecimal(totalNutrition.Calories, 1)
	totalNutrition.Proteins = utils.RoundToDecimal(totalNutrition.Proteins, 1)
	totalNutrition.Fats = utils.RoundToDecimal(totalNutrition.Fats, 1)
	totalNutrition.Carbs = utils.RoundToDecimal(totalNutrition.Carbs, 1)

	return totalNutrition
}
