package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateRecipeStepsAuditFactory = shared.NewAuditActionExtendedFactory(
		data.CreateOperation, data.RecipeStepsComponent, &CreateOrReplaceRecipeStepDTO{})

	UpdateRecipeStepsAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.RecipeStepsComponent, &CreateOrReplaceRecipeStepDTO{})

	DeleteRecipeStepsAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.RecipeStepsComponent)
)
