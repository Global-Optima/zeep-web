import type { AdditiveCategoryItemDTO, AdditiveDTO } from '../../additives/models/additives.model'

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
}

export interface StoreAdditiveCategoriesFilter {
	search?: string
	isMultipleSelect?: boolean
}

export interface StoreAdditiveCategoryItemDTO extends AdditiveCategoryItemDTO {
	storeAdditiveId: number
	storePrice: number
	isDefault: boolean
}

export interface StoreAdditiveCategoryDTO {
	id: number
	name: string
	description: string
	additives: StoreAdditiveCategoryItemDTO[]
	isMultipleSelect: boolean
}
