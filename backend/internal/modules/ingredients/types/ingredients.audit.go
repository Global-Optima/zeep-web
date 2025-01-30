package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateIngredientAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.IngredientComponent)

	UpdateIngredientAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.IngredientComponent, &UpdateIngredientDTO{})

	DeleteIngredientAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.IngredientComponent)
)
