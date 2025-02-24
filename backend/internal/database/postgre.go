package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/golang-migrate/migrate/v4"
	postgresMigration "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBHandler struct {
	DB *gorm.DB
}

const (
	maxOpenConns          = 25
	maxIdleConns          = 25
	connMaxLifetime       = time.Hour
	defaultMigrationsPath = "migrations"
)

func InitDB(dsn string) (*DBHandler, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	cfg := config.GetConfig()
	if cfg.IsTest {
		if err := ResetDatabase(db, defaultMigrationsPath); err != nil {
			return nil, fmt.Errorf("failed to reset database: %w", err)
		}
	} else {
		if err := applyMigrations(db, defaultMigrationsPath); err != nil {
			return nil, fmt.Errorf("failed to apply migrations: %w", err)
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB from gorm DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	log.Println("Database connected and migrations applied successfully.")
	return &DBHandler{DB: db}, nil
}

func ResetDatabase(db *gorm.DB, migrationsPath string) error {
	if err := db.Exec("DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;").Error; err != nil {
		return fmt.Errorf("failed to reset schema: %w", err)
	}

	cfg := config.GetConfig()
	if cfg.IsTest {
		sourceURL := determineSourceURL(migrationsPath)
		if sourceURL == "" {
			log.Println("Test environment: no valid migrations folder found; skipping migrations")
			return nil // Skip applying migrations in test mode if none is found.
		}
		migrationsPath = sourceURL[len("file://"):]
	} else {
		if info, err := os.Stat(migrationsPath); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("migrations path does not exist: %s", migrationsPath)
			}
			return fmt.Errorf("failed to check migrations path: %w", err)
		} else if !info.IsDir() {
			return fmt.Errorf("migrations path is not a directory: %s", migrationsPath)
		}
	}

	return applyMigrations(db, migrationsPath)
}

func applyMigrations(db *gorm.DB, migrationsPath string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB from gorm DB: %w", err)
	}

	driver, err := postgresMigration.WithInstance(sqlDB, &postgresMigration.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	sourceURL := determineSourceURL(migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}

	log.Println("Starting migrations...")
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration error: %w", err)
	}
	log.Println("Migrations completed successfully.")
	return nil
}

func determineSourceURL(migrationsPath string) string {
	cfg := config.GetConfig()

	if cfg.IsTest {
		log.Println("Test environment: searching for migrations folder...")

		candidatePaths := []string{
			"../../../" + defaultMigrationsPath, // two levels up
			"../../" + defaultMigrationsPath,    // two levels up
			"../" + defaultMigrationsPath,       // one level up
			"./" + defaultMigrationsPath,        // same directory
		}
		if callerDir, ok := utils.GetCallerDir(2); ok {
			if foundPath := utils.SearchForCandidatePath(callerDir, candidatePaths); foundPath != "" {
				sourceURL := "file://" + foundPath
				log.Printf("Test environment: using migrations folder: %s", foundPath)
				return sourceURL
			}
			log.Printf("Test environment: no candidate migrations folder found; using provided migrations path: %s", migrationsPath)
		} else {
			log.Printf("Test environment: failed to determine caller directory; using provided migrations path: %s", migrationsPath)
		}
	}

	if cfg.IsDevelopment {
		log.Println("Development environment: using provided migrations path")
		return "file://" + migrationsPath
	}

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		log.Printf("failed to get absolute path of migrations: %v", err)
		absPath = migrationsPath // fallback to provided path
	}
	absPath = filepath.ToSlash(absPath)
	return "file://" + absPath
}

type BaseRepository interface {
	GetDB() *gorm.DB
}

type baseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) BaseRepository {
	return &baseRepository{db: db}
}

func (r *baseRepository) GetDB() *gorm.DB {
	return r.db
}
