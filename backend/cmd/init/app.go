package init

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/limiters"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/censor"
	"go.uber.org/zap"
	"log"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/container"

	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
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

func InitializeRouter(dbHandler *database.DBHandler, redisClient *database.RedisClient, storageRepo storage.StorageRepository) *gin.Engine {
	cfg := config.GetConfig()

	router := gin.New()
	router.Use(logger.ZapLoggerMiddleware())
	router.Use(limiters.LimitRequestBody(30 * 1024 * 1024))
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.Server.ClientURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(middleware.RedisMiddleware(redisClient.Client))

	apiRouter := routes.NewRouter(router, "/api", "/v1")

	//apiRouter.CustomerRoutes.Use(middleware.CustomerAuth())
	apiRouter.EmployeeRoutes.Use(middleware.EmployeeAuth())

	storageHandler := storage.NewStorageHandler(storageRepo)                // temp
	storage.RegisterStorageRoutes(apiRouter.EmployeeRoutes, storageHandler) // temp

	appContainer := container.NewContainer(dbHandler, &storageRepo, apiRouter, logger.GetZapSugaredLogger())
	appContainer.MustInitModules()

	return router
}

// Temporary init: for testing purposes
func InitializeStorage(cfg *config.Config, logger *zap.SugaredLogger) storage.StorageRepository {
	storageRepo, err := storage.NewStorageRepository(
		cfg.S3.Endpoint,
		cfg.S3.AccessKey,
		cfg.S3.SecretKey,
		cfg.S3.BucketName,
		logger,
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

	if err := localization.InitLocalizer(); err != nil {
		logger.GetZapSugaredLogger().Fatalf("Failed to initialize localizer: %v", err)
	}

	if err := censor.InitializeCensor(); err != nil {
		logger.GetZapSugaredLogger().Fatalf("Failed to initialize censor: %v", err)
	}

	dbHandler := InitializeDatabase(cfg)

	redisClient := InitializeRedis(cfg)

	storageRepo := InitializeStorage(cfg, logger.GetZapSugaredLogger()) // temp

	// kafkaManager := InitializeKafka(cfg)

	router := InitializeRouter(dbHandler, redisClient, storageRepo) // temp

	return router, cfg
}
