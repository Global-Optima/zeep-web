package types

import (
	// additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

var storeAdditiveCategoryMiniPreload = []utils.LocalizedPreload{
	utils.Translation("NameTranslation"),
	utils.Translation("DescriptionTranslation"),
}

var StoreAdditivePreloadMap = []utils.LocalizedPreload{
	utils.Translation("NameTranslation"),
	utils.Translation("DescriptionTranslation"),
	{Relation: "Unit", Nested: unitTypes.UnitPreloadMap},
	{Relation: "Category", Nested: storeAdditiveCategoryMiniPreload},
}

var StoreAdditiveCategoryPreloadMap = []utils.LocalizedPreload{
	utils.Translation("NameTranslation"),
	utils.Translation("DescriptionTranslation"),
	{Relation: "Additives", Nested: StoreAdditivePreloadMap},
}
