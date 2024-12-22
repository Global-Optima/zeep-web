package categories

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategories(filter *types.CategoriesFilterDTO) ([]data.ProductCategory, error)
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

func (r *categoryRepository) GetCategories(filter *types.CategoriesFilterDTO) ([]data.ProductCategory, error) {
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
	if err := r.db.Delete(&data.ProductCategory{}, id).Error; err != nil {
		return err
	}
	return nil
}
