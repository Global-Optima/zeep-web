package barcode

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type BarcodeRepository interface {
	GenerateAndAssignBarcode(skuID uint) (string, error)
	GetSKUByBarcode(barcode string) (*data.SKU, error)
	GetSKUBeforeDeletion(skuID uint) (*data.SKU, error)
	UpdateSKUBarcode(sku *data.SKU) error
}

type barcodeRepository struct {
	db *gorm.DB
}

func NewBarcodeRepository(db *gorm.DB) BarcodeRepository {
	return &barcodeRepository{db: db}
}

func (r *barcodeRepository) GenerateAndAssignBarcode(skuID uint) (string, error) {
	var sku data.SKU
	err := r.db.First(&sku, skuID).Error
	if err != nil {
		return "", err
	}

	barcode, err := generateUPCBarcode(sku)
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

func (r *barcodeRepository) GetSKUByBarcode(barcode string) (*data.SKU, error) {
	var sku data.SKU
	err := r.db.Preload("Unit").Preload("Package").Where("barcode = ?", barcode).First(&sku).Error
	if err != nil {
		return nil, err
	}
	return &sku, nil
}

func (r *barcodeRepository) GetSKUBeforeDeletion(skuID uint) (*data.SKU, error) {
	var sku data.SKU
	err := r.db.Preload("Unit").Preload("Package").First(&sku, skuID).Error
	if err != nil {
		return nil, err
	}
	return &sku, nil
}

func (r *barcodeRepository) UpdateSKUBarcode(sku *data.SKU) error {
	return r.db.Save(sku).Error
}

func generateUPCBarcode(sku data.SKU) (string, error) {
	manufacturerCode := generateManufacturerCodeFromSupplierID(sku.SupplierID)

	productCode := fmt.Sprintf("%05d", sku.ID)

	baseCode := fmt.Sprintf("0%s%s", manufacturerCode, productCode)

	checkDigit := calculateUPCCheckDigit(baseCode)

	fullBarcode := baseCode + strconv.Itoa(checkDigit)

	return fullBarcode, nil
}

func generateManufacturerCodeFromSupplierID(supplierID uint) string {
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%d", supplierID)))
	hash := hex.EncodeToString(hasher.Sum(nil))

	numericCode, _ := strconv.Atoi(hash[:5])
	return fmt.Sprintf("%05d", numericCode%100000)
}

func calculateUPCCheckDigit(code string) int {
	if len(code) != 11 {
		panic("UPC base code must be exactly 11 digits")
	}

	total := 0
	for i, r := range code {
		digit := int(r - '0')
		if i%2 == 0 {

			total += digit * 3
		} else {

			total += digit
		}
	}

	checkDigit := (10 - (total % 10)) % 10
	return checkDigit
}

func isBarcodeExists(barcode string, tx *gorm.DB) (bool, error) {
	var sku data.SKU
	err := tx.Where("barcode = ?", barcode).First(&sku).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
