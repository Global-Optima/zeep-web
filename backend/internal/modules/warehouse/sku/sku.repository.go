package sku

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/sku/types"
	"gorm.io/gorm"
)

type SKURepository interface {
	GetAllSKUs(filter *types.SKUFilter) ([]data.SKU, error)
	GetSKUByID(skuID uint) (*data.SKU, error)
	CreateSKU(sku *data.SKU) error
	UpdateSKU(sku *data.SKU) error
	DeleteSKU(skuID uint) error
	DeactivateSKU(skuID uint) error
}

type skuRepository struct {
	db *gorm.DB
}

func NewSKURepository(db *gorm.DB) SKURepository {
	return &skuRepository{db: db}
}

func (r *skuRepository) GetAllSKUs(filter *types.SKUFilter) ([]data.SKU, error) {
	var skus []data.SKU
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
		// By default, retrieve only active SKUs
		query = query.Where("is_active = ?", true)
	}

	err := query.Find(&skus).Error
	if err != nil {
		return nil, err
	}

	return skus, nil
}

func (r *skuRepository) GetSKUByID(skuID uint) (*data.SKU, error) {
	var sku data.SKU
	err := r.db.Preload("Unit").Preload("Package").First(&sku, skuID).Error
	if err != nil {
		return nil, err
	}
	return &sku, nil
}

func (r *skuRepository) CreateSKU(sku *data.SKU) error {
	return r.db.Create(sku).Error
}

func (r *skuRepository) UpdateSKU(sku *data.SKU) error {
	return r.db.Save(sku).Error
}

func (r *skuRepository) DeleteSKU(skuID uint) error {
	return r.db.Delete(&data.SKU{}, skuID).Error
}

func (r *skuRepository) DeactivateSKU(skuID uint) error {
	return r.db.Model(&data.SKU{}).Where("id = ?", skuID).Update("is_active", false).Error
}
