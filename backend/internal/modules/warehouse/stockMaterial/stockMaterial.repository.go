package stockMaterial

import (
	"errors"
	"fmt"

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

	IsBarcodeExists(barcode string) (bool, error)
	GetStockMaterialByBarcode(barcode string) (*data.StockMaterial, error)
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
		Preload("Ingredient.Unit")

	if filter == nil {
		filter = &types.StockMaterialFilter{}
	}

	// Apply filters only if they exist
	if filter.Search != nil && *filter.Search != "" {
		search := "%" + *filter.Search + "%"
		query = query.Joins("JOIN stock_material_categories ON stock_material_categories.id = stock_materials.category_id").
			Where("(stock_materials.name ILIKE ? OR stock_materials.description ILIKE ? OR stock_material_categories.name ILIKE ? OR stock_materials.barcode ILIKE ?)", search, search, search, search)
	}

	if filter.LowStock != nil && *filter.LowStock {
		query = query.Where("quantity < safety_stock")
	}

	if filter.IsActive != nil {
		query = query.Where("stock_materials.is_active = ?", *filter.IsActive)
	} else {
		query = query.Where("stock_materials.is_active = ?", true) // Default to active materials
	}

	if filter.IngredientID != nil {
		query = query.Where("stock_materials.ingredient_id = ?", *filter.IngredientID)
	}

	if filter.CategoryID != nil {
		query = query.Where("stock_materials.category_id = ?", *filter.CategoryID)
	}

	if filter.SupplierID != nil {
		query = query.Joins("JOIN supplier_materials ON supplier_materials.stock_material_id = stock_materials.id").
			Where("supplier_materials.supplier_id = ?", *filter.SupplierID)
	}

	if filter.ExpirationInDays != nil {
		query = query.Where("stock_materials.expiration_period_in_days <= ?", *filter.ExpirationInDays)
	}

	query = query.Order("stock_materials.created_at DESC")

	// Apply pagination only if provided
	if filter.Pagination != nil {
		var err error
		query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StockMaterial{})
		if err != nil {
			return nil, err
		}
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
		Preload("Ingredient.Unit").
		First(&stockMaterial, stockMaterialID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrStockMaterialNotFound
		}
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
	var existingBarcode data.StockMaterial
	err := r.db.Model(&data.StockMaterial{}).Where("barcode = ?", stockMaterial.Barcode).First(&existingBarcode).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := r.db.Create(stockMaterial).Error; err != nil {
			return err
		}

		return r.db.Preload("Ingredient").Preload("Unit").Preload("StockMaterialCategory").
			First(stockMaterial, stockMaterial.ID).Error
	}

	return fmt.Errorf("barcode %s already exists", stockMaterial.Barcode)
}

func (r *stockMaterialRepository) CreateStockMaterials(stockMaterials []data.StockMaterial) error {
	return r.db.Create(&stockMaterials).Error
}

func (r *stockMaterialRepository) UpdateStockMaterial(id uint, stockMaterial *data.StockMaterial) error {
	return r.db.Model(&data.StockMaterial{}).Where("id = ?", id).Save(stockMaterial).Error
}

