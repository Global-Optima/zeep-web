package sku

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/sku/types"
	"gorm.io/gorm"
)

type SKURepository interface {
	GetAllSKUs(filter *types.SKUFilter) ([]data.StockMaterial, error)
	GetSKUByID(skuID uint) (*data.StockMaterial, error)
	GetSKUsByIDs(skuIDs []uint) ([]data.StockMaterial, error)
	CreateSKU(sku *data.StockMaterial) error
	CreateSKUs(skus []data.StockMaterial) error
	UpdateSKU(sku *data.StockMaterial) error
	UpdateSKUFields(skuID uint, fields map[string]interface{}) (*data.StockMaterial, error)
	DeleteSKU(skuID uint) error
	DeactivateSKU(skuID uint) error
}

type skuRepository struct {
	db *gorm.DB
}

func NewSKURepository(db *gorm.DB) SKURepository {
	return &skuRepository{db: db}
}

func (r *skuRepository) GetAllSKUs(filter *types.SKUFilter) ([]data.StockMaterial, error) {
	var skus []data.StockMaterial
	query := r.db.Preload("Unit").Preload("Package")

	if filter != nil {
		if filter.Name != nil && *filter.Name != "" {
			query = query.Where("name ILIKE ?", "%"+*filter.Name+"%")
		}
		if filter.Category != nil && *filter.Category != "" {
			query = query.Where("category = ?", *filter.Category)
		}
		if filter.LowStock != nil {
			if *filter.LowStock {
				query = query.Where("quantity < safety_stock")
			}
		}
		if filter.ExpirationFlag != nil {
			query = query.Where("expiration_flag = ?", *filter.ExpirationFlag)
		}
		if filter.IsActive != nil {
			query = query.Where("is_active = ?", *filter.IsActive)
		}
	} else {

		query = query.Where("is_active = ?", true)
	}

	err := query.Find(&skus).Error
	if err != nil {
		return nil, err
	}

	return skus, nil
}

func (r *skuRepository) GetSKUByID(skuID uint) (*data.StockMaterial, error) {
	var sku data.StockMaterial
	err := r.db.Preload("Unit").Preload("Package").First(&sku, skuID).Error
	if err != nil {
		return nil, err
	}
	return &sku, nil
}

func (r *skuRepository) GetSKUsByIDs(skuIDs []uint) ([]data.StockMaterial, error) {
	var skus []data.StockMaterial
	if len(skuIDs) == 0 {
		return skus, nil // Return an empty slice if no IDs are provided
	}

	err := r.db.Where("id IN ?", skuIDs).Find(&skus).Error
	if err != nil {
		return nil, err
	}

	return skus, nil
}

func (r *skuRepository) CreateSKU(sku *data.StockMaterial) error {
	return r.db.Create(sku).Error
}

func (r *skuRepository) CreateSKUs(skus []data.StockMaterial) error {
	return r.db.Create(&skus).Error
}

func (r *skuRepository) UpdateSKU(sku *data.StockMaterial) error {
	return r.db.Save(sku).Error
}

func (r *skuRepository) UpdateSKUFields(skuID uint, fields map[string]interface{}) (*data.StockMaterial, error) {
	var sku data.StockMaterial

	if err := r.db.Preload("Unit").First(&sku, skuID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("SKU not found")
		}
		return nil, err
	}

	if err := r.db.Model(&sku).Updates(fields).Error; err != nil {
		return nil, err
	}

	if err := r.db.Preload("Unit").Preload("Package").First(&sku, skuID).Error; err != nil {
		return nil, err
	}

	return &sku, nil
}

func (r *skuRepository) DeleteSKU(skuID uint) error {
	return r.db.Delete(&data.StockMaterial{}, skuID).Error
}

func (r *skuRepository) DeactivateSKU(skuID uint) error {
	return r.db.Model(&data.StockMaterial{}).Where("id = ?", skuID).Update("is_active", false).Error
}
