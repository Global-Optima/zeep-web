package pdf

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

type PDFReceiptDetails struct {
	OrderID   uint
	StoreID   uint
	OrderDate string
	Total     float64
	Products  []PDFProduct
	FilePath  string
}

type PDFProduct struct {
	ProductName string
	Size        string
	Quantity    int
	Price       float64
	Additives   []PDFAdditive
}

type PDFAdditive struct {
	Name  string
	Price float64
}

func GeneratePDFReceipt(details PDFReceiptDetails) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, fmt.Sprintf("Receipt for Order #%d", details.OrderID))
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Store ID: %d", details.StoreID))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Order Date: %s", details.OrderDate))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Total: %.2f", details.Total))
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Items")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)

	for _, product := range details.Products {
		pdf.Cell(100, 10, fmt.Sprintf("%s (Size: %s)", product.ProductName, product.Size))
		pdf.Cell(20, 10, fmt.Sprintf("x%d", product.Quantity))
		pdf.Cell(40, 10, fmt.Sprintf("%.2f", product.Price))
		pdf.Ln(8)

		if len(product.Additives) > 0 {
			pdf.SetFont("Arial", "I", 10)
			for _, additive := range product.Additives {
				pdf.Cell(120, 10, fmt.Sprintf("  + %s", additive.Name))
				pdf.Cell(40, 10, fmt.Sprintf("%.2f", additive.Price))
				pdf.Ln(6)
			}
			pdf.SetFont("Arial", "", 12)
		}
	}

	// Write PDF to a buffer
	var buffer bytes.Buffer
	err := pdf.Output(&buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buffer.Bytes(), nil
}
