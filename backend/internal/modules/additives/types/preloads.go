package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

var AdditivePreloadMap = []utils.LocalizedPreload{
	{Relation: "NameTranslation", Localized: true},
	{Relation: "DescriptionTranslation", Localized: true},
	{Relation: "Category", Localized: false, Nested: []utils.LocalizedPreload{
		{Relation: "NameTranslation", Localized: true},
		{Relation: "DescriptionTranslation", Localized: true},
	}},
	{Relation: "Unit", Localized: false},
	{Relation: "Ingredients.Ingredient.Unit", Localized: false},
	{Relation: "Ingredients.Ingredient.IngredientCategory", Localized: false, Nested: []utils.LocalizedPreload{
		{Relation: "NameTranslation", Localized: true},
		{Relation: "DescriptionTranslation", Localized: true},
	}},
	{Relation: "AdditiveProvisions.Provision.Unit", Localized: false},
}
