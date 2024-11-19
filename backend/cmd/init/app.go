package init

import (
	"fmt"
	"log"

	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitializeConfig() *config.Config {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	utils.InitTTLFromConfig(cfg)
	return cfg
}

func InitializeDatabase(cfg *config.Config) *database.DBHandler {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	dbHandler, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	return dbHandler
}

func InitializeRedis(cfg *config.Config) *database.RedisClient {
	redisClient, err := database.InitRedis(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	utils.InitCache(redisClient.Client, redisClient.Ctx)

	return redisClient
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

func InitializeRouter(dbHandler *database.DBHandler, redisClient *database.RedisClient, storageRepo storage.StorageRepository) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(middleware.RedisMiddleware(redisClient.Client))

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

	return router
}

// Temporary init: for testing purposes
func InitializeStorage(cfg *config.Config) storage.StorageRepository {
	storageRepo, err := storage.NewStorageRepository(
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3BucketName,
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

	router := InitializeRouter(dbHandler, redisClient, storageRepo) // temp
	return router, cfg
}
