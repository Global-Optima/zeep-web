import type { PaginationParams } from '@/core/utils/pagination.utils'

export interface CreateUnitDTO {
	name: string
	conversionFactor: number
}

export interface UpdateUnitDTO {
	name?: string
	conversionFactor?: number
}

export interface UnitResponse {
	id: number
	name: string
	conversionFactor: number
}

export interface UnitsFilterDTO extends PaginationParams {
	search?: string
}
