package container

import (
	"sync"

	asynqManager "github.com/Global-Optima/zeep-web/backend/internal/asynqTasks"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth/employeeToken"

	"github.com/Global-Optima/zeep-web/backend/api/storage"

	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/container/modules"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
	"github.com/Global-Optima/zeep-web/backend/internal/scheduler"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	once                    sync.Once
	DbHandler               *database.DBHandler
	RedisClient             *database.RedisClient
	AsynqManager            *asynqManager.AsynqManager
	storageRepo             *storage.StorageRepository
	employeeTokenManager    *employeeToken.EmployeeTokenManager
	router                  *routes.Router
	logger                  *zap.SugaredLogger
	Additives               *modules.AdditivesModule
	Audits                  *modules.AuditsModule
	Notifications           *modules.NotificationModule
	Auth                    *modules.AuthModule
	Categories              *modules.CategoriesModule
	Customers               *modules.CustomersModule
	Employees               *modules.EmployeesModule
	Franchisees             *modules.FranchiseesModule
	Ingredients             *modules.IngredientsModule
	IngredientCategories    *modules.IngredientCategoriesModule
	Orders                  *modules.OrdersModule
	Products                *modules.ProductsModule
	Regions                 *modules.RegionsModule
	Stores                  *modules.StoresModule
	StoreStocks             *modules.StoreStockModule
	StoreSynchronizer       *modules.StoreSynchronizerModule
	Suppliers               *modules.SuppliersModule
	StockRequests           *modules.StockRequestsModule
	Warehouses              *modules.WarehousesModule
	StockMaterials          *modules.StockMaterialsModule
	StockMaterialCategories *modules.StockMaterialCategoriesModule
	Units                   *modules.UnitsModule
	Analytics               *modules.AnalyticsModule
}

func NewContainer(dbHandler *database.DBHandler, redisClient *database.RedisClient, storageRepo *storage.StorageRepository, employeeTokenManager *employeeToken.EmployeeTokenManager, router *routes.Router, logger *zap.SugaredLogger) *Container {
	return &Container{
		DbHandler:            dbHandler,
		RedisClient:          redisClient,
		storageRepo:          storageRepo,
		employeeTokenManager: employeeTokenManager,
		router:               router,
		logger:               logger,
	}
}

func (c *Container) mustInit() {
	baseModule := common.NewBaseModule(c.DbHandler.DB, c.router, c.logger)
	cronManager := scheduler.NewCronManager(c.logger)

	var err error
	c.AsynqManager, err = asynqManager.NewAsyncManager(c.RedisClient.Client, c.logger)
	if err != nil {
		c.logger.Fatalf("Failed to create asynq manager: %v", err)
	}

	c.Audits = modules.NewAuditsModule(baseModule)
	c.Franchisees = modules.NewFranchiseesModule(baseModule, c.Audits.Service)
	c.Regions = modules.NewRegionsModule(baseModule, c.Audits.Service)
	c.Notifications = modules.NewNotificationModule(baseModule)
	c.Categories = modules.NewCategoriesModule(baseModule, c.Audits.Service)
	c.Customers = modules.NewCustomersModule(baseModule)
	c.Employees = modules.NewEmployeesModule(baseModule, c.Audits.Service, c.Franchisees.Service, c.Regions.Service, *c.employeeTokenManager)
	c.Ingredients = modules.NewIngredientsModule(baseModule, c.Audits.Service)
	c.Suppliers = modules.NewSuppliersModule(baseModule, c.Audits.Service)
	c.StockMaterials = modules.NewStockMaterialsModule(baseModule, c.Audits.Service)
	c.StockMaterialCategories = modules.NewStockMaterialCategoriesModule(baseModule, c.Audits.Service)
	c.Units = modules.NewUnitsModule(baseModule, c.Audits.Service)
	c.IngredientCategories = modules.NewIngredientCategoriesModule(baseModule, c.Audits.Service)
	c.Warehouses = modules.NewWarehousesModule(baseModule, c.StockMaterials.Repo, c.Notifications.Service, cronManager, c.Regions.Service, c.Franchisees.Service, c.Audits.Service)
	c.Stores = modules.NewStoresModule(baseModule, c.Franchisees.Service, c.Audits.Service)
	c.StoreStocks = modules.NewStoreStockModule(baseModule, c.Ingredients.Service, c.Franchisees.Service, c.Audits.Service, c.Notifications.Service, c.Stores.Service, cronManager)

	c.Additives = modules.NewAdditivesModule(baseModule, c.Audits.Service, c.Franchisees.Service, c.Ingredients.Repo, c.StoreStocks.Repo, *c.storageRepo)
	c.Products = modules.NewProductsModule(baseModule, c.Audits.Service, c.Franchisees.Service, c.Additives.Service, c.Ingredients.Repo, c.Additives.StoreAdditivesModule.Repo, c.StoreStocks.Repo, *c.storageRepo, c.Notifications.Service)
	c.Auth = modules.NewAuthModule(baseModule, c.Customers.Repo, c.Employees.Repo, *c.employeeTokenManager)
	c.Orders = modules.NewOrdersModule(baseModule, c.AsynqManager, c.Products.StoreProductsModule.Repo, c.Additives.StoreAdditivesModule.Repo, c.StoreStocks.Repo, c.Notifications.Service)
	c.StockRequests = modules.NewStockRequestsModule(baseModule, c.Franchisees.Service, c.Regions.Service, c.StockMaterials.Repo, c.Notifications.Service, c.Audits.Service)
	c.StoreSynchronizer = modules.NewStoreSynchronizerSynchronizerModule(baseModule, c.Stores.Repo, c.Additives.StoreAdditivesModule.Repo, c.StoreStocks.Repo)
	c.Analytics = modules.NewAnalyticsModule(baseModule)

	cronManager.Start()
}

func (c *Container) MustInitModules() {
	c.once.Do(c.mustInit)
}

func (c *Container) GetDB() *gorm.DB {
	return c.DbHandler.DB
}
