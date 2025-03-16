package init

import (
	"fmt"
	"log"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/container"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/limiters"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth/employeeToken"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/censor"

	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// Prometheus imports
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	promHttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	promHttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latencies in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"method", "endpoint"},
	)
)

func InitializePrometheus() {
	prometheus.MustRegister(promHttpRequestsTotal)
	prometheus.MustRegister(promHttpRequestDuration)
}

func InitializeConfig() *config.Config {
	cfg := config.LoadConfig()
	return cfg
}

func InitializeDatabase(cfg *config.Config) *database.DBHandler {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSL_Mode,
	)

	dbHandler, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	return dbHandler
}

func InitializeRedis(cfg *config.Config) *database.RedisClient {
	redisClient, err := database.InitRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB, cfg.Redis.Username)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	utils.InitCache(redisClient.Client, redisClient.Ctx)

	return redisClient
}

func InitializeRouter(dbHandler *database.DBHandler, redisClient *database.RedisClient, storageRepo storage.StorageRepository) *gin.Engine {
	cfg := config.GetConfig()

	gin.SetMode(cfg.GinMode)

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

	utils.InitValidators()
	router.Use(middleware.SanitizeMiddleware())

	router.Use(middleware.RedisMiddleware(redisClient.Client))

	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// Record request count
		promHttpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			fmt.Sprintf("%d", c.Writer.Status()),
		).Inc()

		// Record request duration
		promHttpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(time.Since(start).Seconds())
	})

	apiRouter := routes.NewRouter(router, "/api", "/v1")

	employeeTokenManager := employeeToken.NewEmployeeTokenManager(dbHandler.DB)
	apiRouter.EmployeeRoutes.Use(middleware.EmployeeAuth(employeeTokenManager))

	storageHandler := storage.NewStorageHandler(storageRepo)                // temp
	storage.RegisterStorageRoutes(apiRouter.EmployeeRoutes, storageHandler) // temp

	appContainer := container.NewContainer(dbHandler, redisClient, &storageRepo, &employeeTokenManager, apiRouter, logger.GetZapSugaredLogger())
	appContainer.MustInitModules()

	router.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	return router
}

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

	err := logger.InitLogger("error", "logs/application.log", cfg.IsDevelopment)
	if err != nil {
		panic(err)
	}

	if err := utils.InitBarcodeFont(); err != nil {
		log.Fatalf("Failed to initialize Barcode font: %v", err)
	}

	if err := localization.InitLocalizer(nil); err != nil {
		logger.GetZapSugaredLogger().Fatalf("Failed to initialize localizer: %v", err)
	}

	if err := censor.InitCensor(); err != nil {
		logger.GetZapSugaredLogger().Fatalf("Failed to initialize censor: %v", err)
	}

	dbHandler := InitializeDatabase(cfg)

	redisClient := InitializeRedis(cfg)

	storageRepo := InitializeStorage(cfg)

	router := InitializeRouter(dbHandler, redisClient, storageRepo)

	return router, cfg
}
