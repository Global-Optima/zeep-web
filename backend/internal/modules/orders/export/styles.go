package export

import "github.com/tealeg/xlsx"

func setHeadersStyle(headerRow *xlsx.Row) {
	style := xlsx.NewStyle()
	style.Font.Bold = true
	style.Fill.FgColor = "C6C6C6"
	style.Fill.PatternType = "solid"

	for _, cell := range headerRow.Cells {
		cell.SetStyle(style)
	}
}

func setColumnWidths(sheet *xlsx.Sheet) error {
	for i := range len(sheet.Cols) {
		err := sheet.SetColWidth(i, i, 50)
		if err != nil {
			return err
		}
	}

	return nil
}
