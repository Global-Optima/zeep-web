package storeInventoryManager

import (
	"gorm.io/gorm"
)

type StoreInventoryManagerRepository interface {
}

type storeInventoryManagerRepository struct {
	db *gorm.DB
}

func NewStoreInventoryManagerRepository(db *gorm.DB) StoreInventoryManagerRepository {
	return &storeInventoryManagerRepository{db: db}
}
