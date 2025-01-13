import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { AdditiveDTO } from '../../additives/models/additives.model'

export interface CreateStoreAdditiveDTO {
	additiveId: number
	storePrice: number
}

export interface UpdateStoreAdditiveDTO {
	storePrice: number
}

export interface StoreAdditiveDTO extends AdditiveDTO {
	storeAdditiveId: number
	storePrice: number
	id: number
	name: string
	description: string
}

export interface StoreAdditiveCategoryItemDTO {
	storeAdditiveId: number
	storePrice: number
	id: number
	name: string
}

export interface StoreAdditiveCategoryDTO {
	id: number
	name: string
	description: string
	additives: StoreAdditiveCategoryItemDTO[]
	isMultipleSelect: boolean
}

export interface StoreAdditivesFilterDTO extends PaginationParams {
	search?: string
}
