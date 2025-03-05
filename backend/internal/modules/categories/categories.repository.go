package categories

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategories(filter *types.ProductCategoriesFilterDTO) ([]data.ProductCategory, error)
	GetCategoryByID(id uint) (*data.ProductCategory, error)
	CreateCategory(category *data.ProductCategory) (uint, error)
	UpdateCategory(id uint, category *data.ProductCategory) error
	DeleteCategory(id uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetCategories(filter *types.ProductCategoriesFilterDTO) ([]data.ProductCategory, error) {
	var categories []data.ProductCategory
	query := r.db.Model(&data.ProductCategory{})

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.ProductCategory{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) GetCategoryByID(id uint) (*data.ProductCategory, error) {
	var category data.ProductCategory
	if err := r.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) CreateCategory(category *data.ProductCategory) (uint, error) {
	if err := r.db.Create(category).Error; err != nil {
		return 0, err
	}
	return category.ID, nil
}

func (r *categoryRepository) UpdateCategory(id uint, category *data.ProductCategory) error {
	return r.db.Model(&data.ProductCategory{}).
		Where("id = ?", id).
		Updates(category).Error
}

func (r *categoryRepository) DeleteCategory(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := checkCategoryReferences(tx, id); err != nil {
			return err
		}

		if err := tx.Delete(&data.ProductCategory{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func checkCategoryReferences(db *gorm.DB, categoryID uint) error {
	var category data.ProductCategory

	err := db.
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Where(&data.ProductCategory{BaseEntity: data.BaseEntity{ID: categoryID}}).
		First(&category).Error
	if err != nil {
		return err
	}

	if len(category.Products) > 0 {
		return types.ErrCategoryIsInUse
	}

	return nil
}
