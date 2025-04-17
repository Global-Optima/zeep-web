package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/translations"
)

type IngredientsModule struct {
	*common.BaseModule
	Repo    ingredients.IngredientRepository
	Service ingredients.IngredientService
	Handler *ingredients.IngredientHandler
}

func NewIngredientsModule(base *common.BaseModule, auditService audit.AuditService, translationManager translations.TranslationManager) *IngredientsModule {
	repo := ingredients.NewIngredientRepository(base.DB)
	service := ingredients.NewIngredientService(repo, base.Logger)
	handler := ingredients.NewIngredientHandler(service, auditService)

	base.Router.RegisterIngredientRoutes(handler)

	return &IngredientsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
