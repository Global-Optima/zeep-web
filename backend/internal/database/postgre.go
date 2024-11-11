package database

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	postgresMigration "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBHandler struct {
	DB *gorm.DB
}

// InitDB initializes the database connection and applies migrations
func InitDB(dsn string) (*DBHandler, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := applyMigrations(db, "migrations"); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Database connected and migrations applied successfully.")
	return &DBHandler{DB: db}, nil
}

// applyMigrations applies database migrations from the specified path
func applyMigrations(db *gorm.DB, migrationsPath string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sqlDB from gorm DB: %w", err)
	}

	// Optional: Configure connection pooling
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(time.Hour)

	driver, err := postgresMigration.WithInstance(sqlDB, &postgresMigration.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}

	log.Println("Starting migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration error: %w", err)
	}
	log.Println("Migrations completed successfully.")
	return nil
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
