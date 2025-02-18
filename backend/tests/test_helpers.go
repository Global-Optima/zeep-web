package tests

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

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

func LoadTestData(db *gorm.DB) error {
	_, callerFile, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get caller information")
	}

	baseDir := filepath.Dir(callerFile)
	sqlPath := filepath.Join(baseDir, "..", "scripts", "CREATE_TEST_DATA.sql")
	sqlPath = filepath.Clean(sqlPath)

	if _, err := os.Stat(sqlPath); err != nil {
		return fmt.Errorf("SQL file not found at %s: %w", sqlPath, err)
	}

	content, err := os.ReadFile(sqlPath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file %s: %w", sqlPath, err)
	}

	if err := db.Exec(string(content)).Error; err != nil {
		return fmt.Errorf("failed to execute SQL file %s: %w", sqlPath, err)
	}

	log.Printf("Test data inserted successfully from %s", sqlPath)
	return nil
}

func TruncateAllTables(db *gorm.DB) error {
	var tables []string
	if err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to get table names: %w", err)
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table)
		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}
	return nil
}

func BaseFilterWithPagination(page, pageSize int) utils.BaseFilter {
	return utils.BaseFilter{
		Pagination: &utils.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
	}
}
