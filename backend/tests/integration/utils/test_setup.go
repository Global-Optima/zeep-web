package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
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

	cfg := &config.Config{
		DBUser:     getEnv("DB_USER", "defaultuser"),
		DBPassword: getEnv("DB_PASSWORD", "defaultpassword"),
		DBName:     getEnv("DB_NAME", "defaultdb"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     port,
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

func setupDatabase(cfg *config.Config, t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize test database! Details: %v", err)
	}
	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	apiRouter := routes.NewRouter(router, "/api", "/test")

	dbHandler := &database.DBHandler{DB: db}
	apiRouter.RegisterProductRoutes(product.NewProductHandler(product.NewProductService(product.NewProductRepository(dbHandler.DB))))
	apiRouter.RegisterStoresRoutes(stores.NewStoreHandler(stores.NewStoreService(stores.NewStoreRepository(dbHandler.DB))))
	apiRouter.RegisterProductCategoriesRoutes(categories.NewCategoryHandler(categories.NewCategoryService(categories.NewCategoryRepository(dbHandler.DB))))
	apiRouter.RegisterAdditivesRoutes(additives.NewAdditiveHandler(additives.NewAdditiveService(additives.NewAdditiveRepository(dbHandler.DB))))

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
