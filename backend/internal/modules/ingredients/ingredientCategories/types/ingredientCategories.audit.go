package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateIngredientCategoryAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.IngredientCategoryComponent)

	UpdateIngredientCategoryAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.IngredientCategoryComponent, &UpdateIngredientCategoryDTO{})

	DeleteIngredientCategoryAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.IngredientCategoryComponent)
)
