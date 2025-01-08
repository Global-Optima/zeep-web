import type { PaginationParams } from '@/core/utils/pagination.utils'
import type {
	ProductDTO,
	ProductSizeDetailsDTO,
} from '@/modules/kiosk/products/models/product.model'

export interface StoreProductDTO extends ProductDTO {
	storeProductId: number
	storePrice: number
	storeProductSizeCount: number
	isAvailable: boolean
}

export interface StoreProductDetailsDTO extends StoreProductDTO {
	sizes: StoreProductSizeDetailsDTO[]
}

export interface StoreProductSizeDetailsDTO extends ProductSizeDetailsDTO {
	storePrice: number
}

export interface CreateStoreProductDTO {
	productId: number
	isAvailable: boolean
	productSizes?: CreateStoreProductSizeDTO[]
}

export interface CreateStoreProductSizeDTO {
	productSizeID: number
	storePrice?: number
}

export interface UpdateStoreProductDTO {
	isAvailable?: boolean
	productSizes?: UpdateStoreProductSizeDTO[]
}

export interface UpdateStoreProductSizeDTO {
	productSizeID: number
	storePrice?: number
}

export interface StoreProductsFilterDTO extends PaginationParams {
	categoryId?: number
	isAvailable?: boolean
	search?: string
	maxPrice?: number
	minPrice?: number
}

export interface StoreProductSizesFilterDTO extends PaginationParams {
	categoryId?: number
	name?: string
	measure?: string
	search?: string
	isDefault?: boolean
	minSize?: number
	maxSize?: number
}
