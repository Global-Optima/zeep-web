package main

import (
	"fmt"
	"log"

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

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

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

	productRepo := product.NewProductRepository(dbHandler.DB)
	productService := product.NewProductService(productRepo)
	productHandler := product.NewProductHandler(productService)

	storesRepo := stores.NewStoreRepository(dbHandler.DB)
	storesService := stores.NewStoreService(storesRepo)
	storesHandler := stores.NewStoreHandler(storesService)

	productCategoriesRepo := categories.NewCategoryRepository(dbHandler.DB)
	productCategoriesService := categories.NewCategoryService(productCategoriesRepo)
	productCategoriesHandler := categories.NewCategoryHandler(productCategoriesService)

	additivesRepo := additives.NewAdditiveRepository(dbHandler.DB)
	additivesService := additives.NewAdditiveService(additivesRepo)
	additivesHandler := additives.NewAdditiveHandler(additivesService)

	router := gin.Default()

	router.Use(cors.Default())

	apiRouter := routes.NewRouter(router, "/api", "/v1")
	apiRouter.RegisterProductRoutes(productHandler)
	apiRouter.RegisterStoresRoutes(storesHandler)
	apiRouter.RegisterProductCategoriesRoutes(productCategoriesHandler)
	apiRouter.RegisterAdditivesRoutes(additivesHandler)

	port := cfg.ServerPort
	log.Printf("Starting server on port %d...", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
