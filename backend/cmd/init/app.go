package init

import (
	"fmt"
	"log"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/customers"

	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/supplier"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/inventory"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitializeConfig() *config.Config {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}

func InitializeDatabase(cfg *config.Config) *database.DBHandler {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	dbHandler, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	return dbHandler
}

func InitializeRedis(cfg *config.Config) *database.RedisClient {
	redisClient, err := database.InitRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	utils.InitCache(redisClient.Client, redisClient.Ctx)

	return redisClient
}

// func InitializeKafka(cfg *config.Config) *kafka.KafkaManager {
// 	kafkaManager, err := kafka.NewKafkaManager(cfg.Kafka)
// 	if err != nil {
// 		log.Fatalf("failed to initialize kafka instance: %v", err)
// 		return nil
// 	}
// 	return kafkaManager
// }

func InitializeModule[T any, H any](
	dbHandler *database.DBHandler,
	initService func(dbHandler *database.DBHandler) (T, error),
	createHandler func(T) H,
	registerRoutes func(H)) {

	service, err := initService(dbHandler)
	if err != nil {
		log.Fatalf("Error initializing service: %v", err)
		return
	}

	handler := createHandler(service)

	registerRoutes(handler)
}

func InitializeRouter(dbHandler *database.DBHandler, redisClient *database.RedisClient, storageRepo storage.StorageRepository) *gin.Engine {
	cfg := config.GetConfig()

	router := gin.New()
	router.Use(logger.ZapLoggerMiddleware())

	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.Server.ClientURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(middleware.RedisMiddleware(redisClient.Client))

	apiRouter := routes.NewRouter(router, "/api", "/v1")

	storageHandler := storage.NewStorageHandler(storageRepo)        // temp
	storage.RegisterStorageRoutes(apiRouter.Routes, storageHandler) // temp

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (auth.AuthenticationService, error) {
			return auth.NewAuthenticationService(
				auth.NewAuthorizationRepository(dbHandler.DB),
				customers.NewCustomerRepository(dbHandler.DB),
				employees.NewEmployeeRepository(dbHandler.DB),
				logger.GetZapSugaredLogger(),
			), nil
		},
		auth.NewAuthenticationHandler,
		apiRouter.RegisterAuthenticationRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (product.ProductService, error) {
			return product.NewProductService(product.NewProductRepository(dbHandler.DB), logger.GetZapSugaredLogger()), nil
		},
		product.NewProductHandler,
		apiRouter.RegisterProductRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (stores.StoreService, error) {
			return stores.NewStoreService(stores.NewStoreRepository(dbHandler.DB)), nil
		},
		stores.NewStoreHandler,
		apiRouter.RegisterStoresRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (categories.CategoryService, error) {
			return categories.NewCategoryService(categories.NewCategoryRepository(dbHandler.DB)), nil
		},
		categories.NewCategoryHandler,
		apiRouter.RegisterProductCategoriesRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (additives.AdditiveService, error) {
			return additives.NewAdditiveService(additives.NewAdditiveRepository(dbHandler.DB), logger.GetZapSugaredLogger()), nil
		},
		additives.NewAdditiveHandler,
		apiRouter.RegisterAdditivesRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (employees.EmployeeService, error) {
			return employees.NewEmployeeService(employees.NewEmployeeRepository(dbHandler.DB), logger.GetZapSugaredLogger()), nil
		},
		employees.NewEmployeeHandler,
		apiRouter.RegisterEmployeesRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (orders.OrderService, error) {
			return orders.NewOrderService(
				orders.NewOrderRepository(dbHandler.DB),
				product.NewProductRepository(dbHandler.DB),
				additives.NewAdditiveRepository(dbHandler.DB),
				logger.GetZapSugaredLogger(),
			), nil
		},
		orders.NewOrderHandler,
		apiRouter.RegisterOrderRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (supplier.SupplierService, error) {
			return supplier.NewSupplierService(supplier.NewSupplierRepository(dbHandler.DB)), nil
		},
		supplier.NewSupplierHandler,
		apiRouter.RegisterSupplierRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (storeWarehouses.StoreWarehouseService, error) {
			return storeWarehouses.NewStoreWarehouseService(storeWarehouses.NewStoreWarehouseRepository(dbHandler.DB), logger.GetZapSugaredLogger()), nil
		},
		storeWarehouses.NewStoreWarehouseHandler,
		apiRouter.RegisterStoreWarehouseRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (stockMaterial.StockMaterialService, error) {
			return stockMaterial.NewStockMaterialService(stockMaterial.NewStockMaterialRepository(dbHandler.DB)), nil
		},
		stockMaterial.NewStockMaterialHandler,
		apiRouter.RegisterStockMaterialRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (barcode.BarcodeService, error) {
			return barcode.NewBarcodeService(
				barcode.NewBarcodeRepository(dbHandler.DB),
				stockMaterial.NewStockMaterialRepository(dbHandler.DB),
				barcode.NewPrinterService()), nil
		},
		barcode.NewBarcodeHandler,
		apiRouter.RegisterBarcodeRouter,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (inventory.InventoryService, error) {
			return inventory.NewInventoryService(
				inventory.NewInventoryRepository(dbHandler.DB),
				stockMaterial.NewStockMaterialRepository(dbHandler.DB),
				barcode.NewBarcodeRepository(dbHandler.DB),
				inventory.NewPackageRepository(dbHandler.DB)), nil
		},
		inventory.NewInventoryHandler,
		apiRouter.RegisterInventoryRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (warehouse.WarehouseService, error) {
			return warehouse.NewWarehouseService(warehouse.NewWarehouseRepository(dbHandler.DB)), nil
		},
		warehouse.NewWarehouseHandler,
		apiRouter.RegisterWarehouseRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (stockRequests.StockRequestService, error) {
			return stockRequests.NewStockRequestService(stockRequests.NewStockRequestRepository(dbHandler.DB)), nil
		},
		stockRequests.NewStockRequestHandler,
		apiRouter.RegisterStockRequestRoutes,
	)

	return router
}

// Temporary init: for testing purposes
func InitializeStorage(cfg *config.Config) storage.StorageRepository {
	storageRepo, err := storage.NewStorageRepository(
		cfg.S3.Endpoint,
		cfg.S3.AccessKey,
		cfg.S3.SecretKey,
		cfg.S3.BucketName,
	)
	if err != nil {
		log.Fatalf("Failed to initialize storage repository: %v", err)
	}
	return storageRepo
}

func InitializeApp() (*gin.Engine, *config.Config) {
	cfg := InitializeConfig()

	err := logger.InitLoggers("info", "logs/gin.log", "logs/service.log", cfg.IsDevelopment)
	if err != nil {
		panic(err)
	}

	dbHandler := InitializeDatabase(cfg)

	redisClient := InitializeRedis(cfg)

	storageRepo := InitializeStorage(cfg) // temp

	// kafkaManager := InitializeKafka(cfg)

	router := InitializeRouter(dbHandler, redisClient, storageRepo) // temp

	return router, cfg
}
