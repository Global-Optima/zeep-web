// ingredient.dtos.ts

import type { PaginationParams } from '@/core/utils/pagination.utils'

/**
 * DTO for creating a new ingredient.
 */
export interface CreateIngredientDTO {
	name: string
	calories: number
	fat: number
	carbs: number
	proteins: number
	expiresAt?: string 
}

/**
 * DTO for updating an existing ingredient.
 */
export interface UpdateIngredientDTO {
	id: number
	name?: string
	calories?: number
	fat?: number
	carbs?: number
	proteins?: number
	expiresAt?: string
}

/**
 * DTO for the ingredient response.
 */
export interface IngredientResponseDTO {
	id: number
	name: string
	calories: number
	fat: number
	carbs: number
	proteins: number
	expiresAt?: string
}

/**
 * Filter DTO for fetching ingredients with optional conditions.
 */
export interface IngredientFilter extends PaginationParams {
	productSizeId?: number
	name?: string
	minCalories?: number
	maxCalories?: number
}
