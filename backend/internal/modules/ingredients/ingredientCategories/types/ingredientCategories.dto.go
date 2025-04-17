package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/translations"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type CreateIngredientCategoryDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description" binding:"omitempty"`
}

type UpdateIngredientCategoryDTO struct {
	Name        *string `json:"name" binding:"omitempty,min=1"`
	Description *string `json:"description" binding:"omitempty"`
}

type IngredientCategoryResponse struct {
	ID                     uint   `json:"id"`
	Name                   string `json:"name"`
	NameTranslation        string `json:"nameTranslation,omitempty"` // locale‑specific
	Description            string `json:"description"`
	DescriptionTranslation string `json:"descriptionTranslation,omitempty"` // locale‑specific
}

type IngredientCategoryFilter struct {
	utils.BaseFilter
	Search *string `form:"search"` // Search by name or description
}

type IngredientCategoryTranslationDTO struct {
	Name        translations.FieldLocale `json:"name" binding:"omitempty"`
	Description translations.FieldLocale `json:"description" binding:"omitempty"`
}
