package storeAdditives

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies/types"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type StoreAdditiveRepository interface {
	CreateStoreAdditives(storeAdditives []data.StoreAdditive) ([]uint, error)
	GetStoreAdditiveByID(storeID, storeAdditiveID uint) (*data.StoreAdditive, error)
	GetAdditivesListToAdd(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]data.Additive, error)
	GetStoreAdditives(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]data.StoreAdditive, error)
	GetStoreAdditivesByIDs(storeID uint, IDs []uint) ([]data.StoreAdditive, error)
	GetStoreAdditiveCategories(storeID, storeProductSizeID uint, filter *types.StoreAdditiveCategoriesFilter) ([]data.AdditiveCategory, error)
	UpdateStoreAdditive(storeID, storeAdditiveID uint, input *data.StoreAdditive) error
	DeleteStoreAdditive(storeID, storeAdditiveID uint) error

	CloneWithTransaction(tx *gorm.DB) StoreAdditiveRepository
}

type storeAdditiveRepository struct {
	db *gorm.DB
}

func NewStoreAdditiveRepository(db *gorm.DB) StoreAdditiveRepository {
	return &storeAdditiveRepository{db: db}
}

func (r *storeAdditiveRepository) CloneWithTransaction(tx *gorm.DB) StoreAdditiveRepository {
	return &storeAdditiveRepository{db: tx}
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

func (r *storeAdditiveRepository) GetAdditivesListToAdd(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]data.Additive, error) {
	var additives []data.Additive

	query := r.db.
		Preload("Category").
		Preload("Unit").
		Joins("JOIN additive_categories ON additives.additive_category_id = additive_categories.id")

	var err error

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("additives.name LIKE ? OR additives.description LIKE ? OR additives.size LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if filter.MinPrice != nil {
		query = query.Where("additives.base_price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("additives.base_price <= ?", *filter.MaxPrice)
	}

	if filter.CategoryID != nil {
		query = query.Where("additive_categories.id = ?", *filter.CategoryID)
	}

	if filter.ProductSizeID != nil {
		query = query.Where("product_additives.product_size_id = ?", *filter.ProductSizeID)
	}

	query = query.Where("additives.id NOT IN (?)", r.db.Model(&data.StoreAdditive{}).Select("additive_id").Where("store_id = ?", storeID))

	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Additive{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	if err := query.Find(&additives).Error; err != nil {
		return nil, err
	}

	return additives, nil
}

func (r *storeAdditiveRepository) GetStoreAdditiveCategories(
	storeID, storeProductSizeID uint,
	filter *types.StoreAdditiveCategoriesFilter,
) ([]data.AdditiveCategory, error) {

	var categories []data.AdditiveCategory

	subAdditives := r.db.Model(&data.Additive{}).
		Select("additives.id").
		Joins("JOIN product_size_additives ON product_size_additives.additive_id = additives.id").
		Joins("JOIN store_product_sizes ON store_product_sizes.product_size_id = product_size_additives.product_size_id").
		Joins("JOIN store_additives ON store_additives.additive_id = additives.id").
		Where("store_product_sizes.id = ?", storeProductSizeID).
		Where("store_additives.store_id = ?", storeID)

	subCats := r.db.Model(&data.Additive{}).
		Select("additive_category_id").
		Where("id IN (?)", subAdditives).
		Group("additive_category_id").
		Having("COUNT(additives.id) > 0")

	query := r.db.Model(&data.AdditiveCategory{}).
		Where("id IN (?)", subCats).
		Preload("Additives", "id IN (?)", subAdditives).
		Preload("Additives.Unit").
		Preload("Additives.StoreAdditives", "store_id = ?", storeID).
		Preload("Additives.ProductSizeAdditives", func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN store_product_sizes ON store_product_sizes.product_size_id = product_size_additives.product_size_id").
				Where("store_product_sizes.id = ?", storeProductSizeID)
		})

	if filter.IsMultipleSelect != nil {
		query = query.Where("is_multiple_select = ?", *filter.IsMultipleSelect)
	}
	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	if err := query.Find(&categories).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, moduleErrors.ErrNotFound
		}
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
		Preload("Additive.Unit").
		Preload("Additive.Ingredients.Ingredient.Unit").
		Preload("Additive.Ingredients.Ingredient.IngredientCategory").
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
		Joins("JOIN additives ON additives.id = store_additives.additive_id").
		Preload("Additive.Category").
		Preload("Additive.Unit")

	var err error

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("additives.name ILIKE ? OR additives.description ILIKE ?", searchTerm, searchTerm)
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

func (r *storeAdditiveRepository) GetStoreAdditivesByIDs(storeID uint, IDs []uint) ([]data.StoreAdditive, error) {
	var storeAdditives []data.StoreAdditive

	query := r.db.Model(&data.StoreAdditive{}).
		Where("store_id = ?", storeID).
		Preload("Additive.Category").
		Preload("Additive.Unit").
		Where("id IN (?)", IDs)

	if err := query.Find(&storeAdditives).Error; err != nil {
		return nil, err
	}

	return storeAdditives, nil
}

func (r *storeAdditiveRepository) UpdateStoreAdditive(storeID, storeAdditiveID uint, input *data.StoreAdditive) error {
	if input == nil {
		return fmt.Errorf("input cannot be nil")
	}

	res := r.db.Model(&data.StoreAdditive{}).
		Where("store_id = ? AND id = ?", storeID, storeAdditiveID).
		Updates(input)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return types.ErrStoreAdditiveNotFound
	}

	return nil
}

func (r *storeAdditiveRepository) DeleteStoreAdditive(storeID, storeAdditiveID uint) error {

	err := r.db.Where("store_id = ? AND id = ?", storeID, storeAdditiveID).
		Delete(&data.StoreAdditive{}).Error

	return err
}
