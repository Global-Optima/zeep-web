package types

import (
	// additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	productCategoriesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

var ProductSizePreloadMap = []utils.LocalizedPreload{
	{Relation: "Unit"}, // no translations on Unit
}

var ProductPreloadMap = []utils.LocalizedPreload{
	utils.Translation("NameTranslation"),
	utils.Translation("DescriptionTranslation"),

	{Relation: "Category", Nested: productCategoriesTypes.ProductCategoryPreloadMap},

	{Relation: "ProductSizes", Nested: ProductSizePreloadMap},
}
