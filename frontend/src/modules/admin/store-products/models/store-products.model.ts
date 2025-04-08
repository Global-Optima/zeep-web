import type { PaginationParams } from '@/core/utils/pagination.utils'
import type {
	BaseProductDTO,
	BaseProductSizeDTO,
	ProductSizeAdditiveDTO,
	ProductTotalNutrition,
} from '@/modules/kiosk/products/models/product.model'
import type { IngredientsDTO } from '../../ingredients/models/ingredients.model'
import type { ProvisionDTO } from '@/modules/admin/provisions/models/provision.models'

export interface StoreProductDTO extends BaseProductDTO {
	id: number
	productId: number
	storePrice: number
	basePrice: number
	storeProductSizeCount: number
	productSizeCount: number
	isAvailable: boolean
	isOutOfStock: boolean
}

export interface StoreProductSizeDTO extends BaseProductSizeDTO {
	id: number
	totalNutrition: ProductTotalNutrition
	productSizeId: number
	storePrice: number
}

export interface StoreProductDetailsDTO extends StoreProductDTO {
	sizes: StoreProductSizeDetailsDTO[]
}

export interface StoreProductSizeDetailsDTO extends StoreProductSizeDTO {
	additives: ProductSizeAdditiveDTO[]
	ingredients: ProductSizeIngredientDTO[]
}

export interface ProductSizeIngredientDTO {
	quantity: number
	ingredient: IngredientsDTO
}

export interface ProductSizeProvisionDTO {
	volume: number
	provision: ProvisionDTO
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
	isOutOfStock?: boolean
	isAvailable?: boolean
	search?: string
	maxPrice?: number
	minPrice?: number
	storeId?: number
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
