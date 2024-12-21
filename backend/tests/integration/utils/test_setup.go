package utils

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestEnvironment struct {
	Router *gin.Engine
	DB     *gorm.DB
	Config *config.Config
}

func NewTestEnvironment(t *testing.T) *TestEnvironment {
	cfg := loadConfig()
	db := setupDatabase(cfg, t)
	setupRedis(cfg, t)
	router := setupRouter(db)

	truncateAndLoadMockData(db)

	return &TestEnvironment{
		Router: router,
		DB:     db,
		Config: cfg,
	}
}

func loadConfig() *config.Config {
	_, b, _, _ := runtime.Caller(0)
	baseDir := filepath.Join(filepath.Dir(b), "../../..")
	envFilePath := filepath.Join(baseDir, "tests", ".env")

	if _, err := os.Stat(envFilePath); err == nil {
		err := godotenv.Load(envFilePath)
		if err != nil {
			log.Fatalf("Failed to load .env file! Details: %s", err)
		}
		log.Println("Loaded configuration from .env file")
	} else {
		log.Println(".env file not found. Loading configuration from environment variables")
	}

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("Failed to load configuration from environment variables! Details: %s", err)
	}
	return cfg
}

func LoadConfigFromEnv() (*config.Config, error) {
	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Printf("Invalid DB_PORT value; using default 5432. Error: %v", err)
		port = 5432
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		log.Printf("Invalid REDIS_DB value; using default 0. Error: %v", err)
		redisDB = 0
	}

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisHost = validateRedisHost(redisHost)

	redisPort, err := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	if err != nil {
		log.Printf("Invalid REDIS_PORT value; using default 6379. Error: %v", err)
		redisPort = 6379
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     port,
			User:     getEnv("DB_USER", "defaultuser"),
			Password: getEnv("DB_PASSWORD", "defaultpassword"),
			Name:     getEnv("DB_NAME", "defaultdb"),
		},

		Redis: config.RedisConfig{
			Host:     redisHost,
			Port:     redisPort,
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Environment variable %s not set, using default: %s", key, defaultValue)
		return defaultValue
	}
	return value
}

func validateRedisHost(redisHost string) string {
	if redisHost == "" {
		log.Println("REDIS_HOST is empty; using default 'localhost'")
		return "localhost"
	}

	if net.ParseIP(redisHost) != nil {
		return redisHost
	}

	if strings.Contains(redisHost, ":") {
		hostPart := strings.Split(redisHost, ":")[0]
		if hostPart == "" || len(hostPart) > 253 {
			log.Printf("Invalid REDIS_HOST value '%s'; using 'localhost'", redisHost)
			return "localhost"
		}
	}

	return redisHost
}

func setupDatabase(cfg *config.Config, t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize test database! Details: %v", err)
	}
	return db
}

func setupRedis(cfg *config.Config, t *testing.T) *database.RedisClient {
	redisClient, err := database.InitRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		t.Fatalf("Failed to initialize Redis: %v", err)
	}

	utils.InitCache(redisClient.Client, redisClient.Ctx)

	return redisClient
}

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	apiRouter := routes.NewRouter(router, "/api", "/test")

	dbHandler := &database.DBHandler{DB: db}
	apiRouter.RegisterProductRoutes(product.NewProductHandler(product.NewProductService(product.NewProductRepository(dbHandler.DB), logger.GetZapSugaredLogger())))
	apiRouter.RegisterStoresRoutes(stores.NewStoreHandler(stores.NewStoreService(stores.NewStoreRepository(dbHandler.DB))))
	apiRouter.RegisterProductCategoriesRoutes(categories.NewCategoryHandler(categories.NewCategoryService(categories.NewCategoryRepository(dbHandler.DB))))
	apiRouter.RegisterAdditivesRoutes(additives.NewAdditiveHandler(additives.NewAdditiveService(additives.NewAdditiveRepository(dbHandler.DB), logger.GetZapSugaredLogger())))

	return router
}

func truncateAndLoadMockData(db *gorm.DB) {
	truncateTables(db)
	loadMockData(db)
}

func truncateTables(db *gorm.DB) {
	var tables []string
	db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables)
	for _, table := range tables {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table))
	}
}

func loadMockData(db *gorm.DB) {
	_, b, _, _ := runtime.Caller(0)
	baseDir := filepath.Join(filepath.Dir(b), "../../..")
	sqlFilePath := filepath.Join(baseDir, "scripts", "CREATE_TEST_DATA.sql")

	sqlContent, err := os.ReadFile(sqlFilePath)
	if err != nil {
		fmt.Printf("Failed to read mock data file: %v\n", err)
		return
	}

	if err := db.Exec(string(sqlContent)).Error; err != nil {
		fmt.Printf("Error executing mock data SQL: %v\n", err)
	}
}

func (env *TestEnvironment) Close() {
	truncateTables(env.DB)
	mockDB, _ := env.DB.DB()
	if err := mockDB.Close(); err != nil {
		fmt.Printf("Error closing the database connection: %v\n", err)
	}
}
