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
	SubOrders []PDFSubOrder
}

type PDFSubOrder struct {
	ProductName string
	Size        string
	Price       float64
	Status      string
	Additives   []PDFAdditive
}

type PDFAdditive struct {
	Name  string
	Price float64
}

// GeneratePDFReceipt generates a receipt PDF for an order
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
	pdf.Cell(0, 10, "Suborders")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)

	for _, suborder := range details.SubOrders {
		// Display product information
		pdf.Cell(120, 10, fmt.Sprintf("%s (Size: %s)", suborder.ProductName, suborder.Size))
		pdf.Cell(40, 10, fmt.Sprintf("%.2f", suborder.Price))
		pdf.Ln(8)

		// Display suborder status
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(160, 10, fmt.Sprintf("Status: %s", suborder.Status))
		pdf.Ln(6)

		// Display additives, if any
		if len(suborder.Additives) > 0 {
			pdf.SetFont("Arial", "I", 10)
			for _, additive := range suborder.Additives {
				pdf.Cell(120, 10, fmt.Sprintf("  + %s", additive.Name))
				pdf.Cell(40, 10, fmt.Sprintf("%.2f", additive.Price))
				pdf.Ln(6)
			}
			pdf.SetFont("Arial", "", 12)
		}

		pdf.Ln(4) // Add spacing between suborders
	}

	// Write PDF to a buffer
	var buffer bytes.Buffer
	err := pdf.Output(&buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buffer.Bytes(), nil
}
