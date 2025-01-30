package export

import "github.com/tealeg/xlsx"

var KazHeaders = []string{"Тапсырыс нөмірі", "Тапсырыс берушінің аты", "Филиал атауы", "Тапсырыс тауар нөмірі", "Тауар атауы", "Өлшемі", "Бағасы", "Жалпы баға (қоспалардың бағасын қоса)", "Қоспалар", "Тапсырыс күні"}
var RusHeaders = []string{"Номер заказа", "Имя покупателя", "Название филиала", "Номер подзаказа", "Имя продукта", "Размер", "Цена", "Итого (с учетом цены добавок)", "Добавки", "Дата заказа"}
var EngHeaders = []string{"Order ID", "Customer Name", "Store Name", "Suborder ID", "Product Name", "Size", "StorePrice", "Total (with additive prices added)", "Additives", "Order Date"}

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
		err := sheet.SetColWidth(i, i, 30)
		if err != nil {
			return err
		}
	}

	return nil
}
