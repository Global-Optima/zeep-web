import type { PaginationParams } from '@/core/utils/pagination.utils'

export interface StoresFilter extends PaginationParams {
	search?: string
	franchiseeId?: number
}

export interface CreateStoreDTO {
	name: string
	franchiseeId?: number
	warehouseId: number
	facilityAddress: {
		address: string
	}
	isActive: boolean
	contactPhone: string
	contactEmail: string
	storeHours: string
}

export interface UpdateStoreDTO {
	name: string
	franchiseeId?: number | null
	warehouseId?: number
	facilityAddress: {
		address: string
	}
	isActive: boolean
	contactPhone: string
	contactEmail: string
	storeHours: string
}
