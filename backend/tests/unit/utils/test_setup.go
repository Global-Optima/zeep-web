package utils

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func CreateMockRow(columns []string, data ...interface{}) *sqlmock.Rows {
	rows := sqlmock.NewRows(columns)

	driverValues := make([]driver.Value, len(data))
	for i, v := range data {
		driverValues[i] = v
	}

	rows.AddRow(driverValues...)
	return rows
}

func CreateMockRows(columns []string, data [][]interface{}) *sqlmock.Rows {
	rows := sqlmock.NewRows(columns)
	for _, row := range data {
		driverValues := make([]driver.Value, len(row))
		for i, v := range row {
			driverValues[i] = v
		}
		rows.AddRow(driverValues...)
	}
	return rows
}

func SetupProductRepository(t *testing.T) (product.ProductRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open Gorm DB with sqlmock: %v", err)
	}

	repo := product.NewProductRepository(gormDB)
	return repo, mock
}
