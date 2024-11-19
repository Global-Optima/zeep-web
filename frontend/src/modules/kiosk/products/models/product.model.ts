export interface StoreProducts {
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

// src/modules/kiosk/products/models/store-product-details.dto.ts

export interface StoreProductDetailsDTO {
	id: number
	name: string
	description: string
	imageUrl: string
	sizes: ProductSizeDTO[]
	defaultAdditives: AdditiveDTO[]
}

export interface ProductSizeDTO {
	id: number
	name: string
	basePrice: number
	measure: string
}

export interface AdditiveDTO {
	id: number
	name: string
	description: string
	price: number
	imageUrl: string
	categoryId: number
}

export interface AdditiveCategoryDTO {
	id: number
	name: string
	additives: AdditiveDTO[]
	isMultipleSelect: boolean
}
