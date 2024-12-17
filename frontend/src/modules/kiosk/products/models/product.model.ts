import type { AdditiveCategoryItem } from '@/modules/admin/additives/models/additives.model'

export interface ProductsFilter {
	storeId?: number
	categoryId?: number
	searchTerm?: string
	limit?: number
	offset?: number
}

export interface Products {
	id: number
	name: string
	description: string
	imageUrl: string
	category: string
	basePrice: number
}

export interface ProductCategory {
	id: number
	name: string
	description: string
}

export interface StoreProductDetailsDTO {
	id: number
	name: string
	description: string
	imageUrl: string
	sizes: ProductSizeDTO[]
	defaultAdditives: AdditiveCategoryItem[]
}

export interface ProductSizeDTO {
	id: number
	name: string
	basePrice: number
	measure: string
}
