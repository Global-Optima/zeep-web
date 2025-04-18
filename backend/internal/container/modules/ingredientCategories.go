package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/ingredientCategories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/translations"
)

type IngredientCategoriesModule struct {
	*common.BaseModule
	Repo    ingredientCategories.IngredientCategoryRepository
	Service ingredientCategories.IngredientCategoryService
	Handler *ingredientCategories.IngredientCategoryHandler
}

func NewIngredientCategoriesModule(base *common.BaseModule, auditService audit.AuditService, translationManager translations.TranslationManager) *IngredientCategoriesModule {
	repo := ingredientCategories.NewIngredientCategoryRepository(base.DB)
	transactionManager := ingredientCategories.NewTransactionManager(base.DB, repo, translationManager)
	service := ingredientCategories.NewIngredientCategoryService(repo, transactionManager)
	handler := ingredientCategories.NewIngredientCategoryHandler(service, auditService)

	base.Router.RegisterIngredientCategoriesRoutes(handler)

	return &IngredientCategoriesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
