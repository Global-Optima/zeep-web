import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { IngredientsDTO } from '../../ingredients/models/ingredients.model'
import type { UnitDTO } from '../../units/models/units.model'
// Filters
export interface AdditiveCategoriesFilterQuery extends PaginationParams {
	includeEmpty?: boolean
	productSizeId?: number
	isMultipleSelect?: boolean
	search?: string
}

export interface AdditiveFilterQuery extends PaginationParams {
	search?: string
	minPrice?: number
	maxPrice?: number
	categoryId?: number
	productSizeId?: number
	storeId?: number
}

// Base DTOs
export interface BaseAdditiveCategoryDTO {
	id: number
	name: string
	description: string
	isMultipleSelect: boolean
}

export interface BaseAdditiveDTO {
	name: string
	description: string
	basePrice: number
	imageUrl: string
	size: number
	unit: UnitDTO
	category: BaseAdditiveCategoryDTO
}

// Additive DTOs
export interface AdditiveDTO extends BaseAdditiveDTO {
	id: number
}

export interface AdditiveDetailsDTO extends AdditiveDTO {
	ingredients: SelectedDetailedIngredientDTO[]
}

export interface SelectedDetailedIngredientDTO {
	ingredient: IngredientsDTO
	quantity: number
}

export interface BaseAdditiveCategoryItemDTO {
	name: string
	description: string
	basePrice: number
	imageUrl: string
	size: number
	unit: UnitDTO
	categoryId: number
}

export interface AdditiveCategoryItemDTO extends BaseAdditiveCategoryItemDTO {
	id: number
}

export interface AdditiveCategoryDTO extends BaseAdditiveCategoryDTO {
	additives: AdditiveCategoryItemDTO[]
}

// Create and Update DTOs
export interface CreateAdditiveCategoryDTO {
	name: string
	description?: string
	isMultipleSelect: boolean
}

export interface UpdateAdditiveCategoryDTO {
	name?: string
	description?: string
	isMultipleSelect?: boolean
}

export interface UpdateAdditiveDTO {
	name?: string
	description?: string
	basePrice?: number
	imageUrl?: string
	size?: number
	unitId?: number
	additiveCategoryId?: number
	ingredients?: SelectedIngredientDTO[]
}

export interface AdditiveCategoryResponseDTO {
	id: number
	name: string
	description: string
	isMultipleSelect: boolean
}

export interface CreateAdditiveDTO {
	name: string
	description: string
	basePrice: number
	imageUrl?: string
	size: number
	unitId: number
	additiveCategoryId: number
	ingredients: SelectedIngredientDTO[]
	image?: File
}

// Ingredient DTOs
export interface SelectedIngredientDTO {
	ingredientId: number
	quantity: number
}
