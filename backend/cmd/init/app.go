package init

import (
	"fmt"
	"log"
	"time"

	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/kafka"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/supplier"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/inventory"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/sku"
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
	"github.com/Global-Optima/zeep-web/backend/internal/websockets"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
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

func InitializeKafka(cfg *config.Config) *kafka.KafkaManager {
	kafkaManager, err := kafka.NewKafkaManager(cfg.Kafka)
	if err != nil {
		log.Fatalf("failed to initialize kafka instance: %v", err)
		return nil
	}
	return kafkaManager
}

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

func InitializeWebsocket(router *gin.Engine) *websockets.WebSocketHub {
	hub := websockets.GetHubInstance()
	wsGroup := router.Group("/ws")

	ordersHandler := websockets.OrdersWebSocketHandler(hub)
	websockets.RegisterOrderWebsocketRoutes(wsGroup, ordersHandler)

	return hub
}

func InitializeRouter(dbHandler *database.DBHandler, redisClient *database.RedisClient, kafkaManager *kafka.KafkaManager, storageRepo storage.StorageRepository) *gin.Engine {
	cfg := config.GetConfig()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.Server.ClientURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(middleware.RedisMiddleware(redisClient.Client))
	hub := InitializeWebsocket(router)

	apiRouter := routes.NewRouter(router, "/api", "/v1")

	storageHandler := storage.NewStorageHandler(storageRepo)        // temp
	storage.RegisterStorageRoutes(apiRouter.Routes, storageHandler) // temp

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (product.ProductService, error) {
			return product.NewProductService(product.NewProductRepository(dbHandler.DB)), nil
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
			return additives.NewAdditiveService(additives.NewAdditiveRepository(dbHandler.DB)), nil
		},
		additives.NewAdditiveHandler,
		apiRouter.RegisterAdditivesRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (employees.EmployeeService, error) {
			return employees.NewEmployeeService(employees.NewEmployeeRepository(dbHandler.DB)), nil
		},
		employees.NewEmployeeHandler,
		apiRouter.RegisterEmployeesRoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (orders.OrderService, error) {
			return orders.NewOrderService(
				orders.NewOrderRepository(dbHandler.DB),
				orders.NewSubOrderRepository(dbHandler.DB),
				product.NewProductRepository(dbHandler.DB),
				additives.NewAdditiveRepository(dbHandler.DB),
				kafkaManager,
				orders.NewOrdersNotifier(hub)), nil
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
		func(dbHandler *database.DBHandler) (sku.SKUService, error) {
			return sku.NewSKUService(sku.NewSKURepository(dbHandler.DB)), nil
		},
		sku.NewSKUHandler,
		apiRouter.RegisterSKURoutes,
	)

	InitializeModule(
		dbHandler,
		func(dbHandler *database.DBHandler) (barcode.BarcodeService, error) {
			return barcode.NewBarcodeService(
				barcode.NewBarcodeRepository(dbHandler.DB),
				sku.NewSKURepository(dbHandler.DB),
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
				sku.NewSKURepository(dbHandler.DB),
				barcode.NewBarcodeRepository(dbHandler.DB)), nil
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

	dbHandler := InitializeDatabase(cfg)

	redisClient := InitializeRedis(cfg)

	storageRepo := InitializeStorage(cfg) // temp

	kafkaManager := InitializeKafka(cfg)

	router := InitializeRouter(dbHandler, redisClient, kafkaManager, storageRepo) // temp

	return router, cfg
}
