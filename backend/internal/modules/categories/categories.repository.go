package categories

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategories() ([]data.ProductCategory, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetCategories() ([]data.ProductCategory, error) {
	var categories []data.ProductCategory
	err := r.db.
		Order("name ASC").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	print(categories)
	return categories, nil
}
