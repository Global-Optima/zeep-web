package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

type DBHandler struct {
	DB *gorm.DB
}

func InitDB(url string) DBHandler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalf("error occured while opening db conn: %s", err)
	}

	log.Println("database connected: ", url)

	err = db.Migrator().CreateTable(
	// &addresses.CustomerAddress{},
	// &addresses.FacilityAddress{},
	// &products.Additive{},
	// &products.AdditiveCategory{},
	// &products.Category{},
	// &products.Category{},
	// &products.Product{},
	// &products.ProductAdditive{},
	// &products.ProductSize{},
	// &products.RecipeStep{},
	// TODO: add tables to create
	)

	if err != nil {
		log.Fatalf("error occurred while migration: %s", err)
	}

	log.Println("migrations done")

	dbConn = db
	return DBHandler{db}
}

func GetDBHandler() *DBHandler {
	if dbConn == nil {
		log.Fatal("database connection is not initialized")
	}

	return &DBHandler{DB: dbConn}
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
