package main

import (
	"fmt"
	"log"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadConfig()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Name,
	)

	dbHandler := database.InitDB(dsn)

	productRepo := product.NewProductRepository(dbHandler.DB)
	productService := product.NewProductService(productRepo)
	productHandler := product.NewProductHandler(productService)

	router := gin.Default()

	apiRouter := routes.NewRouter(router, "/api", "/v1")
	apiRouter.RegisterProductRoutes(productHandler)

	port := config.AppConfig.Server.Port
	log.Printf("Starting server on port %d...", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
