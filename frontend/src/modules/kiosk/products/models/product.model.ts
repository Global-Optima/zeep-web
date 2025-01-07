import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { AdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import type { IngredientsDTO } from '@/modules/admin/ingredients/models/ingredients.model'

export enum ProductSizeNames {
	S = 'S',
	M = 'M',
	L = 'L',
	XL = 'XL',
}

export enum ProductSizeMeasures {
	ML = 'мл',
	G = 'г',
	PIECE = 'шт',
}

export interface ProductsFilter extends PaginationParams {
	categoryId?: number
	search?: string
}

export interface BaseProductDTO {
	id: number
	name: string
	description: string
	imageUrl: string
	videoUrl: string
	category: ProductCategoryDTO
}

export interface ProductDTO extends BaseProductDTO {
	productSizeCount: number
	basePrice: number
}

export interface ProductSizeIngredientDTO {
	id: number
	name: string
	calories: number
	fat: number
	carbs: number
	proteins: number
}

export interface ProductSizeDTO {
	id: number
	name: string
	basePrice: number
	measure: string
	size: number
	isDefault: boolean
}

export interface ProductSizeDetailsDTO extends ProductSizeDTO {
	productId: number
	additives: ProductSizeDetailsAdditiveDTO[]
	ingredients: IngredientsDTO[]
}

export interface ProductSizeDetailsAdditiveDTO extends AdditiveDTO {
	isDefault: boolean
}

export interface ProductDetailsDTO extends BaseProductDTO {
	sizes: ProductSizeDTO[]
}

export interface CreateProductDTO {
	name: string
	description: string
	imageUrl: string
	categoryId: number
}

export interface SelectedAdditiveTypesDTO {
	additiveId: number
	isDefault: boolean
}

export interface CreateProductSizeDTO {
	productId: number
	name: string
	measure: string
	basePrice: number
	size: number
	isDefault?: boolean
	additives?: SelectedAdditiveTypesDTO[]
	ingredientIds?: number[]
}

export interface UpdateProductDTO {
	name?: string
	description?: string
	imageUrl?: string
	categoryId?: number
}

export interface UpdateProductSizeDTO {
	name?: string
	measure?: string
	basePrice?: number
	size?: number
	isDefault?: boolean
	additives?: SelectedAdditiveTypesDTO[]
	ingredientIds?: number[]
}

export interface ProductCategoryDTO {
	id: number
	name: string
	description: string
}

export interface ProductCategoriesFilterDTO {
	search?: string
	page?: number
	pageSize?: number
}

export interface CreateProductCategoryDTO {
	name: string
	description: string
}

export interface UpdateProductCategoryDTO {
	name?: string
	description?: string
}
