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
	once                  sync.Once
	DbHandler             *database.DBHandler
	router                *routes.Router
	logger                *zap.SugaredLogger
	Additives             *modules.AdditivesModule
	Auth                  *modules.AuthModule
	Categories            *modules.CategoriesModule
	Customers             *modules.CustomersModule
	Employees             *modules.EmployeesModule
	Ingredients           *modules.IngredientsModule
	Orders                *modules.OrdersModule
	Products              *modules.ProductsModule
	Stores                *modules.StoresModule
	StoreWarehouses       *modules.StoreWarehouseModule
	Suppliers             *modules.SuppliersModule
	StockRequests         *modules.StockRequestsModule
	Warehouses            *modules.WarehousesModule
	StockMaterialPackages *modules.StockMaterialPackagesModule
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

	c.Additives = modules.NewAdditivesModule(baseModule)
	c.Categories = modules.NewCategoriesModule(baseModule)
	c.Customers = modules.NewCustomersModule(baseModule)
	c.Employees = modules.NewEmployeesModule(baseModule)
	c.Ingredients = modules.NewIngredientsModule(baseModule)
	c.StockRequests = modules.NewStockRequestsModule(baseModule)
	c.StoreWarehouses = modules.NewStoreWarehouseModule(baseModule)
	c.Stores = modules.NewStoresModule(baseModule)
	c.Suppliers = modules.NewSuppliersModule(baseModule)
	c.Warehouses = modules.NewWarehousesModule(baseModule)
	c.StockMaterialPackages = modules.NewStockMaterialPackagesModule(baseModule)

	c.Products = modules.NewProductsModule(baseModule, c.Ingredients.Repo, c.StoreWarehouses.Repo)
	c.Auth = modules.NewAuthModule(baseModule, c.Customers.Repo, c.Employees.Repo)
	c.Orders = modules.NewOrdersModule(baseModule, c.Products.Repo, c.Additives.Repo)
}

func (c *Container) MustInitModules() {
	c.once.Do(c.mustInit)
}

func (c *Container) GetDB() *gorm.DB {
	c.MustInitModules()
	return c.DbHandler.DB
}
