package stockMaterial

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"sync"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/golang/freetype/truetype"
	"github.com/jung-kurt/gofpdf"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

type StockMaterialService interface {
	GetAllStockMaterials(filter *types.StockMaterialFilter) ([]types.StockMaterialsDTO, error)
	GetStockMaterialByID(stockMaterialID uint) (*types.StockMaterialsDTO, error)
	CreateStockMaterial(req *types.CreateStockMaterialDTO) (*types.StockMaterialsDTO, error)
	UpdateStockMaterial(stockMaterialID uint, req *types.UpdateStockMaterialDTO) error
	DeleteStockMaterial(stockMaterialID uint) error
	DeactivateStockMaterial(stockMaterialID uint) error
	GetStockMaterialBarcode(stockMaterialID uint) ([]byte, error)
	GenerateStockMaterialBarcodePDF(stockMaterialID uint) ([]byte, error)
	GenerateBarcode() (*types.GenerateBarcodeResponse, error)
	RetrieveStockMaterialByBarcode(barcode string) (*types.StockMaterialsDTO, error)
}

type stockMaterialService struct {
	repo StockMaterialRepository
}

func NewStockMaterialService(repo StockMaterialRepository) StockMaterialService {
	return &stockMaterialService{
		repo: repo,
	}
}

var (
	fontFace      font.Face
	fontInitError error
	fontInitOnce  sync.Once
)

func (s *stockMaterialService) GetAllStockMaterials(filter *types.StockMaterialFilter) ([]types.StockMaterialsDTO, error) {
	stockMaterials, err := s.repo.GetAllStockMaterials(filter)
	if err != nil {
		return nil, err
	}

	stockMaterialResponses := make([]types.StockMaterialsDTO, 0)
	for _, stockMaterial := range stockMaterials {
		stockMaterialResponses = append(stockMaterialResponses, *types.ConvertStockMaterialToStockMaterialResponse(&stockMaterial))
	}

	return stockMaterialResponses, nil
}

func (s *stockMaterialService) GetStockMaterialByID(stockMaterialID uint) (*types.StockMaterialsDTO, error) {
	stockMaterial, err := s.repo.GetStockMaterialByID(stockMaterialID)
	if err != nil {
		return nil, err
	}

	if stockMaterial == nil {
		return nil, errors.New("StockMaterial not found")
	}

	stockMaterialResponse := types.ConvertStockMaterialToStockMaterialResponse(stockMaterial)

	return stockMaterialResponse, nil
}

func (s *stockMaterialService) CreateStockMaterial(req *types.CreateStockMaterialDTO) (*types.StockMaterialsDTO, error) {
	stockMaterial := types.ConvertCreateStockMaterialRequestToStockMaterial(req)

	err := s.repo.CreateStockMaterial(stockMaterial)
	if err != nil {
		return nil, err
	}

	stockMaterialResponse := types.ConvertStockMaterialToStockMaterialResponse(stockMaterial)
	return stockMaterialResponse, nil
}

func (s *stockMaterialService) UpdateStockMaterial(stockMaterialID uint, req *types.UpdateStockMaterialDTO) error {
	stockMaterial, err := s.repo.GetStockMaterialByID(stockMaterialID)
	if err != nil {
		return fmt.Errorf("failed to fetch stock material: %w", err)
	}

	if stockMaterial == nil {
		return fmt.Errorf("stock material with ID %d not found", stockMaterialID)
	}

	updatedStockMaterial, err := types.ValidateAndApplyUpdate(stockMaterial, req)
	if err != nil {
		return err
	}

	err = s.repo.UpdateStockMaterial(stockMaterialID, updatedStockMaterial)
	if err != nil {
		return fmt.Errorf("failed to update stock material: %w", err)
	}

	return nil
}

func (s *stockMaterialService) DeleteStockMaterial(stockMaterialID uint) error {
	return s.repo.DeleteStockMaterial(stockMaterialID)
}

