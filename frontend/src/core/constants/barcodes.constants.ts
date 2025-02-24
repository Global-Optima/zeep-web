export const SAVED_BARISTA_BARCODE_SETTINGS_KEY = 'ZEEP_BARCODE_BARISTA_SETTINGS'

/**
 * ENUM: Названия/метки для размеров
 */
export enum BarcodePrinterSizeName {
	SIZE_100x50 = '100x50 мм',
	SIZE_50x25 = '50x25 мм',
	SIZE_80x40 = '80x40 мм',
	SIZE_60x30 = '60x30 мм',
	CUSTOM = 'Свой размер',
}

/**
 * ТИП: Конфигурация принтера (label, width, height)
 */
export interface BarcodePrinterSizeConfig {
	label: BarcodePrinterSizeName
	width: number
	height: number
}

/**
 * ПРЕДОПРЕДЕЛЕННЫЕ РАЗМЕРЫ (Thermal printer standards)
 */
export const PREDEFINED_BARCODE_SIZES: BarcodePrinterSizeConfig[] = [
	{ label: BarcodePrinterSizeName.SIZE_100x50, width: 100, height: 50 },
	{ label: BarcodePrinterSizeName.SIZE_50x25, width: 50, height: 25 },
	{ label: BarcodePrinterSizeName.SIZE_80x40, width: 80, height: 40 },
	{ label: BarcodePrinterSizeName.SIZE_60x30, width: 60, height: 30 },
]
