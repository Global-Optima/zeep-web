import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { TranslationLocalizedField } from '@/modules/kiosk/products/models/product.model'

export interface CreateUnitDTO {
	name: string
	conversionFactor: number
}

export interface UpdateUnitDTO {
	name?: string
	conversionFactor?: number
}

export interface UnitDTO {
	id: number
	name: string
	conversionFactor: number
}

export interface UnitsFilterDTO extends PaginationParams {
	search?: string
}

export interface UnitTranslationsDTO {
	name?: TranslationLocalizedField
}
