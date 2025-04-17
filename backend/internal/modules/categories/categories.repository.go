package categories

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategories(locale data.LanguageCode, filter *types.ProductCategoriesFilterDTO) ([]data.ProductCategory, error)
	GetCategoryByID(id uint) (*data.ProductCategory, error)
	GetTranslatedCategoryByID(locale data.LanguageCode, id uint) (*data.ProductCategory, error)
	CreateCategory(category *data.ProductCategory) (uint, error)
	UpdateCategory(id uint, category *data.ProductCategory) error
	DeleteCategory(id uint) error
	FindRawProductCategoryByID(id uint, category *data.ProductCategory) error
	CloneWithTransaction(tx *gorm.DB) CategoryRepository
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CloneWithTransaction(tx *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: tx,
	}
}

func (r *categoryRepository) FindRawProductCategoryByID(id uint, category *data.ProductCategory) error {
	return r.db.Model(&data.ProductCategory{}).
		Where(&data.ProductCategory{BaseEntity: data.BaseEntity{ID: id}}).
		First(category).Error
}

func (r *categoryRepository) GetCategories(locale data.LanguageCode, filter *types.ProductCategoriesFilterDTO) ([]data.ProductCategory, error) {
	var categories []data.ProductCategory

	q := r.db.Model(&data.ProductCategory{})

	if s := filter.Search; s != "" {
		pattern := "%" + s + "%"
		q = q.Where("name ILIKE ? OR description ILIKE ?", pattern, pattern)
	}

	paged, err := utils.ApplySortedPaginationForModel(
		q, filter.Pagination, filter.Sort, &data.ProductCategory{},
	)
	if err != nil {
		return nil, err
	}

	paged = utils.ApplyLocalizedPreloads(
		paged,
		locale,
		types.ProductCategoryPreloadMap,
	)

	if err := paged.Find(&categories).Error; err != nil {
		return nil, err
	}

	if categories == nil {
		categories = []data.ProductCategory{}
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

func (r *categoryRepository) GetTranslatedCategoryByID(locale data.LanguageCode, id uint) (*data.ProductCategory, error) {
	var category data.ProductCategory

	q := r.db.Model(&data.ProductCategory{}).Where("id = ?", id)

	q = utils.ApplyLocalizedPreloads(q, locale, types.ProductCategoryPreloadMap)

	if err := q.First(&category).Error; err != nil {
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
		Save(category).Error
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
