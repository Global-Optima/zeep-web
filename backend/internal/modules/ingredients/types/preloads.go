package types

import (
	ingredientCategoryTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/ingredientCategories/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

var IngredientPreloadMap = []utils.LocalizedPreload{
	utils.Translation("NameTranslation"),

	{Relation: "Unit"},

	{Relation: "IngredientCategory", Nested: ingredientCategoryTypes.IngredientCategoryPreloadMap},
}
