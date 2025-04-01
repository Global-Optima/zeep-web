// ingredient.dtos.ts

import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { UnitDTO } from '../../units/models/units.model'

/**
 * DTO for creating a new ingredient.
 */
export interface CreateIngredientDTO {
	name: string
	calories: number
	fat: number
	carbs: number
	proteins: number
	categoryId: number
	unitId: number
	expirationInDays: number
	isAllergen: boolean
}

/**
 * DTO for updating an existing ingredient.
 */
export interface UpdateIngredientDTO {
	name?: string
	calories?: number
	fat?: number
	carbs?: number
	proteins?: number
	categoryId?: number
	unitId?: number
	expirationInDays?: number
	isAllergen?: boolean
}

/**
 * DTO for the ingredient response.
 */
export interface IngredientsDTO {
	id: number
	name: string
	calories: number
	fat: number
	carbs: number
	proteins: number
	expirationInDays: number
	isAllergen: boolean
	unit: UnitDTO
	category: IngredientCategoryDTO
}

/**
 * Filter DTO for fetching ingredients with optional conditions.
 */
export interface IngredientFilter extends PaginationParams {
	productSizeId?: number
	name?: string
	minCalories?: number
	maxCalories?: number
	isAllergen?: boolean
}

export interface CreateIngredientCategoryDTO {
	name: string
	description?: string
}

export interface UpdateIngredientCategoryDTO {
	name?: string
	description?: string
}

export interface IngredientCategoryDTO {
	id: number
	name: string
	description: string
}

export interface IngredientCategoryFilter extends PaginationParams {
	search?: string
}
