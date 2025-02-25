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

const (
	defaultDBPort        = 5432
	defaultDBHost        = "localhost"
	defaultDBUser        = "defaultuser"
	defaultDBPassword    = "defaultpassword"
	defaultDBName        = "defaultdb"
	defaultRedisPort     = 6379
	defaultRedisHost     = "localhost"
	defaultRedisDB       = 0
	defaultRedisPassword = ""
	testDataFileName     = "CREATE_TEST_DATA.sql"
	scriptsDir           = "scripts"
	maxHostLength        = 253
)

func LoadConfigFromEnv() (*config.Config, error) {
	cfg := &config.Config{
		Database: loadDBConfig(),
		Redis:    loadRedisConfig(),
	}
	return cfg, nil
}

func loadDBConfig() config.DatabaseConfig {
	port, err := strconv.Atoi(getEnv("DB_PORT", strconv.Itoa(defaultDBPort)))
	if err != nil {
		log.Printf("Invalid DB_PORT value; using default %d. Error: %v", defaultDBPort, err)
		port = defaultDBPort
	}

	return config.DatabaseConfig{
		Host:     getEnv("DB_HOST", defaultDBHost),
		Port:     port,
		User:     getEnv("DB_USER", defaultDBUser),
		Password: getEnv("DB_PASSWORD", defaultDBPassword),
		Name:     getEnv("DB_NAME", defaultDBName),
	}
}

func loadRedisConfig() config.RedisConfig {
	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", strconv.Itoa(defaultRedisDB)))
	if err != nil {
		log.Printf("Invalid REDIS_DB value; using default %d. Error: %v", defaultRedisDB, err)
		redisDB = defaultRedisDB
	}

	redisPort, err := strconv.Atoi(getEnv("REDIS_PORT", strconv.Itoa(defaultRedisPort)))
	if err != nil {
		log.Printf("Invalid REDIS_PORT value; using default %d. Error: %v", defaultRedisPort, err)
		redisPort = defaultRedisPort
	}

	redisHost := getEnv("REDIS_HOST", defaultRedisHost)
	redisHost = validateRedisHost(redisHost)

	return config.RedisConfig{
		Host:     redisHost,
		Port:     redisPort,
		Password: getEnv("REDIS_PASSWORD", defaultRedisPassword),
		DB:       redisDB,
	}
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
		log.Printf("REDIS_HOST is empty; using default '%s'", defaultRedisHost)
		return defaultRedisHost
	}

	if net.ParseIP(redisHost) != nil {
		return redisHost
	}

	if strings.Contains(redisHost, ":") {
		hostPart := strings.Split(redisHost, ":")[0]
		if hostPart == "" || len(hostPart) > maxHostLength {
			log.Printf("Invalid REDIS_HOST value '%s'; using '%s'", redisHost, defaultRedisHost)
			return defaultRedisHost
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
	sqlPath := filepath.Join(baseDir, "..", scriptsDir, testDataFileName)
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
		Sort: &utils.Sort{
			Direction: utils.SORT_DEFAULT_QUERY,
		},
	}
}

func StringPtr(s string) *string {
	return &s
}

func UintPtr(u uint) *uint {
	return &u
}

func FloatPtr(f float64) *float64 {
	return &f
}
