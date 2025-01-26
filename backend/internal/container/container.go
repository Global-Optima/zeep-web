package container

import (
	"fmt"
	"sync"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/container/modules"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	once                    sync.Once
	DbHandler               *database.DBHandler
	router                  *routes.Router
	logger                  *zap.SugaredLogger
	Additives               *modules.AdditivesModule
	Audits                  *modules.AuditsModule
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
	StoreWarehouses         *modules.StoreWarehouseModule
	Suppliers               *modules.SuppliersModule
	StockRequests           *modules.StockRequestsModule
	Warehouses              *modules.WarehousesModule
	StockMaterials          *modules.StockMaterialsModule
	StockMaterialCategories *modules.StockMaterialCategoriesModule
	Barcodes                *modules.BarcodeModule
	Units                   *modules.UnitsModule
	Analytics               *modules.AnalyticsModule
}

func NewContainer(dbHandler *database.DBHandler, router *routes.Router, logger *zap.SugaredLogger) *Container {
	return &Container{
		DbHandler: dbHandler,
		router:    router,
		logger:    logger,
	}
}

func (c *Container) mustInit() {
	var err error
	cfg := config.GetConfig()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	c.DbHandler, err = database.InitDB(dsn)
	if err != nil {
		c.logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	baseModule := common.NewBaseModule(c.DbHandler.DB, c.router, c.logger)

	c.Audits = modules.NewAuditsModule(baseModule)
	c.Franchisees = modules.NewFranchiseesModule(baseModule, c.Audits.Service)
	c.Regions = modules.NewRegionsModule(baseModule, c.Audits.Service)
	c.Categories = modules.NewCategoriesModule(baseModule, c.Audits.Service)
	c.Customers = modules.NewCustomersModule(baseModule)
	c.Employees = modules.NewEmployeesModule(baseModule, c.Audits.Service)
	c.Ingredients = modules.NewIngredientsModule(baseModule, c.Audits.Service)
	c.StoreWarehouses = modules.NewStoreWarehouseModule(baseModule, c.Ingredients.Service, c.Audits.Service)
	c.Stores = modules.NewStoresModule(baseModule, c.Audits.Service)
	c.Suppliers = modules.NewSuppliersModule(baseModule)
	c.StockMaterials = modules.NewStockMaterialsModule(baseModule)
	c.StockMaterialCategories = modules.NewStockMaterialCategoriesModule(baseModule)
	c.Barcodes = modules.NewBarcodeModule(baseModule, c.StockMaterials.Repo)
	c.Units = modules.NewUnitsModule(baseModule)
	c.IngredientCategories = modules.NewIngredientCategoriesModule(baseModule, c.Audits.Service)
	c.Warehouses = modules.NewWarehousesModule(baseModule, c.StockMaterials.Repo, c.Barcodes.Repo)

	c.Additives = modules.NewAdditivesModule(baseModule, c.Audits.Service, c.Ingredients.Repo, c.StoreWarehouses.Repo)
	c.Products = modules.NewProductsModule(baseModule, c.Audits.Service, c.Ingredients.Repo, c.StoreWarehouses.Repo)
	c.Auth = modules.NewAuthModule(baseModule, c.Customers.Repo, c.Employees.Repo)
	c.Orders = modules.NewOrdersModule(baseModule, c.Products.Repo, c.Additives.Repo)
	c.StockRequests = modules.NewStockRequestsModule(baseModule, c.StockMaterials.Repo)
	c.Analytics = modules.NewAnalyticsModule(baseModule)
}

func (c *Container) MustInitModules() {
	c.once.Do(c.mustInit)
}

func (c *Container) GetDB() *gorm.DB {
	c.MustInitModules()
	return c.DbHandler.DB
}
