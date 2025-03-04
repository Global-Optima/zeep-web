package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	mockStorage "github.com/Global-Optima/zeep-web/backend/tests/integration/utils/s3-mock-repository"

	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/container"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestEnvironment struct {
	Router *gin.Engine
	DB     *gorm.DB
	Config *config.Config
	Tokens map[string]string
}

func NewTestEnvironment(t *testing.T) *TestEnvironment {
	cfg := loadConfig()

	if err := logger.InitLogger("debug", "logs/test_application.log", cfg.IsTest); err != nil {
		log.Fatalf("Failed to initialize test loggers: %v", err)
	}

	dbHandler := setupDatabase(cfg, t)
	wd, err := os.Getwd()
	if err != nil {
		logger.GetZapSugaredLogger().Fatalf("Failed to get current working directory: %v", err)
	}
	projectRoot := filepath.Join(wd, "..", "..", "..")
	localizationLanguagesPath := filepath.Join(projectRoot, "internal/localization/languages")
	if err := localization.InitLocalizer(&localizationLanguagesPath); err != nil {
		logger.GetZapSugaredLogger().Fatalf("Failed to initialize localizer: %v", err)
	}

	setupRedis(cfg, t)
	router := setupRouter(dbHandler)

	truncateAndLoadMockData(dbHandler.DB)
	log.Println("Mock data loaded successfully")

	return &TestEnvironment{
		Router: router,
		DB:     dbHandler.DB,
		Config: cfg,
		Tokens: make(map[string]string),
	}
}

func loadConfig() *config.Config {
	var cfg *config.Config

	_, b, _, _ := runtime.Caller(0)
	baseDir := filepath.Join(filepath.Dir(b), "../../../")
	envFilePath := filepath.Join(baseDir, "tests", "test.env")

	if _, err := os.Stat(envFilePath); err == nil {
		cfg, err = config.LoadTestConfig(envFilePath)
		if err != nil {
			log.Fatalf("Failed to load test.env file! Details: %s", err)
		}

		log.Println("Loaded configuration from test.env file")
	} else {
		log.Println("test.env file not found. Loading configuration from environment variables")

		cfg, err = LoadConfigFromEnv()
		if err != nil {
			log.Fatalf("Failed to load configuration from environment variables! Details: %s", err)
		}
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

func setupDatabase(cfg *config.Config, t *testing.T) *database.DBHandler {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name,
	)
	dbHandler, err := database.InitDB(dsn)
	if err != nil {
		t.Fatalf("Failed to initialize test database! Details: %v", err)
	}

	log.Println("Integration tests setup: Database connected")
	return dbHandler
}

func setupRedis(cfg *config.Config, t *testing.T) *database.RedisClient {
	redisClient, err := database.InitRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		t.Fatalf("Failed to initialize Redis: %v", err)
	}

	utils.InitCache(redisClient.Client, redisClient.Ctx)
	return redisClient
}

func setupMockStorage() *storage.StorageRepository {
	storageRepo, err := mockStorage.NewMockStorageRepository()
	if err != nil {
		logger.GetZapSugaredLogger().Fatalf("Failed to initialize mock storage repository: %v", err)
	}
	return &storageRepo
}

func setupRouter(dbHandler *database.DBHandler) *gin.Engine {
	router := gin.New()
	utils.InitValidators()
	router.Use(middleware.SanitizeMiddleware())
	router.Use(logger.ZapLoggerMiddleware())

	apiRouter := routes.NewRouter(router, "/api", "/test")
	apiRouter.EmployeeRoutes.Use(middleware.EmployeeAuth())

	storageRepo := setupMockStorage()

	testContainer := container.NewContainer(dbHandler, storageRepo, apiRouter, logger.GetZapSugaredLogger())
	testContainer.MustInitModules()

	return router
}

func truncateAndLoadMockData(db *gorm.DB) {
	if err := truncateTables(db); err != nil {
		log.Fatalf("Failed to truncate tables: %v", err)
	}
	loadMockData(db)
}

func truncateTables(db *gorm.DB) error {
	var tables []string
	if err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to get table names: %w", err)
	}
	// Build a list of tables to truncate (excluding schema_migrations).
	var toTruncate []string
	for _, table := range tables {
		if table == "schema_migrations" {
			continue
		}
		toTruncate = append(toTruncate, table)
	}
	// Execute a single TRUNCATE statement for all tables
	if len(toTruncate) > 0 {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;",
			strings.Join(toTruncate, ", "))
		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("failed to truncate tables: %w", err)
		}
	}
	return nil
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
	if err := truncateTables(env.DB); err != nil {
		log.Printf("Failed to truncate tables: %v", err)
	}
	mockDB, _ := env.DB.DB()
	if err := mockDB.Close(); err != nil {
		fmt.Printf("Error closing the database connection: %v\n", err)
	}
}
