package barcode

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type BarcodeRepository interface {
	GenerateAndAssignBarcode(stockMaterialID uint) (string, error)
	AssignBarcode(barcode string, stockMaterialID uint) error
	GetStockMaterialByBarcode(barcode string) (*data.StockMaterial, error)
	GetStockMaterialBeforeDeletion(stockMaterialID uint) (*data.StockMaterial, error)
	UpdateStockMaterialBarcode(stockMaterial *data.StockMaterial) error

	GetSupplierMaterialByStockMaterialID(stockMaterialID uint) (*data.SupplierMaterial, error)
}

type barcodeRepository struct {
	db *gorm.DB
}

func NewBarcodeRepository(db *gorm.DB) BarcodeRepository {
	return &barcodeRepository{db: db}
}

func (r *barcodeRepository) GenerateAndAssignBarcode(stockMaterialID uint) (string, error) {
	var stockMaterial data.StockMaterial
	err := r.db.First(&stockMaterial, stockMaterialID).Error
	if err != nil {
		return "", err
	}

	supplierMaterial, err := r.GetSupplierMaterialByStockMaterialID(stockMaterialID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch supplier for StockMaterial: %w", err)
	}
	if supplierMaterial == nil {
		return "", errors.New("supplier not found for the given StockMaterial")
	}

	barcode, err := utils.GenerateUPCBarcode(stockMaterial, supplierMaterial.SupplierID)
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

	stockMaterial.Barcode = barcode
	err = r.db.Save(&stockMaterial).Error
	if err != nil {
		return "", err
	}

	return barcode, nil
}

func (r *barcodeRepository) AssignBarcode(barcode string, stockMaterialID uint) error {
	var stockMaterial data.StockMaterial
	err := r.db.First(&stockMaterial, stockMaterialID).Error
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

	stockMaterial.Barcode = barcode
	err = r.db.Save(&stockMaterial).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *barcodeRepository) GetStockMaterialByBarcode(barcode string) (*data.StockMaterial, error) {
	var stockMaterial data.StockMaterial
	err := r.db.Preload("Unit").Preload("StockMaterialCategory").Preload("Package").Where("barcode = ?", barcode).First(&stockMaterial).Error
	if err != nil {
		return nil, err
	}
	return &stockMaterial, nil
}

func (r *barcodeRepository) GetStockMaterialBeforeDeletion(stockMaterialID uint) (*data.StockMaterial, error) {
	var stockMaterial data.StockMaterial
	err := r.db.Preload("Unit").Preload("Package").First(&stockMaterial, stockMaterialID).Error
	if err != nil {
		return nil, err
	}
	return &stockMaterial, nil
}

func (r *barcodeRepository) UpdateStockMaterialBarcode(stockMaterial *data.StockMaterial) error {
	return r.db.Save(stockMaterial).Error
}

func isBarcodeExists(barcode string, tx *gorm.DB) (bool, error) {
	var stockMaterial data.StockMaterial
	err := tx.Where("barcode = ?", barcode).First(&stockMaterial).Error
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
