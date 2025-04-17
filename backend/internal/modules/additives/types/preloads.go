package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

var AdditivePreloadMap = []utils.LocalizedPreload{
	utils.Translation("NameTranslation"),
	utils.Translation("DescriptionTranslation"),

	{Relation: "Category", Localized: false, Nested: []utils.LocalizedPreload{
		utils.Translation("NameTranslation"),
		utils.Translation("DescriptionTranslation"),
	}},

	{Relation: "Unit"},

	{Relation: "Ingredients.Ingredient.Unit"},
	{Relation: "Ingredients.Ingredient.IngredientCategory", Nested: []utils.LocalizedPreload{
		utils.Translation("NameTranslation"),
		utils.Translation("DescriptionTranslation"),
	}},

	{Relation: "AdditiveProvisions.Provision.Unit"},
}

var AdditiveCategoryPreloadMap = []utils.LocalizedPreload{
	utils.Translation("NameTranslation"),
	utils.Translation("DescriptionTranslation"),

	{Relation: "Additives", Localized: false, Nested: []utils.LocalizedPreload{
		utils.Translation("NameTranslation"),
		utils.Translation("DescriptionTranslation"),
	}},
}
