import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { IngredientsDTO } from '../../ingredients/models/ingredients.model'
import type { UnitDTO } from '../../units/models/units.model'
import type { ProvisionDTO } from '../../provisions/models/provision.models'
import type { SelectedProvisionDTO } from '@/modules/kiosk/products/models/product.model'
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
	name: string
	description: string
	isMultipleSelect: boolean
	isRequired: boolean
}

export interface AdditiveCategoryDTO extends BaseAdditiveCategoryDTO {
	id: number
}

export interface BaseAdditiveDTO {
	name: string
	description: string
	basePrice: number
	imageUrl: string
	size: number
	unit: UnitDTO
	category: AdditiveCategoryDTO
	machineId: string
}

// Additive DTOs
export interface AdditiveDTO extends BaseAdditiveDTO {
	id: number
}

export interface AdditiveDetailsDTO extends AdditiveDTO {
	ingredients: SelectedDetailedIngredientDTO[]
	provisions: SelectedDetailedProvisionsDTO[]
}

export interface SelectedDetailedIngredientDTO {
	ingredient: IngredientsDTO
	quantity: number
}

export interface SelectedDetailedProvisionsDTO {
	provision: ProvisionDTO
	volume: number
}

export interface AdditiveCategoryDetailsDTO extends AdditiveCategoryDTO {
	additivesCount: number
}

// Create and Update DTOs
export interface CreateAdditiveCategoryDTO {
	name: string
	description?: string
	isMultipleSelect: boolean
	isRequired: boolean
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
	size?: number
	unitId?: number
	additiveCategoryId?: number
	machineId?: string
	ingredients?: SelectedIngredientDTO[]
	provisions?: SelectedProvisionDTO[]
	image?: File
	deleteImage: boolean
}

export interface CreateAdditiveDTO {
	name: string
	description?: string
	basePrice: number
	imageUrl?: string
	size: number
	unitId: number
	additiveCategoryId: number
	machineId: string
	ingredients: SelectedIngredientDTO[]
	provisions?: SelectedProvisionDTO[]
	image?: File
}

// Ingredient DTOs
export interface SelectedIngredientDTO {
	ingredientId: number
	quantity: number
}
