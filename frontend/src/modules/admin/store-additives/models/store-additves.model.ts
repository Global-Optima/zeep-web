import type {
  AdditiveDTO,
  BaseAdditiveDTO,
} from '../../additives/models/additives.model'

export interface CreateStoreAdditiveDTO {
	additiveId: number
	storePrice: number
}

export interface UpdateStoreAdditiveDTO {
	storePrice: number
}

export interface StoreAdditiveDTO extends AdditiveDTO {
	additiveId: number
	storePrice: number
  isOutOfStock: boolean
}

export interface StoreAdditiveCategoriesFilter {
	search?: string
	isMultipleSelect?: boolean
}

export interface StoreAdditiveCategoryItemDTO extends BaseAdditiveDTO {
	id: number
	additiveId: number
	storePrice: number
  isOutOfStock: boolean
	isDefault: boolean
  isHidden: boolean
}

export interface StoreAdditiveCategoryDTO {
	id: number
	name: string
	description: string
	additives: StoreAdditiveCategoryItemDTO[]
	isMultipleSelect: boolean
}
