package tests

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/container"
	"github.com/Global-Optima/zeep-web/backend/internal/container/modules"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	testContainer *container.Container
	once          sync.Once
)

func NewTestContainer() *container.Container {
	once.Do(func() {
		var cfg *config.Config
		var err error

		cfg, err = config.LoadTestConfig()
		if err != nil {
			log.Println("failed to load test configuration from file, trying to load from env...")
			cfg, err = LoadConfigFromEnv()
			if err != nil {
				log.Fatalf("Failed to load test configuration: %v", err)
			}
		}

		if err := logger.InitLoggers("debug", "logs/test_gin.log", "logs/test_service.log", cfg.IsDevelopment); err != nil {
			log.Fatalf("Failed to initialize test loggers: %v", err)
		}
		sugarLog := logger.GetZapSugaredLogger()

		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
		)
		dbHandler, err := database.InitDB(dsn)
		if err != nil {
			sugarLog.Fatalf("Failed to initialize test database: %v", err)
		}

		r := gin.New()
		r.Use(logger.ZapLoggerMiddleware())
		r.Use(gin.Recovery())

		apiRouter := routes.NewRouter(r, "/api", "/test")

		testContainer = container.NewContainer(dbHandler, apiRouter, sugarLog)
		testContainer.MustInitModules()

		time.Sleep(100 * time.Millisecond)
	})
	return testContainer
}

func GetTestDB() *gorm.DB {
	return NewTestContainer().GetDB()
}

func GetOrdersModule() *modules.OrdersModule {
	return NewTestContainer().Orders
}

func GetStockRequestsModule() *modules.StockRequestsModule {
	return NewTestContainer().StockRequests
}

func GetStoreStocksModule() *modules.StoreStockModule {
	return NewTestContainer().StoreStocks
}
