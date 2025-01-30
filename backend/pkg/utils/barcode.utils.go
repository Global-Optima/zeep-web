package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strconv"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/golang/freetype/truetype"
	"golang.org/x/exp/rand"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
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

// GenerateRandomEAN13 generates a random 13-digit barcode following EAN-13 standards.
func GenerateRandomEAN13() string {
	rand.Seed(uint64(time.Now().UnixNano()))

	// Generate a random 12-digit base
	baseCode := ""
	for i := 0; i < 12; i++ {
		baseCode += strconv.Itoa(rand.Intn(10)) // Random digit 0-9
	}

	// Calculate the check digit
	checkDigit := CalculateEAN13CheckDigit(baseCode)

	// Combine base and check digit
	fullBarcode := baseCode + strconv.Itoa(checkDigit)

	return fullBarcode
}

// CalculateEAN13CheckDigit computes the last digit (checksum) for a valid EAN-13 barcode.
func CalculateEAN13CheckDigit(code string) int {
	if len(code) != 12 {
		panic("EAN-13 base code must be exactly 12 digits")
	}

	total := 0
	for i, r := range code {
		digit := int(r - '0')
		if i%2 == 0 {
			total += digit // Odd-positioned digits (1st, 3rd, 5th, etc.) are weighted as 1
		} else {
			total += digit * 3 // Even-positioned digits (2nd, 4th, etc.) are weighted as 3
		}
	}

	return (10 - (total % 10)) % 10
}

func GenerateBarcodeImage(barcodeData string) ([]byte, error) {
	// Generate barcode
	bcode, err := code128.Encode(barcodeData)
	if err != nil {
		return nil, err
	}

	// Scale barcode
	scaledBcode, err := barcode.Scale(bcode, 150, 70)
	if err != nil {
		return nil, err
	}

	// Create a new blank image (barcode + text area)
	finalImg := image.NewRGBA(image.Rect(0, 0, 150, 90))
	draw.Draw(finalImg, finalImg.Bounds(), image.White, image.Point{}, draw.Src)

	// Draw the barcode on the image
	draw.Draw(finalImg, scaledBcode.Bounds().Add(image.Pt(0, 0)), scaledBcode, image.Point{}, draw.Over)

	fontData, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	fontFace := truetype.NewFace(fontData, &truetype.Options{
		Size: 12,
	})

	// Add text below barcode
	col := color.Black
	d := &font.Drawer{
		Dst:  finalImg,
		Src:  image.NewUniform(col),
		Face: fontFace,
	}

	textBounds := d.MeasureString(barcodeData)
	textWidth := textBounds.Ceil()
	textX := (finalImg.Bounds().Dx() - textWidth) / 2 // Center horizontally
	textY := 85                                       // Position below the barcode

	d.Dot = fixed.Point26_6{
		X: fixed.I(textX),
		Y: fixed.I(textY),
	}

	d.DrawString(barcodeData)

	var buf bytes.Buffer
	err = png.Encode(&buf, finalImg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
