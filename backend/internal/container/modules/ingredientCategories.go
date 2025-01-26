package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/ingredientCategories"
)

type IngredientCategoriesModule struct {
	*common.BaseModule
	Repo    ingredientCategories.IngredientCategoryRepository
	Service ingredientCategories.IngredientCategoryService
	Handler *ingredientCategories.IngredientCategoryHandler
}

func NewIngredientCategoriesModule(base *common.BaseModule) *IngredientCategoriesModule {
	repo := ingredientCategories.NewIngredientCategoryRepository(base.DB)
	service := ingredientCategories.NewIngredientCategoryService(repo)
	handler := ingredientCategories.NewIngredientCategoryHandler(service)

	base.Router.RegisterIngredientCategoriesRoutes(handler)

	return &IngredientCategoriesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
