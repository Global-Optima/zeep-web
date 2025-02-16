package tests

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

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