func (r *stockMaterialRepository) UpdateStockMaterialFields(stockMaterialID uint, fields types.UpdateStockMaterialDTO) (*data.StockMaterial, error) {
	var stockMaterial data.StockMaterial

	// Check if the stock material exists
	if err := r.db.Select("id").Where("id = ?", stockMaterialID).First(&stockMaterial).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrStockMaterialNotFound
		}
		return nil, err
	}

	// Create a partial StockMaterial struct with only the fields to update
	updateModel := data.StockMaterial{}

	// Only set fields that are provided in the DTO
	if fields.Name != nil {
		updateModel.Name = *fields.Name
	}
	if fields.Description != nil {
		updateModel.Description = *fields.Description
	}
	if fields.SafetyStock != nil {
		updateModel.SafetyStock = *fields.SafetyStock
	}
	if fields.UnitID != nil {
		updateModel.UnitID = *fields.UnitID
	}
	if fields.CategoryID != nil {
		updateModel.CategoryID = *fields.CategoryID
	}
	if fields.IngredientID != nil {
		updateModel.IngredientID = *fields.IngredientID
	}
	if fields.Barcode != nil {
		updateModel.Barcode = *fields.Barcode
	}
	if fields.ExpirationPeriodInDays != nil {
		updateModel.ExpirationPeriodInDays = *fields.ExpirationPeriodInDays
	}
	if fields.IsActive != nil {
		updateModel.IsActive = *fields.IsActive
	}
	if fields.Size != nil {
		updateModel.Size = *fields.Size
	}

	// Perform the update using the partial struct
	if err := r.db.Model(&stockMaterial).Updates(updateModel).Error; err != nil {
		return nil, err
	}

	// Reload the updated stock material with necessary preloads
	if err := r.db.Preload("Unit").
		Preload("StockMaterialCategory").
		Preload("Ingredient").
		Preload("Ingredient.IngredientCategory").
		Preload("Ingredient.Unit").
		First(&stockMaterial, stockMaterialID).Error; err != nil {
		return nil, err
	}

	return &stockMaterial, nil
}

func (r *stockMaterialRepository) DeleteStockMaterial(stockMaterialID uint) error {
	var stockMaterial data.StockMaterial

	if err := r.db.First(&stockMaterial, stockMaterialID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.ErrStockMaterialNotFound
		}
		return err
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := checkStockMaterialReferences(tx, stockMaterialID); err != nil {
			return err
		}
		if err := tx.Delete(&data.StockMaterial{}, stockMaterialID).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func checkStockMaterialReferences(db *gorm.DB, stockMaterialID uint) error {
	var warehouseStock data.WarehouseStock
	if err := db.Where("stock_material_id = ?", stockMaterialID).First(&warehouseStock).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	} else {
		return types.ErrStockMaterialInUse
	}

	var sri data.StockRequestIngredient
	if err := db.Where("stock_material_id = ?", stockMaterialID).First(&sri).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	} else {
		return types.ErrStockMaterialInUse
	}

	var supplierMaterial data.SupplierMaterial
	if err := db.Where("stock_material_id = ?", stockMaterialID).First(&supplierMaterial).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	} else {
		return types.ErrStockMaterialInUse
	}

	var swdm data.SupplierWarehouseDeliveryMaterial
	if err := db.Where("stock_material_id = ?", stockMaterialID).First(&swdm).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	} else {
		return types.ErrStockMaterialInUse
	}

	return nil
}

func (r *stockMaterialRepository) DeactivateStockMaterial(stockMaterialID uint) error {
	var stockMaterial data.StockMaterial
	err := r.db.Select("is_active").Where("id = ?", stockMaterialID).First(&stockMaterial).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("stock material with id %d not found", stockMaterialID)
		}
		return err
	}

	if !stockMaterial.IsActive {
		return nil // Already inactive, no need to update
	}

	result := r.db.Model(&data.StockMaterial{}).Where("id = ?", stockMaterialID).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *stockMaterialRepository) PopulateStockMaterial(stockMaterialID uint, stockMaterial *data.StockMaterial) error {
	return r.db.Preload("Ingredient").Preload("StockMaterialCategory").First(stockMaterial, "id = ?", stockMaterialID).Error
}

func (r *stockMaterialRepository) IsBarcodeExists(barcode string) (bool, error) {
	var stockMaterial data.StockMaterial
	err := r.db.Where("barcode = ?", barcode).First(&stockMaterial).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *stockMaterialRepository) GetStockMaterialByBarcode(barcode string) (*data.StockMaterial, error) {
	var stockMaterial data.StockMaterial
	err := r.db.
		Preload("Unit").
		Preload("Ingredient").
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory").
		Preload("StockMaterialCategory").
		Where("barcode = ?", barcode).
		First(&stockMaterial).Error
	if err != nil {
		return nil, err
	}
	return &stockMaterial, nil
}
