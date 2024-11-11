package additives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type AdditiveRepository interface {
	GetAdditivesByStoreAndProduct(storeID uint, productID uint) ([]data.AdditiveCategory, error)
}

type additiveRepository struct {
	db *gorm.DB
}

func NewAdditiveRepository(db *gorm.DB) AdditiveRepository {
	return &additiveRepository{db: db}
}

func (r *additiveRepository) GetAdditivesByStoreAndProduct(storeID uint, productID uint) ([]data.AdditiveCategory, error) {
	var additiveCategories []data.AdditiveCategory

	// Step 1: Get all additive IDs associated with the product
	var additiveIDs []uint

	// Get additive IDs from DefaultProductAdditive
	var defaultAdditiveIDs []uint
	err := r.db.Model(&data.DefaultProductAdditive{}).
		Where("product_id = ?", productID).
		Pluck("additive_id", &defaultAdditiveIDs).Error
	if err != nil {
		return nil, err
	}

	// Get additive IDs from ProductAdditive via ProductSize
	var productAdditiveIDs []uint
	err = r.db.Model(&data.ProductAdditive{}).
		Joins("JOIN product_sizes ps ON ps.id = product_additives.product_size_id").
		Where("ps.product_id = ?", productID).
		Pluck("additive_id", &productAdditiveIDs).Error
	if err != nil {
		return nil, err
	}

	// Combine additive IDs
	additiveIDMap := make(map[uint]struct{})
	for _, id := range defaultAdditiveIDs {
		additiveIDMap[id] = struct{}{}
	}
	for _, id := range productAdditiveIDs {
		additiveIDMap[id] = struct{}{}
	}
	for id := range additiveIDMap {
		additiveIDs = append(additiveIDs, id)
	}

	if len(additiveIDs) == 0 {
		// No additives associated with the product
		return []data.AdditiveCategory{}, nil
	}

	// Step 2: Fetch additive categories with additives filtered by store and product
	err = r.db.
		Model(&data.AdditiveCategory{}).
		Preload("Additives", func(db *gorm.DB) *gorm.DB {
			return db.
				Select("additives.*, sa.price as store_price").
				Joins("JOIN store_additives sa ON sa.additive_id = additives.id AND sa.store_id = ?", storeID).
				Where("additives.id IN (?)", additiveIDs)
		}).
		Find(&additiveCategories).Error

	if err != nil {
		return nil, err
	}

	return additiveCategories, nil
}
