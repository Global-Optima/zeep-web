package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
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

func NewProductsModule(
	base *common.BaseModule,
	auditService audit.AuditService,
	franchiseeService franchisees.FranchiseeService,
	ingredientRepo ingredients.IngredientRepository,
	storeWarehouseRepo storeWarehouses.StoreWarehouseRepository,
	notificationService notifications.NotificationService,
) *ProductsModule {
	repo := product.NewProductRepository(base.DB)
	service := product.NewProductService(repo, notificationService, base.Logger)
	handler := product.NewProductHandler(service, auditService)

	recipeModule := NewRecipeModule(base, auditService)
	storeProductsModule := NewStoreProductsModule(base, auditService, service, franchiseeService, repo, ingredientRepo, storeWarehouseRepo)

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

func NewRecipeModule(base *common.BaseModule, auditService audit.AuditService) *RecipeModule {
	repo := recipes.NewRecipeRepository(base.DB)
	service := recipes.NewRecipeService(repo, base.Logger)
	handler := recipes.NewRecipeHandler(service, auditService)

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
	auditService audit.AuditService,
	productService product.ProductService,
	franchiseeService franchisees.FranchiseeService,
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
	handler := storeProducts.NewStoreProductHandler(service, productService, franchiseeService, auditService, base.Logger)

	base.Router.RegisterStoreProductRoutes(handler)

	return &StoreProductsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
