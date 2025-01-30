package export

import (
	"bytes"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/tealeg/xlsx"
)

func GenerateSalesExcel(data []types.OrderExportDTO) ([]byte, error) {
	file := xlsx.NewFile()

	ordersSheet, err := file.AddSheet("Orders List")
	if err != nil {
		return nil, err
	}

	addOrdersSheet(ordersSheet, data, file)

	buffer := bytes.NewBuffer(nil)
	if err := file.Write(buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func addOrdersSheet(sheet *xlsx.Sheet, data []types.OrderExportDTO, file *xlsx.File) {
	headerRow := sheet.AddRow()
	headers := []string{"Order ID", "Store Name", "Customer Name", "Total", "Status", "Order Date"}
	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	err := setColumnWidths(sheet)
	if err != nil {
		logger.GetZapSugaredLogger().Errorln(err.Error())
		return
	}
	setHeadersStyle(headerRow)

	for _, order := range data {
		row := sheet.AddRow()
		addCell(row, fmt.Sprintf("%d", order.ID))
		addCell(row, order.StoreName)
		addCell(row, order.CustomerName)
		addCell(row, fmt.Sprintf("%.2f", order.Total))
		addCell(row, order.Status)
		addCell(row, order.CreatedAt.Format("2006-01-02 15:04:05"))

		productSheetName := fmt.Sprintf("Order_%d_Products", order.ID)
		productSheet, _ := file.AddSheet(productSheetName)
		addProductSheet(productSheet, order.Suborders, file)
	}
}

func addProductSheet(sheet *xlsx.Sheet, suborders []types.SuborderDTO, file *xlsx.File) {
	headerRow := sheet.AddRow()
	headers := []string{"Product Name", "Size", "StorePrice"}
	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	err := setColumnWidths(sheet)
	if err != nil {
		logger.GetZapSugaredLogger().Errorln(err.Error())
		return
	}
	setHeadersStyle(headerRow)

	for _, suborder := range suborders {
		row := sheet.AddRow()
		addCell(row, suborder.ProductSize.ProductName)
		addCell(row, suborder.ProductSize.SizeName)
		addCell(row, fmt.Sprintf("%.2f", suborder.Price))

		if len(suborder.Additives) != 0 {
			additiveSheetName := fmt.Sprintf("Suborder_%d_Additives", suborder.ID)
			additiveSheet, _ := file.AddSheet(additiveSheetName)
			addAdditiveSheet(additiveSheet, suborder.Additives)
		}
	}
}

func addAdditiveSheet(sheet *xlsx.Sheet, additives []types.SuborderAdditiveDTO) {
	headerRow := sheet.AddRow()
	headers := []string{"Additive Name", "StorePrice"}
	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	err := setColumnWidths(sheet)
	if err != nil {
		logger.GetZapSugaredLogger().Errorln(err.Error())
		return
	}
	setHeadersStyle(headerRow)

	for _, additive := range additives {
		row := sheet.AddRow()
		addCell(row, additive.Additive.Name)
		addCell(row, fmt.Sprintf("%.2f", additive.Price))
	}
}

func addCell(row *xlsx.Row, value string) {
	cell := row.AddCell()
	cell.Value = value
}
