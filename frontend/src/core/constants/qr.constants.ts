export const SAVED_BARISTA_QR_SETTINGS_KEY = 'ZEEP_QR_BARISTA_SETTINGS'

/**
 * ENUM: Названия/метки для размеров (QR Standard Sizes)
 */
export enum QRPrinterSizeName {
	SIZE_100x100 = '100x100 мм', // Large QR for product labeling
	SIZE_80x80 = '80x80 мм', // Standard for POS & logistics
	SIZE_60x60 = '60x60 мм', // Common for packaging
	SIZE_40x40 = '40x40 мм', // Compact for receipts & tickets
	CUSTOM = 'Свой размер',
}

/**
 * ТИП: Конфигурация принтера (label, width, height)
 */
export interface QRPrinterSizeConfig {
	label: QRPrinterSizeName
	width: number
	height: number
}

/**
 * ПРЕДОПРЕДЕЛЕННЫЕ РАЗМЕРЫ (QR Code Label Standards)
 * Based on common QR label printing dimensions.
 */
export const PREDEFINED_QR_SIZES: QRPrinterSizeConfig[] = [
	{ label: QRPrinterSizeName.SIZE_100x100, width: 100, height: 100 }, // For detailed QR data
	{ label: QRPrinterSizeName.SIZE_80x80, width: 80, height: 80 }, // Logistics/warehouse
	{ label: QRPrinterSizeName.SIZE_60x60, width: 60, height: 60 }, // Product packaging
	{ label: QRPrinterSizeName.SIZE_40x40, width: 40, height: 40 }, // POS receipts/tickets
]
