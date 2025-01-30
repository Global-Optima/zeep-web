package barcode

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type BarcodeRepository interface {
	AssignBarcode(barcode string, stockMaterialID uint) error
	GetStockMaterialByBarcode(barcode string) (*data.StockMaterial, error)
	GetStockMaterialBeforeDeletion(stockMaterialID uint) (*data.StockMaterial, error)
	UpdateStockMaterialBarcode(stockMaterial *data.StockMaterial) error

	GetSupplierMaterialByStockMaterialID(stockMaterialID uint) (*data.SupplierMaterial, error)

	GetBarcodesForStockMaterials(stockMaterialIDs []uint) ([]data.StockMaterial, error)
	GetBarcodeForStockMaterial(stockMaterialID uint) (*data.StockMaterial, error)
	IsBarcodeExists(barcode string) (bool, error)
}

type barcodeRepository struct {
	db *gorm.DB
}

func NewBarcodeRepository(db *gorm.DB) BarcodeRepository {
	return &barcodeRepository{db: db}
}

func (r *barcodeRepository) AssignBarcode(barcode string, stockMaterialID uint) error {
	var stockMaterial data.StockMaterial
	err := r.db.First(&stockMaterial, stockMaterialID).Error
	if err != nil {
		return err
	}

	exists, err := r.IsBarcodeExists(barcode)
	if err != nil {
		return fmt.Errorf("failed to check barcode uniqueness: %w", err)
	}
	if exists {
		return errors.New("generated barcode already exists")
	}

	stockMaterial.Barcode = barcode
	err = r.db.Save(&stockMaterial).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *barcodeRepository) GetStockMaterialByBarcode(barcode string) (*data.StockMaterial, error) {
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

func (r *barcodeRepository) GetStockMaterialBeforeDeletion(stockMaterialID uint) (*data.StockMaterial, error) {
	var stockMaterial data.StockMaterial
	err := r.db.Preload("Unit").First(&stockMaterial, stockMaterialID).Error
	if err != nil {
		return nil, err
	}
	return &stockMaterial, nil
}

func (r *barcodeRepository) UpdateStockMaterialBarcode(stockMaterial *data.StockMaterial) error {
	return r.db.Save(stockMaterial).Error
}

func (r *barcodeRepository) IsBarcodeExists(barcode string) (bool, error) {
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

func (r *barcodeRepository) GetSupplierMaterialByStockMaterialID(stockMaterialID uint) (*data.SupplierMaterial, error) {
	var supplierMaterial data.SupplierMaterial
	err := r.db.Where("stock_material_id = ?", stockMaterialID).First(&supplierMaterial).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &supplierMaterial, nil
}

func (r *barcodeRepository) GetBarcodesForStockMaterials(stockMaterialIDs []uint) ([]data.StockMaterial, error) {
	var stockMaterials []data.StockMaterial
	err := r.db.Select("id, barcode").
		Where("id IN ?", stockMaterialIDs).
		Find(&stockMaterials).Error
	if err != nil {
		return nil, err
	}
	return stockMaterials, nil
}

func (r *barcodeRepository) GetBarcodeForStockMaterial(stockMaterialID uint) (*data.StockMaterial, error) {
	var stockMaterial data.StockMaterial
	err := r.db.Select("id, barcode").
		Where("id = ?", stockMaterialID).
		First(&stockMaterial).Error
	if err != nil {
		return nil, err
	}
	return &stockMaterial, nil
}
