package stockMaterial

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type StockMaterialRepository interface {
	GetAllStockMaterials(filter *types.StockMaterialFilter) ([]data.StockMaterial, error)
	GetStockMaterialByID(stockMaterialID uint) (*data.StockMaterial, error)
	GetStockMaterialsByIDs(stockMaterialIDs []uint) ([]data.StockMaterial, error)
	CreateStockMaterial(stockMaterial *data.StockMaterial) error
	CreateStockMaterials(stockMaterials []data.StockMaterial) error
	UpdateStockMaterial(id uint, stockMaterial *data.StockMaterial) error
	UpdateStockMaterialFields(stockMaterialID uint, fields types.UpdateStockMaterialDTO) (*data.StockMaterial, error)
	DeleteStockMaterial(stockMaterialID uint) error
	DeactivateStockMaterial(stockMaterialID uint) error

	PopulateStockMaterial(stockMaterialID uint, stockMaterial *data.StockMaterial) error
}

type stockMaterialRepository struct {
	db *gorm.DB
}

func NewStockMaterialRepository(db *gorm.DB) StockMaterialRepository {
	return &stockMaterialRepository{db: db}
}

func (r *stockMaterialRepository) GetAllStockMaterials(filter *types.StockMaterialFilter) ([]data.StockMaterial, error) {
	var stockMaterials []data.StockMaterial
	query := r.db.Model(&data.StockMaterial{}).
		Preload("Unit").
		Preload("StockMaterialCategory").
		Preload("Ingredient").
		Preload("Ingredient.IngredientCategory").
		Preload("Ingredient.Unit").
		Preload("Packages").Preload("Packages.Unit")

	query = query.Where("is_active = ?", true)

	if filter != nil {
		if filter.Search != nil && *filter.Search != "" {
			search := "%" + *filter.Search + "%"
			query = query.Joins("JOIN stock_material_categories ON stock_material_categories.id = stock_materials.category_id").
				Where("(stock_materials.name ILIKE ? OR stock_materials.description ILIKE ? OR stock_material_categories.name ILIKE ?)", search, search, search)
		}

		if filter.LowStock != nil && *filter.LowStock {
			query = query.Where("quantity < safety_stock")
		}

		if filter.IsActive != nil {
			query = query.Where("is_active = ?", *filter.IsActive)
		}

		if filter.IngredientID != nil {
			query = query.Where("ingredient_id = ?", *filter.IngredientID)
		}

		if filter.CategoryID != nil {
			query = query.Where("category_id = ?", *filter.CategoryID)
		}

		if filter.ExpirationInDays != nil {
			query = query.Where("expiration_period_in_days <= ?", *filter.ExpirationInDays)
		}
	}

	query = query.Order("created_at DESC")

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StockMaterial{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&stockMaterials).Error; err != nil {
		return nil, err
	}

	return stockMaterials, nil
}

func (r *stockMaterialRepository) GetStockMaterialByID(stockMaterialID uint) (*data.StockMaterial, error) {
	var stockMaterial data.StockMaterial
	err := r.db.Preload("Unit").
		Preload("StockMaterialCategory").
		Preload("Ingredient").
		Preload("Ingredient.IngredientCategory").
		Preload("Ingredient.Unit").Preload("Packages").Preload("Packages.Unit").
		First(&stockMaterial, stockMaterialID).Error
	if err != nil {
		return nil, err
	}
	return &stockMaterial, nil
}

func (r *stockMaterialRepository) GetStockMaterialsByIDs(stockMaterialIDs []uint) ([]data.StockMaterial, error) {
	var stockMaterials []data.StockMaterial
	if len(stockMaterialIDs) == 0 {
		return stockMaterials, nil // Return an empty slice if no IDs are provided
	}

	err := r.db.Where("id IN ?", stockMaterialIDs).Find(&stockMaterials).Error
	if err != nil {
		return nil, err
	}

	return stockMaterials, nil
}

func (r *stockMaterialRepository) CreateStockMaterial(stockMaterial *data.StockMaterial) error {
	return r.db.Create(stockMaterial).Error
}

func (r *stockMaterialRepository) CreateStockMaterials(stockMaterials []data.StockMaterial) error {
	return r.db.Create(&stockMaterials).Error
}

func (r *stockMaterialRepository) UpdateStockMaterial(id uint, stockMaterial *data.StockMaterial) error {
	return r.db.Model(&data.StockMaterial{}).Where("id = ?", id).Updates(stockMaterial).Error
}

func (r *stockMaterialRepository) UpdateStockMaterialFields(stockMaterialID uint, fields types.UpdateStockMaterialDTO) (*data.StockMaterial, error) {
	var stockMaterial data.StockMaterial

	if err := r.db.Preload("Unit").First(&stockMaterial, stockMaterialID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("StockMaterial not found")
		}
		return nil, err
	}

	if err := r.db.Model(&stockMaterial).Updates(fields).Error; err != nil {
		return nil, err
	}

	if err := r.db.Preload("Unit").Preload("Packages").First(&stockMaterial, stockMaterialID).Error; err != nil {
		return nil, err
	}

	return &stockMaterial, nil
}

func (r *stockMaterialRepository) DeleteStockMaterial(stockMaterialID uint) error {
	return r.db.Delete(&data.StockMaterial{}, stockMaterialID).Error
}

func (r *stockMaterialRepository) DeactivateStockMaterial(stockMaterialID uint) error {
	return r.db.Model(&data.StockMaterial{}).Where("id = ?", stockMaterialID).Update("is_active", false).Error
}

func (r *stockMaterialRepository) PopulateStockMaterial(stockMaterialID uint, stockMaterial *data.StockMaterial) error {
	return r.db.Preload("Ingredient").Preload("StockMaterialCategory").First(stockMaterial, "id = ?", stockMaterialID).Error
}
