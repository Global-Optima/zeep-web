package export

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/tealeg/xlsx"
)

func setGreyRowStyle(row *xlsx.Row) {
	style := xlsx.NewStyle()
	style.Fill.FgColor = "D3D3D3"
	style.Fill.PatternType = "solid"

	for _, cell := range row.Cells {
		cell.SetStyle(style)
	}
}

func GenerateSalesExcelV2(data []types.OrderExportDTO, headers []string) ([]byte, error) {
	file := xlsx.NewFile()

	ordersSheet, err := file.AddSheet("Все заказы")
	if err != nil {
		return nil, err
	}

	addOrdersData(ordersSheet, data, headers)

	buffer := bytes.NewBuffer(nil)
	if err := file.Write(buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func addOrdersData(sheet *xlsx.Sheet, data []types.OrderExportDTO, headers []string) {
	headerRow := sheet.AddRow()
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
	for i, order := range data {
		for _, suborder := range order.Suborders {
			row := sheet.AddRow()
			row.AddCell().Value = fmt.Sprintf("%d", order.ID)
			row.AddCell().Value = order.CustomerName
			row.AddCell().Value = order.StoreName
			row.AddCell().Value = fmt.Sprintf("%d", suborder.ID)
			row.AddCell().Value = suborder.ProductSize.ProductName
			row.AddCell().Value = suborder.ProductSize.SizeName
			row.AddCell().Value = fmt.Sprintf("%.2f KZT", suborder.Price)

			total := suborder.Price
			var additivesDetails []string
			for _, additive := range suborder.Additives {
				total += additive.Price
				additivesDetails = append(additivesDetails, fmt.Sprintf("ID: %d %s(%.2f) - %.2f KZT", additive.Additive.ID, additive.Additive.Name, additive.Additive.Size, additive.Price))
			}
			row.AddCell().Value = fmt.Sprintf("%.2f KZT", total)
			row.AddCell().Value = strings.Join(additivesDetails, "\n")
			row.AddCell().Value = order.CreatedAt.Format("2006-01-02 15:04:05")
		}

		if i < len(data)-1 {
			greyRow := sheet.AddRow()
			for j := 0; j < 10; j++ {
				greyRow.AddCell()
			}
			setGreyRowStyle(greyRow)
		}
	}
}
