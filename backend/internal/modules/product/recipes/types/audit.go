package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateProductAuditFactory = shared.NewAuditActionExtendedFactory(
		data.CreateOperation, data.RecipeStepsComponent, &CreateOrReplaceRecipeStepDTO{})

	UpdateProductAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.RecipeStepsComponent, &CreateOrReplaceRecipeStepDTO{})

	DeleteProductAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.RecipeStepsComponent)
)
