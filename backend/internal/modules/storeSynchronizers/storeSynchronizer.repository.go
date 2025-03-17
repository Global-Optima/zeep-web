package storeSynchronizers

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StoreSynchronizeRepository interface {
	GetNotSynchronizedAdditiveIngredientsIDs(storeID uint, lastSync time.Time) ([]uint, error)
	GetNotSynchronizedProductSizeWithAdditivesIngredients(storeID uint, lastSync time.Time) ([]uint, error)
	GetNotSynchronizedProductSizesAdditivesIDs(storeID uint, lastSync time.Time) ([]uint, error)
}

type storeSynchronizeRepository struct {
	db *gorm.DB
}

func NewStoreSynchronizeRepository(db *gorm.DB) StoreSynchronizeRepository {
	return &storeSynchronizeRepository{db: db}
}

func (r *storeSynchronizeRepository) CloneWithTransaction(tx *gorm.DB) StoreSynchronizeRepository {
	return &storeSynchronizeRepository{
		db: tx,
	}
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
		Where("additive_ingredients.created_at > ? OR product_size_additives.created_at > ?", lastSync, lastSync).
		Pluck("additive_ingredients.ingredient_id", &notSynchronizedProductSizesAdditivesIDs).Error
	if err != nil {
		return nil, err
	}

	return notSynchronizedProductSizesAdditivesIDs, nil
}

func (r *storeSynchronizeRepository) GetNotSynchronizedProductSizeWithAdditivesIngredients(storeID uint, lastSync time.Time) ([]uint, error) {
	var notSynchronizedIngredients []uint

	err := r.db.Model(&data.ProductSizeIngredient{}).
		Distinct("product_size_ingredients.ingredient_id").
		Joins("JOIN product_sizes ON product_sizes.id = product_size_ingredients.product_size_id").
		Joins("JOIN store_product_sizes ON product_sizes.id = store_product_sizes.product_size_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Where("store_products.store_id = ?", storeID).
		Where("product_size_ingredients.created_at > ?", lastSync).
		Pluck("product_size_ingredients.ingredient_id", &notSynchronizedIngredients).Error
	if err != nil {
		return nil, err
	}

	return notSynchronizedIngredients, nil
}

func (r *storeSynchronizeRepository) GetNotSynchronizedProductSizesAdditivesIDs(storeID uint, lastSync time.Time) ([]uint, error) {
	var notSynchronizedProductSizesAdditivesIDs []uint
	err := r.db.Model(&data.ProductSizeAdditive{}).
		Distinct("product_size_additives.additive_id").
		Joins("JOIN product_sizes ON product_sizes.id = product_size_additives.product_size_id").
		Joins("JOIN store_product_sizes ON product_sizes.id = store_product_sizes.product_size_id").
		Joins("JOIN store_products ON store_products.id = store_product_sizes.store_product_id").
		Where("store_products.store_id = ?", storeID).
		Where("product_size_additives.created_at > ?", lastSync).
		Pluck("product_size_additives.additive_id", &notSynchronizedProductSizesAdditivesIDs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []uint{}, nil
		}
		return nil, fmt.Errorf("failed to fetch product sizes: %w", err)
	}

	return notSynchronizedProductSizesAdditivesIDs, nil
}
