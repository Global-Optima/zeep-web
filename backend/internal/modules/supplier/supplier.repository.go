package supplier

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/supplier/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type SupplierRepository interface {
	CreateSupplier(supplier *data.Supplier) error
	GetSupplierByID(id uint) (*data.Supplier, error)
	UpdateSupplier(id uint, fields *data.Supplier) error
	DeleteSupplier(id uint) error
	GetAllSuppliers(filter types.SuppliersFilter) ([]data.Supplier, error)
	ExistsByContactPhone(phone string) (bool, error)

	UpsertSupplierMaterials(supplierID uint, newMaterials []data.SupplierMaterial) error
	GetMaterialsBySupplier(supplierID uint) ([]data.SupplierMaterial, error)
}

type supplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &supplierRepository{db}
}

func (r *supplierRepository) CreateSupplier(supplier *data.Supplier) error {
	return r.db.Create(supplier).Error
}

func (r *supplierRepository) GetSupplierByID(id uint) (*data.Supplier, error) {
	var supplier data.Supplier
	err := r.db.First(&supplier, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, types.ErrSupplierNotFound
	}
	return &supplier, err
}

func (r *supplierRepository) UpdateSupplier(id uint, dto *data.Supplier) error {
	if err := r.db.Model(&data.Supplier{}).Where("id = ?", id).Updates(dto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.ErrSupplierNotFound
		}
		return err
	}

	return nil
}

func (r *supplierRepository) DeleteSupplier(id uint) error {
	return r.db.Delete(&data.Supplier{}, id).Error
}

func (r *supplierRepository) GetAllSuppliers(filter types.SuppliersFilter) ([]data.Supplier, error) {
	var suppliers []data.Supplier

	query := r.db.Model(&data.Supplier{})

	if filter.Search != nil && *filter.Search != "" {
		search := "%" + strings.ToLower(*filter.Search) + "%"
		query = query.Where(
			"LOWER(name) ILIKE ? OR LOWER(city) ILIKE ? OR LOWER(address) ILIKE ?",
			search,
			search,
			search,
		)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Supplier{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&suppliers).Error
	return suppliers, err
}

func (r *supplierRepository) ExistsByContactPhone(phone string) (bool, error) {
	var count int64
	err := r.db.Model(&data.Supplier{}).Where("contact_phone = ?", phone).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *supplierRepository) UpsertSupplierMaterials(supplierID uint, newMaterials []data.SupplierMaterial) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existingMaterials []data.SupplierMaterial
		if err := tx.Where("supplier_id = ?", supplierID).Find(&existingMaterials).Error; err != nil {
			return fmt.Errorf("failed to fetch existing materials: %w", err)
		}

		existingMap := make(map[uint]data.SupplierMaterial)
		for _, existing := range existingMaterials {
			existingMap[existing.StockMaterialID] = existing
		}

		processedIDs := make(map[uint]bool)

		for _, newMaterial := range newMaterials {
			if existing, exists := existingMap[newMaterial.StockMaterialID]; exists {
				existing.SupplierID = supplierID
				if err := tx.Save(&existing).Error; err != nil {
					return fmt.Errorf("failed to update existing material: %w", err)
				}

				if err := r.upsertSupplierPrice(tx, existing.ID, newMaterial.SupplierPrices[0]); err != nil {
					return err
				}
			} else {
				newMaterial.SupplierID = supplierID
				if err := tx.Create(&newMaterial).Error; err != nil {
					return fmt.Errorf("failed to create new material: %w", err)
				}

				if err := r.upsertSupplierPrice(tx, newMaterial.ID, newMaterial.SupplierPrices[0]); err != nil {
					return err
				}
			}
			processedIDs[newMaterial.StockMaterialID] = true
		}

		for _, existing := range existingMaterials {
			if !processedIDs[existing.StockMaterialID] {
				if err := tx.Where("supplier_material_id = ?", existing.ID).Delete(&data.SupplierPrice{}).Error; err != nil {
					return fmt.Errorf("failed to delete old prices: %w", err)
				}
				if err := tx.Delete(&existing).Error; err != nil {
					return fmt.Errorf("failed to delete old material: %w", err)
				}
			}
		}

		return nil
	})
}

func (r *supplierRepository) upsertSupplierPrice(tx *gorm.DB, supplierMaterialID uint, price data.SupplierPrice) error {
	var existingPrice data.SupplierPrice
	err := tx.Where("supplier_material_id = ?", supplierMaterialID).
		First(&existingPrice).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			price.SupplierMaterialID = supplierMaterialID
			if err := tx.Create(&price).Error; err != nil {
				return fmt.Errorf("failed to create new price: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to fetch existing price: %w", err)
	}

	existingPrice.BasePrice = price.BasePrice
	if err := tx.Save(&existingPrice).Error; err != nil {
		return fmt.Errorf("failed to update existing price: %w", err)
	}

	return nil
}

func (r *supplierRepository) GetMaterialsBySupplier(supplierID uint) ([]data.SupplierMaterial, error) {
	var materials []data.SupplierMaterial
	err := r.db.Preload("StockMaterial").
		Preload("StockMaterial.Unit").
		Preload("StockMaterial.StockMaterialCategory").
		Preload("StockMaterial.Ingredient").
		Preload("StockMaterial.Ingredient.Unit").
		Preload("StockMaterial.Ingredient.IngredientCategory").
		Preload("SupplierPrices").
		Where("supplier_id = ?", supplierID).
		Find(&materials).Error
	return materials, err
}
