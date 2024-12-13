package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func GenerateUPCBarcode(sku data.StockMaterial, supplierID uint) (string, error) {
	manufacturerCode := GenerateManufacturerCodeFromSupplierID(supplierID)

	productCode := fmt.Sprintf("%05d", sku.ID)

	baseCode := fmt.Sprintf("0%s%s", manufacturerCode, productCode)

	checkDigit := CalculateUPCCheckDigit(baseCode)

	fullBarcode := baseCode + strconv.Itoa(checkDigit)

	return fullBarcode, nil
}

func GenerateManufacturerCodeFromSupplierID(supplierID uint) string {
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%d", supplierID)))
	hash := hex.EncodeToString(hasher.Sum(nil))

	numericCode, _ := strconv.Atoi(hash[:5])
	return fmt.Sprintf("%05d", numericCode%100000)
}

func CalculateUPCCheckDigit(code string) int {
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
