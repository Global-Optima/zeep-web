package storeAdditives

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies/types"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type StoreAdditiveRepository interface {
	CreateStoreAdditives(storeAdditives []data.StoreAdditive) ([]uint, error)
	GetStoreAdditiveByID(storeID, storeAdditiveID uint) (*data.StoreAdditive, error)
	GetStoreAdditives(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]data.StoreAdditive, error)
	GetStoreAdditiveCategories(storeID, productSizeID uint, filter *types.StoreAdditiveCategoriesFilter) ([]data.AdditiveCategory, error)
	UpdateStoreAdditive(storeID, storeAdditiveID uint, input *data.StoreAdditive) error
	DeleteStoreAdditive(storeID, storeAdditiveID uint) error
}

type storeAdditiveRepository struct {
	db *gorm.DB
}

func NewStoreAdditiveRepository(db *gorm.DB) StoreAdditiveRepository {
	return &storeAdditiveRepository{db: db}
}

func (r *storeAdditiveRepository) CreateStoreAdditives(storeAdditives []data.StoreAdditive) ([]uint, error) {
	if len(storeAdditives) == 0 {
		return nil, nil
	}

	if err := r.db.Create(&storeAdditives).Error; err != nil {
		return nil, err
	}

	var ids []uint
	for _, sa := range storeAdditives {
		ids = append(ids, sa.ID)
	}

	return ids, nil
}

func (r *storeAdditiveRepository) GetStoreAdditiveCategories(storeID, productSizeID uint, filter *types.StoreAdditiveCategoriesFilter) ([]data.AdditiveCategory, error) {
	var categories []data.AdditiveCategory

	subquery := r.db.Model(&data.Additive{}).
		Select("additive_category_id").
		Joins("JOIN store_additives ON store_additives.additive_id = additives.id").
		Where("store_additives.store_id = ?", storeID).
		Where("EXISTS (SELECT 1 FROM product_size_additives WHERE product_size_additives.additive_id = additives.id AND product_size_additives.product_size_id = ?)", productSizeID).
		Group("additive_category_id").
		Having("COUNT(additives.id) > 0")

	query := r.db.Model(&data.AdditiveCategory{}).
		Preload("Additives", "id IN (SELECT additive_id FROM store_additives WHERE store_id = ?)", storeID).
		Preload("Additives.StoreAdditives", "store_id = ?", storeID).
		Preload("Additives.ProductSizeAdditives", "product_size_id = ?", productSizeID).
		Where("id IN (?)", subquery)

	if filter.IsMultipleSelect != nil && *filter.IsMultipleSelect {
		query = query.Where("product_size_additives.is_default = ?", *filter.IsMultipleSelect)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"name ILIKE ? OR description ILIKE ?",
			searchTerm, searchTerm,
		)
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *storeAdditiveRepository) GetStoreAdditiveByID(storeID, storeAdditiveID uint) (*data.StoreAdditive, error) {
	var storeAdditive *data.StoreAdditive

	err := r.db.Model(&data.StoreAdditive{}).
		Where("store_id = ?", storeID).
		Where("id = ?", storeAdditiveID).
		Preload("Additive.Category").
		First(&storeAdditive).Error

	if err != nil {
		return nil, err
	}

	return storeAdditive, nil
}

func (r *storeAdditiveRepository) GetStoreAdditives(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]data.StoreAdditive, error) {
	var storeAdditives []data.StoreAdditive

	query := r.db.Model(&data.StoreAdditive{}).
		Where("store_id = ?", storeID).
		Preload("Additive.Category")

	var err error

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("additives.name LIKE ? OR additives.description LIKE ? OR additives.size LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if filter.MinPrice != nil {
		query = query.Where("store_additives.price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("store_additives.price <= ?", *filter.MaxPrice)
	}

	if filter.CategoryID != nil {
		query = query.Where("additive_categories.id = ?", *filter.CategoryID)
	}
	if filter.ProductSizeID != nil {
		query = query.Where("product_additives.product_size_id = ?", *filter.ProductSizeID)
	}

	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StoreAdditive{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&storeAdditives).Error; err != nil {
		return nil, err
	}

	return storeAdditives, nil
}

func (r *storeAdditiveRepository) UpdateStoreAdditive(storeID, storeAdditiveID uint, input *data.StoreAdditive) error {
	if input == nil {
		return fmt.Errorf("input cannot be nil")
	}

	err := r.db.Model(&data.StoreAdditive{}).
		Where("store_id = ? AND id = ?", storeID, storeAdditiveID).
		Updates(input).Error

	return err
}

func (r *storeAdditiveRepository) DeleteStoreAdditive(storeID, storeAdditiveID uint) error {

	err := r.db.Where("store_id = ? AND id = ?", storeID, storeAdditiveID).
		Delete(&data.StoreAdditive{}).Error

	return err
}
