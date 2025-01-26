package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/recipes"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
)

type ProductsModule struct {
	*common.BaseModule
	Repo          product.ProductRepository
	Service       product.ProductService
	Handler       *product.ProductHandler
	Recipe        *RecipeModule
	StoreProducts *StoreProductsModule
}

func NewProductsModule(base *common.BaseModule, ingredientRepo ingredients.IngredientRepository, storeWarehouseRepo storeWarehouses.StoreWarehouseRepository) *ProductsModule {
	repo := product.NewProductRepository(base.DB)
	service := product.NewProductService(repo, base.Logger)
	handler := product.NewProductHandler(service)

	recipeModule := NewRecipeModule(base)
	storeProductsModule := NewStoreProductsModule(base, repo, ingredientRepo, storeWarehouseRepo)

	base.Router.RegisterProductRoutes(handler)

	return &ProductsModule{
		BaseModule:    base,
		Repo:          repo,
		Service:       service,
		Handler:       handler,
		Recipe:        recipeModule,
		StoreProducts: storeProductsModule,
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

	base.Router.RegisterRecipeRoutes(handler)

	return &RecipeModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}

type StoreProductsModule struct {
	*common.BaseModule
	Repo        storeProducts.StoreProductRepository
	ProductRepo product.ProductRepository
	Service     storeProducts.StoreProductService
	Handler     *storeProducts.StoreProductHandler
}

func NewStoreProductsModule(
	base *common.BaseModule,
	productRepo product.ProductRepository,
	ingredientRepo ingredients.IngredientRepository,
	storeWarehouseRepo storeWarehouses.StoreWarehouseRepository,
) *StoreProductsModule {
	repo := storeProducts.NewStoreProductRepository(base.DB)
	service := storeProducts.NewStoreProductService(
		repo,
		productRepo,
		ingredientRepo,
		storeProducts.NewTransactionManager(base.DB, repo, storeWarehouseRepo),
		base.Logger)
	handler := storeProducts.NewStoreProductHandler(service)

	base.Router.RegisterStoreProductRoutes(handler)

	return &StoreProductsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
