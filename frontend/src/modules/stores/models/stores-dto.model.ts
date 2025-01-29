import type { PaginationParams } from '@/core/utils/pagination.utils'

export interface StoresFilter extends PaginationParams {
	searchTerm?: string
	isFranchise?: boolean
}

export interface CreateStoreDTO {
	name: string
	isFranchise: boolean
	facilityAddress: {
		address: string
		longitude?: number
		latitude?: number
	}
	contactPhone: string
	contactEmail: string
	storeHours: string
}

export interface UpdateStoreDTO {
	name: string
	isFranchise: boolean
	facilityAddress: {
		address: string
		longitude?: number
		latitude?: number
	}
	contactPhone: string
	contactEmail: string
	storeHours: string
}
