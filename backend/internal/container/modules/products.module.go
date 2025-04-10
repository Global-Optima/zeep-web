package modules

import (
	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/recipes"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/technicalMap"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
)

type ProductsModule struct {
	*common.BaseModule
	Repo                product.ProductRepository
	Service             product.ProductService
	Handler             *product.ProductHandler
	Recipe              *RecipeModule
	StoreProductsModule *StoreProductsModule
	TechnicalMapModule  *ProductTechMapModule
}

func NewProductsModule(
	base *common.BaseModule,
	auditService audit.AuditService,
	franchiseeService franchisees.FranchiseeService,
	ingredientRepo ingredients.IngredientRepository,
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	storageRepo storage.StorageRepository,
	notificationService notifications.NotificationService,
) *ProductsModule {
	repo := product.NewProductRepository(base.DB)
	service := product.NewProductService(repo, notificationService, storageRepo, base.Logger)
	handler := product.NewProductHandler(service, auditService, base.Logger)

	recipeModule := NewRecipeModule(base, auditService)
	storeProductsModule := NewStoreProductsModule(base, auditService, service, franchiseeService, repo, ingredientRepo, storeAdditiveRepo, storeStockRepo, storageRepo)
	technicalMapModule := NewProductTechMapModule(base)

	base.Router.RegisterProductRoutes(handler, technicalMapModule.Handler)

	return &ProductsModule{
		BaseModule:          base,
		Repo:                repo,
		Service:             service,
		Handler:             handler,
		Recipe:              recipeModule,
		StoreProductsModule: storeProductsModule,
		TechnicalMapModule:  technicalMapModule,
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
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	storageRepo storage.StorageRepository,
) *StoreProductsModule {
	repo := storeProducts.NewStoreProductRepository(base.DB)
	service := storeProducts.NewStoreProductService(
		repo,
		productRepo,
		ingredientRepo,
		storeAdditiveRepo,
		storageRepo,
		storeProducts.NewTransactionManager(base.DB, repo, storeAdditiveRepo, storeStockRepo, ingredientRepo),
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

type ProductTechMapModule struct {
	*common.BaseModule
	Repo    technicalMap.TechnicalMapRepository
	Service technicalMap.TechnicalMapService
	Handler *technicalMap.TechnicalMapHandler
}

func NewProductTechMapModule(base *common.BaseModule) *ProductTechMapModule {
	repo := technicalMap.NewTechnicalMapRepository(base.DB)
	service := technicalMap.NewTechnicalMapService(repo, base.Logger)
	handler := technicalMap.NewTechnicalMapHandler(service)

	return &ProductTechMapModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
