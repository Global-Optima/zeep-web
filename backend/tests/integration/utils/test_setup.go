package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/database"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
	"github.com/gin-gonic/gin"
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

	cfg, err := config.LoadConfig(envFilePath)
	if err != nil {
		log.Fatalf("Failed to load testing config! Details: %s", err)
	}
	return cfg
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
