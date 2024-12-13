package barcode

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type BarcodeRepository interface {
	GenerateAndAssignBarcode(skuID uint) (string, error)
	AssignBarcode(barcode string, skuID uint) error
	GetSKUByBarcode(barcode string) (*data.StockMaterial, error)
	GetSKUBeforeDeletion(skuID uint) (*data.StockMaterial, error)
	UpdateSKUBarcode(sku *data.StockMaterial) error

	GetSupplierMaterialByStockMaterialID(stockMaterialID uint) (*data.SupplierMaterial, error)
}

type barcodeRepository struct {
	db *gorm.DB
}

func NewBarcodeRepository(db *gorm.DB) BarcodeRepository {
	return &barcodeRepository{db: db}
}

func (r *barcodeRepository) GenerateAndAssignBarcode(skuID uint) (string, error) {
	var sku data.StockMaterial
	err := r.db.First(&sku, skuID).Error
	if err != nil {
		return "", err
	}

	supplierMaterial, err := r.GetSupplierMaterialByStockMaterialID(skuID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch supplier for SKU: %w", err)
	}
	if supplierMaterial == nil {
		return "", errors.New("supplier not found for the given SKU")
	}

	barcode, err := utils.GenerateUPCBarcode(sku, supplierMaterial.SupplierID)
	if err != nil {
		return "", err
	}

	exists, err := isBarcodeExists(barcode, r.db)
	if err != nil {
		return "", fmt.Errorf("failed to check barcode uniqueness: %w", err)
	}
	if exists {
		return "", errors.New("generated barcode already exists")
	}

	sku.Barcode = barcode
	err = r.db.Save(&sku).Error
	if err != nil {
		return "", err
	}

	return barcode, nil
}

func (r *barcodeRepository) AssignBarcode(barcode string, skuID uint) error {
	var sku data.StockMaterial
	err := r.db.First(&sku, skuID).Error
	if err != nil {
		return err
	}

	exists, err := isBarcodeExists(barcode, r.db)
	if err != nil {
		return fmt.Errorf("failed to check barcode uniqueness: %w", err)
	}
	if exists {
		return errors.New("generated barcode already exists")
	}

	sku.Barcode = barcode
	err = r.db.Save(&sku).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *barcodeRepository) GetSKUByBarcode(barcode string) (*data.StockMaterial, error) {
	var sku data.StockMaterial
	err := r.db.Preload("Unit").Preload("Package").Where("barcode = ?", barcode).First(&sku).Error
	if err != nil {
		return nil, err
	}
	return &sku, nil
}

func (r *barcodeRepository) GetSKUBeforeDeletion(skuID uint) (*data.StockMaterial, error) {
	var sku data.StockMaterial
	err := r.db.Preload("Unit").Preload("Package").First(&sku, skuID).Error
	if err != nil {
		return nil, err
	}
	return &sku, nil
}

func (r *barcodeRepository) UpdateSKUBarcode(sku *data.StockMaterial) error {
	return r.db.Save(sku).Error
}

func isBarcodeExists(barcode string, tx *gorm.DB) (bool, error) {
	var sku data.StockMaterial
	err := tx.Where("barcode = ?", barcode).First(&sku).Error
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
