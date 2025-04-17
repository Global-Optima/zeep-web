package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

func ConvertToIngredientCategoryResponse(category *data.IngredientCategory) *IngredientCategoryResponse {
	return &IngredientCategoryResponse{
		ID:                     category.ID,
		Name:                   category.Name,
		NameTranslation:        utils.FirstTranslation(category.NameTranslation),
		Description:            category.Description,
		DescriptionTranslation: utils.FirstTranslation(category.DescriptionTranslation),
	}
}