func (s *stockMaterialService) DeactivateStockMaterial(stockMaterialID uint) error {
	return s.repo.DeactivateStockMaterial(stockMaterialID)
}

func (s *stockMaterialService) GetStockMaterialBarcode(stockMaterialID uint) ([]byte, error) {
	stockMaterial, err := s.repo.GetStockMaterialByID(stockMaterialID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock material: %w", err)
	}

	barcodeImage, err := utils.GenerateBarcodeImage(stockMaterial.Barcode)
	if err != nil {
		return nil, fmt.Errorf("failed to generate barcode image: %w", err)
	}

	return barcodeImage, nil
}

func initFont() {
	fontInitOnce.Do(func() {
		f, err := truetype.Parse(goregular.TTF)
		if err != nil {
			fontInitError = fmt.Errorf("failed to parse goregular font: %w", err)
			return
		}

		// Adjust Size as needed. This is a typical 12pt usage.
		fontFace = truetype.NewFace(f, &truetype.Options{
			Size: 12,
		})
	})
}

func (s *stockMaterialService) GenerateStockMaterialBarcodePDF(stockMaterialID uint) ([]byte, error) {
	stockMaterial, err := s.repo.GetStockMaterialByID(stockMaterialID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock material: %w", err)
	}

	// Initialize font first
	initFont()
	if fontInitError != nil {
		return nil, fontInitError
	}

	// 2. Encode the barcode data using Code-128
	barcodeData := stockMaterial.Barcode
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
	var imgBuf bytes.Buffer
	if err := png.Encode(&imgBuf, finalImg); err != nil {
		return nil, fmt.Errorf("failed to encode final PNG: %w", err)
	}

	// 8. Convert PNG to PDF
	//    Instead of using "A4", we define a small custom page size:
	//
	//    1 pixel = ~0.264583 mm
	//    labelWidth (300 px)  ~ 79.375 mm
	//    labelHeight (120 px) ~ 31.75 mm
	//
	//    Add a small margin around your label if needed.
	const marginMM = 5.0
	imgWidthMM := float64(labelWidth) * 0.264583
	imgHeightMM := float64(labelHeight) * 0.264583

	// Create a new PDF with custom page size:
	// Use gofpdf.NewCustom so we can specify exact dimensions.
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size: gofpdf.SizeType{
			Wd: imgWidthMM + (2 * marginMM),
			Ht: imgHeightMM + (2 * marginMM),
		},
		FontDirStr: "",
	})
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	// Register and place the image
	imgOptions := gofpdf.ImageOptions{
		ImageType: "PNG",
		ReadDpi:   true,
	}
	pdf.RegisterImageOptionsReader("barcode.png", imgOptions, &imgBuf)

	// Place image with the margin as an offset
	pdf.ImageOptions(
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
	if err := pdf.Output(&pdfBuf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfBuf.Bytes(), nil
}

func (s *stockMaterialService) GenerateBarcode() (*types.GenerateBarcodeResponse, error) {
	maxRetries := 10
	var barcode string
	var exists bool
	var err error

	for i := 0; i < maxRetries; i++ {
		barcode = utils.GenerateRandomEAN13()

		exists, err = s.repo.IsBarcodeExists(barcode)
		if err != nil {
			return nil, fmt.Errorf("failed to check uniqueness of barcode: %w", err)
		}

		if !exists {
			break
		}
	}

	if exists {
		return nil, fmt.Errorf("failed to generate unique barcode after %d attempts", maxRetries)
	}

	response := types.ToGenerateBarcodeResponse(barcode)
	return &response, nil
}

func (s *stockMaterialService) RetrieveStockMaterialByBarcode(barcode string) (*types.StockMaterialsDTO, error) {
	stockMaterial, err := s.repo.GetStockMaterialByBarcode(barcode)
	if err != nil {
		return nil, err
	}
	if stockMaterial == nil {
		return nil, errors.New("StockMaterial not found with the provided barcode")
	}

	return types.ConvertStockMaterialToStockMaterialResponse(stockMaterial), nil
}
