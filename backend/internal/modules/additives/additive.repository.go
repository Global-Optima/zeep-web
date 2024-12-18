package additives

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type AdditiveRepository interface {
	GetAdditiveByID(additiveID uint) (*data.Additive, error)
	GetAdditives(filter types.AdditiveFilterQuery) ([]data.Additive, error)
	CreateAdditive(additive *data.Additive) error
	UpdateAdditive(additive *data.Additive) error
	DeleteAdditive(additiveID uint) error

	GetAdditiveCategories(filter types.AdditiveCategoriesFilterQuery) ([]data.AdditiveCategory, error)
	CreateAdditiveCategory(category *data.AdditiveCategory) error
	UpdateAdditiveCategory(category *data.AdditiveCategory) error
	DeleteAdditiveCategory(categoryID uint) error
	GetAdditiveCategoryByID(categoryID uint) (*data.AdditiveCategory, error)
}

type additiveRepository struct {
	db *gorm.DB
}

func NewAdditiveRepository(db *gorm.DB) AdditiveRepository {
	return &additiveRepository{db: db}
}

func (r *additiveRepository) GetAdditiveCategories(filter types.AdditiveCategoriesFilterQuery) ([]data.AdditiveCategory, error) {
	var categories []data.AdditiveCategory

	query := r.db.Preload("Additives")

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"additive_categories.name LIKE ? OR additive_categories.description LIKE ?",
			searchTerm, searchTerm,
		)
	}

	if filter.ProductSizeId != nil {
		query = query.Preload("Additives", "EXISTS (SELECT 1 FROM product_additives WHERE product_additives.additive_id = additives.id AND product_additives.product_size_id = ?)", *filter.ProductSizeId)
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *additiveRepository) GetAdditives(filter types.AdditiveFilterQuery) ([]data.Additive, error) {
	var additives []data.Additive

	query := r.db.
		Preload("Category").
		Joins("JOIN additive_categories ON additives.additive_category_id = additive_categories.id")

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

	var err error
	query, err = utils.ApplyPagination(query, filter.Pagination, &data.Additive{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&additives).Error; err != nil {
		return nil, err
	}

	return additives, nil
}

func (r *additiveRepository) GetAdditiveByID(additiveID uint) (*data.Additive, error) {
	var additive data.Additive
	err := r.db.Where("id = ?", additiveID).First(&additive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("additive with ID %d not found", additiveID)
		}
		return nil, fmt.Errorf("failed to fetch additive with ID %d: %w", additiveID, err)
	}
	return &additive, nil
}

func (r *additiveRepository) CreateAdditive(additive *data.Additive) error {
	return r.db.Create(additive).Error
}

func (r *additiveRepository) UpdateAdditive(additive *data.Additive) error {
	return r.db.Save(additive).Error
}

func (r *additiveRepository) DeleteAdditive(additiveID uint) error {
	return r.db.Where("id = ?", additiveID).Delete(&data.Additive{}).Error
}

func (r *additiveRepository) CreateAdditiveCategory(category *data.AdditiveCategory) error {
	return r.db.Create(category).Error
}

func (r *additiveRepository) UpdateAdditiveCategory(category *data.AdditiveCategory) error {
	return r.db.Save(category).Error
}

func (r *additiveRepository) DeleteAdditiveCategory(categoryID uint) error {
	return r.db.Where("id = ?", categoryID).Delete(&data.AdditiveCategory{}).Error
}

func (r *additiveRepository) GetAdditiveCategoryByID(categoryID uint) (*data.AdditiveCategory, error) {
	var category data.AdditiveCategory
	err := r.db.Preload("Additives").First(&category, categoryID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}
