package utils

import (
	"fmt"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func GenerateUPCBarcode(sku data.StockMaterial, supplierID uint) (string, error) {
	manufacturerCode := fmt.Sprintf("%05d", supplierID%100000)

	productCode := fmt.Sprintf("%05d", sku.ID%100000)

	baseCode := fmt.Sprintf("0%s%s", manufacturerCode, productCode)

	checkDigit := CalculateUPCCheckDigit(baseCode)

	fullBarcode := baseCode + strconv.Itoa(checkDigit)

	if len(fullBarcode) != 12 {
		return "", fmt.Errorf("invalid barcode length: %s", fullBarcode)
	}

	return fullBarcode, nil
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

	return (10 - (total % 10)) % 10
}
