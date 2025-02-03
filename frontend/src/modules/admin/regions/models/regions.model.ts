import type { PaginationParams } from '@/core/utils/pagination.utils'

export interface RegionDTO {
	id: number
	name: string
}

export interface AllRegionsFilter {
	search: string
}

export interface CreateRegionDTO {
	name: string
}

export interface UpdateRegionDTO {
	name?: string
}

export interface RegionFilterDTO extends PaginationParams {
	search?: string
}
