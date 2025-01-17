package stockMaterialCategory

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type StockMaterialCategoryRepository interface {
	Create(category *data.StockMaterialCategory) error
	GetByID(id uint) (*data.StockMaterialCategory, error)
	GetAll() ([]data.StockMaterialCategory, error)
	Update(id uint, updates data.StockMaterialCategory) error
	Delete(id uint) error
}

type stockMaterialCategoryRepository struct {
	db *gorm.DB
}

func NewStockMaterialCategoryRepository(db *gorm.DB) StockMaterialCategoryRepository {
	return &stockMaterialCategoryRepository{db: db}
}

func (r *stockMaterialCategoryRepository) Create(category *data.StockMaterialCategory) error {
	return r.db.Create(category).Error
}

func (r *stockMaterialCategoryRepository) GetByID(id uint) (*data.StockMaterialCategory, error) {
	var category data.StockMaterialCategory
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *stockMaterialCategoryRepository) GetAll() ([]data.StockMaterialCategory, error) {
	var categories []data.StockMaterialCategory
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *stockMaterialCategoryRepository) Update(id uint, updates data.StockMaterialCategory) error {
	return r.db.Model(&data.StockMaterialCategory{}).Where("id = ?", id).Updates(updates).Error
}

func (r *stockMaterialCategoryRepository) Delete(id uint) error {
	return r.db.Delete(&data.StockMaterialCategory{}, id).Error
}
