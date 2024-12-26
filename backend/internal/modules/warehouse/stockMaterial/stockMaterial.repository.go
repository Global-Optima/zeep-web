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
	CreateSupplierMaterial(supplierMaterial *data.SupplierMaterial) error
	UpdateStockMaterial(stockMaterial *data.StockMaterial) error
	UpdateStockMaterialFields(stockMaterialID uint, fields types.UpdateStockMaterialDTO) (*data.StockMaterial, error)
	DeleteStockMaterial(stockMaterialID uint) error
	DeactivateStockMaterial(stockMaterialID uint) error
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
		Preload("Package")

	query = query.Where("is_active = ?", true)

	if filter != nil {
		if filter.Search != nil && *filter.Search != "" {
			search := "%" + *filter.Search + "%"
			query = query.Where("(name ILIKE ? OR category ILIKE ?)", search, search)
		}

		if filter.LowStock != nil && *filter.LowStock {
			query = query.Where("quantity < safety_stock")
		}

		if filter.ExpirationFlag != nil {
			query = query.Where("expiration_flag = ?", *filter.ExpirationFlag)
		}

		if filter.IsActive != nil {
			query = query.Where("is_active = ?", *filter.IsActive)
		}
	}

	var err error
	query, err = utils.ApplyPagination(query, filter.Pagination, &data.StockMaterial{})
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
	err := r.db.Preload("Unit").Preload("Package").First(&stockMaterial, stockMaterialID).Error
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

func (r *stockMaterialRepository) CreateSupplierMaterial(supplierMaterial *data.SupplierMaterial) error {
	return r.db.Create(supplierMaterial).Error
}

func (r *stockMaterialRepository) UpdateStockMaterial(stockMaterial *data.StockMaterial) error {
	return r.db.Save(stockMaterial).Error
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

	if err := r.db.Preload("Unit").Preload("Package").First(&stockMaterial, stockMaterialID).Error; err != nil {
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
