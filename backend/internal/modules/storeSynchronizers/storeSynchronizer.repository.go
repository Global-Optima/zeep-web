package storeSynchronizers

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type StoreSynchronizeRepository interface {
	GetNotSynchronizedAdditiveIngredientsIDs(storeID uint, lastSync time.Time) ([]uint, error)
	GetNotSynchronizedProductSizesIDs(storeID uint, lastSync time.Time) ([]uint, error)
	GetNotSynchronizedProductSizeIngredientsIDs(storeID uint, lastSync time.Time) ([]uint, error)
	GetNotSynchronizedProductSizesAdditivesIDs(storeID uint, lastSync time.Time) ([]uint, error)
	GetNotSynchronizedAdditivesIDs(storeID uint, lastSync time.Time) ([]uint, error)
	GetNotSynchronizedProductSizesProvisionsIDs(storeID uint, lastSync time.Time) ([]uint, error)
	GetNotSynchronizedAdditiveProvisionIDs(storeID uint, lastSync time.Time) ([]uint, error)
}

type storeSynchronizeRepository struct {
	db *gorm.DB
}

func NewStoreSynchronizeRepository(db *gorm.DB) StoreSynchronizeRepository {
	return &storeSynchronizeRepository{db: db}
}

func (r *storeSynchronizeRepository) GetNotSynchronizedAdditiveIngredientsIDs(storeID uint, lastSync time.Time) ([]uint, error) {
	var notSynchronizedProductSizesAdditivesIDs []uint
	err := r.db.Model(&data.AdditiveIngredient{}).
		Distinct("additive_ingredients.ingredient_id").
		Joins("JOIN additives ON additive_ingredients.additive_id = additives.id").
		Joins("JOIN product_size_additives ON product_size_additives.additive_id = additives.id").
		Joins("JOIN product_sizes ON product_sizes.id = product_size_additives.product_size_id").
		Joins("JOIN store_product_sizes ON store_product_sizes.product_size_id = product_sizes.id").
		Joins("JOIN store_products ON store_product_sizes.store_product_id = store_products.id").
		Where("store_products.store_id = ?", storeID).
		Where("additive_ingredients.created_at > ? OR product_size_additives.created_at > ? OR additive_ingredients.updated_at > ? OR product_size_additives.updated_at > ?",
			lastSync, lastSync, lastSync, lastSync).
		Pluck("additive_ingredients.ingredient_id", &notSynchronizedProductSizesAdditivesIDs).Error
	if err != nil {
		return nil, err
	}

	return notSynchronizedProductSizesAdditivesIDs, nil
}

func (r *storeSynchronizeRepository) GetNotSynchronizedProductSizeIngredientsIDs(storeID uint, lastSync time.Time) ([]uint, error) {
	var notSynchronizedIngredients []uint

	err := r.db.Model(&data.ProductSizeIngredient{}).
		Distinct("product_size_ingredients.ingredient_id").
		Joins("JOIN product_sizes ON product_sizes.id = product_size_ingredients.product_size_id").
		Joins("JOIN store_product_sizes ON product_sizes.id = store_product_sizes.product_size_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Where("store_products.store_id = ?", storeID).
		Where("product_size_ingredients.created_at > ? OR product_size_ingredients.updated_at > ?", lastSync, lastSync).
		Pluck("product_size_ingredients.ingredient_id", &notSynchronizedIngredients).Error
	if err != nil {
		return nil, err
	}

	return notSynchronizedIngredients, nil
}

func (r *storeSynchronizeRepository) GetNotSynchronizedProductSizesIDs(storeID uint, lastSync time.Time) ([]uint, error) {
	var productSizeIDs []uint
	err := r.db.Model(&data.ProductSize{}).
		Distinct("product_sizes.id").
		Joins("JOIN store_product_sizes ON product_sizes.id = store_product_sizes.product_size_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Where("store_products.store_id = ?", storeID).
		Where("additives_updated_at > ? OR provisions_updated_at > ?", lastSync, lastSync).
		Pluck("product_sizes.id", &productSizeIDs).Error
	if err != nil {
		return nil, err
	}

	return productSizeIDs, nil
}

func (r *storeSynchronizeRepository) GetNotSynchronizedProductSizesAdditivesIDs(storeID uint, lastSync time.Time) ([]uint, error) {
	var notSynchronizedProductSizesAdditivesIDs []uint
	err := r.db.Model(&data.ProductSizeAdditive{}).
		Distinct("product_size_additives.additive_id").
		Joins("JOIN product_sizes ON product_sizes.id = product_size_additives.product_size_id").
		Joins("JOIN store_product_sizes ON product_sizes.id = store_product_sizes.product_size_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Where("store_products.store_id = ?", storeID).
		Where("product_size_additives.created_at > ? OR product_size_additives.updated_at > ?",
			lastSync, lastSync).
		Pluck("product_size_additives.additive_id", &notSynchronizedProductSizesAdditivesIDs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product size additives: %w", err)
	}

	return notSynchronizedProductSizesAdditivesIDs, nil
}

func (r *storeSynchronizeRepository) GetNotSynchronizedProductSizesProvisionsIDs(storeID uint, lastSync time.Time) ([]uint, error) {
	var notSynchronizedProductSizesProvisionsIDs []uint
	err := r.db.Model(&data.ProductSizeProvision{}).
		Distinct("product_size_provisions.provision_id").
		Joins("JOIN product_sizes ON product_sizes.id = product_size_provisions.product_size_id").
		Joins("JOIN store_product_sizes ON product_size_provisions.product_size_id = store_product_sizes.product_size_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Where("store_products.store_id = ?", storeID).
		Where("product_size_provisions.created_at > ? OR product_size_provisions.updated_at > ? OR product_sizes.provisions_updated_at > ?",
			lastSync, lastSync, lastSync).
		Pluck("product_size_provisions.provision_id", &notSynchronizedProductSizesProvisionsIDs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product size provisions: %w", err)
	}

	return notSynchronizedProductSizesProvisionsIDs, nil
}

func (r *storeSynchronizeRepository) GetNotSynchronizedAdditivesIDs(storeID uint, lastSync time.Time) ([]uint, error) {
	var notSynchronizedAdditives []uint
	err := r.db.Model(&data.Additive{}).
		Distinct("additives.id").
		Joins("JOIN store_additives ON store_additives.additive_id = additives.id").
		Where("store_additives.store_id = ?", storeID).
		Where("additives.ingredients_updated_at > ? OR additives.provisions_updated_at > ?",
			lastSync, lastSync).
		Pluck("additive.id", &notSynchronizedAdditives).Error
	if err != nil {
		return nil, err
	}

	return notSynchronizedAdditives, nil
}

func (r *storeSynchronizeRepository) GetNotSynchronizedAdditiveProvisionIDs(storeID uint, lastSync time.Time) ([]uint, error) {
	var notSynchronizedAdditiveProvisionsIDs []uint
	err := r.db.Model(&data.AdditiveProvision{}).
		Distinct("additive_provisions.provision_id").
		Joins("JOIN additives ON additive_provisions.additive_id = additives.id").
		Joins("JOIN store_additives ON store_additives.additive_id = additive_provisions.additive_id").
		Where("store_additives.store_id = ?", storeID).
		Where("additive_provisions.created_at > ? OR additive_provisions.updated_at > ? OR additives.provisions_updated_at > ?",
			lastSync, lastSync, lastSync).
		Pluck("additive_provisions.provision_id", &notSynchronizedAdditiveProvisionsIDs).Error
	if err != nil {
		return nil, err
	}

	return notSynchronizedAdditiveProvisionsIDs, nil
}
