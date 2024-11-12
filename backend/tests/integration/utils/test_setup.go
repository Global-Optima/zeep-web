package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TestDBConfig TestingDatabaseConfig

type TestingDatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func LoadTestConfig() {
	_, b, _, _ := runtime.Caller(0)
	baseDir := filepath.Join(filepath.Dir(b), "../../..")

	configPath := filepath.Join(baseDir, "tests", "config.test.yml")
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading test config file: %v", err)
	}

	if err := viper.Sub("testing-database").Unmarshal(&TestDBConfig); err != nil {
		log.Fatalf("Unable to decode test database config into struct: %v", err)
	}
}

func LoadSchema(db *gorm.DB, schemaFile string) error {
	_, b, _, _ := runtime.Caller(0)
	baseDir := filepath.Join(filepath.Dir(b), "../../..")

	schemaPath := filepath.Join(baseDir, schemaFile)

	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	if err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error; err != nil {
		return fmt.Errorf("failed to reset schema: %w", err)
	}

	if err := db.Exec(string(schema)).Error; err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

func SetupTestDatabase(t *testing.T) *gorm.DB {
	LoadTestConfig()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		TestDBConfig.Host, TestDBConfig.Port, TestDBConfig.User, TestDBConfig.Password, TestDBConfig.Name,
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize GORM: %v", err)
	}

	if err := LoadSchema(gormDB, "scripts/CREATE_TABLES.sql"); err != nil {
		t.Fatalf("Failed to load schema: %v", err)
	}

	return gormDB
}

func SetupProductRepository(t *testing.T) (product.ProductRepository, *gorm.DB) {
	db := SetupTestDatabase(t)
	return product.NewProductRepository(db), db
}

func SetupProductService(t *testing.T) (product.ProductService, *gorm.DB) {
	repo, db := SetupProductRepository(t)
	return product.NewProductService(repo), db
}

func SetupProductHandler(t *testing.T) (*product.ProductHandler, *gorm.DB) {
	service, db := SetupProductService(t)
	return product.NewProductHandler(service), db
}

func SetupTestRouter(handler *product.ProductHandler) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/stores/:store_id/products", handler.GetStoreProducts)
		api.GET("/stores/:store_id/products/search", handler.SearchStoreProducts)
		api.GET("/stores/:store_id/products/:product_id", handler.GetStoreProductDetails)
	}
	return router
}
