package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strconv"
	"sync"
	"time"

	"github.com/jung-kurt/gofpdf"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/golang/freetype/truetype"
	"golang.org/x/exp/rand"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

const (
	labelWidth  = 300
	labelHeight = 120
	marginMM    = 5.0
)

var (
	fontFace     font.Face
	fontInitOnce sync.Once
)

// InitBarcodeFont ensures the font is parsed only once (thread-safe).
func InitBarcodeFont() error {
	var fontError error
	fontInitOnce.Do(func() {
		f, err := truetype.Parse(goregular.TTF)
		if err != nil {
			fontError = fmt.Errorf("failed to parse goregular font: %w", err)
			return
		}
		fontFace = truetype.NewFace(f, &truetype.Options{Size: 12})
	})
	return fontError
}

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

func GenerateBarcodeImage(barcodeData string) (*bytes.Buffer, error) {
	bcode, err := code128.Encode(barcodeData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode barcode data %q: %w", barcodeData, err)
	}

	// 3. Scale barcode to fit exactly within the image dimensions
	//    Suppose we want the final label to be 300px wide x 120px tall
	const labelWidth = 300
	const labelHeight = 120
	scaledBcode, err := barcode.Scale(bcode, labelWidth, labelHeight-20) // leave space for text
	if err != nil {
		return nil, fmt.Errorf("failed to scale barcode: %w", err)
	}

	// 4. Create a blank image with a white background
	finalImg := image.NewRGBA(image.Rect(0, 0, labelWidth, labelHeight))
	draw.Draw(finalImg, finalImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// 5. Draw barcode on the image
	bcBounds := scaledBcode.Bounds()
	draw.Draw(finalImg, bcBounds, scaledBcode, image.Point{}, draw.Over)

	// 6. Draw the text under the barcode
	drawer := &font.Drawer{
		Dst:  finalImg,
		Src:  image.NewUniform(color.Black),
		Face: fontFace, // Make sure fontFace is defined somewhere
	}
	textWidth := drawer.MeasureString(barcodeData).Ceil()
	textX := (labelWidth - textWidth) / 2
	textY := labelHeight - 5

	drawer.Dot = fixed.Point26_6{
		X: fixed.I(textX),
		Y: fixed.I(textY),
	}
	drawer.DrawString(barcodeData)

	// 7. Encode image as PNG in memory
	imgBuf := &bytes.Buffer{}
	if err := png.Encode(imgBuf, finalImg); err != nil {
		return imgBuf, fmt.Errorf("failed to encode final PNG: %w", err)
	}

	return imgBuf, nil
}

// ConvertImageToPDF converts an image to a PDF with custom dimensions.
func ConvertImageToPDF(imgBuf *bytes.Buffer) ([]byte, error) {
	// 8. Convert PNG to PDF
	//    Instead of using "A4", we define a small custom page size:
	//
	//    1 pixel = ~0.264583 mm
	//    labelWidth (300 px)  ~ 79.375 mm
	//    labelHeight (120 px) ~ 31.75 mm
	//
	//    Add a small margin around your label if needed.
	imgWidthMM := float64(labelWidth) * 0.264583
	imgHeightMM := float64(labelHeight) * 0.264583

	// Create a new PDF with custom page size:
	// Use gofpdf.NewCustom so we can specify exact dimensions.
	pdfFile := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size: gofpdf.SizeType{
			Wd: imgWidthMM + (2 * marginMM),
			Ht: imgHeightMM + (2 * marginMM),
		},
		FontDirStr: "",
	})
	pdfFile.SetAutoPageBreak(false, 0)
	pdfFile.AddPage()

	// Register and place the image
	imgOptions := gofpdf.ImageOptions{
		ImageType: "PNG",
		ReadDpi:   true,
	}
	pdfFile.RegisterImageOptionsReader("barcode.png", imgOptions, imgBuf)

	// Place image with the margin as an offset
	pdfFile.ImageOptions(
		"barcode.png",
		marginMM,    // x-pos
		marginMM,    // y-pos
		imgWidthMM,  // image width in mm
		imgHeightMM, // image height in mm
		false,       // flow: false = free-floating
		imgOptions,
		0,
		"",
	)

	// 9. Output PDF as byte slice
	var pdfBuf bytes.Buffer
	if err := pdfFile.Output(&pdfBuf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfBuf.Bytes(), nil
}

func GenerateBarcodePDF(barcodeData string) ([]byte, error) {
	imgBuf, err := GenerateBarcodeImage(barcodeData)
	if err != nil {
		return nil, err
	}
	pdfBuf, err := ConvertImageToPDF(imgBuf)
	if err != nil {
		return nil, err
	}

	return pdfBuf, nil
}
