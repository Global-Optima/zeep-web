package database

import (
	"fmt"
	"log"
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

var migrationsPath = "migrations"

func InitDB(dsn string) (*DBHandler, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	cfg := config.GetConfig()
	if cfg.IsTest {
		if err := ResetDatabase(db, migrationsPath); err != nil {
			return nil, fmt.Errorf("failed to reset database: %w", err)
		}
	}

	if err := applyMigrations(db, migrationsPath); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB from gorm DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected and migrations applied successfully.")
	return &DBHandler{DB: db}, nil
}

func ResetDatabase(db *gorm.DB, migrationsPath string) error {
	if err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error; err != nil {
		return fmt.Errorf("failed to reset schema: %w", err)
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
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		log.Fatalf("failed to get absolute path of migrations: %v", err)
	}

	absPath = filepath.ToSlash(absPath)

	sourceURL := "file://" + absPath
	if cfg.IsDevelopment {
		sourceURL = "file://" + migrationsPath
	}

	if cfg.IsTest {
		candidatePaths := []string{
			"../../migrations", // two levels up
			"../migrations",    // one level up
			"migrations",       // same directory
		}
		if callerDir, ok := utils.GetCallerDir(2); ok {
			if foundPath := utils.SearchForCandidatePath(callerDir, candidatePaths); foundPath != "" {
				sourceURL = "file://" + foundPath
				log.Printf("Test environment: using migrations folder: %s", foundPath)
			} else {
				log.Printf("Test environment: no candidate migrations folder found; using default: %s", sourceURL)
			}
		} else {
			log.Printf("Test environment: failed to determine caller directory; using default migrations path: %s", sourceURL)
		}
	}

	return sourceURL
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
