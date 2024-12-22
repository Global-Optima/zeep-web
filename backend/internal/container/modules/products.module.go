package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/recipes"
)

type ProductsModule struct {
	*common.BaseModule
	Repo    product.ProductRepository
	Service product.ProductService
	Handler *product.ProductHandler
	Recipe  *RecipeModule
}

func NewProductsModule(base *common.BaseModule) *ProductsModule {
	repo := product.NewProductRepository(base.DB)
	service := product.NewProductService(repo, base.Logger)
	handler := product.NewProductHandler(service)

	recipeModule := NewRecipeModule(base)

	base.Router.RegisterProductRoutes(handler)
	base.Router.RegisterRecipeRoutes(recipeModule.Handler)

	return &ProductsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
		Recipe:     recipeModule,
	}
}

type RecipeModule struct {
	*common.BaseModule
	Repo    recipes.RecipeRepository
	Service recipes.RecipeService
	Handler *recipes.RecipeHandler
}

func NewRecipeModule(base *common.BaseModule) *RecipeModule {
	repo := recipes.NewRecipeRepository(base.DB)
	service := recipes.NewRecipeService(repo, base.Logger)
	handler := recipes.NewRecipeHandler(service)

	return &RecipeModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
