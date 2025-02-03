import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { RegionDTO } from '@/modules/admin/regions/models/regions.model'

export interface WarehouseDTO {
	id: number
	name: string
	facilityAddress: {
		id: number
		address: string
		longitude: number
		latitude: number
	}
	region: RegionDTO
	createdAt: Date
}

export interface WarehouseFilter extends PaginationParams {
	search?: string
	regionId?: number
}
