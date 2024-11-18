package init

import (
	"fmt"
	"log"

	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
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

func InitializeRouter(dbHandler *database.DBHandler, storageRepo storage.StorageRepository) *gin.Engine {

	router := gin.Default()
	router.Use(cors.Default())

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

	storageRepo := InitializeStorage(cfg) // temp

	router := InitializeRouter(dbHandler, storageRepo) // temp
	return router, cfg
}
