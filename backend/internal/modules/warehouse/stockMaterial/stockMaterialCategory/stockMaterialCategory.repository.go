package stockMaterialCategory

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialCategory/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type StockMaterialCategoryRepository interface {
	Create(category *data.StockMaterialCategory) error
	GetByID(id uint) (*data.StockMaterialCategory, error)
	GetAll(filter types.StockMaterialCategoryFilter) ([]data.StockMaterialCategory, error)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, moduleErrors.ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *stockMaterialCategoryRepository) GetAll(filter types.StockMaterialCategoryFilter) ([]data.StockMaterialCategory, error) {
	var categories []data.StockMaterialCategory

	query := r.db.Model(&data.StockMaterialCategory{})

	if filter.Search != nil && *filter.Search != "" {
		search := "%" + *filter.Search + "%"
		query = query.Where("LOWER(name) ILIKE ? OR LOWER(description) ILIKE ?", search, search)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StockMaterialCategory{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	err = query.Find(&categories).Error
	return categories, err
}

func (r *stockMaterialCategoryRepository) Update(id uint, updates data.StockMaterialCategory) error {
	return r.db.Model(&data.StockMaterialCategory{}).Where("id = ?", id).Updates(updates).Error
}

func (r *stockMaterialCategoryRepository) Delete(id uint) error {
	return r.db.Delete(&data.StockMaterialCategory{}, id).Error
}
